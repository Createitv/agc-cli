# 官网接口注册表

官网内置 Go 注册表导出的全部接口，无需启动本地 REST 服务即可浏览。启动本地服务后，页面优先使用 `/api/v1/endpoints` 返回的有效接口定义。

修改 `pkg/domain` 中的接口契约后，在此目录运行：

```bash
npm run registry:generate
npm test -- --run
npm run build
npm run check:seo
```

将更新后的 `src/endpoint-registry.json` 与契约变更一起提交。生成过程只读取公开接口定义，不鉴权、不发送华为 API 请求。参数表展示已登记的字段；完整请求、响应、权限及业务条件请使用每个接口的参考文档链接。

## 中英文 SEO 与 AEO

默认中文路径 `/`，英文路径 `/en/`。每个 API 家族有独立参考路径 `/reference/{family}/` 和 `/en/reference/{family}/`。构建时生成全部 28 个静态 HTML 页面，浏览器再接管交互；无 JavaScript 的读取端也能看到产品介绍、FAQ 和该家族的接口列表。

`src/siteContent.ts` 是双语项目定义、FAQ、元信息和 JSON-LD 的共用来源。不要把注册范围写成生产验证范围，不添加未经验证的版本、评分、用户数或性能数据。页面 FAQ 与结构化数据必须一致。

构建同时生成双语 Markdown、`api-reference.json` 和带语言对应关系的 sitemap。`public/llms.txt` 提供辅助入口，`robots.txt` 允许 OAI-SearchBot；这些文件与结构化数据不保证收录、排名或模型推荐。ChatGPT 实测证据保存在仓库 `docs/verification/seo-aeo/`，区分上线前基线、搜索引用和显式提供 URL 的读取测试。

### ChatGPT 读取兼容性

2026-10-07 的实际测试发现：ChatGPT 可以读取中英文 HTML 页面，但读取 `overview.zh-CN.md` 时报告 `Unsupported content-type: text/markdown`。因此 `overview.zh-CN.md` 和 `overview.en.md` 保留 Markdown 内容和文件名，通过 `public/_headers` 以 `text/plain; charset=utf-8` 返回。变更后需要实际读回响应头，并在 ChatGPT 中重测，不能仅凭文件上传成功认定可读取。

同日部署后复测：两个介绍文件均返回 HTTP 200、`text/plain; charset=utf-8`，ChatGPT 实际读取并引用了中英文 Profile 参数说明。完整错误与成功读回应保留在 `docs/verification/seo-aeo/chatgpt-direct-initial.json` 和 `chatgpt-compatibility-fixed.json` 中；不把显式提供 URL 后的读取成功当作搜索排名提升。
