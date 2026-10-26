export type Endpoint = {
  id: string;
  familyId: string;
  name: string;
  description: string;
  method: string;
  path: string;
  command: string;
  parameters?: { name: string; in: string; required: boolean; type?: string; description?: string }[];
  body?: string;
  sourceUrl?: string;
  direction?: string;
};

export function EndpointList({ endpoints, language }: { endpoints: Endpoint[]; language: 'en' | 'zh' }) {
  const zh = language === 'zh';
  return (
    <section className="endpointList" aria-label={zh ? '完整接口列表' : 'Complete endpoint list'}>
      <h4>{zh ? '全部接口' : 'All endpoints'} <span>{endpoints.length}</span></h4>
      <p className="endpointNote">{zh ? '展开接口查看已登记参数与命令；完整协议、响应字段及权限以官方文档为准。' : 'Expand an endpoint for registered parameters and commands. Consult its reference for full schemas, responses, and permissions.'}</p>
      {endpoints.map(endpoint => (
        <details className="endpointEntry" key={endpoint.id}>
          <summary>
            <span className={`httpMethod method-${endpoint.method.toLowerCase()}`}>{endpoint.method}</span>
            <span className="endpointHeading"><b>{endpoint.name}</b><code>{endpoint.id}</code><code>{endpoint.path}</code></span>
          </summary>
          <div className="endpointContent">
            <p>{endpoint.description}</p>
            <p className="endpointDirection">{endpoint.direction === 'inbound-callback' ? (zh ? '华为回调到开发者服务' : 'Inbound callback to your service') : endpoint.direction === 'local' ? (zh ? '本地集成接口' : 'Local integration') : (zh ? '开发者调用华为 API' : 'Developer-to-Huawei API')}</p>
            <div className="endpointParameters">
              <table>
                <caption>{zh ? '已登记参数' : 'Registered parameters'}</caption>
                <thead><tr>{(zh ? ['参数', '位置', '必填', '类型', '说明'] : ['Parameter', 'Location', 'Required', 'Type', 'Description']).map(label => <th key={label}>{label}</th>)}</tr></thead>
                <tbody>{(endpoint.parameters ?? []).map(parameter => <tr key={`${parameter.in}-${parameter.name}`}>
                  <td><code>{parameter.name}</code></td><td>{parameter.in}</td><td>{parameter.required ? (zh ? '是' : 'Yes') : (zh ? '否' : 'No')}</td><td>{parameter.type ?? '—'}</td><td>{parameter.description ?? '—'}</td>
                </tr>)}</tbody>
              </table>
              {!endpoint.parameters?.length && <p>{zh ? '注册表未列出参数，请查看接口参考。' : 'Parameters are not listed in the registry. See the reference.'}</p>}
            </div>
            {endpoint.body && <p>{zh ? '请求体' : 'Request body'}: <code>{endpoint.body}</code></p>}
            <p>{zh ? '查看接口定义' : 'Inspect definition'}</p>
            <pre>{endpoint.command} --pretty</pre>
            <p>{zh ? '本地 REST 定义地址' : 'Local REST definition'}</p>
            <pre>/api/v1/{endpoint.familyId}/endpoints/{endpoint.id}</pre>
            {endpoint.direction !== 'inbound-callback' && <>
              <p>{zh ? '调用预演（先填写所需参数）' : 'Dry run (supply required parameters first)'}</p>
              <pre>{endpoint.command} --invoke --dry-run=true{endpoint.body ? ' --body body.json' : ''}</pre>
            </>}
            {endpoint.sourceUrl && <a href={endpoint.sourceUrl.startsWith('https://') ? endpoint.sourceUrl : `https://github.com/Createitv/agc-cli/blob/main/${endpoint.sourceUrl}`} target="_blank" rel="noreferrer">{zh ? '查看接口参考文档 ↗' : 'Read endpoint reference ↗'}</a>}
          </div>
        </details>
      ))}
    </section>
  );
}
