#!/usr/bin/env python3
"""Normalize reviewed Chinese Header/Body sections without flattening models."""
import json
import re
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference/gameplay'

def field(row,loc,source):
    name,required,typ,description=row
    kind='string' if typ.startswith('String') else 'integer' if typ in ('Integer','Long') else 'array' if typ=='List<UserId>' else None
    if kind is None or required not in ('是','否'):raise ValueError(f'Unreviewed gameplay field {name}:{typ}')
    p={'name':name,'in':loc,'type':kind,'required':required=='是','description':description,'sourceUrl':source}
    length=re.fullmatch(r'String\((\d+)\)',typ)
    if length:p['maxLength']=int(length[1])
    if typ=='Long':p['format']='int64'
    if name=='userIdType':p['enum']=['1','2','3','4','5','6']
    if name=='pageNo':p['minimum']=1
    if name=='pageSize':p['maximum']=100
    if kind=='array':p.update(itemsType='object')
    if name=='method':p['enum']=['external.hms.gs.syncGameResources']
    return p

for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:continue
    doc=json.loads(path.read_text());tables=[t for t in doc['tables'] if t['section']=='请求参数'];fields=[]
    if path.stem=='resources':
        if 'application/x-www-form-urlencoded' not in doc['examples'][0]['text']:raise ValueError('Resource wire format needs review')
        fields=[field(row,'body',tables[0]['sourceUrl']) for row in tables[0]['rows'][1:]]
        children=[]
        for row in tables[1]['rows'][1:]:
            child=field(row,'body',tables[1]['sourceUrl'])
            if child['name']=='type':child['enum']=['1','2']
            if child['name']=='operType':child['enum']=['0','1']
            children.append(child)
        for p in fields:
            if p['name']=='resources':
                p['encodedJson']={'type':'array','itemsType':'object','maxItems':50,'itemsProperties':children,'sourceUrl':tables[1]['sourceUrl'],'itemsRules':[{'id':'upload-resource-url','kind':'required_if','field':'body.url','when':{'field':'body.operType','equals':'0'},'sourceUrl':tables[1]['sourceUrl'],'message':'url is required when uploading a resource.','suggestion':'Supply a publicly downloadable resource URL with its file extension.'}]}

    else:
        if [row[0] for row in tables[0]['rows'][1:]]!=['appId','timestamp','nonce','sign']:raise ValueError('Callback header section needs review')
        example=next(x['text'] for x in doc['examples'] if x['section']=='请求示例')
        if 'headers:' not in example or 'Content-Type: application/json' not in example:raise ValueError('Callback Header/Body example needs review')
        fields=[field(row,'header',tables[0]['sourceUrl']) for row in tables[0]['rows'][1:]]+[field(row,'body',tables[1]['sourceUrl']) for row in tables[1]['rows'][1:]]
        if path.stem=='basic-info':
            for p in fields:
                if p['name']=='userIds':p['itemsProperties']=[field(row,'body',tables[2]['sourceUrl']) for row in tables[2]['rows'][1:]]
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(fields,ensure_ascii=False,indent=2)+'\n')
