"""用显式授权和专用应用验证签名上传；不记录签名地址或授权头。"""
import argparse
import base64
import datetime
import hashlib
import json
import pathlib
import subprocess
import tempfile
import uuid

ROOT = pathlib.Path(__file__).resolve().parents[1]


def invoke(binary, profile, endpoint, extra):
    proc = subprocess.run([str(binary), "--profile", profile, "upload", endpoint, "--invoke", "--dry-run=false", *extra], capture_output=True, text=True, timeout=65)
    if proc.returncode:
        raise RuntimeError(endpoint + " failed; raw output withheld")
    result = json.loads(proc.stdout)["data"]
    if result.get("dryRun") is not False or not 200 <= result.get("statusCode", 0) < 300:
        raise RuntimeError(endpoint + " did not return live HTTP success")
    return result


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", required=True, type=pathlib.Path)
    parser.add_argument("--profile", required=True)
    parser.add_argument("--app-id", required=True)
    args = parser.parse_args()
    fixture = json.loads((ROOT / "docs/verification/dedicated-test-app.json").read_text())
    if fixture.get("appId") != args.app_id or fixture.get("readbackVerified") is not True:
        raise SystemExit("Upload requires the dedicated test application")
    run = datetime.datetime.now(datetime.timezone.utc).strftime("%Y%m%dT%H%M%SZ") + "-" + uuid.uuid4().hex[:8]
    report = {"runId": run, "flow": "signed-file-upload", "profile": args.profile, "appId": args.app_id, "status": "running", "steps": [], "retention": "Unlinked test object retained; no object deletion API is registered. Application metadata is not modified."}
    output = ROOT / "docs/verification/lifecycles" / (run + "-signed-file-upload.json")
    try:
        content = base64.b64decode("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jRL0AAAAASUVORK5CYII=")
        report.update(contentLength=len(content), sha256=hashlib.sha256(content).hexdigest())
        with tempfile.TemporaryDirectory(prefix="agc-upload-") as directory:
            image = pathlib.Path(directory) / "fixture.png"
            image.write_bytes(content)
            query = {"appId": args.app_id, "suffix": "png", "fileName": "agc-cli-" + run + ".png", "contentLength": len(content)}
            acquired = invoke(args.binary, args.profile, "upload-url-new", [item for k, v in query.items() for item in ("--query", f"{k}={v}")])
            body = acquired.get("body", {})
            info = body.get("urlInfo", {})
            if str(body.get("ret", {}).get("code")) != "0" or info.get("method", "").upper() != "PUT" or not info.get("objectId") or not info.get("url", "").startswith("https://"):
                raise RuntimeError("Invalid upload address contract")
            report["objectId"] = info["objectId"]
            report["steps"].append({"family": "upload", "endpoint": "upload-url-new", "status": "passed", "httpStatus": acquired["statusCode"], "businessCode": 0})
            extra = ["--param", "uploadUrl=" + info["url"], "--file", str(image)]
            headers = info.get("headers", {})
            if not isinstance(headers, dict):
                raise RuntimeError("Unsupported upload headers format")
            for key, value in headers.items():
                extra += ["--header", str(key) + "=" + str(value)]
            uploaded = invoke(args.binary, args.profile, "upload-file-new", extra)
            etag = uploaded.get("etag", "").strip('"').lower()
            if etag != hashlib.md5(content).hexdigest():
                raise RuntimeError("Server ETag differs from fixture MD5")
            report["steps"].append({"family": "upload", "endpoint": "upload-file-new", "status": "passed", "httpStatus": uploaded["statusCode"], "assertion": "server ETag matches uploaded fixture MD5", "etag": etag})
            report["status"] = "verified-upload"
    except Exception as error:
        report.update(status="incomplete", failure=str(error) if isinstance(error, RuntimeError) else type(error).__name__)
    finally:
        report["checkedAt"] = datetime.datetime.now(datetime.timezone.utc).isoformat()
        output.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n")
        print(report["status"], report.get("failure", ""), str(output.relative_to(ROOT)))
    return 0 if report["status"] == "verified-upload" else 1


if __name__ == "__main__":
    raise SystemExit(main())
