#!/usr/bin/env python3
"""Import Computer Use captures from browser-capture-reference.mjs's schema.

No network, credentials or hidden browser state. Re-run with the JSON array saved
from the browser's rendered documents. Semantic overrides below are source-reviewed;
this is not a general-purpose inference of Huawei API contracts.
"""
import copy,json,re,sys
from pathlib import Path
ROOT=Path(__file__).resolve().parents[1]
D=json.loads(Path(sys.argv[1]).read_text())
IDS={
 'agcapi-cloudhost-newapp':'hosting/create-version','agcapi-cloudhost-modifyversion':'hosting/modify-files','agcapi-cloudhost-uploadurl':'hosting/upload-url','agcapi-cloudhost-uploadversion':'hosting/upload-version','agcapi-cloudhost-archivedversion':'hosting/archive-version','agcapi-cloudhost-releaseversion':'hosting/release-version','agcapi-cloudhost-queryversion':'hosting/list-versions','agcapi-cloudhost-queryfile':'hosting/list-version-files','agcapi-cloudhost-deleteversion':'hosting/delete-version','agcapi-cloudhost-createnvir':'hosting/create-environment','agcapi-cloudhost-queryenvir':'hosting/list-environments','agcapi-cloudhost-deletenvir':'hosting/delete-environment','agcapi-cloudhost-createsite':'hosting/create-site','agcapi-cloudhost-querysite':'hosting/list-sites','agcapi-cloudhost-refreshsite':'hosting/update-site','agcapi-cloudhost-deletesite':'hosting/delete-site',
 'server-rest-import':'cloud-auth/import-users','server-rest-export':'cloud-auth/export-users','server-rest-verify':'cloud-auth/verify-token','server-rest-revoke':'cloud-auth/revoke-token',
 'client-rest-auth-obtaintoken':'cloud-auth/client-token','client-rest-auth-verifycode':'cloud-auth/request-code','client-rest-auth-register':'cloud-auth/register','client-rest-auth-signin':'cloud-auth/sign-in-password','client-rest-auth-usersigini':'cloud-auth/sign-in-code','client-rest-auth-modify-phone-email':'cloud-auth/update-account','client-rest-auth-resetpwd':'cloud-auth/reset-password','client-rest-auth-modifypwd':'cloud-auth/update-password','client-rest-auth-userlink':'cloud-auth/link-user','client-rest-auth-unlink':'cloud-auth/unlink-user','client-rest-auth-signout':'cloud-auth/sign-out','client-rest-auth-userdelete':'cloud-auth/delete-user','client-rest-auth-reauthenticate':'cloud-auth/reauthenticate',
 'clouddb-obtaintoken-rest':'clouddb/client-token','clouddb-signin-rest':'clouddb/sign-in','clouddb-query-data-rest':'clouddb/query-data','clouddb-update-data-rest':'clouddb/update-data','clouddb-delete-data-rest':'clouddb/delete-data',
 'storage-restapi-uploadfile':'storage/upload-file','storage-restapi-getfile':'storage/download-file','storage-restapi-deletefile':'storage/delete-file','storage-restapi-getmetadata':'storage/get-metadata','storage-restapi-changemetadata':'storage/update-metadata','storage-restapi-listfile':'storage/list-files',
 'agc-serverless-comm-api-obtainuserid':'templates/user-id','agc-serverless-get-id':'templates/player-id'}
def save(path,obj):
 p=ROOT/path;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(json.dumps(obj,ensure_ascii=False,indent=2)+'\n')
def param(row,loc,src):
 name=row[0]; req,typ=(row[1],row[2]) if row[1] in ('M','O') else (row[2],row[1]);desc=row[3]; low=typ.lower()
 p=dict(name=name,in_=loc,type='string',required=req=='M',description=desc,sourceUrl=src)
 p['in']=p.pop('in_')
 model=None
 if low.startswith(('list<','array<')):
  child=re.search(r'<([^>]+)>',typ)[1]; p['type']='array';p['itemsType']='string' if child.lower().startswith('string') else 'object'
  if p['itemsType']=='object':model=child
  m=re.search(r'>\((\d+)\)',typ)
  if m:p['maxItems']=int(m[1])
 elif low in ('integer','int','long','short','byte','int32','int64'):p['type']='integer'
 elif low in ('double','float'):p['type']='number'
 elif low=='boolean':p['type']='boolean'
 elif low=='bytearray':p['type']='' # Wire shape is not illustrated; do not guess base64 vs numeric array.
 elif low=='json':p['type']='object' # opStore's explicit request example is an object; recorded below.
 elif not low.startswith('string'):
  p['type']='object';model=typ
 m=re.fullmatch(r'String\((\d+)\)',typ,re.I)
 if m:p['maxLength']=int(m[1])
 return p,model
manifest=[]
for doc in D:
 slug=doc['sourceUrl'].rsplit('/',1)[1];key=re.sub(r'-\d+$','',slug);target=IDS.get(key)
 if not target:
  save(Path('pkg/domain/reference/supplemental')/('serverless-'+slug+'.json'),doc);continue
 family,id=target.split('/');doc['captureTool']='Computer Use; scripts/browser-capture-reference.mjs-compatible merged-cell capture'
 req=[t for t in doc['tables'] if t['section']=='请求参数'];heads=[t for t in req if t['location']=='header'];body=[t for t in req if t['location'] in ('body','')];others=[t for t in req if t['location'] in ('query','path')]
 ps=[]
 for t in heads+others:
  for row in t['rows'][1:]:
   p,_=param(row,t['location'],t['sourceUrl'])
   # Credential headers are generated/checked by the authorization layer, not
   # unconditional input fields before it has had a chance to resolve a profile.
   if p['name'].lower() in ('authorization','client_id'):p['required']=False
   if not any(v['name']==p['name'] and v['in']==p['in'] for v in ps):ps.append(p)
 if body:
  pending=[]
  for row in body[0]['rows'][1:]:
   p,m=param(row,'body',body[0]['sourceUrl']);ps.append(p)
   if m:pending.append((p,m))
  # These endpoint pages inline their request models in pre-order. Never read
  # response tables to infer request requiredness.
  idx=1
  def attach(p,m):
   global idx
   if idx>=len(body):return
   table=body[idx];idx+=1;children=[];nexts=[]
   for row in table['rows'][1:]:
    c,cm=param(row,'body',table['sourceUrl']);children.append(c)
    if cm:nexts.append((c,cm))
   p['itemsProperties' if p['type']=='array' else 'properties']=children
   for c,cm in nexts:attach(c,cm)
  for p,m in pending:
   # delete-data omits MsgInfo's inline table; retain the reviewed query/update model.
   if family=='clouddb' and id=='delete-data' and m=='MsgInfo':continue
   attach(p,m)
 proto=dict(doc['tables'][0]['rows'])
 url=proto['接口URL'].split()[0];method=proto['承载协议'].split()[-1];direction='developer-to-huawei';base='https://connect-api.cloud.huawei.com';auth='server-api-client';ctype='application/json'
 if url.startswith('https://connect-drcn.'):
  base='https://connect-drcn.dbankcloud.cn';path=url.removeprefix(base);auth='client-api-token'
 elif family=='storage':
  base='https://ops-server-drcn.agcstorage.link/v0';path=url.removeprefix('https://{domain}');ctype='';auth='server-api-client'
 elif family=='templates':path='{callbackUrl}';base='';direction='inbound-callback';auth='callback-body';ps.insert(0,dict(name='callbackUrl',**{'in':'path'},type='string',required=True,description='Developer HTTPS authentication callback URL',sourceUrl=doc['sourceUrl']))
 elif id=='upload-version':path='{uploadUrl}';base='';auth='signed-upload';ctype='';ps=[p for p in ps if p['in']!='body'];ps+= [dict(name='uploadUrl',**{'in':'path'},type='string',required=True,sourceUrl=doc['sourceUrl']),dict(name='file',**{'in':'file'},type='string',required=True,sourceUrl=doc['sourceUrl'])]
 else:path=url.removeprefix('https://{domain}')
 if family=='cloud-auth' and id in ('verify-token','revoke-token'):path=path.split('?')[0]
 if family=='hosting' and id=='upload-version':
  for p in ps:
   if p['name']=='Authorization':p['required']=True
 if family=='cloud-auth' and key.startswith('server-rest-'):auth='service-account-or-api-client'
 if id=='client-token':auth='body-client-credentials'
 if family=='storage':
  for n in re.findall(r'{([^}]+)}',path):ps.insert(0,dict(name=n,**{'in':'path'},type='string',required=True,description='Storage instance or object path; percent-encoded exactly once.',sourceUrl=doc['sourceUrl']))
  if id=='upload-file':ps.append(dict(name='file',**{'in':'file'},type='string',required=True,sourceUrl=doc['sourceUrl']))
 # Source conflicts are visible in discovery/docs; do not silently infer new rules.
 conflicts=[]
 if family=='clouddb' and id in ('query-data','update-data','delete-data'):
  conflicts.append('opStore is labelled Json/JSON字符串 in the table, but request examples use an object. The adapter uses the explicit example object; queryCond remains an encoded JSON string. SchemaField.t descriptions enumerate numeric codes while write examples use TYPE_* strings. Delete Schema.fs says one item but its example contains three; no conflicting one-item cap is imposed. FieldValue.ba is labelled byteArray without a wire example; preserve caller JSON rather than guessing base64 or array.')
 if family=='cloud-auth' and id=='export-users':conflicts.append('block is labelled String but examples use integer; users and ret response tables say JSON strings while examples use array/object. Both wire shapes are preserved; no conversion is assumed.')
 if family=='storage' and id=='update-metadata':doc['implementationNotes']=['GET with metadata=update matches the prototype and Java HttpGet example; the short request example omits the metadata query.']
 if conflicts:doc['unresolvedRules']=conflicts
 def walk(params):
  for p in params:
   n=p['name']
   if n=='service':p['enum']=['hosting']
   if n=='contentType' and family=='hosting':p['enum']=['902']
   if n=='environment' and family=='hosting':p['enum']=['0','1']
   if n=='cleanSuffix':p['enum']=['0','1']
   if n=='querySize':p['maximum']=1000
   if n=='sendInterval':p['minimum']=30;p['maximum']=120
   if n=='action':p['enum']=['1001','1002']
   if n=='os' and family=='clouddb':p['maxItems']=1
   if n=='t' and '操作类型' in p.get('description',''):p['enum']=['1']
   if n=='deleteFiles':p['maxItems']=1000
   if n=='fileSize' and family=='hosting':p['minimum']=0;p['maximum']=10*1024**3
   if n=='fileSha256' and family=='hosting':p['pattern']='^[a-fA-F0-9]{64}$'
   if n=='lang' and family=='cloud-auth' and id=='request-code':
    lang_path=ROOT/'pkg/domain/reference/supplemental/serverless-agc-auth-verifycode-langtype-0000002563105507.json'
    if lang_path.exists():p['enum']=[row[0] for row in json.loads(lang_path.read_text())['tables'][0]['rows'][1:]]
   if n=='_v':p['enum']=['4']
   if n=='queryType':p['enum']=[str(i) for i in range(6)]
   if n=='provider' and family in ('cloud-auth','clouddb') and p['type']=='integer':p['enum']=[str(i) for i in range(14)]
   if n in ('autoCreateUser','useJwt'):p['enum']=['0','1']
   if n=='grant_type':p['enum']=['client_credentials']
   if n=='queryCond':p['encodedJson']={'type':'object','properties':[{'name':'queryConditions','in':'body','type':'array','required':True,'itemsType':'object','itemsProperties':[{'name':'conditionType','in':'body','type':'string','required':True},{'name':'fieldName','in':'body','type':'string'},{'name':'value','in':'body'}]}]}
   walk(p.get('properties',[]));walk(p.get('itemsProperties',[]))
 walk(ps)
 if family=='clouddb' and id=='delete-data':
  for p in ps:
   if p['name']=='msgInfo':p['properties']=[dict(name='type',**{'in':'body'},type='integer',required=True,enum=['3'],sourceUrl=doc['sourceUrl']),dict(name='opStore',**{'in':'body'},type='object',required=True,properties=[dict(name='storeName',**{'in':'body'},type='string',required=True)],sourceUrl=doc['sourceUrl'])]
 if family=='clouddb' and id in ('query-data','update-data'):
  for p in ps:
   if p['name']=='msgInfo':
    for c in p.get('properties',[]):
     if c['name']=='type':c['enum']=['5' if id=='query-data' else '3']
     if c['name']=='opStore':c['properties']=[dict(name='storeName',**{'in':'body'},type='string',required=True)]
 if family=='cloud-auth' and id=='export-users':
  for p in ps:
   if p['name']=='block':p['type']='';p['description']+=' Documented String/example integer conflict; accepted without coercion.'
 rules=[]
 if family=='cloud-auth' and id in ('register','request-code'):
  rules.append(dict(id='account_identity',kind='at_least_one',fields=['body.phone','body.email'],sourceUrl=doc['sourceUrl'],message='Supply a phone or email identity.',suggestion='Provide phone with country code, or email.'))
 if family=='cloud-auth' and id=='register':
  rules.append(dict(id='phone_verify_code',kind='required_with',field='body.verifyCode',fields=['body.phone'],sourceUrl=doc['sourceUrl'],message='Phone registration requires verifyCode.',suggestion='Obtain a test verification code before registering.'))
 if family=='cloud-auth' and id=='update-account':
  rules.append(dict(id='new_identity',kind='at_least_one',fields=['body.newPhone','body.newEmail'],sourceUrl=doc['sourceUrl'],message='Supply a new phone or email.',suggestion='Specify newPhone or newEmail and its verification code.'))
 save(Path('pkg/domain/reference')/family/(id+'.rules.json'),rules)
 save(Path('pkg/domain/reference')/family/(id+'.json'),doc)
 save(Path('pkg/domain/reference')/family/(id+'.parameters.json'),ps)
 desc=next((n['text'] for n in doc['notes'] if n['section']=='功能介绍'),doc['title'])
 manifest.append(dict(id=id,familyId=family,name=doc['title'],description=desc,method=method,path=path,command='agc '+family+' '+id,status='registered',source='Huawei AppGallery Connect API reference',sourceUrl=doc['sourceUrl'],officialSlug=slug,direction=direction,defaultBaseUrl=base,authentication=auth,requestContentType=ctype,parameters=ps))
save('pkg/domain/reference/serverless/endpoints.json',manifest)
print('Imported',len(D),'documents;',len(manifest),'operations')
