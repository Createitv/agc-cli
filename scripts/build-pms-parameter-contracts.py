#!/usr/bin/env python3
"""Normalize browser-reviewed PMS root parameter tables.

PMS pages list literal request properties; linked models stay separate.
The multipart resource endpoint is deliberately handled by its file adapter.
Authentication alternatives belong to credential handling, not required inputs.
"""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/pms'
AUTH_FIELDS = {'authorization', 'client_id', 'oauth2token', 'teamid'}
for path in sorted(ROOT.glob('*.json')):
    if path.stem.count('.'):
        continue
    doc = json.loads(path.read_text())
    parameters = {}
    for table in doc['tables']:
        location = table['location']
        if table['section'] != '请求参数' or location not in ('header', 'query', 'body', 'path'):
            continue
        for row in table['rows'][1:]:
            if len(row) != 4 or row[1] not in ('M', 'O'):
                raise ValueError(f'Unexpected parameter row in {path}: {row}')
            name, required, typ, description = row
            if location == 'header' and name.lower() in AUTH_FIELDS:
                continue
            if typ == 'MultipartFile':
                kind = 'file'
            elif typ.startswith('String'):
                kind = 'string'
            elif typ.startswith('Integer') or typ == 'Long':
                kind = 'integer'
            elif typ.startswith('List<'):
                kind = 'array'
            elif any(link['title'] == typ for link in doc['links']):
                kind = 'object'
            else:
                raise ValueError(f'Unreviewed parameter type {typ} in {path}')
            field = {'name': name, 'in': 'file' if kind == 'file' else location, 'type': kind, 'required': required == 'M', 'description': description, 'sourceUrl': table['sourceUrl']}
            if kind == 'array':
                field['itemsType'] = 'string' if typ.startswith('List<String') else 'object'
            length = re.fullmatch(r'String\((\d+)\)', typ)
            if length:
                field['maxLength'] = int(length[1])
            if path.stem == 'resourcemanage-harmonyosnext' and name in ('fileType', 'zoneId'):
                if '固定取值： 1：' not in description:
                    raise ValueError('Resource fixed value requires review')
                field['enum'] = ['1']
            if kind == 'array':
                maximum = re.search(r'列表中不超过(\d+)', description)
                if maximum:
                    field['maxItems'] = int(maximum[1])
            parameters[location + '.' + name] = field
    if not parameters:
        raise ValueError(f'No reviewed parameter tables found in {path}')
    path.with_name(path.stem + '.parameters.json').write_text(json.dumps(list(parameters.values()), ensure_ascii=False, indent=2) + '\n')
