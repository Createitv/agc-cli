import {
  ArrowRight,
  Check,
  ChevronDown,
  ChevronRight,
  CircleDot,
  Clipboard,
  Cloud,
  Code2,
  Download,
  FileJson2,
  Fingerprint,
  Globe2,
  KeyRound,
  Layers3,
  MonitorDown,
  PackageCheck,
  RadioTower,
  Rocket,
  ShieldCheck,
  Terminal,
  Workflow,
  X,
} from 'lucide-react';
import { useEffect, useMemo, useRef, useState } from 'react';

type Language = 'en' | 'zh';
type InstallPlatform = 'unix' | 'windows';

type Capability = {
  id: string;
  name: string;
  description: string;
  command: string;
  restPath: string;
  status: string;
  endpointCount?: number;
};

const installCommands: Record<InstallPlatform, string> = {
  unix: 'curl -fsSL https://createitv.github.io/agc-cli/install.sh | sh',
  windows: 'irm https://createitv.github.io/agc-cli/install.ps1 | iex',
};

const fallbackCapabilities: Capability[] = [
  { id: 'publishing', name: 'Publishing API', description: 'Release metadata, package state, submissions, review withdrawal, and release timing.', command: 'agc publishing', restPath: '/api/v1/publishing', status: 'registered', endpointCount: 14 },
  { id: 'upload', name: 'Upload Management', description: 'APP packages, icons, screenshots, videos, PDFs, and multipart binary uploads.', command: 'agc upload', restPath: '/api/v1/upload', status: 'registered', endpointCount: 6 },
  { id: 'provisioning', name: 'Provisioning', description: 'Certificates, profiles, ACLs, devices, and fingerprints for release workflows.', command: 'agc provisioning', restPath: '/api/v1/provisioning', status: 'registered', endpointCount: 17 },
  { id: 'domains', name: 'Domain Management', description: 'Inspect, precheck, download, and update atomic service domain configuration.', command: 'agc domains', restPath: '/api/v1/domains', status: 'registered', endpointCount: 5 },
  { id: 'testing', name: 'Testing', description: 'Test versions, packages, groups, invitations, feedback, and open testing links.', command: 'agc testing', restPath: '/api/v1/testing', status: 'registered', endpointCount: 27 },
  { id: 'reports', name: 'Reports', description: 'Request and download CSV, Excel, PDF, and filtered date-window reports.', command: 'agc reports', restPath: '/api/v1/reports', status: 'registered', endpointCount: 12 },
  { id: 'projects', name: 'Project Management', description: 'Teams, projects, app summaries, services, certificate fingerprints, and SDK files.', command: 'agc projects', restPath: '/api/v1/projects', status: 'registered', endpointCount: 8 },
  { id: 'comments', name: 'Comments', description: 'Read HarmonyOS and Android review data, ratings, and replies through one surface.', command: 'agc comments', restPath: '/api/v1/comments', status: 'registered', endpointCount: 8 },
  { id: 'pms', name: 'PMS', description: 'HarmonyOS and Android products, subscriptions, promotions, review resources, and prices.', command: 'agc pms', restPath: '/api/v1/pms', status: 'registered', endpointCount: 40 },
  { id: 'gameplay', name: 'Game Playing', description: 'Outbound resource sync plus inbound AppGallery game playing callbacks.', command: 'agc gameplay', restPath: '/api/v1/gameplay', status: 'registered', endpointCount: 8 },
  { id: 'game-items', name: 'Game Item Mall', description: 'Inbound role query and order callbacks configured in AGC.', command: 'agc game-items', restPath: '/api/v1/game-items', status: 'registered', endpointCount: 2 },
  { id: 'resources', name: 'Resource Predownload', description: 'Coordinate resource package versions, files, upload confirmation, and publishing.', command: 'agc resources', restPath: '/api/v1/resources', status: 'registered', endpointCount: 8 },
  { id: 'cicd', name: 'CI/CD Platform', description: 'Connect local Hvigor builds to the same command surface.', command: 'agc cicd', restPath: '/api/v1/cicd', status: 'registered', endpointCount: 1 },
];

const capabilityTranslations: Record<Language, Record<string, { name: string; description: string }>> = {
  en: Object.fromEntries(fallbackCapabilities.map(({ id, name, description }) => [id, { name, description }])),
  zh: {
    publishing: { name: '应用发布', description: '管理发布资料、软件包状态、提交审核、撤回审核与定时发布。' },
    upload: { name: '上传管理', description: '上传应用包、图标、截图、视频、PDF 与 multipart 二进制文件。' },
    provisioning: { name: '证书与设备', description: '管理证书、Provisioning Profile、ACL、测试设备与指纹。' },
    domains: { name: '域名管理', description: '检查、预校验、下载和更新原子化服务域名配置。' },
    testing: { name: '测试服务', description: '管理测试版本、软件包、用户组、邀请、反馈与公开测试链接。' },
    reports: { name: '报表下载', description: '请求并下载 CSV、Excel、PDF 与指定时间范围的运营报表。' },
    projects: { name: '项目管理', description: '管理团队、项目、应用摘要、服务、证书指纹与 SDK 配置文件。' },
    comments: { name: '评论管理', description: '统一读取 HarmonyOS 与 Android 评论、评分并处理回复。' },
    pms: { name: '商品管理', description: '管理 HarmonyOS 与 Android 商品、订阅、促销、审核资料与价格。' },
    gameplay: { name: '游戏联运', description: '处理游戏资源同步以及 AppGallery 游戏联运回调。' },
    'game-items': { name: '游戏道具商城', description: '处理在 AGC 中配置的角色查询与订单回调。' },
    resources: { name: '资源预下载', description: '协调资源包版本、文件上传、上传确认与发布。' },
    cicd: { name: 'CI/CD 平台', description: '把本地 Hvigor 构建接入同一套命令控制面。' },
  },
};

const translations = {
  en: {
    documentTitle: 'agc CLI — AppGallery Connect from one command surface',
    nav: { workflow: 'Workflow', profiles: 'Profiles', registry: 'API registry', quickstart: 'Quick start', language: 'Language', install: 'Install' },
    hero: {
      kicker: 'AppGallery Connect, mapped for the terminal',
      titleLead: 'Move every release through',
      titleAccent: 'one red line.',
      body: 'A Go command center for discovering, dry-running, and invoking AppGallery Connect APIs—with stable JSON for people, CI, and coding agents.',
      start: 'Get started',
      browse: 'Browse 156 interfaces',
      families: 'API families',
      interfaces: 'registered interfaces',
      formats: 'output formats',
      factsLabel: 'Project facts',
    },
    deck: { live: 'local API', demo: 'reference mode', result: 'project profile resolved · dry-run on', label: 'Command output example' },
    rail: [
      { label: 'AUTH', detail: 'PS256 / OAuth' },
      { label: 'PROFILE', detail: 'project-bound' },
      { label: 'PACKAGE', detail: 'upload + status' },
      { label: 'RELEASE', detail: 'invoke by choice' },
    ],
    principles: [
      { title: 'Dry-run first', body: 'Nothing is sent until you opt in.' },
      { title: 'Agent-readable', body: 'JSON envelopes include next actions.' },
      { title: 'One surface', body: 'CLI, local REST, and web share a registry.' },
    ],
    workflow: {
      eyebrow: 'The release rail',
      title: 'Context flows forward. Secrets do not.',
      body: 'Each stage owns one job. Credentials remain local, project context remains portable, and destructive API calls require an explicit switch.',
    },
    profile: {
      eyebrow: 'Automatic profile binding',
      title: 'Choose once per project. Override only when you mean it.',
      body: 'Bind production in .agc/project.json. Commands resolve the profile in a predictable order, while --profile staging remains available for one-off work.',
      items: ['CLI override', 'Project default', 'Active account'],
      safe: 'project safe',
      note: 'No client keys or private keys are stored here.',
    },
    registry: {
      eyebrow: 'Interface registry',
      title: 'Discover the API before writing the request.',
      body: 'Every entry carries its method, path, parameter locations, CLI command, REST route, and dry-run path.',
      families: 'API families',
      interfaces: 'interfaces',
      selected: 'Selected family',
      registered: 'registered',
      discover: 'DISCOVER',
      status: 'STATUS',
      localRest: 'LOCAL REST',
      typed: 'typed interface entries',
      copy: 'Copy discover command',
      copied: 'Copied',
    },
    quickstart: {
      eyebrow: 'Quick start',
      title: 'Three commands to a bound, inspectable workspace.',
      body: 'Use placeholders for your AppGallery Connect values. The second command is what makes the profile automatic for this project.',
      steps: [
        { number: '01', title: 'Save a credential profile', body: 'Service Account is the preferred route. Secrets stay in ~/.agc/credentials.json with restricted permissions.', command: 'agc auth login --service-account-file ~/.agc/service-account.json --name production' },
        { number: '02', title: 'Bind it to the project', body: 'Pin the app, project, package, and default profile once. Every later command resolves that context automatically.', command: 'agc init --app-id <app-id> --project-id <project-id> --package-name com.example.app --default-profile production' },
        { number: '03', title: 'Inspect before invoking', body: 'Endpoint commands are dry-run first. Review the exact method and URL, then opt into the real request.', command: 'agc publishing app-info-query --invoke --query appId=<app-id> --dry-run=false' },
      ],
      copyStep: 'Copy step command',
    },
    console: {
      eyebrow: 'Local command center',
      title: 'The browser is a window into the same registry.',
      body: 'Start the local REST server when you want live capability data. Without it, this site stays useful as a complete reference surface.',
      routesLabel: 'Local server routes',
      routes: ['families + affordances', 'all interface entries', 'machine-readable contract', 'dry-run by default'],
    },
    footer: { tagline: 'AppGallery Connect, from one command surface.', license: 'MIT licensed', stack: 'Go + React', cloud: 'Built for Cloudflare' },
    install: {
      button: 'Install',
      eyebrow: 'Direct install · v0.1.0',
      title: 'Put agc on your command line.',
      body: 'Choose your platform, copy one command, and the installer will fetch the matching binary and verify its SHA-256 checksum.',
      unix: 'macOS / Linux',
      windows: 'Windows',
      terminal: 'INSTALL CHANNEL / STABLE',
      copy: 'Copy command',
      copied: 'Copied',
      note: 'macOS supports Apple silicon and Intel. Linux supports x64 and ARM64. Windows supports x64.',
      checksum: 'SHA-256 verified',
      destination: 'User-writable install path',
      close: 'Close install window',
    },
  },
  zh: {
    documentTitle: 'agc CLI — 用一套命令管理 AppGallery Connect',
    nav: { workflow: '发布流程', profiles: '凭据 Profile', registry: '接口注册表', quickstart: '快速开始', language: '语言', install: '安装' },
    hero: {
      kicker: '为终端而生的 AppGallery Connect 控制面',
      titleLead: '让每一次发布沿着',
      titleAccent: '同一条红线前进。',
      body: '用一套 Go 命令发现、预演并调用 AppGallery Connect API，为开发者、CI 和 Coding Agent 提供稳定 JSON。',
      start: '快速开始',
      browse: '浏览 156 个接口',
      families: 'API 家族',
      interfaces: '已注册接口',
      formats: '输出格式',
      factsLabel: '项目数据',
    },
    deck: { live: '本地 API', demo: '参考模式', result: '已解析项目 Profile · 默认 dry-run', label: '命令输出示例' },
    rail: [
      { label: '鉴权', detail: 'PS256 / OAuth' },
      { label: 'PROFILE', detail: '项目自动绑定' },
      { label: '软件包', detail: '上传与状态' },
      { label: '发布', detail: '确认后调用' },
    ],
    principles: [
      { title: '默认预演', body: '明确确认前不会发送请求。' },
      { title: 'Agent 可读', body: 'JSON 返回下一步可执行操作。' },
      { title: '统一控制面', body: 'CLI、本地 REST 与 Web 共用注册表。' },
    ],
    workflow: {
      eyebrow: '发布轨道',
      title: '上下文向前流动，密钥留在本地。',
      body: '每个阶段只完成一件事。凭据保存在本机，项目上下文可以安全共享，写入类 API 必须显式确认。',
    },
    profile: {
      eyebrow: 'Profile 自动绑定',
      title: '每个项目只选择一次，临时切换时再覆盖。',
      body: '在 .agc/project.json 中绑定 production。命令会按固定优先级解析，也可以用 --profile staging 完成一次性切换。',
      items: ['命令行显式覆盖', '项目默认 Profile', '全局激活账户'],
      safe: '可安全提交',
      note: '项目文件不保存 client key 或私钥。',
    },
    registry: {
      eyebrow: '接口注册表',
      title: '先发现 API，再编写请求。',
      body: '每个条目都包含 method、path、参数位置、CLI 命令、REST 路由和 dry-run 调用方式。',
      families: 'API 家族',
      interfaces: '个接口',
      selected: '当前家族',
      registered: '已注册',
      discover: '发现',
      status: '状态',
      localRest: '本地 REST',
      typed: '个类型化接口条目',
      copy: '复制发现命令',
      copied: '已复制',
    },
    quickstart: {
      eyebrow: '快速开始',
      title: '三条命令，建立可检查的项目工作区。',
      body: '请将占位符替换为 AppGallery Connect 的真实值。第二条命令会让该项目自动选择对应 Profile。',
      steps: [
        { number: '01', title: '保存凭据 Profile', body: '推荐使用 Service Account。密钥会以受限权限保存在 ~/.agc/credentials.json。', command: 'agc auth login --service-account-file ~/.agc/service-account.json --name production' },
        { number: '02', title: '绑定当前项目', body: '一次性保存应用、项目、包名和默认 Profile，后续命令会自动解析这些上下文。', command: 'agc init --app-id <app-id> --project-id <project-id> --package-name com.example.app --default-profile production' },
        { number: '03', title: '调用前先检查', body: '接口命令默认 dry-run。确认 method 与 URL 后，再明确执行真实请求。', command: 'agc publishing app-info-query --invoke --query appId=<app-id> --dry-run=false' },
      ],
      copyStep: '复制步骤命令',
    },
    console: {
      eyebrow: '本地命令中心',
      title: '浏览器看到的，就是 CLI 使用的同一份注册表。',
      body: '需要实时能力数据时启动本地 REST 服务；未启动时，官网仍会作为完整参考界面工作。',
      routesLabel: '本地服务路由',
      routes: ['家族与下一步操作', '全部接口条目', '机器可读契约', '默认 dry-run'],
    },
    footer: { tagline: '用一套命令管理 AppGallery Connect。', license: 'MIT 开源协议', stack: 'Go + React', cloud: '部署于 Cloudflare' },
    install: {
      button: '安装',
      eyebrow: '直接安装 · v0.1.0',
      title: '把 agc 安装到命令行。',
      body: '选择平台并复制一条命令。安装器会下载对应架构的可执行文件，并验证 SHA-256 校验值。',
      unix: 'macOS / Linux',
      windows: 'Windows',
      terminal: '安装通道 / 稳定版',
      copy: '复制命令',
      copied: '已复制',
      note: 'macOS 支持 Apple 芯片和 Intel，Linux 支持 x64 与 ARM64，Windows 支持 x64。',
      checksum: 'SHA-256 校验',
      destination: '用户可写安装目录',
      close: '关闭安装窗口',
    },
  },
} as const;

const railIcons = [KeyRound, Fingerprint, PackageCheck, Rocket];
const principleIcons = [ShieldCheck, FileJson2, Cloud];
const profileCodes = ['--profile staging', '.agc/project.json', '~/.agc/credentials.json'];

function initialLanguage(): Language {
  if (typeof window === 'undefined') return 'en';
  const stored = typeof window.localStorage?.getItem === 'function'
    ? window.localStorage.getItem('agc-language')
    : null;
  if (stored === 'en' || stored === 'zh') return stored;
  return window.navigator.language.toLowerCase().startsWith('zh') ? 'zh' : 'en';
}

export function App() {
  const [language, setLanguage] = useState<Language>(initialLanguage);
  const [languageOpen, setLanguageOpen] = useState(false);
  const [installOpen, setInstallOpen] = useState(false);
  const [installPlatform, setInstallPlatform] = useState<InstallPlatform>('unix');
  const [capabilities, setCapabilities] = useState<Capability[]>(fallbackCapabilities);
  const [endpointCount, setEndpointCount] = useState(156);
  const [mode, setMode] = useState<'demo' | 'live'>('demo');
  const [selected, setSelected] = useState('publishing');
  const [copied, setCopied] = useState('');
  const languagePickerRef = useRef<HTMLDivElement>(null);
  const installTriggerRef = useRef<HTMLButtonElement>(null);
  const closeInstallRef = useRef<HTMLButtonElement>(null);
  const text = translations[language];

  useEffect(() => {
    document.documentElement.lang = language === 'zh' ? 'zh-CN' : 'en';
    document.title = text.documentTitle;
    if (typeof window.localStorage?.setItem === 'function') {
      window.localStorage.setItem('agc-language', language);
    }
  }, [language, text.documentTitle]);

  useEffect(() => {
    if (typeof fetch === 'undefined') return;

    fetch('/api/v1/capabilities')
      .then((response) => response.ok ? response.json() : Promise.reject(new Error('offline')))
      .then((body: { data: Capability[] }) => {
        if (Array.isArray(body.data) && body.data.length > 0) {
          setCapabilities(body.data);
          setMode('live');
        }
      })
      .catch(() => setMode('demo'));

    fetch('/api/v1/endpoints')
      .then((response) => response.ok ? response.json() : Promise.reject(new Error('offline')))
      .then((body: { data: unknown[] }) => {
        if (Array.isArray(body.data)) setEndpointCount(body.data.length);
      })
      .catch(() => setEndpointCount(156));
  }, []);

  useEffect(() => {
    if (!languageOpen) return;
    const handleOutside = (event: PointerEvent) => {
      if (!languagePickerRef.current?.contains(event.target as Node)) setLanguageOpen(false);
    };
    document.addEventListener('pointerdown', handleOutside);
    return () => document.removeEventListener('pointerdown', handleOutside);
  }, [languageOpen]);

  useEffect(() => {
    if (!installOpen) return;
    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';
    closeInstallRef.current?.focus();
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setInstallOpen(false);
    };
    window.addEventListener('keydown', handleEscape);
    return () => {
      document.body.style.overflow = previousOverflow;
      window.removeEventListener('keydown', handleEscape);
      installTriggerRef.current?.focus();
    };
  }, [installOpen]);

  const active = useMemo(
    () => capabilities.find((capability) => capability.id === selected) ?? capabilities[0],
    [capabilities, selected],
  );

  const localizedCapability = (capability: Capability) => (
    capabilityTranslations[language][capability.id] ?? {
      name: capability.name,
      description: capability.description,
    }
  );

  const copyCommand = async (id: string, command: string) => {
    try {
      if (!navigator.clipboard) throw new Error('Clipboard is unavailable');
      await navigator.clipboard.writeText(command);
      setCopied(id);
      window.setTimeout(() => setCopied(''), 1600);
    } catch {
      setCopied('');
    }
  };

  const selectLanguage = (nextLanguage: Language) => {
    setLanguage(nextLanguage);
    setLanguageOpen(false);
  };

  const activeCopy = localizedCapability(active);
  const installCommand = installCommands[installPlatform];

  return (
    <main>
      <nav className="topbar" aria-label={language === 'zh' ? '主导航' : 'Primary navigation'}>
        <a className="brand" href="#top" aria-label="agc-cli home">
          <span className="brandMark" aria-hidden="true"><RadioTower size={18} /></span>
          <span>agc<span className="brandSuffix">cli</span></span>
        </a>
        <div className="navlinks">
          <a href="#workflow">{text.nav.workflow}</a>
          <a href="#profiles">{text.nav.profiles}</a>
          <a href="#registry">{text.nav.registry}</a>
          <a href="#quickstart">{text.nav.quickstart}</a>
        </div>
        <div className="navActions">
          <div className="languagePicker" ref={languagePickerRef}>
            <button
              className="languageButton"
              type="button"
              aria-label={text.nav.language}
              aria-haspopup="menu"
              aria-expanded={languageOpen}
              onClick={() => setLanguageOpen((open) => !open)}
            >
              <Globe2 size={15} />
              <span>{language === 'zh' ? '中' : 'EN'}</span>
              <ChevronDown size={13} />
            </button>
            {languageOpen && (
              <div className="languageMenu" role="menu" aria-label={text.nav.language}>
                <button type="button" role="menuitemradio" aria-checked={language === 'en'} onClick={() => selectLanguage('en')}>
                  <span>EN</span><b>English</b>{language === 'en' && <Check size={14} />}
                </button>
                <button type="button" role="menuitemradio" aria-checked={language === 'zh'} onClick={() => selectLanguage('zh')}>
                  <span>中</span><b>简体中文</b>{language === 'zh' && <Check size={14} />}
                </button>
              </div>
            )}
          </div>
          <button ref={installTriggerRef} className="navInstall" type="button" onClick={() => setInstallOpen(true)}>
            <Download size={16} />
            <span>{text.nav.install}</span>
          </button>
        </div>
      </nav>

      <section className="hero" id="top">
        <div className="heroGrid" aria-hidden="true" />
        <div className="heroCopy">
          <div className="heroKicker"><span className="signalDot" /> {text.hero.kicker}</div>
          <h1>{text.hero.titleLead} <em>{text.hero.titleAccent}</em></h1>
          <p className="lede">{text.hero.body}</p>
          <div className="actions">
            <a className="primary" href="#quickstart">{text.hero.start} <ArrowRight size={18} /></a>
            <a className="secondary" href="#registry"><Layers3 size={18} /> {text.hero.browse}</a>
          </div>
          <div className="heroFacts" aria-label={text.hero.factsLabel}>
            <span><strong>13</strong> {text.hero.families}</span>
            <span><strong>{endpointCount}</strong> {text.hero.interfaces}</span>
            <span><strong>3</strong> {text.hero.formats}</span>
          </div>
        </div>

        <div className="controlDeck" aria-label={text.deck.label}>
          <div className="deckHeader">
            <span className="deckLights"><i /><i /><i /></span>
            <span>agc / production</span>
            <span className="deckMode"><span className={mode === 'live' ? 'pulse live' : 'pulse'} />{mode === 'live' ? text.deck.live : text.deck.demo}</span>
          </div>
          <div className="rail" aria-label={text.workflow.eyebrow}>
            {text.rail.map(({ label, detail }, index) => {
              const Icon = railIcons[index];
              return (
                <div className="railStop" key={label}>
                  <span className="railNode"><Icon size={17} /></span>
                  <span><b>{label}</b><small>{detail}</small></span>
                  {index < text.rail.length - 1 && <ChevronRight size={15} className="railArrow" />}
                </div>
              );
            })}
          </div>
          <div className="terminalBody">
            <p><span className="prompt">$</span> agc auth check --pretty</p>
            <pre>{'{\n  "data": {\n    "name": "production",\n    "mode": "service-account",\n    "active": true\n  },\n  "affordances": {\n    "next": "agc publishing endpoints"\n  }\n}'}</pre>
            <p className="terminalResult"><Check size={14} /> {text.deck.result}</p>
          </div>
        </div>
      </section>

      <section className="principles" aria-label={language === 'zh' ? '产品原则' : 'Product principles'}>
        {text.principles.map(({ title, body }, index) => {
          const Icon = principleIcons[index];
          return <div key={title}><Icon /><span><b>{title}</b>{body}</span></div>;
        })}
      </section>

      <section className="section workflowSection" id="workflow">
        <div className="sectionIntro">
          <p className="eyebrow">{text.workflow.eyebrow}</p>
          <h2>{text.workflow.title}</h2>
          <p>{text.workflow.body}</p>
        </div>
        <div className="workflowGrid">
          {text.rail.map(({ label, detail }, index) => {
            const Icon = railIcons[index];
            return (
              <article className="workflowCard" key={label}>
                <span className="workflowIndex">0{index + 1}</span>
                <Icon size={24} />
                <h3>{label}</h3>
                <p>{detail}</p>
                <span className="workflowLine" />
              </article>
            );
          })}
        </div>
      </section>

      <section className="section profileSection" id="profiles">
        <div className="profileCopy">
          <p className="eyebrow">{text.profile.eyebrow}</p>
          <h2>{text.profile.title}</h2>
          <p>{text.profile.body}</p>
          <ol className="precedence">
            {text.profile.items.map((item, index) => (
              <li key={item}><span>{index + 1}</span><b>{item}</b><code>{profileCodes[index]}</code></li>
            ))}
          </ol>
        </div>
        <div className="profileConfig">
          <div className="configTitle"><Fingerprint size={18} /> .agc/project.json <span>{text.profile.safe}</span></div>
          <pre>{'{\n  "appId": "123456789",\n  "projectId": "987654321",\n  "packageName": "com.example.app",\n  "profile": "production"\n}'}</pre>
          <div className="configNote"><ShieldCheck size={16} /> {text.profile.note}</div>
        </div>
      </section>

      <section className="section registrySection" id="registry">
        <div className="sectionIntro registryIntro">
          <div>
            <p className="eyebrow">{text.registry.eyebrow}</p>
            <h2>{text.registry.title}</h2>
          </div>
          <p>{text.registry.body}</p>
        </div>

        <div className="registryShell">
          <aside className="registryNav" aria-label={text.registry.families}>
            <div className="registryNavTitle"><Workflow size={16} /> {text.registry.families} <span>{capabilities.length}</span></div>
            {capabilities.map((capability) => {
              const capabilityCopy = localizedCapability(capability);
              return (
                <button
                  key={capability.id}
                  className={selected === capability.id ? 'registryItem selected' : 'registryItem'}
                  onClick={() => setSelected(capability.id)}
                  aria-pressed={selected === capability.id}
                >
                  <span className="familyCode">{capability.id.slice(0, 2).toUpperCase()}</span>
                  <span>{capabilityCopy.name}<small>{capability.endpointCount ?? 0} {text.registry.interfaces}</small></span>
                  <ChevronRight size={15} />
                </button>
              );
            })}
          </aside>

          <div className="registryDetail">
            <div className="detailHeader">
              <div>
                <p className="eyebrow">{text.registry.selected} / {active.id}</p>
                <h3>{activeCopy.name}</h3>
              </div>
              <span className="registeredBadge"><CircleDot size={14} /> {text.registry.registered}</span>
            </div>
            <p className="detailDescription">{activeCopy.description}</p>
            <div className="commandTable">
              <div><span>{text.registry.discover}</span><code>{active.command} endpoints --output table</code></div>
              <div><span>{text.registry.status}</span><code>{active.command} status --pretty</code></div>
              <div><span>{text.registry.localRest}</span><code>{active.restPath}</code></div>
            </div>
            <div className="detailFooter">
              <span><Code2 size={16} /> {active.endpointCount ?? 0} {text.registry.typed}</span>
              <button onClick={() => copyCommand('active', active.command + ' endpoints --output table')}>
                {copied === 'active' ? <Check size={15} /> : <Clipboard size={15} />}
                {copied === 'active' ? text.registry.copied : text.registry.copy}
              </button>
            </div>
          </div>
        </div>
      </section>

      <section className="section quickstartSection" id="quickstart">
        <div className="sectionIntro quickIntro">
          <p className="eyebrow">{text.quickstart.eyebrow}</p>
          <h2>{text.quickstart.title}</h2>
          <p>{text.quickstart.body}</p>
        </div>
        <div className="quickSteps">
          {text.quickstart.steps.map((step) => (
            <article className="quickStep" key={step.number}>
              <div className="stepNumber">{step.number}</div>
              <div className="stepCopy">
                <h3>{step.title}</h3>
                <p>{step.body}</p>
              </div>
              <div className="commandBox">
                <code><span>$</span> {step.command}</code>
                <button onClick={() => copyCommand(step.number, step.command)} aria-label={text.quickstart.copyStep + ' ' + step.number}>
                  {copied === step.number ? <Check size={16} /> : <Clipboard size={16} />}
                </button>
              </div>
            </article>
          ))}
        </div>
      </section>

      <section className="section consoleSection" id="console">
        <div className="consoleCopy">
          <p className="eyebrow">{text.console.eyebrow}</p>
          <h2>{text.console.title}</h2>
          <p>{text.console.body}</p>
          <div className="serverCommand"><Terminal size={18} /><code>agc web-server --addr :8421</code></div>
        </div>
        <div className="consoleMap" aria-label={text.console.routesLabel}>
          <div><span>GET</span><code>/api/v1/capabilities</code><small>{text.console.routes[0]}</small></div>
          <div><span>GET</span><code>/api/v1/endpoints</code><small>{text.console.routes[1]}</small></div>
          <div><span>GET</span><code>/api/v1/openapi.json</code><small>{text.console.routes[2]}</small></div>
          <div><span>POST</span><code>/api/v1/:family/endpoints/:id/invoke</code><small>{text.console.routes[3]}</small></div>
        </div>
      </section>

      <footer>
        <div className="footerBrand">
          <span className="brandMark"><RadioTower size={18} /></span>
          <span><b>agccli.app</b><small>{text.footer.tagline}</small></span>
        </div>
        <div className="footerMeta"><span>{text.footer.license}</span><span>{text.footer.stack}</span><span>{text.footer.cloud}</span></div>
      </footer>

      {installOpen && (
        <div className="installBackdrop" onMouseDown={(event) => event.currentTarget === event.target && setInstallOpen(false)}>
          <section className="installDialog" role="dialog" aria-modal="true" aria-labelledby="install-title">
            <div className="installDialogHeader">
              <span className="deckLights"><i /><i /><i /></span>
              <span>{text.install.terminal}</span>
              <button ref={closeInstallRef} type="button" aria-label={text.install.close} onClick={() => setInstallOpen(false)}>
                <X size={17} />
              </button>
            </div>
            <div className="installDialogBody">
              <div className="installGlyph" aria-hidden="true"><MonitorDown size={24} /></div>
              <p className="eyebrow">{text.install.eyebrow}</p>
              <h2 id="install-title">{text.install.title}</h2>
              <p className="installLead">{text.install.body}</p>

              <div className="platformTabs" role="tablist" aria-label={language === 'zh' ? '安装平台' : 'Install platform'}>
                <button type="button" role="tab" aria-selected={installPlatform === 'unix'} onClick={() => setInstallPlatform('unix')}>
                  {text.install.unix}
                </button>
                <button type="button" role="tab" aria-selected={installPlatform === 'windows'} onClick={() => setInstallPlatform('windows')}>
                  {text.install.windows}
                </button>
              </div>

              <div className="installTerminal">
                <span className="installPrompt">{installPlatform === 'windows' ? 'PS>' : '$'}</span>
                <code>{installCommand}</code>
                <button type="button" onClick={() => copyCommand('install', installCommand)}>
                  {copied === 'install' ? <Check size={16} /> : <Clipboard size={16} />}
                  <span>{copied === 'install' ? text.install.copied : text.install.copy}</span>
                </button>
              </div>

              <p className="installNote">{text.install.note}</p>
              <div className="installAssurances">
                <span><ShieldCheck size={15} /> {text.install.checksum}</span>
                <span><Download size={15} /> {text.install.destination}</span>
              </div>
            </div>
          </section>
        </div>
      )}
    </main>
  );
}
