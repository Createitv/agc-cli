#!/usr/bin/env python3
"""Normalize all eight reviewed Android/HarmonyOS comment contracts."""
import json
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference/comments'
for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem:continue
    doc=json.loads(path.read_text());fields={};rules=[]
    for table in doc['tables']:
        loc=table['location']
        if table['section']!='请求参数' or loc not in ('query','path','header','body'):continue
        for name,required,typ,description in table['rows'][1:]:
            if name in ('client_id','Authorization'):continue
            kind='string' if typ.startswith('String') else 'integer' if typ in ('Integer','Long') else None
            if kind is None or required not in ('M','O'):raise ValueError(f'Unreviewed comment type {typ}')
            p={'name':name,'in':loc,'type':kind,'required':required=='M','description':description,'sourceUrl':table['sourceUrl']}
            if typ=='Long':p['format']='int64'
            if name=='limit':p['maximum']=100
            if name=='sort':p['enum']=['0','1','2']
            if name=='requestId':p['maxLength']=64
            if loc=='body' and name=='lang' and path.stem=='com-reviewreply-harmonyos':p['enum']=['zh','en','bo','ug']
            fields[loc+'.'+name]=p
    if path.stem=='com-reviewreply':
        rules.append({'id':'reply-target-exclusive','kind':'at_most_one','fields':['body.toReplyId','body.updateReplyId'],'sourceUrl':fields['body.updateReplyId']['sourceUrl'],'message':'toReplyId and updateReplyId cannot both be supplied.','suggestion':'Choose whether to reply to a user reply or modify an existing developer reply, and provide only its corresponding ID.'})
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
    if rules:path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
