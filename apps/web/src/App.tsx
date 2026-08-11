import {
  ArrowDownRight,
  ArrowRight,
  Check,
  ChevronRight,
  CircleDot,
  Clipboard,
  Cloud,
  Code2,
  FileJson2,
  Fingerprint,
  KeyRound,
  Layers3,
  PackageCheck,
  RadioTower,
  Rocket,
  ShieldCheck,
  Terminal,
  Workflow,
} from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';

type Capability = {
  id: string;
  name: string;
  description: string;
  command: string;
  restPath: string;
  status: string;
  endpointCount?: number;
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

const quickSteps = [
  {
    number: '01',
    title: 'Save a credential profile',
    body: 'Service Account is the preferred route. Secrets stay in ~/.agc/credentials.json with restricted permissions.',
    command: 'agc auth login --service-account-file ~/.agc/service-account.json --name production',
  },
  {
    number: '02',
    title: 'Bind it to the project',
    body: 'Pin the app, project, package, and default profile once. Every later command resolves that context automatically.',
    command: 'agc init --app-id <app-id> --project-id <project-id> --package-name com.example.app --default-profile production',
  },
  {
    number: '03',
    title: 'Inspect before invoking',
    body: 'Endpoint commands are dry-run first. Review the exact method and URL, then opt into the real request.',
    command: 'agc publishing app-info-query --invoke --query appId=<app-id> --dry-run=false',
  },
];

const railSteps = [
  { label: 'AUTH', detail: 'PS256 / OAuth', icon: KeyRound },
  { label: 'PROFILE', detail: 'project-bound', icon: Fingerprint },
  { label: 'PACKAGE', detail: 'upload + status', icon: PackageCheck },
  { label: 'RELEASE', detail: 'invoke by choice', icon: Rocket },
];

export function App() {
  const [capabilities, setCapabilities] = useState<Capability[]>(fallbackCapabilities);
  const [endpointCount, setEndpointCount] = useState(156);
  const [mode, setMode] = useState<'demo' | 'live'>('demo');
  const [selected, setSelected] = useState('publishing');
  const [copied, setCopied] = useState('');

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

  const active = useMemo(
    () => capabilities.find((capability) => capability.id === selected) ?? capabilities[0],
    [capabilities, selected],
  );

  const copyCommand = async (id: string, command: string) => {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(id);
      window.setTimeout(() => setCopied(''), 1600);
    } catch {
      setCopied('');
    }
  };

  return (
    <main>
      <nav className="topbar" aria-label="Primary navigation">
        <a className="brand" href="#top" aria-label="agc-cli home">
          <span className="brandMark" aria-hidden="true"><RadioTower size={18} /></span>
          <span>agc<span className="brandSuffix">cli</span></span>
        </a>
        <div className="navlinks">
          <a href="#workflow">Workflow</a>
          <a href="#profiles">Profiles</a>
          <a href="#registry">API registry</a>
          <a href="#quickstart">Quick start</a>
        </div>
        <a className="navCta" href="#quickstart">Run the first command <ArrowDownRight size={15} /></a>
      </nav>

      <section className="hero" id="top">
        <div className="heroGrid" aria-hidden="true" />
        <div className="heroCopy">
          <div className="heroKicker"><span className="signalDot" /> AppGallery Connect, mapped for the terminal</div>
          <h1>Move every release through <em>one red line.</em></h1>
          <p className="lede">A Go command center for discovering, dry-running, and invoking AppGallery Connect APIs—with stable JSON for people, CI, and coding agents.</p>
          <div className="actions">
            <a className="primary" href="#quickstart">Get started <ArrowRight size={18} /></a>
            <a className="secondary" href="#registry"><Layers3 size={18} /> Browse 156 interfaces</a>
          </div>
          <div className="heroFacts" aria-label="Project facts">
            <span><strong>13</strong> API families</span>
            <span><strong>{endpointCount}</strong> registered interfaces</span>
            <span><strong>3</strong> output formats</span>
          </div>
        </div>

        <div className="controlDeck" aria-label="Command output example">
          <div className="deckHeader">
            <span className="deckLights"><i /><i /><i /></span>
            <span>agc / production</span>
            <span className="deckMode"><span className={mode === 'live' ? 'pulse live' : 'pulse'} />{mode === 'live' ? 'local API' : 'reference mode'}</span>
          </div>
          <div className="rail" aria-label="Release workflow">
            {railSteps.map(({ label, detail, icon: Icon }, index) => (
              <div className="railStop" key={label}>
                <span className="railNode"><Icon size={17} /></span>
                <span><b>{label}</b><small>{detail}</small></span>
                {index < railSteps.length - 1 && <ChevronRight size={15} className="railArrow" />}
              </div>
            ))}
          </div>
          <div className="terminalBody">
            <p><span className="prompt">$</span> agc auth check --pretty</p>
            <pre>{`{
  "data": {
    "name": "production",
    "mode": "service-account",
    "active": true
  },
  "affordances": {
    "next": "agc publishing endpoints"
  }
}`}</pre>
            <p className="terminalResult"><Check size={14} /> project profile resolved · dry-run on</p>
          </div>
        </div>
      </section>

      <section className="principles" aria-label="Product principles">
        <div><ShieldCheck /><span><b>Dry-run first</b>Nothing is sent until you opt in.</span></div>
        <div><FileJson2 /><span><b>Agent-readable</b>JSON envelopes include next actions.</span></div>
        <div><Cloud /><span><b>One surface</b>CLI, local REST, and web share a registry.</span></div>
      </section>

      <section className="section workflowSection" id="workflow">
        <div className="sectionIntro">
          <p className="eyebrow">The release rail</p>
          <h2>Context flows forward. Secrets do not.</h2>
          <p>Each stage owns one job. Credentials remain local, project context remains portable, and destructive API calls require an explicit switch.</p>
        </div>
        <div className="workflowGrid">
          {railSteps.map(({ label, detail, icon: Icon }, index) => (
            <article className="workflowCard" key={label}>
              <span className="workflowIndex">0{index + 1}</span>
              <Icon size={24} />
              <h3>{label}</h3>
              <p>{detail}</p>
              <span className="workflowLine" />
            </article>
          ))}
        </div>
      </section>

      <section className="section profileSection" id="profiles">
        <div className="profileCopy">
          <p className="eyebrow">Automatic profile binding</p>
          <h2>Choose once per project. Override only when you mean it.</h2>
          <p>Bind <code>production</code> in <code>.agc/project.json</code>. Commands resolve the profile in a predictable order, while <code>--profile staging</code> remains available for one-off work.</p>
          <ol className="precedence">
            <li><span>1</span><b>CLI override</b><code>--profile staging</code></li>
            <li><span>2</span><b>Project default</b><code>.agc/project.json</code></li>
            <li><span>3</span><b>Active account</b><code>~/.agc/credentials.json</code></li>
          </ol>
        </div>
        <div className="profileConfig">
          <div className="configTitle"><Fingerprint size={18} /> .agc/project.json <span>project safe</span></div>
          <pre>{`{
  "appId": "123456789",
  "projectId": "987654321",
  "packageName": "com.example.app",
  "profile": "production"
}`}</pre>
          <div className="configNote"><ShieldCheck size={16} /> No client keys or private keys are stored here.</div>
        </div>
      </section>

      <section className="section registrySection" id="registry">
        <div className="sectionIntro registryIntro">
          <div>
            <p className="eyebrow">Interface registry</p>
            <h2>Discover the API before writing the request.</h2>
          </div>
          <p>Every entry carries its method, path, parameter locations, CLI command, REST route, and dry-run path.</p>
        </div>

        <div className="registryShell">
          <aside className="registryNav" aria-label="API families">
            <div className="registryNavTitle"><Workflow size={16} /> API families <span>{capabilities.length}</span></div>
            {capabilities.map((capability) => (
              <button
                key={capability.id}
                className={selected === capability.id ? 'registryItem selected' : 'registryItem'}
                onClick={() => setSelected(capability.id)}
                aria-pressed={selected === capability.id}
              >
                <span className="familyCode">{capability.id.slice(0, 2).toUpperCase()}</span>
                <span>{capability.name}<small>{capability.endpointCount ?? 0} interfaces</small></span>
                <ChevronRight size={15} />
              </button>
            ))}
          </aside>

          <div className="registryDetail">
            <div className="detailHeader">
              <div>
                <p className="eyebrow">Selected family / {active.id}</p>
                <h3>{active.name}</h3>
              </div>
              <span className="registeredBadge"><CircleDot size={14} /> {active.status}</span>
            </div>
            <p className="detailDescription">{active.description}</p>
            <div className="commandTable">
              <div><span>DISCOVER</span><code>{active.command} endpoints --output table</code></div>
              <div><span>STATUS</span><code>{active.command} status --pretty</code></div>
              <div><span>LOCAL REST</span><code>{active.restPath}</code></div>
            </div>
            <div className="detailFooter">
              <span><Code2 size={16} /> {active.endpointCount ?? 0} typed interface entries</span>
              <button onClick={() => copyCommand('active', `${active.command} endpoints --output table`)}>
                {copied === 'active' ? <Check size={15} /> : <Clipboard size={15} />}
                {copied === 'active' ? 'Copied' : 'Copy discover command'}
              </button>
            </div>
          </div>
        </div>
      </section>

      <section className="section quickstartSection" id="quickstart">
        <div className="sectionIntro quickIntro">
          <p className="eyebrow">Quick start</p>
          <h2>Three commands to a bound, inspectable workspace.</h2>
          <p>Use placeholders for your AppGallery Connect values. The second command is what makes the profile automatic for this project.</p>
        </div>
        <div className="quickSteps">
          {quickSteps.map((step) => (
            <article className="quickStep" key={step.number}>
              <div className="stepNumber">{step.number}</div>
              <div className="stepCopy">
                <h3>{step.title}</h3>
                <p>{step.body}</p>
              </div>
              <div className="commandBox">
                <code><span>$</span> {step.command}</code>
                <button onClick={() => copyCommand(step.number, step.command)} aria-label={`Copy step ${step.number} command`}>
                  {copied === step.number ? <Check size={16} /> : <Clipboard size={16} />}
                </button>
              </div>
            </article>
          ))}
        </div>
      </section>

      <section className="section consoleSection" id="console">
        <div className="consoleCopy">
          <p className="eyebrow">Local command center</p>
          <h2>The browser is a window into the same registry.</h2>
          <p>Start the local REST server when you want live capability data. Without it, this site stays useful as a complete reference surface.</p>
          <div className="serverCommand"><Terminal size={18} /><code>agc web-server --addr :8421</code></div>
        </div>
        <div className="consoleMap" aria-label="Local server routes">
          <div><span>GET</span><code>/api/v1/capabilities</code><small>families + affordances</small></div>
          <div><span>GET</span><code>/api/v1/endpoints</code><small>all interface entries</small></div>
          <div><span>GET</span><code>/api/v1/openapi.json</code><small>machine-readable contract</small></div>
          <div><span>POST</span><code>/api/v1/:family/endpoints/:id/invoke</code><small>dry-run by default</small></div>
        </div>
      </section>

      <footer>
        <div className="footerBrand">
          <span className="brandMark"><RadioTower size={18} /></span>
          <span><b>agccli.app</b><small>AppGallery Connect, from one command surface.</small></span>
        </div>
        <div className="footerMeta"><span>MIT licensed</span><span>Go + React</span><span>Built for Cloudflare</span></div>
      </footer>
    </main>
  );
}
