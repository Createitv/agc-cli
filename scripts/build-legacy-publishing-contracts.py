#!/usr/bin/env python3
"""Normalize reviewed legacy root tables; retain explicit nested model boundaries."""
import json
import re
from pathlib import Path
REPO=Path(__file__).resolve().parents[1]
ROOT=REPO/'pkg/domain/reference/publishing'
current={e['id'] for e in json.loads((REPO/'cmd/agc/command/testdata/publishing-current.json').read_text())['endpoints']}
# The bare submit body and file-list example require further wire-shape review.
PENDING={'app-file-info'}

def field(row,loc,source):
    name,required,typ,description=row
    kind='string' if typ.startswith('String') else 'integer' if typ.startswith('Integer') else 'boolean' if typ=='boolean' else 'array' if typ=='List<DeviceTypeInfo>' else 'object' if typ=='PhasedReleaseInfo' else None
    if kind is None or required not in ('M','O'):raise ValueError(f'Unreviewed legacy field {name}:{typ}')
    p={'name':name,'in':loc,'type':kind,'required':required=='M','description':description,'sourceUrl':source}
    length=re.fullmatch(r'String\((\d+)\)',typ)
    if length:p['maxLength']=int(length[1])
    if kind=='integer' and '取值范围：' in description:
        enum_text=description.split('取值范围：',1)[1]
        enum_text=re.split(r'默认值|注意|说明|如果|该参数',enum_text)[0]
        values=re.findall(r'(?:^|\s)(-?\d+)：',enum_text)
        if values:p['enum']=sorted(set(values),key=int)
    if kind=='array':p['itemsType']='object'
    return p

for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem or path.stem in current or path.stem in PENDING:continue
    doc=json.loads(path.read_text());fields={};rules=[]
    tables=[t for t in doc['tables'] if t['section']=='请求参数']
    for table in tables:
        if table['location'] not in ('header','query','path'):continue
        for row in table['rows'][1:]:
            if row[0] in ('client_id','Authorization','oauth2Token'):continue
            p=field(row,table['location'],table['sourceUrl'])
            if p['name']=='teamId':
                p['required']=False
                rules.append({'id':'oauth-team','kind':'required_with','field':'header.teamId','fields':['header.oauth2Token'],'sourceUrl':table['sourceUrl'],'message':'OAuth authorization requires its selected teamId.','suggestion':'Supply the selected OAuth team ID in the teamId header.'})
            fields[p['in']+'.'+p['name']]=p
    bodies=[t for t in tables if t['location']=='body']
    if bodies:
        for row in bodies[0]['rows'][1:]:
            p=field(row,'body',bodies[0]['sourceUrl'])
            if path.stem=='app-info-update':
                condition=None
                if p['name'] in ('price','priceDetail') and '当isFree为0时，此参数必选' in p['description']:
                    condition='isFree'
                if p['name'] in ('registeredEntity','registeredEntityName') and '当registeredDclType值为0' in p['description'] and '该字段必填' in p['description']:
                    condition='registeredDclType'
                if condition:
                    rules.append({'id':'conditional-'+p['name'],'kind':'required_if','field':'body.'+p['name'],'when':{'field':'body.'+condition,'equals':'0'},'sourceUrl':p['sourceUrl'],'message':p['name']+' is required when '+condition+'=0.','suggestion':'Supply '+p['name']+' for the selected application settings.'})
            if path.stem in ('app-info-update','app-submit-with-file') and p['name'] in ('deviceTypes','phasedReleaseInfo'):
                if len(bodies)!=2:raise ValueError('Nested legacy model needs review')
                key='itemsProperties' if p['type']=='array' else 'properties'
                p[key]=[field(child,'body',bodies[1]['sourceUrl']) for child in bodies[1]['rows'][1:]]
            if path.stem=='app-submit' and p['required']:
                # The page explicitly makes Body release fields conditional on
                # Query releaseType=3; the empty commercial request is its example.
                if not any('当releaseType为3时' in note['text'] for note in doc['notes']):
                    raise ValueError('Conditional legacy submit Body needs review')
                p['required']=False
                rules.append({'id':'phased-'+p['name'],'kind':'required_if','field':'body.'+p['name'],'when':{'field':'query.releaseType','equals':'3'},'sourceUrl':bodies[0]['sourceUrl'],'message':p['name']+' is required for phased releaseType=3.','suggestion':'Supply this documented phased-release field or choose the intended release type.'})
            fields['body.'+p['name']]=p
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
    if rules:path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
