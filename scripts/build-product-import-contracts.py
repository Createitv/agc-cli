#!/usr/bin/env python3
"""Build reviewed ProductImportInfo contracts from browser-captured model tables.

Only explicit types, String(N) lengths and manually reviewed enums/conditions
are encoded. Referenced models with no reviewed schema retain object shape.
"""
import json
import re
from pathlib import Path

REFERENCE = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/supplemental'
MODELS = {'defaultPriceInfo': 'defaultproductpriceinfo', 'languages': 'productlanguageinfo', 'prices': 'productpriceinfo', 'promotion': 'promotioninfo'}


def fields(document):
    result = []
    for name, required, typ, description in document['tables'][0]['rows'][1:]:
        if typ.startswith('String'):
            kind = 'string'
        elif typ.startswith('Integer') or typ == 'Long':
            kind = 'integer'
        elif typ.startswith('List<'):
            kind = 'array'
        else:
            # These named objects are explicit model references in the table.
            if not any(link['title'] == typ for link in document['links']):
                raise ValueError(f'Unreviewed type {typ} in {document["sourceUrl"]}')
            kind = 'object'
        field = {'name': name, 'in': 'body', 'type': kind, 'required': required == 'M', 'description': description, 'sourceUrl': document['sourceUrl']}
        if kind == 'array':
            field['itemsType'] = 'string' if typ.startswith('List<String') else 'object'
        length = re.fullmatch(r'String\((\d+)\)', typ)
        if length:
            field['maxLength'] = int(length[1])
        result.append(field)
    return result


for platform in ('android', 'harmonyosnext'):
    document = json.loads((REFERENCE / f'product-import-{platform}.json').read_text())
    parameters = fields(document)
    for field in parameters:
        name = field['name']
        if name == 'purchaseType':
            field['enum'] = ['consumable', 'auto_subscription', 'non_consumable'] + (['non-renewing-subscription'] if platform == 'harmonyosnext' else [])
        if name == 'status':
            field['enum'] = ['drafts'] if platform == 'harmonyosnext' else ['active', 'inactive']
        if name == 'subPeriodUnit':
            field['enum'] = ['D', 'W', 'M', 'Y']
        if name in ('force', 'autoAddSaleCountry'):
            field['enum'] = ['0', '1']
        if name in MODELS:
            child = json.loads((REFERENCE / f'{MODELS[name]}-{platform}.json').read_text())
            field['itemsProperties' if field['type'] == 'array' else 'properties'] = fields(child)
    rules = [{'id': 'subscription-' + name, 'kind': 'required_if', 'field': 'body.' + name, 'when': {'field': 'body.purchaseType', 'equals': 'auto_subscription'}, 'sourceUrl': document['sourceUrl'], 'message': name + ' is required for auto_subscription', 'suggestion': 'Supply ' + name + ' using the subscription group and period documented by ProductImportInfo.'} for name in ('subGroupId', 'subPeriod', 'subPeriodUnit')]
    (REFERENCE / f'product-import-{platform}.contract.json').write_text(json.dumps({'parameters': parameters, 'rules': rules}, ensure_ascii=False, indent=2) + '\n')
