#!/usr/bin/env python3
"""Read-only live Publishing checks for the explicitly selected Landlady app.
No tokens, credential responses or full business response bodies are persisted.
"""
import argparse
import datetime
import json
import re
from pathlib import Path
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]
APP = '6917618355927751054'
PACKAGE = 'app.landlady.www.landlady'
CASES = [
    ('app-info-query-v3', ['--query', 'appId='+APP]),
    ('appid-list', ['--query', 'packageName='+PACKAGE]),
    ('app-version-list', ['--header', 'appId='+APP, '--field', 'packageName='+PACKAGE]),
    ('obbfile-package-list', ['--query', 'appId='+APP]),
    ('app-info-query', ['--query', 'appId='+APP]),
]


def business_code(body):
    if not isinstance(body, dict): return None
    ret = body.get('ret')
    if isinstance(ret, dict) and 'code' in ret: return ret['code']
    for key in ('rtnCode', 'code'):
        if key in body: return body[key]
    return None


def identities(value):
    out = []
    if isinstance(value, dict):
        for key, item in value.items():
            if key in ('appId', 'appid', 'packageName', 'appName') and isinstance(item, (str, int)):
                out.append({key: item})
            elif isinstance(item, (dict, list)): out.extend(identities(item))
    elif isinstance(value, list):
        for item in value: out.extend(identities(item))
    return out


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--profile', required=True)
    parser.add_argument('--report', type=Path, required=True)
    args = parser.parse_args()
    report = {'checkedAt': datetime.datetime.now(datetime.timezone.utc).isoformat(),
              'profile': args.profile, 'appId': APP, 'packageName': PACKAGE,
              'scope': 'Registered read-only CLI operations; no production mutations', 'results': []}
    with tempfile.TemporaryDirectory(prefix='agc-publishing-read-') as work:
        binary = str(Path(work) / 'agc')
        subprocess.run(['go', 'build', '-o', binary, './cmd/agc'], cwd=ROOT, check=True)
        for endpoint, fields in CASES:
            result = subprocess.run([binary, '--profile', args.profile, '--timeout', '25s',
                'publishing', endpoint, '--invoke', '--dry-run=false', *fields],
                cwd=ROOT, text=True, capture_output=True, timeout=40)
            entry = {'endpoint': endpoint, 'exitCode': result.returncode}
            try:
                envelope = json.loads(result.stdout)
                data = envelope.get('data', {})
                body = data.get('body', {})
                if isinstance(body, str): body = json.loads(body)
                code = business_code(body)
                entry.update(httpStatus=data.get('statusCode'), businessCode=code,
                    status='passed' if result.returncode == 0 and str(code) == '0' else 'failed',
                    responseKeys=list(body) if isinstance(body, dict) else [],
                    applicationIdentity=identities(body))
                if isinstance(body, dict) and 'appids' in body:
                    entry['expectedAppIdPresent'] = APP in str(body['appids'])
                for key in ('versionList', 'packageList', 'packages', 'appIdList'):
                    if isinstance(body, dict) and isinstance(body.get(key), list): entry[key+'Count'] = len(body[key])
            except (ValueError, TypeError):
                entry['status'] = 'failed'
                match = re.search(r'HTTP (\d+)(?: \(([^)]+)\))?', result.stderr)
                if match:
                    entry['httpStatus'] = int(match[1])
                    entry['businessCode'] = match[2]
                # CLI diagnostics can include upstream payloads; do not persist them.
                entry['diagnostic'] = 'CLI failed without a parseable success envelope; raw output suppressed'
            report['results'].append(entry)
            print(endpoint, entry['status'], entry.get('httpStatus'), entry.get('businessCode'), flush=True)
    args.report.write_text(json.dumps(report, ensure_ascii=False, indent=2)+'\n')


if __name__ == '__main__':
    main()
