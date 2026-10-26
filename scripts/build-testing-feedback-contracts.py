#!/usr/bin/env python3
"""Replace logical wrapper labels with browser-reviewed literal root fields."""
import json
import re
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference'
for endpoint,model_name in [('test-api-query-feedback','testing-feedback-list-request'),('test-api-delete-feedback','testing-delete-feedback-request')]:
    doc=json.loads((ROOT/'testing'/f'{endpoint}.json').read_text())
    model=json.loads((ROOT/'supplemental'/f'{model_name}.json').read_text())
    fields={}
    for table in doc['tables']:
        if table['section']=='请求参数' and table['location']=='header':
            for name,required,typ,description in table['rows'][1:]:
                if name=='appId':fields['header.appId']={'name':name,'in':'header','type':'string','required':required=='M','description':description,'sourceUrl':table['sourceUrl'],'maxLength':32}
    for name,required,typ,description in model['tables'][0]['rows'][1:]:
        kind='string' if typ.startswith('String') else 'integer' if typ=='Integer' else 'array' if typ.startswith('List<String') or name=='delFeedBackInfos' else None
        if kind is None:raise ValueError(f'Unreviewed feedback type {typ}')
        p={'name':name,'in':'body','type':kind,'required':required=='M','description':description,'sourceUrl':model['sourceUrl']}
        if kind=='array':
            p['itemsType']='object' if name=='delFeedBackInfos' else 'string'
            length=re.fullmatch(r'List<String\((\d+)\)>',typ)
            if length:p['itemsMaxLength']=int(length[1])
        if name=='deleteMode':p['enum']=['1','2']
        if name=='delFeedBackInfos':
            if not any(re.search(r'"delFeedBackInfos"\s*:\s*\[',x['text']) for x in doc['examples']):raise ValueError('Deletion array requires a reviewed request example')
            child=json.loads((ROOT/'supplemental/testing-delete-feedback-info.json').read_text())
            p['itemsProperties']=[{'name':row[0],'in':'body','type':'string','required':row[1]=='M','maxLength':int(re.fullmatch(r'String\((\d+)\)',row[2])[1]),'description':row[3],'sourceUrl':child['sourceUrl']} for row in child['tables'][0]['rows'][1:]]
        fields['body.'+name]=p
    (ROOT/'testing'/f'{endpoint}.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
