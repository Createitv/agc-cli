#!/usr/bin/env python3
"""Normalize the two reviewed form callback parameter tables."""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/game-items'
for endpoint in ('propapi-getroleinfo', 'propapi-order'):
    doc = json.loads((ROOT / f'{endpoint}.json').read_text())
    fields = []
    for table in doc['tables']:
        if table['section'] != '请求参数':
            continue
        for name, required, typ, description in table['rows'][1:]:
            if required not in ('M', 'O') or not (typ.startswith('String') or typ in ('Int', 'Long')):
                raise ValueError(f'Unreviewed callback field: {endpoint}, {name}, {typ}')
            field = {'name': name, 'in': 'body', 'required': required == 'M', 'type': 'string' if typ.startswith('String') else 'integer', 'description': description, 'sourceUrl': table['sourceUrl']}
            length = re.fullmatch(r'String\((\d+)\)', typ)
            if length:
                field['maxLength'] = int(length[1])
            if name == 'userIdType':
                field['enum'] = ['1', '2', '3', '4', '5', '6']
            fields.append(field)
    if not fields:
        raise ValueError(f'Missing reviewed parameter table: {endpoint}')
    (ROOT / f'{endpoint}.parameters.json').write_text(json.dumps(fields, ensure_ascii=False, indent=2) + '\n')
