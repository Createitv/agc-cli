#!/usr/bin/env python3
"""Normalize reviewed resource tables using their actual request-model nesting."""
import json
import re
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference/resources'
NESTED={'respackapi-addversion':('resourceVersion','properties'),'respackapi-modifyversion':('resourceVersion','properties'),'respackapi-uploadapplyfor':('filePartApplyInfo','itemsProperties'),'respackapi-uploadconfirm':('filePartConfirmInfos','itemsProperties')}

def field(row,loc,source):
    name,required,typ,description=row
    kind='string' if typ=='String' else 'integer' if typ.startswith('Int') or typ=='Long' else 'array' if typ.startswith('List<') else 'object' if typ in ('GWOpPackageVersion','GWOpPackageVersionEntity') else None
    if kind is None or required not in ('M','O'):raise ValueError(f'Unreviewed resource field {name}:{typ}')
    p={'name':name,'in':loc,'type':kind,'required':required=='M','description':description,'sourceUrl':source}
    maximum=re.search(r'(?:最大长度为?|最长)(\d+)个?字符',description)
    if kind=='string' and maximum:p['maxLength']=int(maximum[1])
    if name=='fileName' and '名称不能包含特殊字符' in description:
        p['pattern']=r'^[^|:/\\*?"<>]+$'
    if kind=='array':p['itemsType']='object'
    if name in ('isDraft','uploadResult'):p['enum']=['0','1']
    if name=='fileType':p['enum']=['0']
    if kind=='integer':
        bounds=re.search(r'取值范围[：:]?\[(\d+),(\d+)\]',description)
        if bounds:p.update(minimum=int(bounds[1]),maximum=int(bounds[2]))
        if name=='version' and '不超过9位' in description:p['maximum']=999999999
    if name=='filePartApplyInfo':p['maxItems']=1000
    return p

for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:continue
    doc=json.loads(path.read_text());fields={};rules=[]
    tables=[t for t in doc['tables'] if t['section']=='请求参数']
    for table in tables:
        if table['location'] not in ('header','query','path'):continue
        for row in table['rows'][1:]:
            if row[0] in ('client_id','Authorization'):continue
            p=field(row,table['location'],table['sourceUrl']);fields[p['in']+'.'+p['name']]=p
    bodies=[t for t in tables if t['location']=='body']
    if path.stem in NESTED:
        if len(bodies)!=3:raise ValueError(f'Resource model tables need review: {path}')
        # Table 0 names a logical request model. Table 1 is the JSON root;
        # table 2 defines its object/list member, as shown by request examples.
        name,property_key=NESTED[path.stem]
        if not any('"'+name+'"' in x['text'] for x in doc['examples'] if x['section']=='请求示例'):raise ValueError('Missing reviewed literal root request example')
        for row in bodies[1]['rows'][1:]:
            p=field(row,'body',bodies[1]['sourceUrl'])
            if row[0]==name:p[property_key]=[field(child,'body',bodies[2]['sourceUrl']) for child in bodies[2]['rows'][1:]]
            fields['body.'+p['name']]=p
    if path.stem=='respackapi-queryfile':
        if '取值范围：[1,500]' not in fields['query.fileIds']['description']:raise ValueError('Resource CSV limit needs review')
        rules.append({'id':'resource-file-id-count','kind':'csv_count','field':'query.fileIds','minItems':1,'maxItems':500,'sourceUrl':fields['query.fileIds']['sourceUrl'],'message':'fileIds requires 1 to 500 nonempty comma-separated file IDs.','suggestion':'Supply at most 500 file IDs separated by commas without empty entries, or select a versionId.'})
        rules.append({'id':'resource-file-identity','kind':'at_least_one','fields':['query.versionId','query.fileIds'],'sourceUrl':fields['query.fileIds']['sourceUrl'],'message':'Supply versionId or fileIds to select resource files.','suggestion':'Provide the resource version ID or a comma-separated file ID list.'})
        if '选择一个传入' not in fields['query.fileIds']['description']:raise ValueError('Resource identity selection needs review')
        rules.append({'id':'resource-file-identity-exclusive','kind':'at_most_one','fields':['query.versionId','query.fileIds'],'sourceUrl':fields['query.fileIds']['sourceUrl'],'message':'Select either versionId or fileIds, as documented.','suggestion':'Keep one resource-file selector and remove the other.'})
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
    if rules:path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
