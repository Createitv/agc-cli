#!/usr/bin/env python3
"""Normalize reviewed Project API tables without combining auth alternatives."""
import json
import re
from pathlib import Path
ROOT = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/projects'
for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:
        continue
    doc=json.loads(path.read_text())
    fields={}
    rules=[]
    for table in doc['tables']:
        loc=table['location']
        if table['section']!='请求参数' or loc not in ('header','query','path','body'):
            continue
        for name,required,typ,description in table['rows'][1:]:
            if name.lower()=='authorization':
                continue
            if name.lower()=='oauth2token' and path.stem not in ('getteamlist','appbriefinfo','addappfingerprint','queryappfingerprint'):
                continue
            if typ not in ('String','Integer','List<String>') or required not in ('M','O'):
                raise ValueError(f'Unreviewed project field: {path}, {name}, {typ}')
            p={'name':name,'in':loc,'type':{'String':'string','Integer':'integer','List<String>':'array'}[typ],'required':required=='M','description':description,'sourceUrl':table['sourceUrl']}
            if typ=='List<String>':
                p['itemsType']='string'
            if path.stem=='queryprojectlist' and name=='teamId':
                p['required']=False
                rules.append({'id':'oauth-team','kind':'required_with','field':'header.teamId','fields':['header.oauth2Token'],'sourceUrl':table['sourceUrl'],'message':'teamId is required with OAuth authorization.','suggestion':'Supply the selected OAuth team ID in the teamId header.'})
            if name=='fromRecCount': p['minimum']=1
            if name=='maxReqCount': p['maximum']=1000
            if name=='queryFlag': p['enum']=['0','1']
            if path.stem=='appbriefinfo' and name=='appType':
                values=re.findall(r'(?:^|\s)(\d+)：',description)
                if not values:raise ValueError('Application types require a reviewed enumeration')
                p['enum']=list(dict.fromkeys(values))
            if name=='addCertFingerprints':
                p['maxItems']=5
                p['itemsPattern']='^[^:]*$'
            fields[loc+'.'+name]=p
    if path.stem=='appbriefinfo':
        source=fields['query.appId']['sourceUrl']
        if '最大100个' not in fields['query.packageNames']['description'] or '以appId为准' not in fields['query.packageNames']['description']:raise ValueError('Package list limit and precedence need review')
        rules.append({'id':'application-package-count','kind':'csv_count','field':'query.packageNames','minItems':1,'maxItems':100,'unlessPresent':'query.appId','sourceUrl':fields['query.packageNames']['sourceUrl'],'message':'packageNames accepts at most 100 nonempty comma-separated package names.','suggestion':'Use at most 100 package names, or select applications with appId, which takes precedence.'})
        rules.append({'id':'application-identity','kind':'at_least_one','fields':['query.appId','query.packageNames'],'sourceUrl':source,'message':'Supply appId or packageNames to select applications.','suggestion':'Provide a comma-separated appId list or packageNames list; appId takes precedence when both exist.'})
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
    if rules:
        path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
