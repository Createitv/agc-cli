#!/usr/bin/env python3
"""Build Provisioning root and linked model contracts from reviewed pages."""
import json
import re
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference/provisioning'

def parameter(name,required,typ,description,location,source):
    kind='string' if typ.startswith('String') else 'integer' if typ.startswith('Integer') else 'array' if typ.startswith('List<') else None
    if kind is None or required not in ('M','O'):
        raise ValueError(f'Unreviewed type {typ}: {source}')
    p={'name':name,'in':location,'type':kind,'required':required=='M','description':description,'sourceUrl':source}
    length=re.fullmatch(r'String\((\d+)\)',typ)
    if length:p['maxLength']=int(length[1])
    if kind=='array':
        p['itemsType']='string' if typ.startswith('List<String') else 'object'
        length=re.fullmatch(r'List<String\((\d+)\)>',typ)
        if length:p['itemsMaxLength']=int(length[1])
        bounds=re.search(r'数组长度为(\d+)~(\d+)',description)
        maximum=re.search(r'(?:数组长度不超过|最大长度)(\d+)',description)
        if bounds:p.update(minItems=int(bounds[1]),maxItems=int(bounds[2]))
        elif maximum:p['maxItems']=int(maximum[1])
    if kind=='integer':
        bounds=re.search(r'取值范围：\[(\d+),(\d+)\]',description)
        minimum=re.search(r'(?:最小值：|最小为)(\d+)',description)
        if bounds:p.update(minimum=int(bounds[1]),maximum=int(bounds[2]))
        elif minimum:p['minimum']=int(minimum[1])
    return p

for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:continue
    doc=json.loads(path.read_text());fields={};rules=[]
    for table in doc['tables']:
        loc=table['location']
        if table['section']!='请求参数' or loc not in ('header','query','body','path'):continue
        for name,required,typ,description in table['rows'][1:]:
            if name.lower() in ('authorization','client_id','oauth2token'):continue
            if name=='aclPermissionList' and path.stem=='provision-api-applyacl':
                # The table's singular model label contradicts its actual JSON example.
                if not any(re.search(r'"aclPermissionList"\s*:\s*\[',x['text']) for x in doc['examples']):
                    raise ValueError('ACL request array needs a reviewed wire example')
                typ='List<OpenApplyPermissionInfo>'
            p=parameter(name,required,typ,description,loc,table['sourceUrl'])
            if name=='teamId':
                p['required']=False
                rules.append({'id':'oauth-team','kind':'required_with','field':'header.teamId','fields':['header.oauth2Token'],'sourceUrl':table['sourceUrl'],'message':'OAuth authorization requires its selected teamId.','suggestion':'Supply the selected OAuth team ID in the teamId header.'})
            if name=='certType':p['enum']=['1','2','3','4']
            if name=='provisionType':p['enum']=['1','2','3','6']
            if name=='order':p['enum']=['1','2']
            if name=='lang':p['enum']=['zh','en','ru']
            model_name='add-device-info' if typ=='List<AddDeviceInfo>' else 'acl-apply-permission-info' if typ=='List<OpenApplyPermissionInfo>' else None
            if model_name:
                model=json.loads((ROOT.parent/'supplemental'/f'{model_name}.json').read_text())
                children=[parameter(*row,'body',model['sourceUrl']) for row in model['tables'][0]['rows'][1:]]
                for child in children:
                    if child['name']=='deviceType':child['enum']=['1','2','4','5','8','19']
                p['itemsProperties']=children
            fields[loc+'.'+name]=p
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
    if rules:path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
