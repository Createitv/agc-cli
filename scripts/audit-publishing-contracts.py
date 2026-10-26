#!/usr/bin/env python3
"""Compare offline CLI contracts with independently captured Huawei parameter tables.

Does not read credentials or call Huawei. Prints JSON; use --cli-json to inspect
an existing `agc publishing endpoints` result, or invokes the local Go CLI.
"""
import argparse
import json
from pathlib import Path
import subprocess

ROOT = Path(__file__).resolve().parents[1]
ALIASES = {
    '更新应用基本信息': 'app-info-update', '更新应用语言信息': 'language-info-update',
    '更新应用文件信息': 'app-file-info', '更新分阶段发布': 'phased-release',
    '查询软件包编译状态': 'query-aabfile', '提交发布': 'app-submit',
}
AUTH = {'Authorization', 'client_id', 'teamId', 'oauth2Token'}


def kind(value):
    if value.startswith('String'): return 'string'
    if value.startswith('Integer'): return 'integer'
    if value.startswith('List<'): return 'array'
    return 'object'


def audit(reference, local):
    results = []
    for official in reference['endpoints']:
        exact = [e for e in local if (e['method'], e['path']) == (official['method'], official['path'])]
        # User and privacy agreements share method/path while using separate
        # models and commands. Preserve the official document's identity.
        exact.sort(key=lambda e: e.get('sourceUrl') != official.get('sourceUrl'))
        legacy = next((e for e in local if e['id'] == ALIASES.get(official['name'])), None)
        matched = exact[0] if exact else legacy
        status = 'path-matched' if exact else ('legacy-path' if legacy else 'missing')
        issues = []
        reviewed = []
        if exact:
            known = {(p['in'], p['name']): p for p in matched.get('parameters', [])}
            for p in official['parameters']:
                if p['in'] == 'header' and p['name'] in AUTH: continue
                actual = known.get((p['in'], p['name']))
                label = p['name'] + '@' + p['in']
                if actual is None:
                    issues.append(label + ': absent from declared contract')
                    continue
                required, expected_type = p['required'], kind(p['type'])
                # 表格与说明/示例冲突时，保留精确的官方原文依据。
                if official['name'] == '查询应用所有版本类型的版本列表' and label == 'appId@header':
                    required = False
                    reviewed.append({'field': label, 'reason': 'appId和packageName至少传一个', 'sourceUrl': 'https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-query-brief-info-list-0000002236201254#section2357705386'})
                if official['name'] == '新建用户协议' and label == 'agreementInfos@body':
                    expected_type = 'array'
                    reviewed.append({'field': label, 'reason': '请求示例使用对象数组，说明规定数组长度不超过100', 'sourceUrl': 'https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-create-user-agreement-0000002271160581#section1945115411389'})
                if actual['required'] != required:
                    issues.append(label + ': requiredness differs')
                if actual.get('type') != expected_type:
                    issues.append(label + ': type missing or differs')
        results.append({**official, 'command': matched['id'] if matched else None,
                        'localMethod': matched['method'] if matched else None,
                        'localPath': matched['path'] if matched else None,
                        'status': status, 'parameterIssues': issues, 'reviewedExceptions': reviewed})
    return {'checkedAt': reference['checkedAt'], 'scope': 'method, path, top-level parameter location/requiredness/type; authentication alternatives and nested model constraints require separate validation',
            'summary': {s: sum(e['status'] == s for e in results) for s in ('path-matched', 'legacy-path', 'missing')},
            'endpoints': results}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--cli-json', type=Path)
    args = parser.parse_args()
    reference = json.loads((ROOT / 'docs/verification/publishing-browser-reference-20261007.json').read_text())
    raw = args.cli_json.read_text() if args.cli_json else subprocess.check_output(
        ['go', 'run', './cmd/agc', 'publishing', 'endpoints'], cwd=ROOT, text=True)
    print(json.dumps(audit(reference, json.loads(raw)['data']), ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
