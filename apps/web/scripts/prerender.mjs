import { createServer } from 'vite';
import { readFile, writeFile, mkdir } from 'node:fs/promises';
import { resolve } from 'node:path';

const server = await createServer({ server: { middlewareMode: true }, appType: 'custom' });
try {
  const { renderPage, markdown } = await server.ssrLoadModule('/src/prerender.tsx');
  const entries = JSON.parse(await readFile('src/endpoint-registry.json','utf8'));
  const families = [...new Set(entries.map(e => e.familyId))];
  const template = await readFile('dist/index.html','utf8');
  const paths = [];
  for (const language of ['zh','en']) {
    for (const family of [undefined,...families]) {
      const path = `${language === 'en' ? '/en/' : '/'}${family ? `reference/${family}/` : ''}`;
      const { head, body } = renderPage(language,family);
      const html = template.replace('lang="zh-CN"',`lang="${language === 'zh'?'zh-CN':'en'}"`).replace(/<!--seo:start-->[\s\S]*?<!--seo:end-->/,()=>head).replace('<div id="root"></div>',()=>`<div id="root">${body}</div>`);
      const target = resolve('dist','.'+path); await mkdir(target,{recursive:true}); await writeFile(resolve(target,'index.html'),html);
      paths.push(path);
    }
    await writeFile(`dist/overview.${language === 'zh'?'zh-CN':'en'}.md`, markdown(language));
  }
  await writeFile('dist/api-reference.json',JSON.stringify({data:entries},null,2)+'\n');
  const sitemap = `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">\n${paths.map(path=>{
    const rest=path.replace(/^\/en\//,'/');
    return `<url><loc>https://agccli.app${path}</loc><lastmod>2026-10-07</lastmod><xhtml:link rel="alternate" hreflang="zh-CN" href="https://agccli.app${rest}"/><xhtml:link rel="alternate" hreflang="en" href="https://agccli.app/en${rest}"/><xhtml:link rel="alternate" hreflang="x-default" href="https://agccli.app${rest}"/></url>`;
  }).join('\n')}\n</urlset>\n`;
  await writeFile('dist/sitemap.xml',sitemap);
  console.log(`Prerendered ${paths.length} bilingual pages and ${entries.length} endpoint contracts.`);
} finally { await server.close(); }
