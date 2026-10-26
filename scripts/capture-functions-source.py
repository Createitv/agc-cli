#!/usr/bin/env python3
"""Capture an allowlist of public vendor source files, never scratch or sessions."""
import argparse,hashlib,json
from pathlib import Path
p=argparse.ArgumentParser(description=__doc__);p.add_argument('--package-root',type=Path,required=True);p.add_argument('--checked-at',required=True);a=p.parse_args();root=Path(__file__).resolve().parents[1]
meta=json.loads((a.package_root/'package.json').read_text());assert meta['name']=='@hw-agconnect/agc-cli' and meta['version']=='1.12.0','Review a changed vendor version before changing the native contract'
shared=['commands/cloudfunction/api/functionApi.js','commands/cloudfunction/api/layerApi.js','commands/cloudfunction/option/options.js','commands/cloudfunction/option/verify.js','commands/cloudfunction/option/parse.js','commands/cloudfunction/util/tools.js','modules/request/edge.js']
commands={'create-function':'createFunctionCmd.js','update-function-code':'updateFunctionCodeCmd.js','update-function-config':'updateFunctionConfigCmd.js','publish-function':'publishFunctionCmd.js','update-and-publish-function':'updateAndPublishFunctionCmd.js','invoke':'invokeCmd.js','list-functions':'listFunctionsCmd.js','query-function-version':'queryFunctionVersionCmd.js','list-function-versions':'listFunctionVersionsCmd.js'}
for op,command in commands.items():
 evidence=[]
 for relative in shared+['commands/cloudfunction/cmd/function/'+command]:
  raw=(a.package_root/relative).read_bytes();evidence.append({'packagePath':relative,'sha256':hashlib.sha256(raw).hexdigest(),'source':raw.decode()})
 doc={'title':op,'sourceUrl':'https://www.npmjs.com/package/@hw-agconnect/agc-cli/v/1.12.0','checkedAt':a.checked_at,'status':'vendor-source-reviewed','package':meta['name'],'version':meta['version'],'sourceFiles':evidence,'notes':[{'text':'Public vendor CLI source is evidence of HTTP behavior, not a published REST specification or live validation. Explicit OAuth session required. Native Go preserves TLS verification, revision workflows and original file bytes; it does not read vendor scratch or invoke Node.'}],'limitations':['The reviewed vendor package advertises older runtimes; current cloud runtime availability needs dedicated-resource verification.','Function trigger/alias/layer lifecycle and template deployment are not implied by these nine function workflows.','Management responses omit environment/session values; invocation application results remain unchanged.','Native responses above 16 MiB return an explicit local-limit error.']}
 target=root/'pkg/domain/reference/functions'/f'{op}.json';target.parent.mkdir(exist_ok=True);target.write_text(json.dumps(doc,ensure_ascii=False,indent=2)+'\n')
print('Captured 9 native function workflows with public vendor source hashes.')
