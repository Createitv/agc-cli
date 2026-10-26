#!/usr/bin/env python3
"""Normalize the five browser-reviewed domain endpoint root tables."""
import json
import re
from pathlib import Path
ROOT = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/domains'
for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:
        continue
    doc = json.loads(path.read_text())
    fields = {}
    for table in doc['tables']:
        location = table['location']
        if table['section'] != '请求参数' or location not in ('header', 'query', 'body'):
            continue
        for name, required, typ, description in table['rows'][1:]:
            if name.lower() in ('authorization', 'client_id'):
                continue
            kind = 'string' if typ.startswith('String') else 'array' if typ == 'List<DomainInfo>' else 'object' if typ == 'DomainInfo' else None
            if kind is None or required not in ('M', 'O'):
                raise ValueError(f'Unreviewed domain field: {path}, {name}, {typ}')
            field = {'name': name, 'in': location, 'type': kind, 'required': required == 'M', 'description': description, 'sourceUrl': table['sourceUrl']}
            length = re.fullmatch(r'String\((\d+)\)', typ)
            if length:
                field['maxLength'] = int(length[1])
            if kind == 'array':
                field['itemsType'] = 'object'
                maximum = re.search(r'集合长度最大(\d+)', description)
                if maximum:
                    field['maxItems'] = int(maximum[1])
            if name == 'category':
                field['enum'] = ['server', 'business']
            if kind in ('array', 'object'):
                model_name = 'domain-info-server' if kind == 'array' else 'domain-info-business'
                model = json.loads((ROOT.parent / 'supplemental' / (model_name+'.json')).read_text())
                nested = []
                for child_name, child_required, child_type, child_description in model['tables'][0]['rows'][1:]:
                    # errorCode/errorMsg describe responses, not submitted inputs.
                    if child_name in ('errorCode', 'errorMsg'):
                        continue
                    child = {'name':child_name,'in':'body','type':'string','required':child_required=='M','description':child_description,'sourceUrl':model['sourceUrl']}
                    child['maxLength'] = int(re.fullmatch(r'String\((\d+)\)',child_type)[1])
                    if child_name == 'type':
                        child['enum'] = ['webView'] if path.stem == 'domain-api-precheck' else ['httpRequest','upload','download','webSocket','webView']
                    nested.append(child)
                field['itemsProperties' if kind == 'array' else 'properties'] = nested
            fields[location+'.'+name] = field
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()), ensure_ascii=False, indent=2)+'\n')
