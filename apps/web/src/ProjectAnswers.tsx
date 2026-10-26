import { content, pagePath, updated, type SiteLanguage } from './siteContent';

export function ProjectAnswers({ language, families }: { language: SiteLanguage; families: string[] }) {
  const c = content[language];
  return <section className="section projectAnswers" id="faq">
    <h2>{c.heading}</h2><p className="projectDefinition">{c.definition}</p>
    <p className="projectSources">{c.source}: <a href="https://github.com/Createitv/agc-cli">GitHub</a> · <a href={`https://github.com/Createitv/agc-cli/blob/main/docs/CLI_USAGE${language === 'en' ? '.en' : ''}.md`}>{language === 'zh' ? '使用指南' : 'Usage guide'}</a> · <a href="https://github.com/Createitv/agc-cli/issues">Issues</a></p>
    <h3>{language === 'zh' ? '支持项目持续维护' : 'Support continued development'}</h3>
    <p>{language === 'zh' ? '支持将用于接口适配、跨平台验证、问题修复和中英文文档维护。GitHub Sponsors 已提交，正在等待审核，尚未开放收款。你也可以通过反馈问题和贡献代码支持项目。' : 'Support helps maintain API contracts, cross-platform validation, fixes, and bilingual documentation. GitHub Sponsors has been submitted for review; payments are not available yet. You can already help by reporting issues and contributing code.'} <a href="https://github.com/Createitv/agc-cli/issues">{language === 'zh' ? '反馈问题' : 'Report an issue'}</a></p>
    <h3>{c.faq}</h3>
    <div className="projectFaq">{c.faqs.map(({q,a}) => <details key={q}><summary>{q}</summary><p>{a}</p></details>)}</div>
    <h3>{language === 'zh' ? '按 API 家族浏览完整参考' : 'Browse the complete API reference by family'}</h3>
    <nav className="familyLinks" aria-label={language === 'zh' ? '接口参考页面' : 'Endpoint reference pages'}>{families.map(family => <a key={family} href={pagePath(language,family)}>{family}</a>)}</nav>
    <p className="contentUpdated">{c.updated}: <time dateTime={updated}>{updated}</time></p>
  </section>;
}
