#!/usr/bin/env python3
"""Render the reviewed Publishing pages as a CLI reference, without network/auth.

This is documentation rendering, not contract validation. Keep source tables
intact, including conditional-requiredness remarks and authentication branches.
"""
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def cell(value):
    return value.replace('|', r'\|').replace('\n', '<br>')


def render():
    audit = json.loads((ROOT / 'docs/verification/publishing-contract-audit-current.json').read_text())
    output = [
        '# Publishing API 逐项 CLI 参考', '',
        '本文从浏览器保存的官方正文生成，每项说明用途、传输位置、参数和前置条件。'
        '生成命令：`python3 scripts/render-publishing-guide.py`。', '',
        '流程和请求文件示例见 [Publishing 使用与核验](PUBLISHING_WORKFLOWS.md)。'
        '表格保留官方必选标记；条件要求和表格/示例冲突需要同时阅读说明。'
        '本文件不代表所有模型规则已实现或全部线上操作已验收。', '',
        '## 参数位置', '',
        '`query` 使用 `--query key=value`，`header` 使用 `--header key=value`，'
        '`body` 使用 JSON 文件 `--body request.json` 或已声明类型的 `--field key=value`。'
        'JSON 整数不加引号，对象和数组按实际模型展开。', '',
        '实际调用选择 `--profile PROFILE`，先使用 `--invoke --dry-run` 查看请求计划；'
        '只有准备实际调用时才使用 `--invoke --dry-run=false`。'
        'dry-run 不证明资源状态、权限或审核结果。', '',
    ]
    for endpoint in audit['endpoints']:
        command = endpoint['command']
        doc = json.loads((ROOT / 'pkg/domain/reference/publishing' / (command + '.json')).read_text())
        output += [f"## {doc['title']}", '',
                   f"命令：`agc publishing {command}`。请求：`{endpoint['method']} {endpoint['path']}`。", '',
                   f"[官方正文]({doc['sourceUrl']})；官网更新：{doc.get('updatedAt', '未标明')}；"
                   f"浏览器读取：{doc.get('checkedAt', '未标明')}。", '']
        for note in doc.get('notes', []):
            if note['section'] in ('功能介绍', '使用约束'):
                output += [note['text'], '']
        for table in doc.get('tables', []):
            if table.get('section') != '请求参数':
                continue
            label = table.get('location') or '参数'
            if label == 'header':
                label += '（' + table.get('authenticationBranch', '认证方式以正文说明为准') + '）'
            output += [f"### {label}", '', f"[参数来源]({table['sourceUrl']})", '']
            rows = table['rows']
            if not rows:
                continue
            output += ['| ' + ' | '.join(cell(x) for x in rows[0]) + ' |',
                       '| ' + ' | '.join('---' for _ in rows[0]) + ' |']
            output += ['| ' + ' | '.join(cell(x) for x in row) + ' |' for row in rows[1:]]
            output += ['']
        remarks = [note for note in doc.get('notes', []) if note['section'] == '请求参数'
                   and not note['text'].endswith('方式：')]
        if remarks:
            output += ['### 请求说明', '']
            for note in remarks:
                output += [note['text'], '']
        output += ['### 完整模型及示例', '', '```sh',
                   f'agc publishing {command} --docs --pretty', '```', '',
                   '该命令返回官方请求/响应示例及已捕获的关联模型。'
                   '业务响应中的 `ret.code` 应与 HTTP 状态分别检查；'
                   '修改后使用相应查询接口确认目标状态。', '']
    (ROOT / 'docs/PUBLISHING_API_REFERENCE.md').write_text('\n'.join(output), encoding='utf-8')


if __name__ == '__main__':
    render()
