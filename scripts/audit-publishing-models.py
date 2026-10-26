#!/usr/bin/env python3
"""Audit request-model fields against browser-read source tables, offline.

Never read user authorization or call Huawei. This checks nested shape, scalar
types, requiredness and explicit string/array bounds. It does not certify
resource-dependent rules, signatures, roles, dates, files or real operations.
"""
import argparse
from datetime import datetime, timezone
import json
from pathlib import Path
import re
import subprocess

ROOT = Path(__file__).resolve().parents[1]


def kind(value):
    if value.startswith('String'): return 'string'
    if value.startswith('Integer'): return 'integer'
    if value.lower() == 'boolean': return 'boolean'
    if value.startswith('List<'): return 'array'
    return 'object'


def audit_model(model, fields, path, documents, endpoint_id, stack=()):
    issues = []
    actual = {field['name']: field for field in fields}
    for table in model.get('tables', []):
        for row in table.get('rows', []):
            if len(row) < 4 or row[1] not in ('M', 'O'): continue
            name, required, typ, description = row[:4]
            field_path = path + '.' + name
            source = table.get('sourceUrl', model['sourceUrl'])

            def issue(code, expected=None, got=None):
                issues.append({'field': field_path, 'code': code, 'sourceUrl': source,
                               'expected': expected, 'actual': got})

            field = actual.get(name)
            if field is None:
                issue('missing_field', typ)
                continue
            expected_kind = kind(typ)
            if field.get('type') != expected_kind:
                issue('type', expected_kind, field.get('type'))
            expected_required = required == 'M'
            # TagQualInfo explicitly makes taskId required when modifying.
            if endpoint_id == 'qualification-review-update' and name == 'taskId':
                expected_required = True
            if field.get('required', False) != expected_required:
                issue('requiredness', expected_required, field.get('required', False))
            match = re.fullmatch(r'String\((\d+)\)', typ)
            if match and field.get('maxLength') != int(match[1]):
                issue('max_length', int(match[1]), field.get('maxLength'))
            match = re.search(r'数组长度不超过(\d+)', description)
            if match and field.get('maxItems') != int(match[1]):
                issue('max_items', int(match[1]), field.get('maxItems'))
            match = re.search(r'数组长度为(\d+)~(\d+)', description)
            if match:
                for key, value in [('minItems', int(match[1])), ('maxItems', int(match[2]))]:
                    if field.get(key, 0) != value: issue(key, value, field.get(key, 0))
            if expected_kind == 'array':
                item_type = typ[5:-1]
                item_kind = kind(item_type)
                if field.get('itemsType') != item_kind:
                    issue('items_type', item_kind, field.get('itemsType'))
                match = re.fullmatch(r'String\((\d+)\)', item_type)
                if match and field.get('itemsMaxLength') != int(match[1]):
                    issue('item_max_length', int(match[1]), field.get('itemsMaxLength'))
                nested_fields = field.get('itemsProperties', [])
                nested_path = field_path + '[]'
            else:
                item_type = typ
                nested_fields = field.get('properties', [])
                nested_path = field_path
            if kind(item_type) != 'object': continue
            # Resolve by the table's actual model link, never by sidebar or
            # previous/next labels; user/privacy model names can be identical.
            link = next((link for link in model.get('links', []) if link['title'] == item_type), None)
            if link is None:
                issue('missing_model_link', item_type)
                continue
            url = link['sourceUrl'].split('#', 1)[0]
            nested = documents.get(url)
            if nested is None:
                issue('model_not_captured', url)
                continue
            if url in stack:
                issue('recursive_model_requires_review', url)
                continue
            issues += audit_model(nested, nested_fields, nested_path, documents, endpoint_id, (*stack, url))
    return issues


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--cli-json', type=Path)
    args = parser.parse_args()
    raw = args.cli_json.read_text() if args.cli_json else subprocess.check_output(
        ['go', 'run', './cmd/agc', 'publishing', 'endpoints'], cwd=ROOT, text=True)
    local = {endpoint['id']: endpoint for endpoint in json.loads(raw)['data']}
    documents = {}
    for file in sorted((ROOT / 'pkg/domain/reference/supplemental').glob('*.json')):
        doc = json.loads(file.read_text())
        if doc.get('status') != 'browser-read' or 'tables' not in doc: continue
        url = doc.get('sourceUrl')
        if url not in documents or doc.get('checkedAt', '') > documents[url].get('checkedAt', ''):
            documents[url] = doc
    coverage = json.loads((ROOT / 'docs/verification/publishing-contract-audit-current.json').read_text())
    results = []
    for operation in coverage['endpoints']:
        endpoint = local[operation['command']]
        doc = json.loads((ROOT / 'pkg/domain/reference/publishing' / (endpoint['id'] + '.json')).read_text())
        roots = []
        for table in doc.get('tables', []):
            if table.get('section') != '请求参数' or table.get('location') != 'body': continue
            for row in table.get('rows', []):
                if len(row) < 4 or row[1] not in ('M', 'O'): continue
                typ = row[2]
                is_list = typ.startswith('List<') or (endpoint['id'] == 'user-agreement-create' and row[0] == 'agreementInfos')
                model_type = typ[5:-1] if typ.startswith('List<') else typ
                if kind(model_type) != 'object': continue
                roots.append((row[0], model_type, is_list))
        issues = []
        for name, model_type, is_list in roots:
            field = next((field for field in endpoint['parameters'] if field['in'] == 'body' and field['name'] == name), {})
            link = next((link for link in doc['links'] if link['title'] == model_type), None)
            model = documents.get(link['sourceUrl'].split('#', 1)[0]) if link else None
            if model is None:
                issues.append({'field': 'body.' + name, 'code': 'model_not_captured', 'expected': model_type})
                continue
            nested = field.get('itemsProperties' if is_list else 'properties', [])
            issues += audit_model(model, nested, 'body.' + name + ('[]' if is_list else ''), documents, endpoint['id'])
        results.append({'endpoint': endpoint['id'], 'sourceUrl': doc['sourceUrl'], 'requestModelRoots': len(roots), 'issues': issues})
    print(json.dumps({'checkedAt': datetime.now(timezone.utc).isoformat(),
                      'scope': 'nested request-model field shape, type, requiredness and explicit string/array bounds; other rules and live operations require separate validation',
                      'summary': {'endpoints': len(results), 'endpointsWithModelIssues': sum(bool(r['issues']) for r in results), 'issues': sum(len(r['issues']) for r in results)},
                      'endpoints': results}, ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
