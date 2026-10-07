export type SiteLanguage = 'zh' | 'en';
export const siteOrigin = 'https://agccli.app';
export const updated = '2026-10-07';
export const content = {
  zh: {
    title: 'agc-cli：鸿蒙与华为 AppGallery Connect 开源命令行工具',
    description: 'agc-cli 是面向华为 AppGallery Connect 的 Go 开源 CLI。浏览 13 个 API 家族、156 个接口，使用凭据 Profile、dry-run 和 JSON 输出，将应用管理接入终端、脚本与 AI Agent。',
    heading: 'agc-cli 是什么，适合谁使用？',
    definition: 'agc-cli 是 Createitv 开源的华为 AppGallery Connect 命令行工具，安装后的命令名是 agc。它面向鸿蒙和华为应用开发者，将接口发现、凭据 Profile、请求预演与 API 调用组织成可复用命令，适合需要维护应用资料、测试用户、商品和报表的个人开发者与团队。',
    source: '项目与使用指南', faq: '常见问题', updated: '内容更新',
    faqs: [
      { q: 'agc-cli 是华为官方工具吗？', a: '不是。agc-cli 是 Createitv 维护的独立开源项目，使用 MIT 许可证。它基于华为 AppGallery Connect API 参考文档组织接口，不代表华为官方产品或官方支持服务。' },
      { q: '如何在 macOS、Windows 和 Linux 安装 agc-cli？', a: 'macOS：brew tap createitv/tap，然后 brew install agc-cli。Windows：scoop bucket add createitv https://github.com/Createitv/scoop-bucket，然后 scoop install agc-cli。Linux：从 GitHub Releases 下载适合架构的二进制。安装后运行 agc version 检查版本。' },
      { q: 'agc-cli 免费吗，需要付费账号吗？', a: 'agc-cli 使用 MIT 许可证，可免费使用和修改。调用华为 API 需要自己的 AppGallery Connect 账号、凭据和相应权限；华为平台服务的费用和限制以华为说明为准。浏览官网接口和使用 CLI 查看接口定义无需登录。' },
      { q: '如何开始第一次应用信息查询？', a: '先用 agc auth login --service-account-file ~/.agc/service-account.json --name production 保存凭据，再用 agc init --app-id YOUR_APP_ID --default-profile production 绑定项目。运行 agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=zh-CN 预演请求，确认后增加 --dry-run=false 发送请求。接口要求 client_id 时显式补上请求头。' },
      { q: '156 个接口是否都完成了生产验证？', a: '156 表示注册表中的接口条目数量，分属 13 个 API 家族，不代表每个接口都已完成生产验证。注册表展示已登记参数和调用入口；完整请求、响应、权限与业务前置条件以每个接口链接的参考文档为准。' },
      { q: '默认 dry-run 会修改我的应用吗？', a: '不会。带 --invoke 的接口命令默认构建并预览请求，不发送远程请求。真实调用需显式设置 --dry-run=false；调用前应检查目标应用、所选凭据和参数。' },
      { q: '可以把 agc-cli 用在脚本或 AI Agent 中吗？', a: '可以。默认 JSON 输出适合脚本处理，agc capabilities 和 agc endpoints 提供能力发现；agc web-server 提供本地 REST 接口，agc openapi 导出接口契约。命令模板帮助导航，但不能替代业务状态检查或审核判断。' },
      { q: '项目 Profile 会自动填入所有接口参数吗？', a: '项目配置可以绑定默认凭据 Profile，临时切换时使用 --profile。应用 ID、语言、请求头和请求体等接口参数仍需按接口要求显式提供。不要将 Service Account 私钥或凭据文件提交到 Git。' },
    ],
  },
  en: {
    title: 'agc-cli: Open-source Huawei AppGallery Connect CLI for HarmonyOS',
    description: 'Discover 13 API families and 156 endpoints with agc-cli, an open-source Go CLI for Huawei AppGallery Connect. Use credential profiles, dry runs, and JSON output in terminal, script, and AI agent workflows.',
    heading: 'What is agc-cli, and who is it for?',
    definition: 'agc-cli is an open-source Go command-line tool for Huawei AppGallery Connect, maintained by Createitv. Its installed command is agc. It organizes API discovery, credential profiles, request previews, and API invocation into reusable commands for HarmonyOS and Huawei app developers who manage app metadata, testers, products, and reports.',
    source: 'Project and guides', faq: 'Frequently asked questions', updated: 'Content updated',
    faqs: [
      { q: 'Is agc-cli an official Huawei tool?', a: 'No. agc-cli is an independent open-source project maintained by Createitv under the MIT license. It organizes endpoints using Huawei AppGallery Connect API references and is not an official Huawei product or support service.' },
      { q: 'How do I install agc-cli on macOS, Windows, or Linux?', a: 'On macOS, run brew tap createitv/tap, then brew install agc-cli. On Windows, run scoop bucket add createitv https://github.com/Createitv/scoop-bucket, then scoop install agc-cli. On Linux, download a binary for your architecture from GitHub Releases. Run agc version to verify installation.' },
      { q: 'Is agc-cli free, and do I need an account?', a: 'agc-cli is free to use and modify under the MIT license. Calling Huawei APIs requires your own AppGallery Connect account, credentials, and permissions. Huawei service pricing and limits remain subject to Huawei terms. Browsing the endpoint reference or inspecting CLI definitions does not require login.' },
      { q: 'How do I query app information for the first time?', a: 'Save credentials with agc auth login --service-account-file ~/.agc/service-account.json --name production, then bind the project with agc init --app-id YOUR_APP_ID --default-profile production. Preview with agc publishing app-info-query --invoke --query appId=YOUR_APP_ID --query lang=en-US. Add --dry-run=false after checking the request. Supply a client_id header explicitly when required.' },
      { q: 'Are all 156 endpoints production-verified?', a: 'No. The registry contains 156 entries across 13 API families; this is registry coverage, not production verification of every endpoint. Registered parameters and command surfaces are shown here. Consult each endpoint reference for full request and response schemas, permissions, and business prerequisites.' },
      { q: 'Does the default dry run change my app?', a: 'No. Endpoint commands with --invoke preview the request without sending it by default. Real invocation requires --dry-run=false. Check the target app, credential profile, and parameters before executing.' },
      { q: 'Can I use agc-cli with scripts or AI agents?', a: 'Yes. JSON output supports script processing. agc capabilities and agc endpoints expose discovery; agc web-server provides local REST routes, and agc openapi exports the contract. Command templates support navigation but do not replace business-state checks or review decisions.' },
      { q: 'Does a project profile fill all endpoint parameters automatically?', a: 'A project configuration can select its default credential profile, with --profile for explicit overrides. App IDs, languages, headers, and request bodies still need to be supplied according to each endpoint. Keep Service Account private keys and credential files out of Git.' },
    ],
  },
};

export function localeFromPath(path: string): SiteLanguage { return /^\/en(?:\/|$)/.test(path) ? 'en' : 'zh'; }
export function familyFromPath(path: string) { return path.match(/\/reference\/([^/]+)/)?.[1]; }
export function pagePath(language: SiteLanguage, family?: string) { return `${language === 'en' ? '/en/' : '/'}${family ? `reference/${family}/` : ''}`; }
export function pageMeta(language: SiteLanguage, family?: string) {
  const c = content[language];
  return { title: family ? `${family} API ${language === 'zh' ? '接口参考' : 'endpoint reference'} | agc-cli` : c.title, description: family ? `${c.description} ${family} API: ${language === 'zh' ? '请求方法、路径、参数与 CLI 命令。' : 'Methods, paths, parameters, and CLI commands.'}` : c.description, url: siteOrigin + pagePath(language, family) };
}
export function structuredData(language: SiteLanguage, family?: string) {
  const c = content[language]; const meta = pageMeta(language, family);
  return { '@context': 'https://schema.org', '@graph': [
    { '@type': 'WebSite', '@id': `${siteOrigin}/#website`, name: 'agc-cli', url: `${siteOrigin}/`, inLanguage: ['zh-CN', 'en'] },
    { '@type': 'WebPage', '@id': `${meta.url}#webpage`, url: meta.url, name: meta.title, description: meta.description, inLanguage: language === 'zh' ? 'zh-CN' : 'en', dateModified: updated, isPartOf: { '@id': `${siteOrigin}/#website` } },
    { '@type': 'SoftwareApplication', '@id': `${siteOrigin}/#software`, name: 'agc-cli', alternateName: 'agc', url: `${siteOrigin}/`, description: c.definition, applicationCategory: 'DeveloperApplication', operatingSystem: 'macOS, Windows, Linux', license: 'https://github.com/Createitv/agc-cli/blob/main/LICENSE', downloadUrl: 'https://github.com/Createitv/agc-cli/releases', offers: { '@type': 'Offer', price: '0', priceCurrency: 'USD' } },
    { '@type': 'SoftwareSourceCode', name: 'agc-cli', codeRepository: 'https://github.com/Createitv/agc-cli', programmingLanguage: 'Go', author: { '@type': 'Person', name: 'Createitv', url: 'https://github.com/Createitv' } },
    { '@type': 'FAQPage', '@id': `${meta.url}#faq`, inLanguage: language === 'zh' ? 'zh-CN' : 'en', mainEntity: c.faqs.map(({q,a}) => ({ '@type': 'Question', name: q, acceptedAnswer: { '@type': 'Answer', text: a } })) },
  ] };
}
