#!/usr/bin/env python3
"""Reproduce reviewed management contracts from public AGCLI source + CUA captures.

Usage: python3 scripts/import-management-reference.py PACKAGE_DIR CAPTURE_JSON
No network, user sessions, package execution or contract inference.
"""
import copy
import hashlib
import json
import pathlib
import sys

ROOT = pathlib.Path(__file__).resolve().parents[1]
PACKAGE = pathlib.Path(sys.argv[1])
PAGES = json.loads(pathlib.Path(sys.argv[2]).read_text())
SOURCE = 'https://www.npmjs.com/package/@hw-agconnect/agc-cli/v/1.12.0'
assert json.loads((PACKAGE/'package.json').read_text())['version'] == '1.12.0'

def save(path, obj):
    p=ROOT/path; p.parent.mkdir(parents=True,exist_ok=True)
    p.write_text(json.dumps(obj,ensure_ascii=False,indent=2)+'\n')

def p(name, where='body', typ='string', required=True, **kw):
    return dict(name=name, **{'in':where}, type=typ, required=required, sourceUrl=SOURCE, **kw)

def s(name, **kw): return p(name, **kw)
def obj(name, children, **kw): return p(name,typ='object',properties=children,**kw)
def arr(name, children=None, items='object', **kw):
    return p(name,typ='array',itemsType=items,**({'itemsProperties':children} if children else {}),**kw)
def optional(name,where='query',typ='string',**kw):return p(name,where,typ,False,**kw)
def path(name, **kw):return p(name,'path',**kw)

defs=[]
def add(family,id,method,url,params,files,auth='agcli-management-oauth',notes=None,pages=None):
    evidence=[]
    for name in files+['modules/request/edge.js']+(['modules/request/RequestService.js'] if auth=='agcli-management-cookie' else []):
        raw=(PACKAGE/name).read_bytes()
        evidence.append(dict(packagePath=name,sha256=hashlib.sha256(raw).hexdigest(),source=raw.decode()))
    doc=dict(title=id,sourceUrl=SOURCE,checkedAt='2026-10-26',status='vendor-source-reviewed',package='@hw-agconnect/agc-cli',version='1.12.0',sourceFiles=evidence,
             notes=notes or [], reviewedPages=[d for d in PAGES if any(x in d['sourceUrl'] for x in (pages or []))])
    save(pathlib.Path('pkg/domain/reference')/family/(id+'.json'),doc)
    defs.append(dict(id=id,familyId=family,name=id,description='Reviewed native management transport; explicit user session required. '+ ' '.join(notes or []),method=method,path=url,command=f'agc {family} {id}',authentication=auth,direction='developer-to-huawei',status='registered',source='Huawei official AGCLI 1.12.0 source; not a public REST specification',sourceUrl=SOURCE,parameters=params,documentation=dict(title=id,sourceUrl=SOURCE,checkedAt='2026-10-26',status='vendor-source-reviewed',readCommand=f'agc {family} {id} --docs --pretty')))

cookie='agcli-management-cookie'
cache=['commands/cloudcache/util/constants.js','commands/cloudcache/util/Validation.js','commands/cloudcache/cmd/infoUsageStatistics.js']
def cache_op(id,method,url,params,file,notes=None):
    add('cache',id,method,url,params,cache+['commands/cloudcache/api/'+file],cookie,notes,pages=['cloudcache-'])
cache_op('execute-command','POST','/agc/edge/scs/scs-service/v1/eval',[s('command',maxLength=4096)],'execCommand.js')
cache_op('connection-info','GET','/agc/edge/scs/scs-service/v1/cache-invoke-address',[],'cacheInvokeAdress.js', ['projectId header comes from session.'])
cache_op('usage','POST','/agc/edge/spms/v1/query/metric/scs-metric',[p('views','query',enum=['query_cache_metric_trend']),p('startTime','query','integer',minimum=12960000,maximum=2158876800000),p('endTime','query','integer',minimum=12960000,maximum=2158876800000)],'getInfoOfUsageStatistics.js',['Epoch timestamps are milliseconds; interval at most three days, per source cmd/infoUsageStatistics.js.'])
cache_op('enable-service','POST','/agc/edge/cpms/project-service/v1/services/order',[p('projectId','query'),s('serviceName',enum=['cloudCache'])],'enableCache.js',['Does not configure a password or prove the service is ready.'])
for id,method,file in [('set-password','POST','setPassword.js'),('reset-password','PUT','changePassword.js')]:
    cache_op(id,method,'/agc/edge/scs/scs-service/v1/cache-password',[s('password',minLength=8,maxLength=32)],file,['Password is write-only; resetting a used cache password can disrupt applications.'])

name_pattern=r'^[A-Za-z](?:[A-Za-z0-9_]{0,28}[A-Za-z0-9])?$'
field=[s('fieldName',pattern=name_pattern,maxLength=30),s('fieldType',enum=['String','Boolean','Byte','Short','Integer','Long','Float','Double','ByteArray','Text','Date','IntAutoIncrement','LongAutoIncrement'])]
field += [optional(n,'body','boolean') for n in ['belongPrimaryKey','notNull','isNeedEncrypt','isSensitive']]
field += [p('defaultValue',required=False,typ='')]
index=[s('indexName',pattern=r'^[A-Za-z][A-Za-z0-9_]{0,29}$'),arr('indexList',[s('fieldName'),s('sortType',enum=['ASC','DESC'])],minItems=1,maxItems=5)]
permission=[s('role',enum=['World','Authenticated','Creator','Administrator']),arr('rights',items='string',itemsEnum=['Read','Upsert','Delete'],required=False)]
object_fields=[s('objectTypeName',pattern=name_pattern,maxLength=30),arr('fields',field,minItems=1),arr('indexes',index,required=False,maxItems=16),arr('permissions',permission,required=False)]
db=['commands/clouddb/util/constants.js','commands/clouddb/model/PermissionRole.js','commands/clouddb/model/PermissionRight.js','commands/clouddb/model/ObjectField.js','commands/clouddb/model/ObjectIndexField.js','commands/clouddb/cmd/createObject.js','commands/clouddb/cmd/modifyObject.js','commands/clouddb/cmd/createZone.js','commands/clouddb/cmd/deleteZone.js']
def dbop(id,method,url,params,file,notes=None):
    add('clouddb',id,method,'/agc/edge/clouddb/clouddbservice/v1/'+url,params,db+['commands/clouddb/api/'+file],cookie,notes,pages=['agcconsole-objecttypes','aboutclouddb'])
dbop('schema-versions','GET','schemaVersion',[],'getSchemeVersions.js')
dbop('get-object-type','GET','objectType',[p('schemaVersion','query','integer'),p('objectTypeName','query',pattern=name_pattern)],'getObjectTypeAPI.js')
dbop('create-object-type','POST','objectType',[obj('objectType',object_fields)],'createObjectTypes.js',['Creation is asynchronous; query process-status and independently read back schema before reporting completion.'])
dbop('update-object-type','PUT','objectType',[obj('objectType',object_fields),obj('originObjectType',copy.deepcopy(object_fields))],'modifyObjectTypes.js',['originObjectType must be the exact prior definition; this is not a patch or a migration of existing data.'])
dbop('delete-object-type','DELETE','objectType',[s('objectTypeName',pattern=name_pattern)],'deleteObjectType.js')
# Exported import format differs from the editor wrapper, including its permission map.
dbop('import-object-types','POST','importObjectType',[p('schemaVersion',typ='integer'),arr('objectTypes',[x for x in object_fields if x['name']!='permissions'],minItems=1),arr('permissions',[s('objectTypeName'),arr('permissions',permission)],required=False)],'importObjectTypes.js',['Raw JSON document, not a JSON string wrapper; only new object types may be imported; local size limit 20MB.'])
defs[-1]['requestContentType']='application/x-www-form-urlencoded'
dbop('export-object-types','POST','exportObjectType',[p('schemaVersion',typ='integer'),s('type',enum=['json','js','java','objective-c']),optional('jsFileType','body'),optional('javaFileType','body'),optional('packageName','body',maxLength=1024)],'exportObjectTypes.js',['Non-JSON formats return source-encoded Base64; no local extraction or execution.'])
for id,method,file in [('list-zones','GET','listDbZone.js'),('create-zone','POST','createDbZone.js'),('delete-zone','DELETE','deleteDbZone.js')]:
    dbop(id,method,'naturalStore',[] if method=='GET' else [s('cloudDBZoneName',maxLength=20,pattern=r'^[A-Za-z][A-Za-z0-9]{0,19}$')],file)
dbop('process-status','POST','processPercentage',[s('processId')],'checkProcessPertentage.js')

fn='/api/wisefunction/serverless/v1/'
fnname=path('functionName',pattern=r'^[a-z][a-z0-9-]{0,30}[a-z0-9]$')
aliasname=path('aliasName',pattern=r'^[a-z][a-z0-9-]{0,14}[a-z0-9]$')
pagination=[optional('pageSize',typ='integer',minimum=1,maximum=1000),optional('pageIndex',typ='integer',minimum=1)]
def fnop(id,method,url,params,file,notes=None,pages=None):
    files=['commands/cloudfunction/api/'+file,'commands/cloudfunction/option/verify.js']
    if file=='triggerApi.js': files+=['commands/cloudfunction/boundary/triggerService.js','commands/cloudfunction/api/functionApi.js']
    if file=='layerApi.js': files+=['commands/cloudfunction/cmd/layer/createLayerVersionCmd.js']
    add('functions',id,method,fn+url,params,files,notes=notes,pages=pages or ['function-aliasname'])
weights=obj('routingConfig',[],additionalProperties=p('weight',typ='integer',minimum=0,maximum=100),propertyNamesPattern=r'^(\$latest|[1-9][0-9]{0,2}|1000)$',minProperties=1,maxProperties=2,description='One or two version weights totaling 100; explicit complete routing required by this CLI.')
aliasbody=[s('name',pattern=r'^[a-z][a-z0-9-]{0,14}[a-z0-9]$'),s('functionVersion'),optional('description','body',maxLength=64),weights]
for id,method in [('get-alias','GET'),('delete-alias','DELETE')]:fnop(id,method,'functions/{functionName}/aliases/{aliasName}',[fnname,aliasname],'aliasApi.js')
fnop('list-aliases','GET','functions/{functionName}/aliases',[fnname,optional('functionVersion')]+pagination,'aliasApi.js')
fnop('create-alias','POST','functions/{functionName}/aliases',[fnname]+aliasbody,'aliasApi.js',['Provide one or two explicit version weights totaling 100; no implicit selection of latest published version.'])
fnop('update-alias','PUT','functions/{functionName}/aliases/{aliasName}',[fnname,aliasname,s('functionVersion'),s('revisionId'),p('description',required=False,maxLength=64),weights],'aliasApi.js',['Read get-alias first and provide revisionId plus the full intended routing configuration.'])
triggerfile='triggerApi.js'
httpbody=[s('funcId'),s('mode',enum=['NO_PATH']),p('authFlag',typ='boolean'),s('authAlgor',enum=['HDA-SYSTEM']),s('authType',enum=['apigw-client','cloudgw-client','internal-auth']),p('enableUrlDecode',typ='boolean')]
cronbody=[s('funcId'),s('triggerType',enum=['Cron']),s('triggerMode',enum=['Sync','ASync']),p('triggerSwitch',typ='boolean'),obj('spec',[s('schedule'),p('payload',typ='',required=False)])]
fnop('create-http-trigger','POST','trigger',httpbody,triggerfile,['funcId is obtained from query-function-version and must include the intended qualifier; create policy separately if required.'],['httptrigger'])
fnop('create-cron-trigger','POST','triggers',cronbody,triggerfile,['Default AGCLI create switch is false; an explicit triggerSwitch is required. Current guide says switch cannot be changed later, unlike AGCLI update; cloud support is unverified.'],['crontrigger'])
for kind,route in [('http-trigger','trigger'),('trigger','triggers')]:
    for verb,method in [('get','GET'),('delete','DELETE')]:fnop(verb+'-'+kind,method,route+'/{triggerId}',[path('triggerId')],triggerfile,pages=['trigger-'])
    fnop('list-'+('http-triggers' if kind=='http-trigger' else 'triggers'),'GET',route+'/list/{funcId}',[path('funcId')]+pagination,triggerfile,['Listing path is also reviewed in boundary/triggerService.js.'],['trigger-'])
    defs[-1]['description']+=' Endpoint list uses exact funcId, no automatic qualifier discovery.'
fnop('update-cron-trigger','PUT','triggers/{triggerId}',[path('triggerId')]+[x for x in cronbody if x['name']!='triggerType']+[s('revisionId')],triggerfile,['Read get-trigger first; preserve funcId, spec, triggerMode and revisionId.'],['crontrigger'])
fnop('get-function-policy','GET','functions/{functionName}/policy',[fnname,p('qualifier','query')],triggerfile,pages=['trigger-'])
fnop('create-function-policy','POST','functions/{functionName}/policy',[fnname,p('qualifier','query'),s('effect',enum=['Allow']),s('action',enum=['InvokeFunction']),s('sourceProvider'),s('sourceId'),s('revisionId',required=False)],triggerfile,['Read policy first; no automatic grants or silent retries.'],['trigger-'])
layer=path('layerName',pattern=r'^[a-z][a-z0-9-]{0,30}[a-z0-9]$')
fnop('list-layers','GET','layers',[optional('layerName'),optional('ownerFlag',enum=['all','self','other']),optional('compatibleRuntime')]+pagination,'layerApi.js',pages=['function-layer'])
fnop('list-layer-versions','GET','layers/{layerName}/versions',[layer,optional('compatibleRuntime')]+pagination,'layerApi.js',pages=['function-layer'])
fnop('get-layer-version','GET','layers/{layerName}/versions/{layerVersion}',[layer,path('layerVersion',pattern=r'^[1-9][0-9]*$')],'layerApi.js',pages=['function-layer'])
fnop('create-layer-version','POST','layers/{layerName}/versions',[layer,arr('compatibleRuntimes',items='string',minItems=1,itemsEnum=['nodejs14.x','java1.8','python3','custom-runtime']),p('sharedScope',typ='integer',enum=['0','1']),s('zipFile',format='byte'),optional('description','body',maxLength=1024),optional('licenseInfo','body',maxLength=64)],'layerApi.js',['zipFile contains Base64 ZIP bytes; validate compressed/decompressed limits and entry count before authorization. Layer name creates the initial version or appends a version. AGCLI 1.12.0 local ZIP validation caps files below 50MiB; this direct API adapter uses the current browser-reviewed guide limits of 150MiB compressed, 300MiB expanded and 30000 entries. This difference and modern runtime wire values still require cloud verification.'],['function-layer'])

# SAM applications are distinct from gallery templates. These are the exact
# deployment primitives provided by the vendor; codeUri must already be uploaded.
sam_files=['commands/serverlessModel/api/admApi.js','commands/serverlessModel/validator/templateJSON.js','commands/serverlessModel/validator/cloudFunctionValidator.js','commands/serverlessModel/validator/httpTriggerValidator.js','commands/serverlessModel/validator/cronTriggerValidator.js']
prop=[s('codeUri'),s('runtime'),p('memorySize',typ='integer',enum=['500','1000','2000','4000']),p('timeout',typ='integer',minimum=1,maximum=1800),s('handler',maxLength=64,pattern=r'^[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*$'),p('instanceConcurrency',typ='integer',minimum=1,maximum=1000),optional('description','body',maxLength=64,description='Vendor display width <=64; non-ASCII UTF-16 units count twice.'),arr('layers',items='string',required=False,maxItems=5),arr('env',[s('name',maxLength=64,pattern=r'^[-._a-zA-Z][-._a-zA-Z0-9]*$'),s('value',maxLength=2048)],required=False)]
event_fields=[s('name',pattern=r'^[a-z][a-zA-Z0-9-]{0,30}[a-zA-Z0-9]$'),s('type',enum=['http','cron']),obj('properties',[optional('authType','body',enum=['apigw-client','cloudgw-client']),optional('enableUrlDecode','body','boolean'),optional('schedule','body'),p('input',required=False,encodedJson=p('value',typ='',required=False))])]
event_rules=[]
for kind,fields in [('http',['authType','enableUrlDecode']),('cron',['schedule','input'])]:
    for field in fields:
        event_rules.append(dict(id='sam-'+kind+'-'+field,kind='required_if',field='body.properties.'+field,when=dict(field='body.type',equals=kind),sourceUrl=SOURCE,message='SAM '+kind+' events require properties.'+field,suggestion='Supply the conditional event property using the source-defined type.'))
prop += [optional('alias','body',pattern=r'^[a-z][a-z0-9-]{0,14}[a-z0-9]$'),arr('events',event_fields,required=False,itemsRules=event_rules)]
component=[s('name',pattern=r'^[a-z][a-z0-9-]{0,30}[a-z0-9]$'),s('type',enum=['Huawei::AGC::CloudFunctions']),obj('properties',prop)]
sam_body=[s('name'),arr('components',component,minItems=1)]
sam_notes=['SAM CloudFunctions application transport, not deployment of the console gallery templates. codeUri is a pre-uploaded bucket/object key. Creation/application acceptance does not prove ready state; search/get-plan must be read back.']
for id,method,url,params in [
    ('search-applications','GET','/api/adm/v1/applications/search',[p('appName','query')]),
    ('create-application','POST','/api/adm/v1/applications',sam_body),
    ('create-plan','POST','/api/adm/v1/applications/{applicationId}/plans',[path('applicationId')]+sam_body),
    ('get-plan','GET','/api/adm/v1/applications/plans/{planId}',[path('planId')]),
    ('apply-plan','PUT','/api/adm/v1/applications/{applicationId}/plans/{planId}',[path('applicationId'),path('planId')])]:
    add('templates',id,method,url,params,sam_files,notes=sam_notes,pages=['serverless-'])
add('templates','deploy','POST','native:templates/deploy',copy.deepcopy(sam_body)+[s('storageInstance'),optional('waitForReady','body','boolean')],sam_files+['commands/serverlessModel/api/cloudStorageApi.js','commands/serverlessModel/cmd/deploy/create.js','commands/serverlessModel/cmd/deploy/cmd.js','commands/serverlessModel/cmd/deploy/output.js'],notes=['Create-only SAM deployment: validate every local codeUri ZIP, search exact application name, upload code, create application, then independently query readiness. Existing applications are refused; use explicit create-plan/get-plan/apply-plan for updates. Local-only storageInstance and waitForReady are not sent to ADM. Uploaded artifacts may be retained after failure. Not a console gallery deployment.'],pages=['serverless-'])
defs[-1]['direction']='native-workflow'

for id,extra in [('fetch-install',[]),('fetch-periodic',[optional('token',maxLength=2048),optional('params',maxLength=1024)]),('fetch-deferred',[optional('link')])]:
    doc=next(d for d in PAGES if d['sourceUrl'].endswith('cloudfoundation-prefetch-cloud-interdev'))
    save(pathlib.Path('pkg/domain/reference/prefetch')/(id+'.json'),dict(doc,relatedDocuments=[d for d in PAGES if ('cloudfoundation-prefetch' in d['sourceUrl'] or 'cloudfoundation-commandtool' in d['sourceUrl']) and d['sourceUrl']!=doc['sourceUrl']],implementationNotes=['Developer-server GET simulation; does not configure Huawei prefetch or prove device delivery. Never loads developer credentials. Installation response limit follows the overview quota of 2MiB; periodic and deferred responses are limited to 3MiB. The configuration guide gives a generic 3MB server-response limit; that does not override the smaller installation quota.']))
    defs.append(dict(id=id,familyId='prefetch',name=id,method='GET',path='{callbackUrl}',command='agc prefetch '+id,authentication='callback-none',direction='inbound-callback',status='registered',source='Huawei browser-reviewed developer server protocol',sourceUrl=doc['sourceUrl'],parameters=[path('callbackUrl'),p('appId','query')]+extra,documentation=dict(title=doc['title'],sourceUrl=doc['sourceUrl'],checkedAt='2026-10-26',status='browser-read',readCommand='agc prefetch '+id+' --docs --pretty')))

device_doc=next(d for d in PAGES if d['sourceUrl'].endswith('cloudfoundation-commandtool-debug'))
save('pkg/domain/reference/prefetch/device-getcache.json',dict(device_doc,relatedDocuments=[d for d in PAGES if d['sourceUrl'].endswith('cloudfoundation-commandtool-errorcode')],implementationNotes=['Runs hdc -t DEVICE shell cf_prefetch getcache -m BUNDLE with separate argv. A successful process is not proof of a new fetch data timestamp. Debug HAP, supported ROM and an explicitly connected test device are prerequisites.']))
defs.append(dict(id='device-getcache',familyId='prefetch',name='Debug periodic prefetch on an explicit device',method='POST',path='local:prefetch/device-getcache',command='agc prefetch device-getcache',authentication='local-device',direction='local',status='registered',source='Huawei browser-reviewed hdc command',sourceUrl=device_doc['sourceUrl'],parameters=[s('bundleName',pattern=r'^[A-Za-z][A-Za-z0-9_]*(?:\.[A-Za-z][A-Za-z0-9_]*)+$'),s('deviceId',pattern=r'^[A-Za-z0-9._:-]+$'),optional('hdcPath','body')],documentation=dict(title=device_doc['title'],sourceUrl=device_doc['sourceUrl'],checkedAt='2026-10-26',status='browser-read',readCommand='agc prefetch device-getcache --docs --pretty')))
for endpoint in defs:
    endpoint['officialSlug']=('agcli-1.12.0-'+endpoint['id']) if endpoint['authentication'].startswith('agcli-management-') else endpoint['sourceUrl'].rsplit('/',1)[-1]
    if endpoint['familyId']=='prefetch':
        for parameter in endpoint['parameters']: parameter['sourceUrl']=endpoint['sourceUrl']
save('pkg/domain/reference/management/endpoints.json',defs)
save('pkg/domain/reference/management/browser-review.json',PAGES)
print(f'Imported {len(defs)} source-reviewed operations and {len(PAGES)} rendered pages')
