#!/usr/bin/env python3
"""Normalize literal Testing request tables; logical wrappers need model review."""
import json
import re
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]/'pkg/domain/reference/testing'
WRAPPERS={'test-api-create-open-test-link','test-api-delete-feedback','test-api-query-feedback'}

def extension_fields(model_name):
    model=json.loads((ROOT.parent/'supplemental'/f'{model_name}.json').read_text())
    result=[]
    for name,required,typ,description in model['tables'][0]['rows'][1:]:
        kind='string' if typ.startswith('String') else 'integer' if typ.startswith('Integer') else 'array' if typ.startswith('List<') else 'object' if typ=='AppReviewContact' else None
        if kind is None:raise ValueError(f'Unreviewed extension type {typ}')
        p={'name':name,'in':'body','type':kind,'required':required=='M','description':description,'sourceUrl':model['sourceUrl']}
        length=re.fullmatch(r'String\((\d+)\)',typ)
        if length:p['maxLength']=int(length[1])
        if kind=='integer' and '取值范围：' in description:
            enum_text=re.split(r'说明|注意|默认值',description.split('取值范围：',1)[1])[0]
            values=re.findall(r'(?:^|\s)(-?\d+)：',enum_text)
            if values:p['enum']=sorted(set(values),key=int)
        if kind=='array':
            p['itemsType']='integer' if typ.startswith('List<Integer') else 'object' if typ=='List<GroupInfo>' else 'string'
            if typ=='List<GroupInfo>':p['itemsProperties']=extension_fields('testing-group-info')
            if p['itemsType']=='integer' and '取值范围：' in description:
                enum_text=re.split(r'说明|注意|默认值',description.split('取值范围：',1)[1])[0]
                values=re.findall(r'(?:^|\s)(-?\d+)：',enum_text)
                if values:p['itemsEnum']=sorted(set(values),key=int)
            item_length=re.fullmatch(r'List<String\((\d+)\)>',typ)
            if item_length:p['itemsMaxLength']=int(item_length[1])
            maximum=re.search(r'数组长度不超过(\d+)',description)
            if maximum:p['maxItems']=int(maximum[1])
        if kind=='object':p['properties']=extension_fields('testing-app-review-contact')
        result.append(p)
    return result
for path in sorted(ROOT.glob('*.json')):
    if '.' in path.stem or path.stem in WRAPPERS-{'test-api-create-open-test-link'}:continue
    doc=json.loads(path.read_text());fields={};rules=[]
    for table in doc['tables']:
        loc=table['location']
        if table['section']!='请求参数' or loc not in ('header','query','body','path'):continue
        # This page calls Body "request" but its example has an unwrapped
        # groupId/limit object. Requiredness is undocumented: review headers only
        # rather than inventing a wrapper or mandatory example fields.
        if path.stem=='test-api-create-open-test-link' and loc=='body':continue
        for name,required,typ,description in table['rows'][1:]:
            if name.lower() in ('authorization','client_id','oauth2token'):continue
            kind='string' if typ.startswith('String') else 'integer' if typ.lower().startswith('integer') else 'array' if typ.startswith('List<') else 'object' if any(x['title']==typ for x in doc['links']) else None
            if kind is None or required not in ('M','O'):raise ValueError(f'Unreviewed Testing field {name}:{typ} in {path}')
            p={'name':name,'in':loc,'type':kind,'required':required=='M','description':description,'sourceUrl':table['sourceUrl']}
            if name=='teamId':
                p['required']=False
                rules.append({'id':'oauth-team','kind':'required_with','field':'header.teamId','fields':['header.oauth2Token'],'sourceUrl':table['sourceUrl'],'message':'OAuth authorization requires its selected teamId.','suggestion':'Supply the selected OAuth team ID in the teamId header.'})
            length=re.fullmatch(r'String\((\d+)\)',typ)
            if length:p['maxLength']=int(length[1])
            if kind=='array':
                p['itemsType']='string' if typ.startswith('List<String') else 'object'
                length=re.fullmatch(r'List<String\((\d+)\)>',typ)
                if length:p['itemsMaxLength']=int(length[1])
                maximum=re.search(r'数组长度不超过(\d+)',description)
                if maximum:p['maxItems']=int(maximum[1])
            if kind=='integer':
                bounds=re.search(r'\[(\d+),(\d+)\]',description)
                if bounds:p.update(minimum=int(bounds[1]),maximum=int(bounds[2]))
                text=description.split('默认值')[0].split('说明')[0]
                enums=re.findall(r'(?:取值范围：|\s)(\d+)[:：]',text)
                if enums:p['enum']=list(dict.fromkeys(enums))
            if path.stem=='test-api-modify-test-version' and name=='languages':
                model=json.loads((ROOT.parent/'supplemental/testing-app-language.json').read_text())
                children=[]
                for child_name,child_required,child_type,child_description in model['tables'][0]['rows'][1:]:
                    child_kind='array' if child_type=='List<DeviceMaterial>' else 'string'
                    child={'name':child_name,'in':'body','type':child_kind,'required':child_required=='M','description':child_description,'sourceUrl':model['sourceUrl']}
                    length=re.fullmatch(r'String\((\d+)\)',child_type)
                    if length:child['maxLength']=int(length[1])
                    minimum=re.search(r'不少于(\d+)个字符',child_description)
                    if minimum:child['minLength']=int(minimum[1])
                    if child_kind=='array':
                        child.update(itemsType='object',maxItems=6)
                        material=json.loads((ROOT.parent/'supplemental/testing-device-material.json').read_text())
                        material_fields=[]
                        for field_name,field_required,field_type,field_description in material['tables'][0]['rows'][1:]:
                            field_kind='integer' if field_type=='Integer' else 'array' if field_type.startswith('List<') else 'string' if field_type.startswith('String(') else None
                            if field_kind is None:raise ValueError(f'Unreviewed DeviceMaterial type: {field_type}')
                            field={'name':field_name,'in':'body','type':field_kind,'required':field_required=='M','description':field_description,'sourceUrl':material['sourceUrl']}
                            if field_kind=='integer':
                                field['enum']=re.findall(r'(\d+)：',field_description)
                            elif field_kind=='array':
                                field.update(itemsType='string',itemsMaxLength=300,maxItems=10)
                                if field_name=='screenShots':field['minItems']=3
                            else:
                                field.update(maxLength=300,minLength=1)
                            material_fields.append(field)
                        child['itemsProperties']=material_fields
                    children.append(child)
                p['itemsProperties']=children
            if path.stem=='test-api-add-test-package' and name=='file':
                model=json.loads((ROOT.parent/'supplemental/testing-package-file-info.json').read_text())
                p['properties']=[]
                for field_name,field_required,field_type,field_description in model['tables'][0]['rows'][1:]:
                    length=re.fullmatch(r'String\((\d+)\)',field_type)
                    if not length:raise ValueError(f'Unreviewed FileInfo type: {field_type}')
                    p['properties'].append({'name':field_name,'in':'body','type':'string','required':field_required=='M','maxLength':int(length[1]),'description':field_description,'sourceUrl':model['sourceUrl']})
            if path.stem=='test-api-modify-test-version' and name in ('deviceTypes','openTestInfo','packagePermissionIntroVideoList'):
                model_name={'deviceTypes':'testing-device-type-info','openTestInfo':'testing-open-test-info','packagePermissionIntroVideoList':'testing-package-permission-intro-video'}[name]
                model=json.loads((ROOT.parent/'supplemental'/f'{model_name}.json').read_text())
                children=[]
                for field_name,field_required,field_type,field_description in model['tables'][0]['rows'][1:]:
                    field_kind='integer' if field_type.startswith('Integer') else 'string' if field_type.startswith('String') else 'object' if field_type=='TestTaskInfo' else None
                    if field_kind is None:raise ValueError(f'Unreviewed {model_name} type: {field_type}')
                    child={'name':field_name,'in':'body','type':field_kind,'required':field_required=='M','description':field_description,'sourceUrl':model['sourceUrl']}
                    length=re.fullmatch(r'String\((\d+)\)',field_type)
                    if length:child['maxLength']=int(length[1])
                    if field_name=='deviceType':
                        child['enum']=re.findall(r'(\d+)：',field_description.split('说明')[0])
                    if field_type=='TestTaskInfo':
                        child['properties']=extension_fields('testing-test-task-info')
                    children.append(child)
                p['properties' if name=='openTestInfo' else 'itemsProperties']=children
            if path.stem=='test-api-modify-test-version' and name=='extInfo':
                p['properties']=extension_fields('testing-version-ext-info')
                source=next(child['sourceUrl'] for child in p['properties'] if child['name']=='aiTypes')
                rules.append({'id':'testing-ai-types','kind':'required_if','field':'body.extInfo.aiTypes','when':{'field':'body.extInfo.isAiGenerate','equals':'1'},'sourceUrl':source,'message':'AI-generated services require aiTypes.','suggestion':'Supply the documented AI service types in extInfo.aiTypes.'})
            if path.stem=='test-api-modify-test-version' and name=='versionVocReq':
                p['itemsProperties']=extension_fields('testing-version-voc-req')
            fields[loc+'.'+name]=p
    if path.stem=='test-api-delete-invite-code':
        example=next(x for x in doc['examples'] if x['section']=='请求示例')
        if 'invitationCodeId=' not in example['text']:
            raise ValueError('Invitation deletion query needs a reviewed wire example')
        # The parameter table omits this existing verified required identifier;
        # the official request example confirms its Query location.
        fields['query.invitationCodeId']={'name':'invitationCodeId','in':'query','type':'string','required':True,'description':'待删除的邀请码ID；参数表缺漏，Query位置由官方请求示例确认。','sourceUrl':example['sourceUrl']}
    path.with_name(path.stem+'.parameters.json').write_text(json.dumps(list(fields.values()),ensure_ascii=False,indent=2)+'\n')
    if rules:path.with_name(path.stem+'.rules.json').write_text(json.dumps(rules,ensure_ascii=False,indent=2)+'\n')
