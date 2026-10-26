import assert from 'node:assert/strict';
import {test} from 'node:test';
import {JSDOM} from '../apps/web/node_modules/jsdom/lib/api.js';
import {captureReference} from './browser-capture-reference.mjs';

test('captures request examples split across application-kind subsections',async()=>{
  const dom=new JSDOM(`<div class="layout-right"><h1>更新隐私政策协议</h1>
    <div><h2 id="request">请求示例</h2></div>
    <div><h3>HarmonyOS应用</h3><pre>PUT /agreement\n{"agreementInfos":[{"lang":"zh"}]}</pre></div>
    <div><h3>元服务</h3><pre>PUT /agreement\n{"agreementInfos":[{"lang":"en"}]}</pre></div>
    <div><h2 id="response">响应示例</h2><pre>{"ret":{"code":0}}</pre></div>
    </div>`);
  // jsdom models the DOM, but has no layout engine to supply innerText.
  Object.defineProperty(dom.window.HTMLElement.prototype,'innerText',{get(){return this.textContent;}});
  globalThis.document=dom.window.document;
  const tab={url:async()=> 'https://official/fixture',playwright:{
    getByRole:()=>({first:()=>({waitFor:async()=>{}})}),evaluate:async(fn,arg)=>fn(arg),
  }};
  try {
    const doc=await captureReference(tab,'https://official/fixture');
    assert.equal(doc.examples.length,2);
    assert.match(doc.examples[0].text,/HarmonyOS应用/);
    assert.match(doc.examples[0].text,/"lang":"zh"/);
    assert.match(doc.examples[0].text,/元服务/);
    assert.match(doc.examples[0].text,/"lang":"en"/);
    assert.doesNotMatch(doc.examples[0].text,/"ret"/);
    assert.match(doc.examples[1].text,/"ret"/);
  } finally {delete globalThis.document;dom.window.close();}
});
