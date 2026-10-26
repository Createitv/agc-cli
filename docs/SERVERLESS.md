# 云开发（Serverless）

本工作区接入八类服务的 102 项操作：49 个官方文档明确的 HTTP/开发者回调入口、42 个依据华为 AGCLI 1.12.0 源码实现的管理传输、10 个原生工作流，以及 1 个设备调试入口。SDK 方法、控制台功能和未公开的资源管理接口不自动算作已支持的 HTTP API。新增管理、SAM 部署及预加载的用法与边界见[管理扩展](SERVERLESS_MANAGEMENT.md)。

所有操作默认 dry-run；增加 `--invoke --dry-run=false` 才会发送。`--docs --pretty` 显示本地保存的完整官方参数表、例子、段落链接、错误码、前提文档及已记录冲突。命令、REST 和 OpenAPI 使用同一领域契约。

| 命令族 | 操作数 | 范围 |
| --- | ---: | --- |
| `agc cloud-auth` | 17 | 客户端令牌、验证码、注册、密码/验证码登录、更新身份/密码、关联/取消关联、登出、注销、重认证；服务端导入/导出、验证/撤销用户凭据 |
| `agc clouddb` | 16 | 客户端数据操作，以及 schema、存储区与异步进度管理；两者鉴权不同 |
| `agc storage` | 6 | 原始文件上传下载、删除、元数据读写、分页列表 |
| `agc hosting` | 16 | 站点四项、环境三项、版本九项（含签名上传） |
| `agc functions` | 29 | 原有 9 个工作流，以及别名、HTTP/定时触发器、调用策略和层管理传输 |
| `agc templates` | 8 | 2 个认证回调、5 个 SAM 应用/计划管理传输、1 个本地 ZIP 到新应用的部署工作流；不等同于控制台模板库 |
| `agc cache` | 6 | 命令执行、连接信息、三天内用量、开通、初设密码与重置密码 |
| `agc prefetch` | 4 | 安装/周期性/跳链安装的开发者服务器 GET 模拟，以及显式设备的 `cf_prefetch getcache` 调试 |

## 发现、文档和请求

```sh
./bin/agc capabilities
./bin/agc storage endpoints --pretty
./bin/agc storage upload-file --docs --pretty
./bin/agc clouddb query-data --docs --pretty
./bin/agc functions create-function --docs --pretty
./bin/agc openapi > /tmp/agc-openapi.json
```

公开 REST 使用 `--param`（路径）、`--query`、`--header`、`--body` 或 `--field`。文件传输使用 `--file`；下载用 `--out` 保存原始字节。REST 的二进制下载响应用 `bodyBase64` 无损承载。

```sh
# 预演文件上传，确认后才增加 --dry-run=false。
./bin/agc --profile dedicated-test storage upload-file --invoke \
  --param bucket_name=TEST_BUCKET --param object_name=checks/hello.txt \
  --header productId=TEST_PROJECT --file ./hello.txt

# 预演站点查询。requestid 必须唯一，service 固定为 hosting。
./bin/agc --profile dedicated-test hosting list-sites --invoke \
  --header productId=TEST_PROJECT --header requestid=UNIQUE_REQUEST_ID \
  --header service=hosting --field 'site={}'

# 数据库查询：queryCond 是 JSON 字符串，opStore 按官方例子是对象。
./bin/agc --profile dedicated-test clouddb query-data --invoke \
  --header productId=TEST_PROJECT --query _v=4 --body ./query.json
```

## 鉴权边界

- 认证服务端的四项操作支持 Service Account 或 API 客户端；API 客户端分支需要 `client_id`，不得同时强制两种分支。
- 客户端认证和 CloudDB 的客户端网关令牌与发布接口令牌分开获取；自动令牌获取使用官方记录的中国站点。其他未核实的客户端区域应显式提供令牌和对应地址。
- 云存储和云托管使用 API 客户端；存储实例域名对应的 Token 区域必须一致。CLI 拒绝把 Service Account 自动用于这些客户端专用服务。
- 云托管上传使用返回的 `uploadInfo.url` 和整组签名头，`--param uploadUrl=... --file ...`；不能加入 AGC Bearer。上传前生成签名和归档/发布状态均需分别读回。
- 模板认证接口是华为向开发者服务器发送的 JSON 回调。`templates` 的调用仅用于模拟，不读取凭据、不转发 AGC Bearer，不证明华为已实际调用或模板已部署。

## 原生云函数

源码依据为 `@hw-agconnect/agc-cli@1.12.0`，不是公开 REST 规范。源码路径和 SHA256 保存在 `--docs` 中。运行不依赖 Node 子进程，不读取官方 CLI 的加密 scratch，不复用发布凭据，也不复现官方包关闭 TLS 验证的行为。

显式 OAuth 会话 JSON 存在 Git 仓库外，权限建议 `0600`。字段为 `uid`、`teamId`、`productId`、`accessToken`、RFC3339 格式的 `expiresAt` 和项目的 `projectBaseUrl`。这些值需来自开发者自己授权的有效会话，不要将真实值写入文档或报告。CLI 会拒绝过期会话和项目不一致；dry-run 不读取会话。

```sh
./bin/agc functions create-function --invoke --name test-function \
  --session-file /PRIVATE_PATH/session.json --project-id TEST_PROJECT \
  --runtime nodejs14.x --code-path ./function.zip
./bin/agc functions update-function-code --invoke --name test-function \
  --session-file /PRIVATE_PATH/session.json --code-path ./function.zip
./bin/agc functions update-function-config --invoke --name test-function \
  --session-file /PRIVATE_PATH/session.json --memory 1000
./bin/agc functions publish-function --invoke --name test-function \
  --session-file /PRIVATE_PATH/session.json --version-description 'test release'
./bin/agc functions invoke --invoke --name test-function \
  --session-file /PRIVATE_PATH/session.json --payload '{"request":"test"}'
```

创建时内存、超时、并发默认值分别为 500、55、100；运行时默认值来自所审阅官方 CLI（nodejs14.x），不代表华为当前仍允许创建这个运行时。其他运行时及上下限见 `--docs` 和 OpenAPI，云端可用性需实测。

更新先读取 `$latest` 的 revision，保留未提供的环境变量、层和 handler；冲突停止，不自动重试。发布失败会明确报告更新可能已生效。管理接口返回的环境值和用户会话不进入响应；函数调用的业务结果按原始 JSON 保留，不把业务 code 或 uid 当作管理错误或凭据。HTTP 状态与执行状态独立返回。

REST 路径仍为 `/api/v1/functions/endpoints/OP/invoke`，请求根字段对应函数输入（`name`、`sessionFile`、`codePath`、`payload` 等）；默认 `dryRun:true`。不能套用通用接口的 `body`、`token` 封装。

## 数据与文档冲突

已保存官方表格和请求例子。CloudDB 的 `opStore` 文本称 JSON 字符串，实际例子是对象；适配器采用实际例子的对象。`queryCond` 明确保留为字符串，并验证解码后的条件及聚合类型。更新与删除共用 upsert 路径，`naturalbase_deleted` 分别必须为 false/true，列顺序与值数量必须匹配。删除表格对字段个数的描述与例子冲突，没有静默应用会拒绝官方例子的单字段上限。

认证导出的 `block` 表格为字符串、例子为整数；响应 `users`/`ret` 也存在字符串与对象/数组描述冲突。适配器保留实际类型，不自行转换。云存储元数据更新的 GET 与 Java 示例一致。

## 验证与限制

离线验证涵盖每项操作的 CLI HTTP 收发、REST 收发、dry-run、发现/文档和 OpenAPI，另有独立来源断言、嵌套字段、条件规则、文件字节、签名头、区域鉴权、revision 工作流与失败场景测试。模拟服务器通过不代表华为真实业务通过。

真实验证必须显式指定 profile、项目和专用资源；未指定时不选择活动业务账号、不发送验证码、不部署业务函数、不发布站点。报告中每项 live 状态为未执行，直至用专用资源执行并完成相应读回。当前未取得可复核管理协议的能力包括控制台模板库部署/升级、预加载云端绑定配置、层删除，以及非 HTTP/定时类型触发器的创建；这些没有被猜测为可调用接口。详见[管理扩展](SERVERLESS_MANAGEMENT.md)。

开发验证：`make ci`、`go test -race ./...`、`make build`。本次工作区二进制是 `./bin/agc`；不代表系统安装的 `agc` 已升级，也不代表已发版或 GitHub Actions 已通过。
