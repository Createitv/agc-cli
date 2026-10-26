#!/usr/bin/env python3
"""Read current official Publishing endpoints, independent of CLI registration.
Requires explicit profile. Uses fixed Landlady identity. Never invokes mutations.
"""
import argparse
import datetime
import json
from pathlib import Path
import subprocess
import urllib.error
import urllib.parse
import urllib.request

ROOT = Path(__file__).resolve().parents[1]
APP = '6917618355927751054'
PACKAGE = 'app.landlady.www.landlady'
READ_NAMES = {'查询应用包名对应的appid', '查询团队账号下的成员角色列表', '查询应用信息',
 '查询应用软件包信息', '查询软件包编译状态', '查询应用所有版本类型的版本列表',
 '查询应用未处理的问题工单', '查询需要上传的资质文件信息', '查询资质审核任务详情',
 '查询用户协议列表', '查询隐私政策协议列表'}


def collect(value, names):
    found = []
    if isinstance(value, dict):
        for key, item in value.items():
            if key in names and isinstance(item, (str, int)): found.append(str(item))
            elif isinstance(item, (dict, list)): found.extend(collect(item, names))
    elif isinstance(value, list):
        for item in value: found.extend(collect(item, names))
    return list(dict.fromkeys(found))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--profile', required=True)
    parser.add_argument('--report', type=Path, required=True)
    args = parser.parse_args()
    auth = subprocess.run(['go', 'run', './cmd/agc', '--profile', args.profile, 'auth', 'token'],
        cwd=ROOT, capture_output=True, text=True)
    if auth.returncode: raise SystemExit('Token generation failed; raw authorization output suppressed')
    token = json.loads(auth.stdout)['data']['access_token']
    refs = json.loads((ROOT/'docs/verification/publishing-browser-reference-20261007.json').read_text())['endpoints']
    refs = {e['name']: e for e in refs if e['name'] in READ_NAMES}
    report = {'checkedAt': datetime.datetime.now(datetime.timezone.utc).isoformat(),
        'profile': args.profile, 'appId': APP, 'packageName': PACKAGE,
        'scope': 'Direct official HTTP read checks; success DOES NOT establish CLI support', 'results': []}

    def call(name, query=None, body=None):
        ref = refs[name]
        assert name in READ_NAMES and ref['method'] in ('GET', 'POST')
        query = query or {}
        url = 'https://connect-api.cloud.huawei.com'+ref['path']
        if query: url += '?'+urllib.parse.urlencode(query)
        headers = {'Authorization': 'Bearer '+token, 'appId': APP}
        data = None
        if body is not None: headers['Content-Type']='application/json'; data=json.dumps(body).encode()
        entry = {'name': name, 'method': ref['method'], 'path': ref['path'], 'query': query, 'body': body}
        decoded = {}
        try:
            request = urllib.request.Request(url, data=data, headers=headers, method=ref['method'])
            try:
                with urllib.request.urlopen(request, timeout=25) as response:
                    status=response.status; raw=response.read()
            except urllib.error.HTTPError as error: status=error.code; raw=error.read()
            try: decoded=json.loads(raw)
            except (ValueError, UnicodeError): decoded={}
            ret=decoded.get('ret',{}) if isinstance(decoded,dict) else {}
            code=ret.get('code',decoded.get('rtnCode',decoded.get('code'))) if isinstance(decoded,dict) else None
            entry.update(httpStatus=status,businessCode=code,
                status='passed' if status<400 and str(code)=='0' else 'failed',
                responseKeys=list(decoded) if isinstance(decoded,dict) else [])
            if name=='查询应用信息':entry['appNames']=collect(decoded,{'appName'})
            if name=='查询应用包名对应的appid': entry['expectedAppIdPresent']=APP in str(decoded.get('appids',[]))
        except Exception as error:entry.update(status='transport-error',errorType=type(error).__name__)
        report['results'].append(entry)
        print(name,entry['status'],entry.get('httpStatus'),entry.get('businessCode'),flush=True)
        return decoded

    info=call('查询应用信息',{'appId':APP})
    call('查询应用包名对应的appid',{'packageName':PACKAGE})
    versions=call('查询应用所有版本类型的版本列表',body={'packageName':PACKAGE})
    call('查询团队账号下的成员角色列表',{'fromPage':1,'pageSize':10,'languageCode':'zh_cn'})
    tags=call('查询需要上传的资质文件信息')
    call('查询应用未处理的问题工单',body={'pageNo':1,'pageSize':10})
    call('查询用户协议列表',body={'name':'','type':2,'fromRecCount':0,'maxReqCount':10})
    call('查询隐私政策协议列表',body={'name':'','type':1,'fromRecCount':0,'maxReqCount':10})
    packages=collect(versions,{'pkgId','packageId'})+collect(info,{'pkgId','packageId'})
    if packages:
        call('查询应用软件包信息',{'packageId':packages[0]})
        call('查询软件包编译状态',{'appId':APP,'pkgIds':','.join(packages[:5])})
    else:
        for name in ('查询应用软件包信息','查询软件包编译状态'):
            report['results'].append({'name':name,'status':'blocked','reason':'Version response supplied no real package ID; no invented fixture sent'})
    tagids=collect(tags,{'tagId'})
    if tagids:call('查询资质审核任务详情',{'tagIds':','.join(tagids[:3])})
    else:report['results'].append({'name':'查询资质审核任务详情','status':'blocked','reason':'No real qualification tag ID obtained'})
    args.report.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')


if __name__=='__main__':main()
