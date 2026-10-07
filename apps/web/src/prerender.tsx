import { renderToString } from 'react-dom/server';
import { App } from './App';
import { content, pageMeta, structuredData, pagePath, siteOrigin, type SiteLanguage } from './siteContent';

const escape = (value: string) => value.replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]!));
export function renderPage(language: SiteLanguage, family?: string) {
  const meta = pageMeta(language,family);
  const head = `<title>${escape(meta.title)}</title>
<meta name="description" content="${escape(meta.description)}" />
<meta name="robots" content="index, follow, max-image-preview:large" />
<meta name="author" content="Createitv" />
<link rel="canonical" href="${meta.url}" />
${(['zh-CN','en','x-default'] as const).map(lang => `<link rel="alternate" hreflang="${lang}" href="${siteOrigin + pagePath(lang === 'en' ? 'en':'zh',family)}" />`).join('\n')}
<meta property="og:type" content="website" />
<meta property="og:site_name" content="agc-cli" />
<meta property="og:url" content="${meta.url}" />
<meta property="og:title" content="${escape(meta.title)}" />
<meta property="og:description" content="${escape(meta.description)}" />
<meta property="og:locale" content="${language === 'zh'?'zh_CN':'en_US'}" />
<meta property="og:locale:alternate" content="${language === 'zh'?'en_US':'zh_CN'}" />
<meta property="og:image" content="${siteOrigin}/social-cover.png" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:title" content="${escape(meta.title)}" />
<meta name="twitter:description" content="${escape(meta.description)}" />
<meta name="twitter:image" content="${siteOrigin}/social-cover.png" />
<script id="seo-schema" type="application/ld+json">${JSON.stringify(structuredData(language,family)).replace(/</g,'\\u003c')}</script>`;
  return { head, body: renderToString(<App renderLanguage={language} renderFamily={family} />) };
}
export function markdown(language: SiteLanguage) {
  const c=content[language];
  return `# ${c.title}\n\n${c.definition}\n\n${c.faqs.map(({q,a})=>`## ${q}\n\n${a}`).join('\n\n')}\n\n## Links\n\n- Website: ${siteOrigin + pagePath(language)}\n- Source: https://github.com/Createitv/agc-cli\n- API contracts: ${siteOrigin}/api-reference.json\n`;
}
