"""Probe selected read endpoints with the active CLI credential; never persist tokens or response bodies."""
import argparse, concurrent.futures, datetime, json, pathlib, subprocess, urllib.request, urllib.error, urllib.parse
ROOT = pathlib.Path(__file__).resolve().parents[1]
REPORT = ROOT / 'docs/verification/live-api-status.json'
APP = '6917618355927751054'
PROJECT = '101653523865229627'
TEAM = '260086000226330263'
QUERIES = {
 ('upload','upload-url-new'): {'appId':APP,'suffix':'png','fileName':'agc-api-probe.png','contentLength':'1'},
 ('publishing','appid-list'): {'packageName':'app.landlady.www.landlady'},
 ('publishing','obbfile-package-list'): {'appId':APP},
 ('provisioning','provision-api-query-device'): {'pageNum':'1','pageSize':'10'},
 ('provisioning','provision-api-query-provision'): {'appId':APP,'pageNum':'1','pageSize':'10'},
 ('provisioning','provision-api-get-fingerprints'): {'appId':APP},
 ('provisioning','provision-api-getacl'): {'appId':APP},
 ('domains','domain-api-get-domain'): {'appId':APP,'category':'0'},
 ('domains','domain-api-get-domain-config'): {'appId':APP},
 ('testing','test-api-get-test-grouplist'): {'appId':APP,'pageNum':'1','pageSize':'10'},
 ('testing','test-api-query-test-user'): {'appId':APP,'pageNum':'1','pageSize':'10'},
 ('projects','appbriefinfo'): {'projectId':PROJECT,'pageNum':'1','pageSize':'10'},
 ('projects','queryappfingerprint'): {},
 ('comments','com-rating-harmonyos'): {'appId':APP},
 ('comments','comapi-getreviews-harmonyos'): {'appId':APP,'pageNum':'1','pageSize':'10'},
 ('pms','basicapplication-harmonyosnext'): {'appId':APP},
 ('provisioning','provision-api-getaclapplystatus'): {'appId':APP},
 ('domains','domain-api-download-domain-config'): {'appId':APP},
 ('projects','getconfigfile'): {'appId':APP,'projectId':PROJECT},
 ('projects','queryservice'): {'appID':APP,'projectId':PROJECT},
 ('comments','com-rating'): {'appId':APP},
 ('comments','comapi-getreviews'): {'appId':APP,'pageNum':'1','pageSize':'10'},
}

# Report URL acquisition is a read operation; these date parameters are candidate
# parameters except for appdownloadexport, whose names match the official example.
for entry in json.loads(REPORT.read_text())['endpoints']:
 if entry['family']=='reports' and entry['method']=='GET':
  QUERIES[('reports',entry['endpoint'])]={'appId':APP,'language':'zh-CN','startTime':'20261001','endTime':'20261006'}

QUERIES[('reports','iapexport')]['currency']='CNY'
READ_POST_BODIES={
 ('provisioning','provision-api-query-cent'): {'pageNum':1,'pageSize':10},
 ('provisioning','provision-api-eligible-acl'): {},
 ('testing','test-api-query-feedback'): {'pageNum':1,'pageSize':10},
 ('testing','test-api-query-feedback-dimension'): {'appId':APP},
}
for platform in ['android','harmonyosnext']:
 for operation in ['bygetproductinfo','bygetpromotioninfo','getproductgroup']:
  READ_POST_BODIES[('pms',operation+'-'+platform)]={'appId':APP,'pageNum':1,'pageSize':10}
for key in READ_POST_BODIES:
 QUERIES[key]={}
QUERIES[('testing','test-api-query-feedback-dimension')]={'appId':APP}

def main():
 parser=argparse.ArgumentParser(description=__doc__)
 parser.add_argument('--only',action='append',help='Endpoint ID to probe, repeatable; otherwise probe untested candidates only')
 parser.add_argument('--workers',type=int,default=8,help='Concurrent read requests (1-16, default 8)')
 parser.add_argument('--profile', required=True, help='Explicit credential profile for live requests')
 args=parser.parse_args()
 if not 1 <= args.workers <= 16:
  parser.error('--workers must be between 1 and 16')
 result=subprocess.run(['go','run','./cmd/agc','--profile',args.profile,'auth','token'],cwd=ROOT,capture_output=True,text=True)
 if result.returncode: raise SystemExit('Local token creation failed; diagnostic output suppressed to protect credentials.')
 token=json.loads(result.stdout)['data']['access_token']
 report=json.loads(REPORT.read_text())
 # 专用群组只能用于所属应用；缺少真实 ID 时不要发送已知必然失败的请求。
 group_file=ROOT/'docs/verification/landlady-test-group.json'
 group=json.loads(group_file.read_text()) if group_file.exists() else {}
 group_id=group.get('groupId') if group.get('appId')==APP else None
 def probe(e):
  key=(e['family'],e['endpoint'])
  path=e['path'].replace('{appId}',APP)
  url='https://connect-api.cloud.huawei.com'+path
  if QUERIES[key]: url+='?'+urllib.parse.urlencode(QUERIES[key])
  headers={'Authorization':'Bearer '+token,'appId':APP,**({'teamId':TEAM} if key==('reports','orderanalysisexport') else {})}
  if key==('testing','test-api-query-test-user'):
   if not group_id: return e,{'status':'blocked','nextRequirement':'Create or select a dedicated test group for this app.'}
   headers['groupId']=group_id
  body=None
  if key in READ_POST_BODIES:
   body=json.dumps(READ_POST_BODIES[key]).encode()
   headers['Content-Type']='application/json'
  req=urllib.request.Request(url,data=body,method=e['method'],headers=headers)
  try:
   with urllib.request.urlopen(req,timeout=20) as r: status=r.status; data=r.read()
  except urllib.error.HTTPError as err: status=err.code; data=err.read()
  except Exception as err:
   return e,{'status':'transport-error','errorType':type(err).__name__}
  code=None; msg=None
  try:
   body=json.loads(data)
   if isinstance(body,dict):
    ret=body.get('ret',{})
    if isinstance(ret,dict): code=ret.get('code'); msg=ret.get('msg')
    if code is None: code=body.get('code',body.get('rtnCode'));  msg=body.get('message',body.get('msg',body.get('rtnDesc',body.get('rtnMsg'))))
  except (ValueError,UnicodeError): pass
  state='http-error' if status>=400 else 'response-needs-review'
  if status<400 and code is not None:
   state='verified-read' if str(code)=='0' else 'business-error'
  return e,{'status':state,'httpStatus':status,'businessCode':code,'businessMessage':str(msg)[:300] if msg else None,'responseBytes':len(data),'testedQuery':QUERIES[key],'testedBody':READ_POST_BODIES.get(key),'testedHeaders':{k:v for k,v in headers.items() if k in ('appId','teamId','groupId')},'checkedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'verificationScope':'Transport probe with candidate parameters; success requires explicit business code zero. Does not prove other accounts/platforms.'}
 known={e['endpoint'] for e in report['endpoints'] if (e['family'],e['endpoint']) in QUERIES}
 if args.only and set(args.only)-known:
  parser.error('Unknown or non-read candidate endpoint: '+','.join(sorted(set(args.only)-known)))
 selected=[e for e in report['endpoints'] if (e['family'],e['endpoint']) in QUERIES and ((args.only and e['endpoint'] in args.only) or (not args.only and e['status']=='not-tested'))]
 with concurrent.futures.ThreadPoolExecutor(max_workers=args.workers) as pool:
  for e,updates in pool.map(probe,selected):
   previous={k:v for k,v in e.items() if k not in ['attempts','family','endpoint','method','path']}
   e.setdefault('attempts',[]).append(previous)
   e.update(updates)
   print(e['family'],e['endpoint'],e['status'],e.get('httpStatus'),e.get('businessCode'),e.get('businessMessage'),flush=True)
 report['checkedAt']=datetime.datetime.now(datetime.timezone.utc).isoformat()
 REPORT.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
if __name__=='__main__': main()
