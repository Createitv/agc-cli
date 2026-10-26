"""专用 .app 的真实签名、上传和异步解析；不提交审核或发布。

证书、Profile 与指纹仅使用本次创建的资源。私钥和密码不写入报告。
上传后的未关联软件包保留，签名资源在结束时清理并逐项读回。
"""
import argparse
import io
import datetime
import hashlib
import importlib.util
import json
import pathlib
import re
import subprocess
import tempfile
import time
import urllib.request
import uuid
import zipfile
import os
import shutil
import ssl

_certificate_spec = importlib.util.spec_from_file_location("certificate_fingerprint", pathlib.Path(__file__).with_name("certificate-fingerprint.py"))
_certificate_module = importlib.util.module_from_spec(_certificate_spec)
_certificate_spec.loader.exec_module(_certificate_module)
certificate_fingerprint_for_csr = _certificate_module.certificate_fingerprint_for_csr


ROOT = pathlib.Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("lifecycle", pathlib.Path(__file__).with_name("live-lifecycle.py"))
lifecycle = importlib.util.module_from_spec(spec)
spec.loader.exec_module(lifecycle)
JAVA = "/Applications/DevEco-Studio.app/Contents/jbr/Contents/Home/bin/java"
SIGNER = "/Applications/DevEco-Studio.app/Contents/sdk/default/openharmony/toolchains/lib/hap-sign-tool.jar"
PACKER = "/Applications/DevEco-Studio.app/Contents/sdk/default/openharmony/toolchains/lib/app_packing_tool.jar"


def inspect_app_modules(source, expected_bundle):
    """逐个核对内层 HAP 的应用身份，外层 pack.info 不能替代内层检查。"""
    report = []
    with zipfile.ZipFile(source) as archive:
        for item in archive.infolist():
            if not item.filename.endswith('.hap'):
                continue
            with zipfile.ZipFile(io.BytesIO(archive.read(item))) as hap:
                config = json.loads(hap.read('module.json'))
                app = config.get('app', {})
                if app.get('bundleName') != expected_bundle:
                    raise RuntimeError('Inner HAP bundle does not match dedicated application')
                report.append({'moduleName': config.get('module', {}).get('name'),
                               'bundleMatchesDedicated': True, 'debug': app.get('debug'),
                               'minAPIVersion': app.get('minAPIVersion'),
                               'targetAPIVersion': app.get('targetAPIVersion')})
    if not report:
        raise RuntimeError('Dedicated APP contains no HAP modules')
    return report


def retain_signed_artifact(source, directory):
    """仅复制最终 APP，不复制私钥、密码文件或临时签名目录。"""
    if source.suffix != ".app":
        raise RuntimeError("Only final APP artifacts can be retained")
    directory.mkdir(mode=0o700, parents=True, exist_ok=True)
    if directory.is_symlink():
        raise RuntimeError("Artifact directory must not be a symlink")
    target = directory / source.name
    with source.open("rb") as reader, os.fdopen(os.open(target, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600), "wb") as writer:
        shutil.copyfileobj(reader, writer)
    digest = hashlib.sha256(target.read_bytes()).hexdigest()
    if digest != hashlib.sha256(source.read_bytes()).hexdigest():
        raise RuntimeError("Retained signed APP byte verification failed")
    return {"path": str(target), "sha256": digest, "bytes": target.stat().st_size}



def run_sdk(arguments, label, tool=SIGNER):
    result = subprocess.run([JAVA, "-jar", tool, *arguments], capture_output=True, text=True, timeout=90)
    output = result.stdout + result.stderr
    # SDK 部分工具报错仍退出 0；错误日志不能当作签名/打包成功。
    if result.returncode or re.search(r"\[Error\]|ERROR:\s*\d+", output, re.IGNORECASE):
        match = re.search(r"(?:ERROR:\s*|code:\s*)(\d+)", output, re.IGNORECASE)
        raise RuntimeError(label + " failed" + ("; SDK code " + match[1] if match else ""))


def download(url, target):
    if not isinstance(url, str) or not url.startswith("https://"):
        raise RuntimeError("Missing HTTPS signing resource download")
    with urllib.request.urlopen(url, timeout=30) as response:
        content = response.read(1024 * 1024)
        if response.status != 200 or not content:
            raise RuntimeError("Signing resource download failed")
    target.write_bytes(content)


def package_state(body, package_id):
    states = [p for p in body.get("pkgStateList", []) if p.get("pkgId") == package_id]
    if len(states) != 1 or type(states[0].get("successStatus")) is not int or states[0]["successStatus"] not in (0, 1, 2):
        raise RuntimeError("Missing unambiguous owned package parse state")
    return states[0]["successStatus"]


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", required=True, type=pathlib.Path)
    parser.add_argument("--profile", required=True)
    parser.add_argument("--compile-profile", required=True)
    parser.add_argument("--sign-compatible-version", choices=("24", "sdk-default"), default="24", help="Controlled signing comparison; sdk-default omits the override (installed SDK default is 9).")
    parser.add_argument("--app-id", required=True)
    parser.add_argument("--artifact-dir", type=pathlib.Path, default=pathlib.Path(tempfile.gettempdir()) / "agc-signed-package-artifacts")
    parser.add_argument("--output-dir", type=pathlib.Path, default=ROOT / "docs/verification/lifecycles")
    args = parser.parse_args()
    args.binary = args.binary.resolve()
    fixture = json.loads((ROOT / "docs/verification/dedicated-test-app.json").read_text())
    if not lifecycle.is_dedicated_app(args.app_id, fixture):
        raise SystemExit("Verified dedicated application required")
    builds = [json.loads(p.read_text()) for p in args.output_dir.glob("*-dedicated-app-package-build.json")]
    builds = [b for b in builds if b.get("status") == "verified-unsigned-app" and b.get("appId") == args.app_id and b.get("bundleName") == fixture["packageName"]]
    if not builds:
        raise SystemExit("Verified dedicated unsigned APP build required")
    build = max(builds, key=lambda b: b["checkedAt"])
    artifact = build["artifacts"][0]
    source = pathlib.Path(artifact["path"])
    if hashlib.sha256(source.read_bytes()).hexdigest() != artifact["sha256"]:
        raise SystemExit("Build artifact changed since verification")
    with zipfile.ZipFile(source) as archive:
        if json.loads(archive.read("pack.info"))["summary"]["app"]["bundleName"] != fixture["packageName"]:
            raise SystemExit("APP bundle does not match dedicated application")
    module_inspection = inspect_app_modules(source, fixture["packageName"])
    run_id = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ") + "-" + uuid.uuid4().hex[:8]
    flow = lifecycle.Flow(args, "signed-test-package", run_id)
    flow.ledger["sourceModuleInspection"] = module_inspection
    flow.ledger["signingCompatibleVersion"] = args.sign_compatible_version
    compatibility_args = [] if args.sign_compatible_version == "sdk-default" else ["-compatibleVersion", args.sign_compatible_version]
    flow.save()
    headers = {"appId": args.app_id}
    cert_id = profile_id = fingerprint = None
    fingerprint_added = False
    with tempfile.TemporaryDirectory(prefix="agc-package-signing-") as directory:
        tmp = pathlib.Path(directory)
        try:
            # 缓存的 fixture 不能证明应用现在可访问；先读回，避免为失效应用创建全局证书。
            flow.call("publishing", "app-info-query-v3", query={"appId": args.app_id}, assert_fn=lambda b: b.get("appInfo", {}).get("parentType") == 13 and str(b["appInfo"].get("projectId")) == str(fixture.get("projectId")), assertion="dedicated application live preflight matches fixture project")
            flow.call("provisioning", "provision-api-get-fingerprints", headers=headers, assertion="dedicated signing API live preflight")
            password = uuid.uuid4().hex
            store, csr = tmp / "private.p12", tmp / "public.csr"
            run_sdk(["generate-keypair", "-keyAlias", "agcverify", "-keyAlg", "ECC", "-keySize", "NIST-P-256", "-keystoreFile", str(store), "-keystorePwd", password, "-keyPwd", password], "Key generation")
            store.chmod(0o600)
            run_sdk(["generate-csr", "-keyAlias", "agcverify", "-keyPwd", password, "-subject", "C=CN,O=AGC CLI Verification,CN=AGC CLI Verification", "-signAlg", "SHA256withECDSA", "-keystoreFile", str(store), "-keystorePwd", password, "-outFile", str(csr)], "CSR generation")
            certificate = flow.call("provisioning", "provision-api-apply-cent", body={"csr": csr.read_text(), "certName": "AGC CLI 软件包 " + run_id[-8:], "certType": 2}, assert_fn=lambda b: bool(b.get("certInfo", {}).get("id")), assertion="new dedicated release certificate ID")
            cert_id = certificate["certInfo"]["id"]
            flow.own("certId", cert_id)
            cert_file = tmp / "public.cer"
            download(certificate["certInfo"].get("certDownloadUrl"), cert_file)
            fingerprint = certificate_fingerprint_for_csr(cert_file, csr)
            flow.ledger["signerCertificateMatchesCSR"] = True
            flow.save()
            initial = flow.call("provisioning", "provision-api-get-fingerprints", headers=headers)
            if fingerprint in json.dumps(initial):
                raise RuntimeError("New certificate unexpectedly already registered")
            flow.own("fingerprint", fingerprint)
            flow.call("provisioning", "provision-api-add-fingerprints", headers=headers, body={"fingerprintList": [fingerprint]})
            fingerprint_added = True
            flow.call("provisioning", "provision-api-get-fingerprints", headers=headers, assert_fn=lambda b: fingerprint in json.dumps(b), assertion="new package signer fingerprint read back")
            provision = flow.call("provisioning", "provision-api-apply-provision", body={"provisionName": "AGC CLI 软件包 " + run_id[-8:], "provisionType": 2, "certId": cert_id, "appId": args.app_id}, assert_fn=lambda b: bool(b.get("provisionInfo", {}).get("id")), assertion="new dedicated release Profile ID")
            profile_id = provision["provisionInfo"]["id"]
            flow.own("provisionId", profile_id)
            profile_file = tmp / "release.p7b"
            download(provision["provisionInfo"].get("provisionDownloadUrl"), profile_file)
            signed = tmp / ("agc-cli-" + run_id + ".app")
            hap_metadata = []
            signed_haps = []
            with zipfile.ZipFile(source) as original:
                for item in original.infolist():
                    content = original.read(item.filename)
                    if item.filename.endswith(".hap"):
                        if pathlib.PurePosixPath(item.filename).name != item.filename:
                            raise RuntimeError("Unexpected non-root HAP entry")
                        unsigned_hap, signed_hap = tmp / ("unsigned-" + item.filename), tmp / item.filename
                        unsigned_hap.write_bytes(content)
                        run_sdk(["sign-app", "-mode", "localSign", "-keyAlias", "agcverify", "-keyPwd", password, "-appCertFile", str(cert_file), "-profileFile", str(profile_file), "-signAlg", "SHA256withECDSA", "-keystoreFile", str(store), "-keystorePwd", password, "-inFile", str(unsigned_hap), "-outFile", str(signed_hap), *compatibility_args], "HAP signing")
                        run_sdk(["verify-app", "-inFile", str(signed_hap), "-outCertChain", str(tmp / "verified-chain.cer"), "-outProfile", str(tmp / "verified-profile.p7b")], "HAP signature verification")
                        content = signed_hap.read_bytes()
                        hap_metadata.append({"fileName": item.filename, "sha256": hashlib.sha256(content).hexdigest(), "bytes": len(content), "sdkSignatureVerified": True})
                        signed_haps.append(str(signed_hap))
                    elif item.filename in ("pack.info", "pac.json"):
                        (tmp / item.filename).write_bytes(content)
            if not hap_metadata:
                raise RuntimeError("APP contains no signed HAP")
            # 源 APP 内的 pack.info 已由 Hvigor 生成。保留它，避免重写已经签名的 HAP。
            run_sdk(["--mode", "app", "--pack-info-path", str(tmp / "pack.info"), "--hap-path", ",".join(signed_haps), "--out-path", str(signed), "--force", "true", "--pac-json-path", str(tmp / "pac.json"), "--replace-pack-info", "false"], "APP packing", tool=PACKER)
            with zipfile.ZipFile(signed) as packed:
                if packed.testzip() is not None or json.loads(packed.read("pack.info"))["summary"]["app"]["bundleName"] != fixture["packageName"]:
                    raise RuntimeError("Packed APP failed independent bundle/ZIP verification")
                for hap in hap_metadata:
                    content = packed.read(hap["fileName"])
                    packed_hap = tmp / ("packed-" + hap["fileName"])
                    packed_hap.write_bytes(content)
                    # 打包工具可能调整包内元信息；必须对最终字节重新验签，不能沿用打包前的结果。
                    run_sdk(["verify-app", "-inFile", str(packed_hap), "-outCertChain", str(tmp / "packed-chain.cer"), "-outProfile", str(tmp / "packed-profile.p7b")], "Packed HAP signature verification")
                    hap["beforePackingSHA256"] = hap["sha256"]
                    hap["sha256"] = hashlib.sha256(content).hexdigest()
                    hap["bytes"] = len(content)
                    hap["finalSDKSignatureVerified"] = True
            flow.ledger["signedPackage"] = {"bundleName": fixture["packageName"], "sha256": hashlib.sha256(signed.read_bytes()).hexdigest(), "bytes": signed.stat().st_size, "haps": hap_metadata}
            flow.ledger["signedArtifact"] = retain_signed_artifact(signed, args.artifact_dir.absolute())
            flow.save()
            address = flow.call("upload", "upload-url-new", query={"appId": args.app_id, "suffix": "app", "fileName": signed.name, "contentLength": signed.stat().st_size})["urlInfo"]
            if address.get("method", "").upper() != "PUT" or not address.get("url", "").startswith("https://"):
                raise RuntimeError("Invalid package upload address")
            command = [str(args.binary), "--profile", args.profile, "upload", "upload-file-new", "--invoke", "--dry-run=false", "--param", "uploadUrl=" + address["url"], "--file", str(signed)]
            for key, value in address.get("headers", {}).items():
                command += ["--header", key + "=" + str(value)]
            upload = subprocess.run(command, capture_output=True, text=True, timeout=90)
            if upload.returncode:
                raise RuntimeError("Signed package upload failed; output withheld")
            uploaded = json.loads(upload.stdout)["data"]
            if uploaded.get("dryRun") is not False or not 200 <= uploaded.get("statusCode", 0) < 300 or uploaded.get("etag", "").strip('"').lower() != hashlib.md5(signed.read_bytes()).hexdigest():
                raise RuntimeError("Signed package upload byte verification failed")
            flow.own("objectId", address["objectId"])
            package = flow.call("testing", "test-api-add-test-package", query={"appId": args.app_id}, body={"distributeMode": 1, "file": {"fileName": signed.name, "objectId": address["objectId"]}}, assert_fn=lambda b: isinstance(b.get("pkgVersion"), list) and len(b["pkgVersion"]) == 1 and isinstance(b["pkgVersion"][0], str), assertion="new package ID returned for signed uploaded APP")
            package_id = package["pkgVersion"][0]
            flow.own("pkgId", package_id)
            for attempt in range(40):
                body = flow.call("publishing", "package-compile-status-v3", query={"appId": args.app_id, "pkgIds": package_id})
                state = package_state(body, package_id)
                flow.ledger["steps"][-1]["packageId"] = package_id
                flow.ledger["steps"][-1]["parseState"] = state
                row = next(s for s in body["pkgStateList"] if s["pkgId"] == package_id)
                error_code = str(row.get("errorCode", ""))
                if re.fullmatch(r"[A-Za-z0-9_.-]{1,80}", error_code):
                    flow.ledger["steps"][-1]["errorCode"] = error_code
                flow.save()
                if state == 0:
                    flow.ledger["steps"][-1]["assertion"] = "new package parsed successfully for dedicated app"
                    flow.ledger["parseVerified"] = True
                    break
                if state == 2:
                    raise RuntimeError("Server package parser rejected signed APP" + ("; code " + error_code if error_code else ""))
                time.sleep(3)
            if not flow.ledger.get("parseVerified"):
                raise RuntimeError("Package parser remains pending after bounded polling")
            try:
                flow.call("publishing", "query-aabfile", query={"appId": args.app_id, "pkgIds": package_id}, profile=args.compile_profile, assert_fn=lambda b: package_state(b, package_id) == 0, assertion="legacy compile query confirms owned parsed package")
            except RuntimeError:
                flow.ledger["legacyCompileVerified"] = False
        except Exception as error:
            flow.ledger["failure"] = str(error) if isinstance(error, RuntimeError) else type(error).__name__
        finally:
            cleanup_ok = True
            operations = []
            if profile_id:
                operations.append(("provision-api-delete-provision", {"query": {"id": profile_id}}, "provision-api-query-provision", {"headers": headers}, lambda b: not any(p.get("id") == profile_id for p in b.get("provisionList", [])), "new release Profile removed"))
            if fingerprint_added:
                operations.append(("provision-api-delete-fingerprints", {"headers": headers, "body": {"fingerprintList": [fingerprint]}}, "provision-api-get-fingerprints", {"headers": headers}, lambda b: fingerprint not in json.dumps(b), "new package signer fingerprint removed"))
            if cert_id:
                operations.append(("provision-api-delete-cent", {"body": {"certIds": [cert_id]}}, "provision-api-query-cent", {"body": {"pageNum": 1, "pageSize": 100}}, lambda b: not any(c.get("id") == cert_id for c in b.get("certList", [])), "new release certificate removed"))
            for endpoint, mutation, query_endpoint, query, assertion, description in operations:
                try:
                    flow.call("provisioning", endpoint, **mutation)
                    flow.call("provisioning", query_endpoint, **query, assert_fn=assertion, assertion=description)
                except Exception:
                    cleanup_ok = False
            flow.ledger["signingResourceCleanupVerified"] = cleanup_ok
    flow.ledger["secretMaterialDeleted"] = True
    flow.ledger["retention"] = "Uploaded unlinked APP/object and package retained; signing resources removed; no review or publication request made"
    flow.ledger["status"] = "verified-package" if flow.ledger.get("parseVerified") and flow.ledger.get("signingResourceCleanupVerified") else "incomplete"
    flow.save()
    print(flow.ledger["status"], flow.ledger.get("failure", ""), str(flow.output.relative_to(ROOT)))
    return int(flow.ledger["status"] != "verified-package")


if __name__ == "__main__":
    raise SystemExit(main())
