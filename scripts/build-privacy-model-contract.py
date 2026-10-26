#!/usr/bin/env python3
"""Compile captured privacy models by exact official links, never model-name guesses."""
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1] / 'pkg/domain/reference/supplemental'
documents = {}
for path in ROOT.glob('*.json'):
    doc = json.loads(path.read_text())
    if doc.get('status') == 'browser-read':
        documents[doc['sourceUrl'].split('#')[0]] = doc


def fields(url, stack=()):
    if url in stack:
        raise ValueError(f'Recursive model requires explicit review: {url}')
    doc = documents[url]
    tables = [t for t in doc['tables'] if t['rows'] and t['rows'][0][0] == '参数名称']
    if len(tables) != 1:
        raise ValueError(f'Model requires one reviewed parameter table: {url}')
    result = []
    for name, required, typ, description in tables[0]['rows'][1:]:
        if required not in ('M', 'O'):
            raise ValueError(f'Unknown requiredness: {name}')
        p = dict(name=name, **{'in': 'body'}, required=required == 'M',
                 description=description, sourceUrl=tables[0]['sourceUrl'])
        array = re.fullmatch(r'List<(.+)>', typ)
        inner = array[1] if array else typ
        length = re.fullmatch(r'String(?:\((\d+)\))?', inner)
        if length:
            kind = 'string'
            if length[1]:
                p['itemsMaxLength' if array else 'maxLength'] = int(length[1])
        elif inner.startswith('Integer'):
            kind = 'integer'
        elif inner == 'boolean':
            kind = 'boolean'
        else:
            kind = 'object'
            links = {x['sourceUrl'].split('#')[0] for x in doc['links']
                     if x['title'] == inner or x['title'].startswith(inner + '（')}
            if len(links) != 1 or next(iter(links)) not in documents:
                raise ValueError(f'Uncaptured or ambiguous model {inner} in {url}')
            p['itemsProperties' if array else 'properties'] = fields(next(iter(links)), stack + (url,))
        p['type'] = 'array' if array else kind
        if array:
            p['itemsType'] = kind
            maximum = re.search(r'数组长度不超过(\d+)', description)
            interval = re.search(r'数组长度为(\d+)~(\d+)', description)
            if maximum:
                p['maxItems'] = int(maximum[1])
            if interval:
                p.update(minItems=int(interval[1]), maxItems=int(interval[2]))
        if kind == 'integer' and '取值范围：' in description:
            text = re.split(r'默认值|说明|注意', description.split('取值范围：', 1)[1])[0]
            values = re.findall(r'(?:^|\s)(-?\d+)：', text)
            if values:
                p['itemsEnum' if array else 'enum'] = sorted(set(values), key=int)
        if name == 'lang' and doc['title'] == 'AgreementInfo（隐私政策协议）':
            p['enum'] = re.findall(r'(?:^|\s)([a-zA-Z-]+)：', description.split('取值范围：', 1)[1])
        result.append(p)
    return result


source = 'https://developer.huawei.com/consumer/cn/doc/doccenter-submission/agc-help-publish-api-data-agreementinfo-privacy-0000002294845630'
(ROOT / 'privacy-agreement-info.contract.json').write_text(
    json.dumps({'sourceUrl': source, 'parameters': fields(source)}, ensure_ascii=False, indent=2) + '\n')
