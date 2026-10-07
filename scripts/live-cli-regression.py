"""Replay proven read requests through the actual CLI concurrently without logging secrets."""
import argparse,concurrent.futures,datetime,json,pathlib,subprocess,time
ROOT=pathlib.Path(__file__).resolve().parents[1]
REPORT=ROOT/'docs/verification/live-api-status.json'
BINARY='/tmp/agc-live-verified'
def main():
 parser=argparse.ArgumentParser(description=__doc__)
 parser.add_argument('--binary',default=BINARY)
 parser.add_argument('--only',action='append')
 args=parser.parse_args()
 report=json.loads(REPORT.read_text())
 selected=[e for e in report['endpoints'] if e['status']=='verified-read' and e.get('testedQuery') is not None and (not args.only or e['endpoint'] in args.only)]
 started=time.monotonic()
 def check(e):
  cmd=[args.binary,e['family'],e['endpoint'],'--invoke','--dry-run=false']
  if '{appId}' in e['path']: cmd+=['--param','appId=6917618355927751054']
  for location,option in [('testedQuery','--query'),('testedHeaders','--header')]:
   for key,value in e.get(location,{}).items(): cmd+=[option,f'{key}={value}']
  body=e.get('testedBody')
  if body is not None:
   # Existing approved query-only POST bodies use simple scalar types.
   import tempfile
   with tempfile.NamedTemporaryFile(mode='w',suffix='.json') as f:
    json.dump(body,f);f.flush();cmd+=['--body',f.name]
    r=subprocess.run(cmd,cwd=ROOT,capture_output=True,text=True,timeout=65)
  else: r=subprocess.run(cmd,cwd=ROOT,capture_output=True,text=True,timeout=65)
  evidence={'checkedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'exitCode':r.returncode}
  if r.returncode==0:
   d=json.loads(r.stdout)['data']; evidence.update(httpStatus=d.get('statusCode'),dryRun=d.get('dryRun'))
  # stderr and stdout may contain signed URLs or application fields, so neither is saved.
  return e,evidence
 with concurrent.futures.ThreadPoolExecutor(max_workers=8) as pool:
  for e,evidence in pool.map(check,selected):
   e['cliEvidence']=evidence
   print(e['family'],e['endpoint'],'CLI_EXIT',evidence['exitCode'],'HTTP',evidence.get('httpStatus'),flush=True)
 report.setdefault('cliRegressionRuns',[]).append(report.get('cliRegressionSummary',{}))
 report['cliRegressionSummary']={'checkedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'workers':8,'elapsedSeconds':round(time.monotonic()-started,2),'tested':len(selected),'success':sum(e['cliEvidence']['exitCode']==0 for e in selected)}
 REPORT.write_text(json.dumps(report,ensure_ascii=False,indent=2)+'\n')
 print(report['cliRegressionSummary'])
if __name__=='__main__':main()
