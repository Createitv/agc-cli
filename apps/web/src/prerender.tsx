import { renderToString } from 'react-dom/server';
import { App } from './App';
import { agentSkillNames, content, pageMeta, structuredData, pagePath, siteOrigin, updated, type SiteLanguage } from './siteContent';

export const contentUpdated = updated;

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

export function agentGuide(language: SiteLanguage) {
  const zh = language === 'zh';
  const title = zh ? 'agc-cli Agent 使用指南' : 'agc-cli Agent Usage Guide';
  const intro = zh
    ? '本文档帮助 Codex、Claude 和其他 Coding Agent 在不猜测命令或泄露凭据的前提下使用 agc-cli。agc-cli 是 Createitv 维护的独立开源项目，不是华为官方产品。'
    : 'This guide helps Codex, Claude, and other coding agents use agc-cli without guessing commands or exposing credentials. agc-cli is an independent open-source project maintained by Createitv; it is not an official Huawei product.';
  const tasks = zh
    ? [['检查身份与权限', 'agc auth whoami --pretty'], ['发现能力', 'agc capabilities --pretty'], ['查看接口', 'agc <family> endpoints --output table'], ['读取接口文档', 'agc <family> <endpoint> --docs --pretty'], ['预演请求', 'agc <family> <endpoint> --invoke ...'], ['明确执行', '在核对目标、凭据和参数后增加 --dry-run=false']]
    : [['Check identity and access', 'agc auth whoami --pretty'], ['Discover capabilities', 'agc capabilities --pretty'], ['List endpoints', 'agc <family> endpoints --output table'], ['Read endpoint docs', 'agc <family> <endpoint> --docs --pretty'], ['Preview a request', 'agc <family> <endpoint> --invoke ...'], ['Invoke explicitly', 'Add --dry-run=false only after checking the target, profile, and parameters']];
  return `# ${title}\n\nLast updated: ${updated}\n\n${intro}\n\n## ${zh ? '推荐工作流' : 'Recommended workflow'}\n\n${tasks.map(([label,command],index)=>`${index+1}. **${label}** — \`${command}\``).join('\n')}\n\n## ${zh ? '安全边界' : 'Safety boundaries'}\n\n- ${zh ? '默认 dry-run 只构建请求，不发送远程变更；真实调用需要显式设置 `--dry-run=false`。' : 'The default dry run builds a request without sending a remote change; real invocation requires an explicit `--dry-run=false`.'}\n- ${zh ? '不要把 Service Account 私钥、OAuth 密钥、凭据响应或 `~/.agc/credentials.json` 提交到 Git。' : 'Never commit Service Account private keys, OAuth secrets, credential responses, or `~/.agc/credentials.json` to Git.'}\n- ${zh ? '注册表覆盖不等于每个接口已经完成生产验证；权限、状态和完整模型以接口文档链接为准。' : 'Registry coverage does not mean every endpoint is production-verified; follow the linked endpoint references for permissions, state, and complete models.'}\n- ${zh ? '上传、提交审核、发布和删除类操作必须由用户明确授权，并在执行后读取远程状态。' : 'Uploads, review submissions, releases, and deletion operations require explicit user authorization and remote-state verification.'}\n\n## Agent Skills\n\n${agentSkillNames.map(name=>`- [${name}](https://github.com/Createitv/agc-cli/tree/main/skills/${name})`).join('\n')}\n\n${zh ? '安装全部 Skills' : 'Install all Skills'}: \`agc skills install --agent codex --global\`\n\n## ${zh ? '机器可读资源' : 'Machine-readable resources'}\n\n- [llms.txt](${siteOrigin}/llms.txt)\n- [llms-full.txt](${siteOrigin}/llms-full.txt)\n- [Skills index](${siteOrigin}/skills.md)\n- [API contracts](${siteOrigin}/api-reference.json)\n- [Open-source repository](https://github.com/Createitv/agc-cli)\n`;
}

export function skillsMarkdown() {
  return `# agc-cli Agent Skills\n\nLast updated: ${updated}\n\nSix task-oriented Skills help coding agents use agc-cli safely and consistently. Install all Skills with \`agc skills install --agent codex --global\`, or select one with \`agc skills install --name <skill> --agent codex --global\`.\n\n${agentSkillNames.map(name=>`## ${name}\n\n- Source: https://github.com/Createitv/agc-cli/tree/main/skills/${name}\n- Direct SKILL.md: https://raw.githubusercontent.com/Createitv/agc-cli/main/skills/${name}/SKILL.md\n- Install: \`agc skills install --name ${name} --agent codex --global\``).join('\n\n')}\n\nThe Skills describe workflows and safety checks. They do not grant Huawei permissions, supply credentials, or prove that a remote action succeeded.\n`;
}

export function llmsFull() {
  return `${markdown('en')}\n\n---\n\n${agentGuide('en')}\n\n---\n\n${skillsMarkdown()}\n\n---\n\n${markdown('zh')}\n\n---\n\n${agentGuide('zh')}\n`;
}
