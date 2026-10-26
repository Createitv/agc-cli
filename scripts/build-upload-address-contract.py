#!/usr/bin/env python3
"""Normalize the browser-reviewed OBS upload address query table."""
import json
import re
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference/upload'
doc=json.loads((ROOT/'upload-url-new.json').read_text())
fields=[]
for table in doc['tables']:
    if table['section']!='请求参数' or table['location']!='query':continue
    for name,required,typ,description in table['rows'][1:]:
        kind='string' if typ.startswith('String') else 'integer' if typ=='Integer' else None
        if kind is None:raise ValueError(f'Unreviewed upload type {typ}')
        p={'name':name,'in':'query','type':kind,'required':required=='M','description':description,'sourceUrl':table['sourceUrl']}
        length=re.fullmatch(r'String\((\d+)\)',typ)
        if length:p['maxLength']=int(length[1])
        if name=='releaseType':p['enum']=['1']
        if name=='chineseMainlandFlag':p['enum']=['0','1']
        fields.append(p)
(ROOT/'upload-url-new.parameters.json').write_text(json.dumps(fields,ensure_ascii=False,indent=2)+'\n')
