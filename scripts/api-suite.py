"""为全部接口生成测试矩阵，核对官方参数，并通过 CLI 并发复测只读功能。

离线模式使用占位参数，不读取真实授权。在线模式必须指定 profile 和应用；
更新、删除、发布和回调需要专用资源闭环，不因本地通过而标记为远端成功。
"""
import argparse
import concurrent.futures
import datetime
import json
import pathlib
import re
import subprocess
import tempfile
import time
import urllib.request
from html.parser import HTMLParser

ROOT = pathlib.Path(__file__).resolve().parents[1]
READ_POST = {
    ("publishing", "app-version-list"),
    ("provisioning", "provision-api-query-cent"),
    ("provisioning", "provision-api-eligible-acl"),
    ("testing", "test-api-query-feedback"),
    ("testing", "test-api-query-feedback-dimension"),
    *(("pms", operation + "-" + platform) for platform in ("android", "harmonyosnext")
      for operation in ("bygetproductinfo", "bygetpromotioninfo", "getproductgroup")),
}
DOC_API = "https://svc-drcn.developer.huawei.com/community/servlet/consumer/cn/documentPortal/getCenterDocument"


class ParameterParser(HTMLParser):
    def __init__(self):
        super().__init__()
        self.request = False
        self.location = None
        self.capture = None
        self.buffer = []
        self.row = []
        self.parameters = []

    def handle_starttag(self, tag, attrs):
        if tag in ("h2", "h3", "h4", "h5", "h6", "td", "th"):
            self.capture = tag
            self.buffer = []
        if tag == "tr":
            self.row = []

    def handle_data(self, data):
        if self.capture:
            self.buffer.append(data)

    def handle_endtag(self, tag):
        if tag == self.capture:
            value = " ".join("".join(self.buffer).split())
            if tag.startswith("h") and tag != "th":
                value = re.sub(r"^\[h\d\]", "", value)
                if value in ("请求参数", "Request Parameters"):
                    self.request = True
                elif value in ("响应参数", "Response Parameters", "请求示例", "Request Example"):
                    self.request = False
                if value.lower() in ("header", "query", "body", "path"):
                    self.location = value.lower()
            else:
                self.row.append(value)
            self.capture = None
        if tag == "tr" and self.request and self.location and len(self.row) >= 3:
            name, required, kind = self.row[:3]
            if re.fullmatch(r"[A-Za-z_][A-Za-z0-9_.\[\]-]*", name) and required in ("M", "O"):
                item = {"name": name, "in": self.location, "required": required == "M", "type": kind}
                if item not in self.parameters:
                    self.parameters.append(item)


def extract_parameters(html):
    """提取请求参数的位置、必填性与类型；不把响应字段混入请求参数。"""
    parser = ParameterParser()
    parser.feed(html)
    return parser.parameters


def response_status(http_status, body):
    """同时检查 HTTP 状态和业务码；没有业务码时保留待判读状态。"""
    if http_status is None:
        return "invalid-response"
    if not 200 <= http_status < 300:
        return "http-error"
    if not isinstance(body, dict):
        return "response-needs-review"
    codes = []
    for obj, keys in [(body, ("rtnCode", "code", "businessCode")),
                      (body.get("ret"), ("code",)), (body.get("error"), ("errorCode",)),
                      (body.get("errorDetail"), ("code",))]:
        if isinstance(obj, dict):
            codes.extend(str(obj[k]) for k in keys if obj.get(k) is not None)
    if any(c != "0" for c in codes):
        return "business-error"
    return "verified-read" if codes else "response-needs-review"


def read_fixture(endpoint, previous, app, project, group):
    """只构造已明确为只读的请求，群组 ID 必须属于指定应用。"""
    key = (endpoint["familyId"], endpoint["id"])
    if key == ("publishing", "app-version-list"):
        return {"headers": {"appId": app}, "query": {}, "params": {}, "body": {}}
    if key == ("testing", "test-api-query-test-user"):
        if group.get("appId") != app or not group.get("groupId"):
            return None
        return {"headers": {"appId": app, "groupId": group["groupId"]}, "query": {}, "params": {}}
    if key == ("publishing", "app-info-query"):
        return {"query": {"appId": app, "lang": "zh-CN"}, "headers": {}, "params": {}}
    if key == ("projects", "queryprojectdetail"):
        return {"params": {"projectId": project}, "headers": {}, "query": {}} if project else None
    if key == ("projects", "queryprojectlist"):
        return {"params": {}, "headers": {}, "query": {}}
    if endpoint.get("direction", "developer-to-huawei") != "developer-to-huawei":
        return None
    if endpoint.get("method") != "GET" and key not in READ_POST:
        return None
    if previous.get("testedQuery") is None:
        return None
    fixture = {"query": dict(previous["testedQuery"]), "headers": dict(previous.get("testedHeaders", {})), "params": {}}
    if "{appId}" in endpoint.get("path", ""):
        fixture["params"]["appId"] = app
    for values in (fixture["query"], fixture["headers"]):
        for name in ("appId", "appID"):
            if name in values:
                values[name] = app
        if "projectId" in values:
            values["projectId"] = project
    if key in (("comments", "com-rating-harmonyos"), ("comments", "comapi-getreviews-harmonyos")):
        end = int(time.time() * 1000)
        fixture["query"].update(beginTime=str(end - 7 * 86400000), endTime=str(end), countries="CN")
    if previous.get("testedBody") is not None:
        fixture["body"] = dict(previous["testedBody"])
        for name, value in (("appId", app), ("appID", app), ("projectId", project)):
            if name in fixture["body"]:
                fixture["body"][name] = value
    return fixture


def invoke_cli(binary, endpoint, fixture, profile, dry_run):
    """运行真实 CLI，仅记录状态与业务码，不保存完整响应或认证信息。"""
    cmd = [str(binary)]
    if profile:
        cmd += ["--profile", profile]
    cmd += [endpoint["familyId"], endpoint["id"], "--invoke", "--dry-run=" + str(dry_run).lower()]
    if dry_run:
        cmd += ["--base-url", "https://example.invalid", "--token", "offline-placeholder"]
    for location, option in (("params", "--param"), ("query", "--query"), ("headers", "--header")):
        for key, value in fixture.get(location, {}).items():
            cmd += [option, str(key) + "=" + str(value)]
    try:
        with tempfile.TemporaryDirectory(prefix="agc-suite-") as tmp:
            if "body" in fixture:
                body_file = pathlib.Path(tmp) / "body.json"
                body_file.write_text(json.dumps(fixture["body"]))
                cmd += ["--body", str(body_file)]
            proc = subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True, timeout=65)
        if proc.returncode:
            # CLI 错误可能包含业务原始数据或签名地址，只保留状态码。
            match = re.search(r"agc endpoint returned HTTP (\d{3})(?: \(([^)]+)\))?", proc.stderr)
            if match:
                http = int(match[1])
                result = {"status": "http-error" if not 200 <= http < 300 else "business-error", "exitCode": proc.returncode, "httpStatus": http}
                if match[2] and re.fullmatch(r"[A-Za-z0-9_.-]{1,80}", match[2]):
                    result["businessCode"] = match[2]
                return result
            return {"status": "cli-error", "exitCode": proc.returncode}
        data = json.loads(proc.stdout)["data"]
        if dry_run:
            return {"status": "passed" if data.get("dryRun") is True else "invalid-response", "exitCode": 0}
        body = data.get("body")
        result = {"status": response_status(data.get("statusCode"), body), "exitCode": 0, "httpStatus": data.get("statusCode")}
        if isinstance(body, dict):
            result["assertions"] = {"businessCodeZero": result["status"] == "verified-read"}
            for field in ("groups", "testerInfo", "data", "list"):
                if isinstance(body.get(field), list):
                    result["collection"] = {"field": field, "count": len(body[field])}
        return result
    except subprocess.TimeoutExpired:
        return {"status": "timeout"}
    except (OSError, ValueError, KeyError, TypeError):
        return {"status": "invalid-response"}


def offline_fixture(endpoint):
    """占位值仅用于离线请求构造，不可发送到华为服务器。"""
    fixture = {"params": {}, "query": {}, "headers": {}}
    location_keys = {"path": "params", "query": "query", "header": "headers", "body": "body", "file": "body"}
    for p in endpoint.get("parameters", []):
        if p.get("required"):
            value = "https://example.invalid/resource" if p["name"] in ("uploadUrl", "callbackUrl") else "fixture"
            if p["in"] == "body" and p.get("type") in ("integer", "number"):
                value = 1
            if p["in"] == "body" and p.get("type") == "array":
                value = ["fixture"]
            if p["in"] == "body" and p.get("type") == "object":
                value = {"key": "fixture"}
            fixture.setdefault(location_keys[p["in"]], {})[p["name"]] = value
    return fixture


def fetch_catalog():
    """读取新版官方目录，只接受同一 API 分类下唯一匹配的中文标题。"""
    body = {"level2NodeAlias": "submission", "centerPrefix": "doccenter", "language": "cn"}
    req = urllib.request.Request(DOC_API.replace("getCenterDocument", "getCenterCatalogTree"), json.dumps(body).encode(), {"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=20) as response:
            root = json.load(response)
    except Exception:
        return {}
    candidates = {}
    def walk(value, ancestors=()):
        if isinstance(value, dict):
            names = ancestors + (value.get("nodeName", ""),)
            if value.get("relateDocument"):
                candidates.setdefault(value.get("nodeName"), []).append({"slug": value["relateDocument"], "ancestors": names})
            for child in value.values():
                walk(child, names)
        elif isinstance(value, list):
            for child in value:
                walk(child, ancestors)
    walk(root)
    return candidates


def fetch_contract(endpoint, catalog=None):
    """读取官方文档；目录迁移匹配仅供核对，不自动改写生产接口。"""
    if endpoint.get("direction") == "local":
        return {"status": "local-integration", "parameters": endpoint.get("parameters", [])}
    filename = endpoint.get("officialSlug", "")
    body = {"fileName": filename, "level2NodeAlias": "submission", "centerPrefix": "doccenter", "language": "cn"}
    try:
        req = urllib.request.Request(DOC_API, json.dumps(body).encode(), {"Content-Type": "application/json"})
        with urllib.request.urlopen(req, timeout=20) as resp:
            result = json.load(resp)
        value = result.get("value", {})
        html = value.get("content", {}).get("content", "")
        if result.get("code") != 0 or not html:
            family_names = {"publishing": "Publishing API参考", "upload": "Upload Management API参考", "testing": "Testing API参考", "provisioning": "Provisioning API参考", "domains": "Domain Management API参考", "reports": "Reports API参考"}
            family_name = family_names.get(endpoint["familyId"])
            matches = [c for c in (catalog or {}).get(endpoint.get("name"), []) if family_name and family_name in c["ancestors"]]
            if len(matches) == 1 and matches[0]["slug"] != filename:
                current = dict(endpoint, officialSlug=matches[0]["slug"])
                found = fetch_contract(current)
                found["referenceResolution"] = "Unique same-family title in current official catalog; check platform and path before replacing production contract"
                return found
            return {"status": "unavailable", "reason": "Official document not found at current center route"}
        urls = re.findall(r"https://connect-api\.cloud\.huawei\.com[^\s<\"&]+", html)
        return {"status": "fetched", "sourceUrl": "https://developer.huawei.com/consumer/cn/doc/doccenter-submission/" + filename,
                "updatedDate": value.get("updatedDate"), "parameters": extract_parameters(html),
                "requestUrls": list(dict.fromkeys(urls))[:3],
                "registeredPathSeen": any(endpoint["path"] in url for url in urls)}
    except Exception as error:
        return {"status": "unavailable", "reason": type(error).__name__}


def requirements(endpoint, previous, fixture):
    """逐接口列出资源、授权、接收服务和清理策略等验收条件。"""
    direction = endpoint.get("direction")
    if direction == "inbound-callback":
        return ["Developer callback receiver", "Official signature verification and event fixture", "Delivery/replay/duplicate processing assertions"]
    if direction == "local":
        return ["Actual Hvigor integration adapter", "Dedicated HarmonyOS project and SDK", "Build artifact readback"]
    if fixture:
        return []
    items = ["Reviewed official parameter contract", "Explicit credential profile accepted by this API"]
    if endpoint["method"] == "GET" or (endpoint["familyId"], endpoint["id"]) in READ_POST:
        items += ["Real resource IDs and complete read fixture"]
    else:
        items += ["Dedicated resource fixture with explicit allowed mutation fields", "Readback assertion", "Dependency sequence and cleanup policy"]
    if previous.get("nextRequirement"):
        items.append(previous["nextRequirement"])
    return items


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=pathlib.Path, default=ROOT / "bin/agc")
    parser.add_argument("--workers", type=int, default=8)
    parser.add_argument("--fetch-docs", action="store_true")
    parser.add_argument("--live", action="store_true")
    parser.add_argument("--profile")
    parser.add_argument("--app-id")
    parser.add_argument("--project-id")
    parser.add_argument("--output", type=pathlib.Path, default=ROOT / "docs/verification/api-suite-results.json")
    args = parser.parse_args()
    if not 1 <= args.workers <= 16:
        parser.error("workers must be between 1 and 16")
    if args.live and (not args.profile or not args.app_id):
        parser.error("live mode requires --profile and --app-id")
    args.binary = args.binary.resolve()
    endpoints = json.loads(subprocess.check_output([str(args.binary), "endpoints"], cwd=ROOT, text=True))["data"]
    history = json.loads((ROOT / "docs/verification/live-api-status.json").read_text())
    previous = {(e["family"], e["endpoint"]): e for e in history["endpoints"]}
    group_file = ROOT / "docs/verification/landlady-test-group.json"
    group = json.loads(group_file.read_text()) if group_file.exists() else {}
    started = time.monotonic()
    catalog = fetch_catalog() if args.fetch_docs else {}

    def check(endpoint):
        old = previous.get((endpoint["familyId"], endpoint["id"]), {})
        fixture = read_fixture(endpoint, old, args.app_id, args.project_id, group) if args.live else None
        record = {"family": endpoint["familyId"], "endpoint": endpoint["id"], "method": endpoint["method"], "path": endpoint["path"],
                  "direction": endpoint.get("direction"), "sourceUrl": endpoint.get("sourceUrl"),
                  "registeredParameters": endpoint.get("parameters", []),
                  "offline": invoke_cli(args.binary, endpoint, offline_fixture(endpoint), None, True),
                  "priorLiveStatus": old.get("status"), "requirements": requirements(endpoint, old, fixture)}
        if args.fetch_docs:
            record["officialContract"] = fetch_contract(endpoint, catalog)
        record["live"] = invoke_cli(args.binary, endpoint, fixture, args.profile, False) if fixture else {"status": "blocked" if args.live else "not-run"}
        # 每条记录明确测试层级，不能用本地模拟代替华为真实业务验收。
        record["testPlan"] = {
            "offlineTransport": ["success-response", "HTTP-200-business-error", "HTTP-429-error", "query-header-body-forwarding"],
            "liveAssertions": ["HTTP success", "explicit business code zero", "expected resource or response schema", "pagination when applicable"],
            "resourceLifecycle": "fixture -> operation -> readback -> cleanup" if not fixture else "read existing scoped resource; do not mutate",
            "parallelRule": "Independent cases may run concurrently; steps touching the same resource must run sequentially",
        }
        if args.fetch_docs and record["officialContract"]["status"] == "unavailable":
            record["requirements"].append("Current official document and complete field/schema verification")
        if record["live"]["status"] not in ("verified-read", "blocked", "not-run"):
            record["requirements"].append("Resolve current live failure: " + record["live"]["status"])
        if fixture:
            record["liveFixture"] = fixture
        return record

    with concurrent.futures.ThreadPoolExecutor(max_workers=args.workers) as pool:
        results = list(pool.map(check, endpoints))
    from collections import Counter
    summary = {"endpoints": len(results), "offline": dict(Counter(r["offline"]["status"] for r in results)),
               "live": dict(Counter(r["live"]["status"] for r in results)),
               "documents": dict(Counter(r.get("officialContract", {}).get("status", "not-run") for r in results))}
    report = {"checkedAt": datetime.datetime.now(datetime.timezone.utc).isoformat(), "workers": args.workers,
              "elapsedSeconds": round(time.monotonic() - started, 2), "profile": args.profile if args.live else None,
              "scope": "All registered endpoints: offline CLI, optional official document extraction, allowlisted live reads only. No mutations.",
              "summary": summary, "endpoints": results}
    args.output.parent.mkdir(parents=True, exist_ok=True)
    args.output.write_text(json.dumps(report, ensure_ascii=False, indent=2) + "\n")
    md = ["# 全接口测试矩阵", "", "本地请求构造与真实功能验收分别统计；未执行的接口保持 blocked/not-run。", "", "```json", json.dumps(summary, ensure_ascii=False, indent=2), "```", "", "| 类别 | 接口 | 本地 | 本轮真实调用 | 待补条件 |", "| --- | --- | --- | --- | --- |"]
    for r in results:
        md.append("| " + " | ".join([r["family"], r["endpoint"], r["offline"]["status"], r["live"]["status"], "; ".join(r["requirements"]).replace("|", "/")]) + " |")
    args.output.with_suffix(".md").write_text("\n".join(md) + "\n")
    print(json.dumps(summary, ensure_ascii=False, indent=2))
    failed = any(r["offline"]["status"] != "passed" for r in results) or any(r["live"]["status"] in ("cli-error", "timeout", "invalid-response", "business-error", "http-error", "response-needs-review") for r in results)
    # 退出码 2 表示真实验收仍有阻塞条件，不把未执行的接口计为全部通过。
    raise SystemExit(1 if failed else 2 if args.live and any(r["live"]["status"] == "blocked" for r in results) else 0)


if __name__ == "__main__":
    main()
