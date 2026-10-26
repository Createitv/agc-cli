#!/usr/bin/env python3
"""Build a per-operation offline matrix from actual go test -json pass events.
No authentication, network calls or raw test output is retained. Live is always
not-run here: only a dedicated-resource lifecycle runner may establish that.
"""
import argparse,datetime,hashlib,json,subprocess
from pathlib import Path
FAMILIES={'functions','cloud-auth','clouddb','storage','hosting','templates'}
def matrix(endpoints,events):
 passed={e.get('Test') for e in events if e.get('Action')=='pass' and e.get('Test')}
 packages={e.get('Package') for e in events if e.get('Action')=='pass' and not e.get('Test')}
 required={'github.com/Createitv/agc-cli/cmd/agc/command','github.com/Createitv/agc-cli/pkg/server'}
 package_ok=required<=packages
 rows=[]
 for endpoint in endpoints:
  f,id=endpoint['familyId'],endpoint['id']
  if f not in FAMILIES:continue
  suffix=f+'/'+id
  row={'family':f,'operation':id,'sourceUrl':endpoint['sourceUrl'],'documentationStatus':endpoint.get('documentation',{}).get('status'),'cliHTTP':'passed' if 'TestEveryEndpointHTTPTransport/'+suffix in passed else 'not-run','restHTTP':'passed' if 'TestEveryServerlessOperationREST/'+suffix in passed else 'not-run','dryRun':'passed' if 'TestEveryRegisteredEndpointCanDryRun/'+suffix in passed else 'not-run','live':{'status':'not-run','reason':'No explicit profile/project and dedicated cloud test resources were selected.'}}
  rows.append(row)
 complete=package_ok and bool(rows) and all(all(r[k]=='passed' for k in ('cliHTTP','restHTTP','dryRun')) for r in rows)
 return {'offlineComplete':complete,'packagesPassed':package_ok,'operationCount':len(rows),'operations':rows,'liveVerified':0}
def main():
 p=argparse.ArgumentParser(description=__doc__);p.add_argument('--binary',default='./bin/agc');p.add_argument('--events',required=True);p.add_argument('--output',required=True);a=p.parse_args();binary=Path(a.binary).resolve()
 endpoints=json.loads(subprocess.check_output([str(binary),'endpoints'],text=True))['data']
 events=[json.loads(line) for line in Path(a.events).read_text().splitlines() if line.strip()]
 report=matrix(endpoints,events);report.update(generatedAt=datetime.datetime.now(datetime.timezone.utc).isoformat(),binary=str(binary),binarySHA256=hashlib.sha256(binary.read_bytes()).hexdigest(),evidence='Real local httptest transports via CLI and REST; no Huawei cloud verification implied.')
 target=Path(a.output);target.parent.mkdir(parents=True,exist_ok=True);target.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
 print(json.dumps({k:report[k] for k in ('offlineComplete','packagesPassed','operationCount','liveVerified')}));return 0 if report['offlineComplete'] else 1
if __name__=='__main__':raise SystemExit(main())
