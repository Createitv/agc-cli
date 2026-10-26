#!/usr/bin/env python3
"""Normalize each browser-reviewed report table and its explicit date rules."""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/reports'
for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:
        continue
    doc = json.loads(path.read_text())
    parameters, rules = [], []
    for table in doc['tables']:
        location = table['location']
        if table['section'] != '请求参数' or location not in ('path', 'query'):
            continue
        for name, required, typ, description in table['rows'][1:]:
            kind = 'array' if typ.startswith('List<String>') else 'integer' if typ.startswith('Integer') else 'string' if typ.startswith('String') else None
            if kind is None or required not in ('M', 'O'):
                raise ValueError(f'Unreviewed report parameter: {path}, {name}, {typ}')
            p = {'name': name, 'in': location, 'type': kind, 'required': required == 'M', 'description': description, 'sourceUrl': table['sourceUrl']}
            if kind == 'array': p['itemsType'] = 'string'
            maximum = re.fullmatch(r'String\((\d+)\)', typ)
            if maximum: p['maxLength'] = int(maximum[1])
            if name == 'language': p['enum'] = ['zh-CN'] if path.stem == 'gamereservationexport' else ['zh-CN', 'en-US', 'ru-RU']
            if name == 'currency': p['enum'] = ['CNY', 'EUR', 'USD']
            if name == 'exportType': p['enum'] = ['CSV', 'EXCEL']
            if name in ('userTag', 'queryType'): p['enum'] = ['0', '1']
            if name == 'startTime':
                days = re.search(r'不超过(\d+)天', description)
                if not days: raise ValueError(f'Missing explicit interval limit: {path}')
                rules.append({'id':'report-date-window','kind':'date_window','fields':['query.startTime','query.endTime'],'maxDays':int(days[1]),'sourceUrl':table['sourceUrl'],'message':'Report dates must be real YYYYMMDD dates in chronological order within the documented query interval.','suggestion':f'Use YYYYMMDD dates with endTime at or after startTime and no more than {days[1]} days apart.'})
            parameters.append(p)
    if any(p['name']=='filterConditionValue' for p in parameters):
        rules.append({'id':'ordered-filter-pairs','kind':'equal_length','fields':['query.filterCondition','query.filterConditionValue'],'sourceUrl':next(p['sourceUrl'] for p in parameters if p['name']=='filterConditionValue'),'message':'Each filterCondition must have a filterConditionValue at the same position.','suggestion':'Repeat both query keys once per filter, preserving matching order.'})
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(parameters,ensure_ascii=False,indent=2)+'\n')
    path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
