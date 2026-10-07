"""通过真实 AGC CLI 验证专用群组、邀请码和测试草稿的资源生命周期。

只清理本次运行创建的资源；不发布应用、不发送邀请通知、不修改现有群组。
每一步写入证据，异常后也保留创建的资源 ID 和清理结果。
"""
import argparse
import concurrent.futures
import datetime
import json
import pathlib
import re
import subprocess
import tempfile
import uuid
import hashlib
import urllib.request
import urllib.parse

ROOT = pathlib.Path(__file__).resolve().parents[1]


def is_owned(ledger, kind, resource_id):
    """删除与修改必须匹配本次运行创建并登记的资源。"""
    return bool(resource_id) and ledger.get("owned", {}).get(kind) == resource_id


def is_dedicated_app(app_id, fixture):
    return bool(app_id) and fixture.get("appId") == app_id and fixture.get("readbackVerified") is True


def group_matches(groups, group_id, expected_name):
    return any(g.get("groupId") == group_id and g.get("groupName") == expected_name for g in groups)


def safe_evidence(body):
    """邀请码、账号信息、下载签名及完整响应均不写入报告。"""
    result = {}
    if isinstance(body.get("ret"), dict):
        result["businessCode"] = body["ret"].get("code")
    else:
        result["businessCode"] = body.get("rtnCode", body.get("code"))
    if body.get("businessCode") is not None and str(body["businessCode"]) != "0":
        result["businessCode"] = body["businessCode"]
    for key in ("appId", "groupId", "versionId", "invitationCodeId"):
        if isinstance(body.get(key), str):
            result[key] = body[key]
    for key in ("certInfo", "provisionInfo"):
        if isinstance(body.get(key), dict) and body[key].get("id"):
            result[key] = {"id": body[key]["id"]}
    return result


def verified_features(ledger):
    """一项变更必须有对应读回断言和资源清理，才算该项闭环成功。"""
    if not ledger.get("cleanupVerified"):
        return {}
    readbacks = {
        "testing/test-api-add-test-group": "new group ID and name exist",
        "testing/test-api-edit-test-group": "same group ID has updated name",
        "testing/test-api-delete-test-group": "temporary group absent after deletion",
        "testing/test-api-add-invite-code": "new code is active; secret omitted",
        "testing/test-api-stop-invite-code": "same code is inactive",
        "testing/test-api-delete-invite-code": "deleted code absent",
        "testing/test-api-add-test-version": "created invitation-test draft visible",
        "testing/test-api-delete-test-version": "temporary draft absent after deletion",
        "provisioning/provision-api-apply-cent": "dedicated debug certificate exists",
        "provisioning/provision-api-add-fingerprints": "new certificate fingerprint read back",
        "provisioning/provision-api-apply-provision": "new Profile exists under dedicated app",
        "provisioning/provision-api-update-provision": "temporary Profile now binds a different real device",
        "provisioning/provision-api-delete-provision": "dedicated Profile absent after cleanup",
        "provisioning/provision-api-delete-fingerprints": "only new fingerprint removed",
        "provisioning/provision-api-delete-cent": "only new certificate removed",
    }
    features = {}
    steps = ledger.get("steps", [])
    for index, step in enumerate(steps):
        key = step["family"] + "/" + step["endpoint"]
        if step["status"] != "passed" or key not in readbacks:
            continue
        for next_step in steps[index + 1:]:
            if next_step.get("assertion") == readbacks[key] and next_step.get("status") == "passed":
                features[key] = {"status": "verified-lifecycle", "operationAt": step.get("checkedAt"), "readbackAt": next_step.get("checkedAt"), "assertion": readbacks[key], "cleanupVerified": True}
                break
    return features


class Flow:
    def __init__(self, args, name, run_id):
        self.args = args
        self.output = args.output_dir / (run_id + "-" + name + ".json")
        self.ledger = {"runId": run_id, "flow": name, "profile": args.profile, "appId": args.app_id, "owned": {}, "steps": [], "status": "running"}
        self.save()

    def save(self):
        self.ledger["updatedAt"] = datetime.datetime.now(datetime.timezone.utc).isoformat()
        self.output.parent.mkdir(parents=True, exist_ok=True)
        temporary = self.output.with_suffix(".tmp")
        temporary.write_text(json.dumps(self.ledger, ensure_ascii=False, indent=2) + "\n")
        temporary.replace(self.output)

    def call(self, family, endpoint, *, headers=None, query=None, body=None, assert_fn=None, assertion=None):
        cmd = [str(self.args.binary), "--profile", self.args.profile, family, endpoint, "--invoke", "--dry-run=false"]
        for values, flag in ((headers or {}, "--header"), (query or {}, "--query")):
            for key, value in values.items():
                cmd += [flag, str(key) + "=" + str(value)]
        step = {"family": family, "endpoint": endpoint, "checkedAt": datetime.datetime.now(datetime.timezone.utc).isoformat(), "assertion": assertion, "status": "running"}
        self.ledger["steps"].append(step)
        self.save()
        try:
            with tempfile.TemporaryDirectory(prefix="agc-lifecycle-") as tmp:
                if body is not None:
                    file = pathlib.Path(tmp) / "body.json"
                    file.write_text(json.dumps(body))
                    cmd += ["--body", str(file)]
                proc = subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True, timeout=65)
            step["exitCode"] = proc.returncode
            if proc.returncode:
                match = re.search(r"agc endpoint returned HTTP (\d{3})(?: \(([^)]+)\))?", proc.stderr)
                step["status"] = "cli-error"
                if match:
                    step["httpStatus"] = int(match[1])
                    step["status"] = "http-error" if int(match[1]) >= 400 else "business-error"
                    if match[2] and re.fullmatch(r"[A-Za-z0-9_.-]{1,80}", match[2]):
                        step["businessCode"] = match[2]
                raise RuntimeError(endpoint + ": " + step["status"])
            data = json.loads(proc.stdout)["data"]
            response = data.get("body")
            if not isinstance(response, dict):
                raise RuntimeError(endpoint + ": no JSON response")
            step.update(safe_evidence(response))
            step["httpStatus"] = data.get("statusCode")
            if data.get("dryRun") is not False or not 200 <= step["httpStatus"] < 300 or str(step["businessCode"]) != "0":
                raise RuntimeError(endpoint + ": no explicit business success")
            if assert_fn and not assert_fn(response):
                step["status"] = "assertion-failed"
                raise RuntimeError(endpoint + ": readback assertion failed")
            step["status"] = "passed"
            print(self.ledger["flow"], endpoint, "passed", flush=True)
            return response
        except Exception as error:
            if step["status"] == "running":
                step["status"] = "timeout" if isinstance(error, subprocess.TimeoutExpired) else "invalid-response"
            raise
        finally:
            self.save()

    def own(self, kind, resource_id):
        if not resource_id:
            raise RuntimeError("create response missing resource ID")
        self.ledger["owned"][kind] = resource_id
        self.save()

    def guard(self, kind, resource_id):
        if not is_owned(self.ledger, kind, resource_id):
            raise RuntimeError("refuse mutation of resource not owned by this run")

    def finish(self):
        self.ledger["verifiedFeatures"] = verified_features(self.ledger)
        self.ledger["status"] = "passed" if all(s["status"] == "passed" for s in self.ledger["steps"]) and self.ledger.get("cleanupVerified") and not self.ledger.get("updateBlocked") else "incomplete"
        self.save()
        return self.ledger


def group_flow(args, run_id):
    flow = Flow(args, "group-and-invitation", run_id)
    h = {"appId": args.app_id}
    group_id = None
    code_id = None
    name = "AGC CLI 验证 " + run_id[-8:]
    try:
        initial = flow.call("testing", "test-api-get-test-grouplist", headers=h, query={"current": 1, "pageSize": 100})
        existing = {g["groupId"] for g in initial.get("groups", [])}
        created = flow.call("testing", "test-api-add-test-group", headers=h, body={"groupName": name, "groupType": 0}, assert_fn=lambda b: bool(b.get("groupId")) and b["groupId"] not in existing, assertion="returned ID is new")
        group_id = created["groupId"]
        flow.own("groupId", group_id)
        flow.call("testing", "test-api-get-test-grouplist", headers=h, query={"current": 1, "pageSize": 100}, assert_fn=lambda b: group_matches(b.get("groups", []), group_id, name), assertion="new group ID and name exist")
        renamed = name + " 已修改"
        flow.guard("groupId", group_id)
        flow.call("testing", "test-api-edit-test-group", headers=h, body={"groupId": group_id, "groupName": renamed})
        flow.call("testing", "test-api-get-test-grouplist", headers=h, query={"current": 1, "pageSize": 100}, assert_fn=lambda b: group_matches(b.get("groups", []), group_id, renamed), assertion="same group ID has updated name")
        flow.call("testing", "test-api-query-test-user", headers={**h, "groupId": group_id}, assert_fn=lambda b: b.get("testerInfo") == [] and b.get("pageInfo", {}).get("totalRecord") == 0, assertion="dedicated empty group has zero testers")
        invite = flow.call("testing", "test-api-add-invite-code", headers=h, body={"groupId": group_id, "invitationCodeValidDays": 30, "invitationCodeInviteLimit": 1}, assert_fn=lambda b: bool(b.get("invitationCodeId")), assertion="new invitation resource ID returned")
        code_id = invite["invitationCodeId"]
        flow.own("invitationCodeId", code_id)
        flow.call("testing", "test-api-get-invite-code", headers=h, query={"groupId": group_id}, assert_fn=lambda b: any(c.get("id") == code_id and c.get("status") == 1 for c in b.get("invitationCodes", [])), assertion="new code is active; secret omitted")
        flow.guard("invitationCodeId", code_id)
        flow.call("testing", "test-api-stop-invite-code", headers=h, body={"invitationCodeId": code_id})
        flow.call("testing", "test-api-get-invite-code", headers=h, query={"groupId": group_id}, assert_fn=lambda b: any(c.get("id") == code_id and c.get("status") == 2 for c in b.get("invitationCodes", [])), assertion="same code is inactive")
        flow.call("testing", "test-api-delete-invite-code", headers=h, query={"invitationCodeId": code_id})
        flow.ledger["invitationDeleted"] = True
        flow.call("testing", "test-api-get-invite-code", headers=h, query={"groupId": group_id}, assert_fn=lambda b: not any(c.get("id") == code_id for c in b.get("invitationCodes", [])), assertion="deleted code absent")
    except Exception as error:
        flow.ledger["failure"] = str(error)
    finally:
        if group_id and is_owned(flow.ledger, "groupId", group_id):
            try:
                # 失败时只清理本次创建的群组，不触碰原有 AGC CLI 接口验证组。
                flow.call("testing", "test-api-delete-test-group", headers=h, query={"groupId": group_id})
                flow.call("testing", "test-api-get-test-grouplist", headers=h, query={"current": 1, "pageSize": 100}, assert_fn=lambda b: not any(g.get("groupId") == group_id for g in b.get("groups", [])), assertion="temporary group absent after deletion")
                flow.ledger["cleanupVerified"] = True
            except Exception as error:
                flow.ledger["cleanupFailure"] = str(error)
    return flow.finish()


def version_flow(args, run_id):
    flow = Flow(args, "test-draft", run_id)
    app_query = {"appId": args.app_id}
    h = {"appId": args.app_id}
    version_id = None
    try:
        initial = flow.call("publishing", "app-version-list", headers=h, body={})
        existing = {v["versionId"] for v in initial.get("versionList", [])}
        response = flow.call("testing", "test-api-add-test-version", query=app_query, body={"releaseType": 6, "testType": 3, "testDesc": "AGC CLI 专用草稿 " + run_id[-8:], "onshelfSelfDetect": 0}, assert_fn=lambda b: bool(b.get("versionId")) and b["versionId"] not in existing, assertion="created draft ID differs from all existing versions")
        version_id = response["versionId"]
        flow.own("versionId", version_id)
        flow.call("publishing", "app-version-list", headers=h, body={}, assert_fn=lambda b: any(v.get("versionId") == version_id and v.get("testType") == 3 for v in b.get("versionList", [])), assertion="created invitation-test draft visible")
    except Exception as error:
        flow.ledger["failure"] = str(error)
    finally:
        if version_id and is_owned(flow.ledger, "versionId", version_id):
            try:
                flow.call("testing", "test-api-delete-test-version", query={**app_query, "versionId": version_id})
                flow.call("publishing", "app-version-list", headers=h, body={}, assert_fn=lambda b: not any(v.get("versionId") == version_id for v in b.get("versionList", [])), assertion="temporary draft absent after deletion")
                flow.ledger["cleanupVerified"] = True
            except Exception as error:
                flow.ledger["cleanupFailure"] = str(error)
    return flow.finish()


def provisioning_flow(args, run_id):
    """新建独立证书和 Profile，只向专用应用添加并移除本次证书指纹。"""
    flow = Flow(args, "certificate-profile-fingerprint", run_id)
    fixture_file = ROOT / "docs/verification/dedicated-test-app.json"
    fixture = json.loads(fixture_file.read_text()) if fixture_file.exists() else {}
    if not is_dedicated_app(args.app_id, fixture):
        flow.ledger.update(status="blocked", failure="Provisioning mutations require the verified dedicated application")
        flow.save()
        return flow.ledger
    h = {"appId": args.app_id}
    cert_id = None
    profile_id = None
    fingerprint = None
    fingerprint_added = False
    with tempfile.TemporaryDirectory(prefix="agc-signing-fixture-") as tmp:
        try:
            initial = flow.call("provisioning", "provision-api-query-cent", body={"pageNum": 1, "pageSize": 100})
            existing_cert_ids = {c["id"] for c in initial.get("certList", [])}
            initial_fingerprints = flow.call("provisioning", "provision-api-get-fingerprints", headers=h)
            devices = flow.call("provisioning", "provision-api-query-device", query={"pageNum": 1, "pageSize": 100}).get("deviceList", [])
            path = pathlib.Path(tmp)
            private = path / "private.pem"
            csr_file = path / "request.csr"
            # 私钥仅存在于临时目录，不进入日志、报告或仓库。
            generated = subprocess.run(["openssl", "req", "-new", "-newkey", "ec", "-pkeyopt", "ec_paramgen_curve:prime256v1", "-nodes", "-keyout", str(private), "-out", str(csr_file), "-sha256", "-subj", "/CN=AGC CLI API validation/O=API Test/C=CN"], capture_output=True, timeout=30)
            if generated.returncode:
                raise RuntimeError("local EC CSR generation failed")
            private.chmod(0o600)
            created = flow.call("provisioning", "provision-api-apply-cent", body={"csr": csr_file.read_text(), "certName": "AGC CLI 验证 " + run_id[-8:], "certType": 1}, assert_fn=lambda b: bool(b.get("certInfo", {}).get("id")) and b["certInfo"]["id"] not in existing_cert_ids, assertion="new EC debug certificate ID")
            info = created["certInfo"]
            cert_id = info["id"]
            flow.own("certId", cert_id)
            flow.call("provisioning", "provision-api-query-cent", body={"pageNum": 1, "pageSize": 100}, assert_fn=lambda b: any(c.get("id") == cert_id and c.get("certType") == 1 for c in b.get("certList", [])), assertion="dedicated debug certificate exists")
            url = info.get("certDownloadUrl")
            if not url or urllib.parse.urlparse(url).scheme != "https":
                raise RuntimeError("missing HTTPS certificate download")
            with urllib.request.urlopen(url, timeout=30) as response:
                cert = response.read(1024 * 1024)
                if response.status != 200 or not cert:
                    raise RuntimeError("certificate download failed")
            cert_file = path / "certificate.cer"
            cert_file.write_bytes(cert)
            parsed = subprocess.run(["openssl", "x509", "-inform", "DER", "-in", str(cert_file), "-noout", "-fingerprint", "-sha256"], capture_output=True, text=True, timeout=20)
            if parsed.returncode:
                parsed = subprocess.run(["openssl", "x509", "-in", str(cert_file), "-noout", "-fingerprint", "-sha256"], capture_output=True, text=True, timeout=20)
            match = re.search(r"=([0-9A-Fa-f:]+)", parsed.stdout)
            if parsed.returncode or not match:
                raise RuntimeError("downloaded certificate is not valid X.509")
            fingerprint = match[1].replace(":", "").upper()
            if fingerprint in json.dumps(initial_fingerprints):
                raise RuntimeError("generated certificate fingerprint unexpectedly already exists")
            flow.ledger["certificateDownload"] = {"bytes": len(cert), "sha256": hashlib.sha256(cert).hexdigest(), "parsedX509": True}
            flow.own("fingerprint", fingerprint)
            flow.call("provisioning", "provision-api-add-fingerprints", headers=h, body={"fingerprintList": [fingerprint]})
            fingerprint_added = True
            flow.call("provisioning", "provision-api-get-fingerprints", headers=h, assert_fn=lambda b: fingerprint in json.dumps(b), assertion="new certificate fingerprint read back")
            device_ids = [d["id"] for d in devices[:1]]
            profile = flow.call("provisioning", "provision-api-apply-provision", body={"provisionName": "AGC CLI 验证 " + run_id[-8:], "provisionType": 1, "certId": cert_id, "appId": args.app_id, "deviceIdList": device_ids}, assert_fn=lambda b: bool(b.get("provisionInfo", {}).get("id")), assertion="new dedicated debug Profile ID")
            profile_id = profile["provisionInfo"]["id"]
            flow.own("provisionId", profile_id)
            flow.call("provisioning", "provision-api-query-provision", headers=h, assert_fn=lambda b: any(p.get("id") == profile_id for p in b.get("provisionList", [])), assertion="new Profile exists under dedicated app")
            if len(devices) >= 2:
                # 修改的是新 Profile，不删除或修改账号中既有真实设备。
                changed_ids = [devices[1]["id"]]
                flow.call("provisioning", "provision-api-update-provision", body={"provisionId": profile_id, "deviceIdList": changed_ids})
                flow.call("provisioning", "provision-api-query-provision", headers=h, assert_fn=lambda b: any(p.get("id") == profile_id and {d.get("id") for d in p.get("deviceList", [])} == set(changed_ids) for p in b.get("provisionList", [])), assertion="temporary Profile now binds a different real device")
            else:
                flow.ledger["updateBlocked"] = "Need a second real device ID to verify a binding change; empty device lists are rejected with 204144776"
        except Exception as error:
            # 外部服务的错误正文不保存，避免报告包含签名地址或账号信息。
            flow.ledger["failure"] = str(error) if isinstance(error, RuntimeError) else type(error).__name__
        finally:
            cleanup_ok = True
            if profile_id and is_owned(flow.ledger, "provisionId", profile_id):
                try:
                    flow.call("provisioning", "provision-api-delete-provision", query={"id": profile_id})
                    flow.call("provisioning", "provision-api-query-provision", headers=h, assert_fn=lambda b: not any(p.get("id") == profile_id for p in b.get("provisionList", [])), assertion="dedicated Profile absent after cleanup")
                except Exception:
                    cleanup_ok = False
            if fingerprint_added and is_owned(flow.ledger, "fingerprint", fingerprint):
                try:
                    flow.call("provisioning", "provision-api-delete-fingerprints", headers=h, body={"fingerprintList": [fingerprint]})
                    flow.call("provisioning", "provision-api-get-fingerprints", headers=h, assert_fn=lambda b: fingerprint not in json.dumps(b), assertion="only new fingerprint removed")
                except Exception:
                    cleanup_ok = False
            if cert_id and is_owned(flow.ledger, "certId", cert_id):
                try:
                    flow.call("provisioning", "provision-api-delete-cent", body={"certIds": [cert_id]})
                    flow.call("provisioning", "provision-api-query-cent", body={"pageNum": 1, "pageSize": 100}, assert_fn=lambda b: not any(c.get("id") == cert_id for c in b.get("certList", [])), assertion="only new certificate removed")
                except Exception:
                    cleanup_ok = False
            flow.ledger["cleanupVerified"] = bool(cert_id) and cleanup_ok
    return flow.finish()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", required=True, type=pathlib.Path)
    parser.add_argument("--profile", required=True)
    parser.add_argument("--app-id", required=True)
    parser.add_argument("--flows", choices=("group", "draft", "both", "provisioning"), default="both")
    parser.add_argument("--output-dir", type=pathlib.Path, default=ROOT / "docs/verification/lifecycles")
    args = parser.parse_args()
    args.binary = args.binary.resolve()
    run_id = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ") + "-" + uuid.uuid4().hex[:8]
    flows = [group_flow, version_flow] if args.flows == "both" else [{"group": group_flow, "draft": version_flow, "provisioning": provisioning_flow}[args.flows]]
    with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:
        reports = list(pool.map(lambda f: f(args, run_id), flows))
    for report in reports:
        print(report["flow"], report["status"], "cleanup", report.get("cleanupVerified", False), report.get("failure", ""))
    raise SystemExit(0 if all(r["status"] == "passed" for r in reports) else 1)


if __name__ == "__main__":
    main()
