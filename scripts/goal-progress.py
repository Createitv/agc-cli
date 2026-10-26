"""从原始接口清单与真实生命周期断言重建进度，不使用报告自报的成功数量。"""
import datetime
import importlib.util
import json
import pathlib
import re

ROOT = pathlib.Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location("lifecycle", pathlib.Path(__file__).with_name("live-lifecycle.py"))
lifecycle = importlib.util.module_from_spec(spec)
spec.loader.exec_module(lifecycle)


def collect(baseline, records):
    original = {(r["family"], r["endpoint"]) for r in baseline["originalBlocked"]}
    sources = {ledger["runId"]: ledger for _, ledger in records if ledger.get("runId")}
    completed = {}
    for filename, ledger in records:
        features = lifecycle.verified_features(ledger)
        if ledger.get("flow") == "domain-precheck-custom-host" and ledger.get("status") == "verified-precheck" and ledger.get("cleanupVerified") is True:
            artifact = ledger.get("verificationFile", {})
            owned = ledger.get("owned", {})
            reads = ledger.get("publicReadbacks", [])
            digest = artifact.get("sha256", "")
            size = artifact.get("bytes")
            valid_file = bool(ledger.get("appId")) and artifact.get("fileName") == ledger["appId"] + ".txt" and artifact.get("downloadHTTPStatus") == 200 and type(size) is int and size > 0 and re.fullmatch(r"[a-f0-9]{64}", digest) and artifact.get("sourceURL") == "https://connect-api.cloud.huawei.com/api/dms/domain-manage/v1/app/domain/verify-file"
            valid_host = owned.get("workerName") and owned.get("hostname") and ledger.get("hostedDomain") == "https://" + owned["hostname"]
            valid_reads = reads and all(r.get("httpStatus") == 200 and r.get("bytes") == size and r.get("sha256") == digest and r.get("matchesDownloadedFile") is True for r in reads)
            passed = [s for s in ledger.get("steps", []) if s.get("family") == "domains" and s.get("endpoint") == "domain-api-precheck" and s.get("status") == "passed" and s.get("exitCode") == 0 and s.get("httpStatus") == 200 and s.get("businessCode") == 0]
            if valid_file and valid_host and valid_reads and passed:
                features["domains/domain-api-precheck"] = {"status": "verified-precheck", "cleanupVerified": True, "assertion": "Owned host served the exact downloaded verification file before successful precheck"}
        if ledger.get("flow") == "multipart-upload" and ledger.get("status") == "verified-part-upload-and-compose":
            steps = ledger.get("steps", [])
            names = ("obbfile-init", "obbfile-uploadurl", "obbfile-upload")
            selected = [next((s for s in steps if s.get("endpoint") == name and s.get("status") == "passed"), {}) for name in names]
            receipts = ledger.get("partReceipts", {})
            source = ledger.get("sourceArtifact", {})
            valid_receipts = len(receipts) > 1 and all(r.get("httpStatus") == 200 and re.fullmatch(r"[a-f0-9]{32}", r.get("expectedMD5", "")) and r.get("etag", "").lower() == r["expectedMD5"] and re.fullmatch(r"[a-f0-9]{64}", r.get("sha256", "")) and type(r.get("bytes")) is int and r["bytes"] > 0 for r in receipts.values())
            valid_steps = all(selected) and all(s.get("httpStatus") == 200 and s.get("businessCode") == 0 for s in selected[:2]) and selected[2].get("partCount") == len(receipts)
            ordered = valid_steps and all(steps.index(selected[i]) < steps.index(selected[i+1]) for i in range(2))
            if ordered and valid_receipts and sum(r["bytes"] for r in receipts.values()) == source.get("bytes") and re.fullmatch(r"[a-f0-9]{64}", source.get("sha256", "")) and ledger.get("owned", {}).get("objectId") and ledger.get("owned", {}).get("nspUploadId"):
                for name, step in zip(names, selected):
                    features["upload/"+name] = {"status": "verified-part-upload", "operationAt": step.get("checkedAt"), "assertion": "Owned multipart session produced signed addresses and exact part ETags", "cleanupVerified": False, "retention": ledger.get("retention"), "limitations": ledger.get("limitations")}
            # compose 未独立下载整文件，不计入原始闭环数。
        if ledger.get("flow") == "legacy-compile-readback" and ledger.get("status") == "verified-read" and ledger.get("exitCode") == 0 and ledger.get("businessCode") == 0 and 200 <= ledger.get("httpStatus", 0) < 300:
            source = sources.get(ledger.get("sourcePackageEvidence"), {})
            package_id, expected = ledger.get("packageId"), ledger.get("expectedParseState")
            states = ledger.get("packageStates", [])
            independent = any(s.get("endpoint") == "package-compile-status-v3" and s.get("status") == "passed" and s.get("packageId") == package_id and s.get("parseState") == expected for s in source.get("steps", []))
            if source.get("appId") == ledger.get("appId") and source.get("owned", {}).get("pkgId") == package_id and package_id and type(expected) is int and expected in (0, 1, 2) and len(states) == 1 and states[0].get("pkgId") == package_id and type(states[0].get("successStatus")) is int and states[0]["successStatus"] == expected and independent:
                features["publishing/query-aabfile"] = {"status": "verified-read", "operationAt": ledger.get("checkedAt"), "assertion": "Exact owned package status agrees with independent V3 readback", "statesVerified": [expected], "sourcePackageEvidence": ledger.get("sourcePackageEvidence"), "retention": "Dedicated package retained; this query performs no mutation"}
        if ledger.get("flow") == "pms-group-create" and ledger.get("status") == "verified-retained" and ledger.get("owned", {}).get("groupId"):
            steps = ledger.get("steps", [])
            for index, step in enumerate(steps):
                if step.get("endpoint") == "addproductgroup-harmonyosnext" and step.get("status") == "passed" and step.get("businessCode") == 0:
                    readbacks = [s for s in steps[index + 1:] if s.get("endpoint") == "getproductgroup-harmonyosnext" and s.get("status") == "passed" and s.get("assertion") == "dedicated empty subscription group read back"]
                    if readbacks:
                        features["pms/addproductgroup-harmonyosnext"] = {"status": "verified-retained", "operationAt": step.get("checkedAt"), "readbackAt": readbacks[0].get("checkedAt"), "assertion": readbacks[0]["assertion"], "cleanupVerified": False, "retention": ledger.get("retention")}
        if ledger.get("flow") == "hvigor-cli-build" and ledger.get("status") == "verified-build" and ledger.get("exitCode") == 0 and ledger.get("assertion") == "CLI generated new HAP; independent ZIP, size and SHA256 readback match":
            artifacts = ledger.get("artifacts", [])
            if artifacts and all(a.get("path", "").endswith(".hap") and a.get("size", 0) > 0 and len(a.get("sha256", "")) == 64 for a in artifacts):
                features["cicd/build-hvigor"] = {"status": "verified-build", "assertion": ledger["assertion"], "operationAt": ledger.get("checkedAt"), "artifacts": artifacts, "retention": "Dedicated local fixture and unsigned HAP retained for further package verification"}
        if ledger.get("flow") == "signed-file-upload" and ledger.get("status") == "verified-upload":
            steps = ledger.get("steps", [])
            uploads = [s for s in steps if s.get("endpoint") == "upload-file-new" and s.get("status") == "passed" and s.get("assertion") == "server ETag matches uploaded fixture MD5" and 200 <= s.get("httpStatus", 0) < 300]
            addresses = [s for s in steps if s.get("endpoint") == "upload-url-new" and s.get("status") == "passed" and s.get("businessCode") == 0]
            if uploads and addresses and ledger.get("objectId") and ledger.get("sha256") and ledger.get("contentLength", 0) > 0 and uploads[0].get("etag"):
                features["upload/upload-file-new"] = {"status": "verified-upload", "assertion": uploads[0]["assertion"], "operationAt": ledger.get("checkedAt"), "cleanupVerified": False, "retention": ledger.get("retention")}
        for key, evidence in features.items():
            family, endpoint = key.split("/", 1)
            if (family, endpoint) in original:
                completed[(family, endpoint)] = {"family": family, "endpoint": endpoint, "file": filename, **evidence}
    return {"checkedAt": datetime.datetime.now(datetime.timezone.utc).isoformat(),
            "originalBlocked": len(original), "verifiedOriginalLifecycles": len(completed),
            "remainingOriginal": len(original) - len(completed),
            "completed": [completed[k] for k in sorted(completed)],
            "scope": "Original 109 endpoint objective remains active. Newly registered helper APIs do not substitute for these endpoints."}


if __name__ == "__main__":
    directory = ROOT / "docs/verification"
    baseline = json.loads((directory / "goal-verification-baseline.json").read_text())
    records = [(str(p.relative_to(ROOT)), json.loads(p.read_text())) for p in sorted((directory / "lifecycles").glob("*.json"))]
    report = collect(baseline, records)
    (directory / "goal-progress.json").write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n")
    print(f"verified {report['verifiedOriginalLifecycles']}/{report['originalBlocked']}; remaining {report['remainingOriginal']}")
