import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { JSDOM } from 'jsdom';

const registry = JSON.parse(await readFile('dist/api-reference.json','utf8')).data;
const families = [...new Set(registry.map(e=>e.familyId))];
const sitemap = await readFile('dist/sitemap.xml','utf8');
let checked=0;
for (const language of ['zh','en']) for (const family of [undefined,...families]) {
  const path=`${language==='en'?'/en/':'/'}${family?`reference/${family}/`:''}`;
  const doc=new JSDOM(await readFile(`dist${path}index.html`,'utf8')).window.document;
  const canonical=`https://agccli.app${path}`;
  assert.equal(doc.documentElement.lang,language==='zh'?'zh-CN':'en');
  assert.equal(doc.querySelector('link[rel="canonical"]').href,canonical);
  assert.equal(doc.querySelector('meta[property="og:url"]').content,canonical);
  assert.equal(doc.querySelectorAll('h1').length,1);
  assert.ok(doc.querySelector('meta[name="description"]').content.length>40);
  assert.ok(doc.querySelector('#root').textContent.includes('AppGallery Connect'));
  for (const lang of ['zh-CN','en','x-default']) assert.ok(doc.querySelector(`link[hreflang="${lang}"]`));
  const graph=JSON.parse(doc.querySelector('#seo-schema').textContent)['@graph'];
  const faq=graph.find(e=>e['@type']==='FAQPage');
  const visibleFaq=[...doc.querySelectorAll('.projectFaq details')];
  assert.equal(visibleFaq.length,faq.mainEntity.length);
  for (const [i,element] of visibleFaq.entries()) {
    assert.equal(element.querySelector('summary').textContent,faq.mainEntity[i].name);
    assert.equal(element.querySelector('p').textContent,faq.mainEntity[i].acceptedAnswer.text);
  }
  assert.equal(doc.querySelectorAll('.endpointEntry').length,registry.filter(e=>e.familyId===(family??'publishing')).length);
  for (const link of doc.querySelectorAll('.familyLinks a')) {
    const target=await readFile(`dist${link.getAttribute('href')}index.html`,'utf8');
    assert.ok(target.includes('endpointEntry'));
  }
  assert.ok(sitemap.includes(`<loc>${canonical}</loc>`));
  checked++;
}
assert.ok((await readFile('dist/robots.txt','utf8')).includes('User-agent: OAI-SearchBot\nAllow: /'));
const headers=await readFile('dist/_headers','utf8');
for (const file of ['overview.zh-CN.md','overview.en.md']) assert.ok(headers.includes(`/${file}\n  Content-Type: text/plain; charset=utf-8`),`${file} must remain readable by clients that reject text/markdown`);
console.log(`SEO checks passed: ${checked} pages, reciprocal locales, visible FAQ/schema parity, ${registry.length} endpoints.`);
