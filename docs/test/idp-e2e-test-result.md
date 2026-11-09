# RAM IdP 端到端测试结果

- 执行时间：2026-11-07 09:32 – 09:42（Asia/Shanghai）
- 执行环境：Go 1.27.3 darwin/amd64；PostgreSQL 127.0.0.1:5432 / `test-make-01`；租户 `1000`
- 测试程序：`test/idp/`
- 说明文档：[idp-e2e-test.md](./idp-e2e-test.md)

## 1. 汇总

| 运行 | 命令 | 结果 | 用时 | 子用例 |
| --- | --- | --- | --- | --- |
| 默认（物理删除，`pg.data.delete=true`） | `go test ./test/idp -run TestIDP -v -count=1` | PASS | 3.118s | 140 通过，0 失败，1 跳过 |
| 逻辑删除（`IDP_TEST_SOFT_DELETE=true`） | `IDP_TEST_SOFT_DELETE=true go test ...` | PASS | 3.481s | 140 通过，0 失败，1 跳过 |
| 实机冒烟（gin `:9981`） | `IDP_LIVE_BASE=http://127.0.0.1:9981 go test .../live-server-smoke` | PASS | 1.982s | 1 通过 |

“跳过”仅为未设置 `IDP_LIVE_BASE` 时的实机冒烟用例；默认两次运行不访问外部端口。

## 2. 分组结果

| 分组 | 通过 | 失败 | 跳过 |
| --- | --- | --- | --- |
| routes-registered | 1 | 0 | 0 |
| schema（库表存在性） | 1 | 0 | 0 |
| identity-provider | 19 | 0 | 0 |
| identity-source | 14 | 0 | 0 |
| identity-source-callback | 14 | 0 | 0 |
| idp-binding | 13 | 0 | 0 |
| idp-credential | 12 | 0 | 0 |
| idp-metadata-cache | 12 | 0 | 0 |
| idp-login（OAuth2/OpenID） | 19 | 0 | 0 |
| cas | 16 | 0 | 0 |
| saml | 9 | 0 | 0 |
| ldap | 10 | 0 | 0 |
| live-server-smoke | 0 | 0 | 1 |
| **合计** | **140** | **0** | **1** |

两种删除模式的分组通过数完全一致。

## 3. 关键验证点

- 6 个管理端控制器全部 12 个标准端点（create/update/detail/enable/disable/state/delete/recovery/physicalDeletion/query/selectPublic 及 provider/source/callback 的 existName）均可达。
- 鉴权中间件生效：无 token 请求返回“Authorization 参数不能为空”。
- 创建→详情→更新→查询→selectPublic→停用/启用→删除 全链路字段值一致。
- 凭证详情不返回 `value`；元数据详情不返回 `metadataRaw`。
- OAuth2：首次登录 `isSignup=true` 并创建账号与绑定；二次登录 `isSignup=false`；绑定中保存 access/refresh token。
- OpenID：`protocol=openid` 时绑定协议落库为 `openid`，id_token claims 解析成功。
- 刷新 token：携带 access + refresh 头成功换发；缺少刷新头返回“刷新token失败”。
- CAS：ST 可签发；1.0 返回 `yes/user`；2.0/3.0 返回 `cas:authenticationSuccess`；票据复用、service 不匹配、非 ST、未知 ST、缺参分别返回预期错误。
- SAML：有证书时 302 跳转到 IdP SSO 且含 `SAMLRequest`；无证书/协议不符/停用/坏回调均返回明确错误。
- LDAP：连接被拒时 `connected=false`；搜索/同步/登录返回连接失败；停用源与缺参被拦截。
- 实机冒烟：用测试进程签发的 PASETO 访问真实 `:9981` 查询接口，返回 `code=200`。

## 4. 数据清理验证

运行结束后抽查（均针对测试前缀）：

| 数据 | 残留 |
| --- | --- |
| `ram_identity_provider`（name LIKE `E2E-%`） | 0 |
| `ram_identity_source`（name LIKE `E2E-%`） | 0 |
| `ram_identity_source_callback`（name LIKE `E2E-%`） | 0 |
| `ram_idp_binding`（external_sub LIKE `e2e-%`） | 0 |
| `ram_account`（account LIKE `e2e_%`） | 0 |
| `ram_account_session_log`（external_sub LIKE `e2e-%`） | 0 |
| `ram_idp_credential`（关联测试 provider） | 0 |
| `ram_idp_metadata_cache`（关联测试 provider） | 0 |

## 5. 缺陷与风险清单

| ID | 严重度 | 现象 | 证据/位置 | 建议 |
| --- | --- | --- | --- | --- |
| D1 | 高（已修复） | 启动建表 runner 未注册 6 张 IdP 表，也未注册 `ram_account_session_log`、`ram_account_auth_mfa`、`ram_account_auth_webauthn`；新库启动后 IdP 接口报 `relation does not exist`（部分服务吞掉错误后返回空列表） | `pkg/middleware/runnerPg/table/a-init-table.go`；修复后 `go run .` 建表，`public` 表数 68 → 77 | 已补充 9 个实体到建表列表。建议补一条启动/CI 校验，确保实体与建表清单同步 |
| D2 | 中 | `IdentityProvider.Vo`、`IdentitySource.Vo` 没有 `no` 字段，而 IdP 登录必须用 `sourceNo`、认证源必须存 provider `No`；管理端无法从接口取得 No | `../../app/models/idp/modRamIdentityProvider`、`modRamIdentitySource/vo.go` | Vo 增加 `no`；或 detail 返回 No；或提供 by-id 转 No 的接口 |
| D3 | 中 | `pg.data.delete=true` 时 `LogicalDeletion` 实际物理删除，`recovery` 必然“数据不存在”，`physicalDeletion` 无实际作用 | 各 `Ram*Service.LogicalDeletion` 中 `if c.sv.Config().Data.Delete` 分支；默认配置 `data/config/app.toml` | 若需要回收站语义，应改为“删除只置 state，physicalDeletion 才删”；否则前端隐藏恢复入口并更名 |
| D4 | 低 | 凭证模型有 `description`，实体/表没有该列：创建/更新描述不落库，详情描述恒为空 | `entityRam/RamIdpCredential.go` 与 `../../app/models/idp/modRamIdpCredential` | 实体补 `description` 列或模型删除该字段 |
| D5 | 高 | CAS 登录路由未挂鉴权中间件，而 `CasService.Login` 依赖 `ctx.Get("ano")`，HTTP 直连永远返回“用户未登录” | `../../app/handler/manage/idp/AuthCas.go`、`authLoginIdp.go` 路由组；`../../app/domain/domainRam/service` | 给 `/cas/login`、`/idp/cas/login` 挂载管理端鉴权中间件，或改用 `holderPg` 读取登录上下文 |
| D6 | 低 | LDAP 登录未校验 `source.Protocol == "ldap"`，非 LDAP 认证源也会走 LDAP 连接 | `../../app/domain/domainRam/service` `Login` | 增加协议校验，避免错误配置产生误导性错误 |
| D7 | 低 | 回调白名单创建只校验 `name`，`callback` 可为空 | `RamIdentitySourceCallbackService.Create` | 补 `callback` 必填与 URL 格式校验 |
| D8 | 中 | `Custom` provider 的 `UserMapping` 未从 `baseConfig` 填充；配置了 `userInfoURL` 的 Custom 源必然报 `cannot find id in userMapping`，只有 id_token 模式可用 | `AccountLoginIdpService.Login` 构造 `idp.ProviderInfo` 时未赋 `UserMapping`；`idp/custom.go` | 增加 `userMapping` 配置解析，或为 Custom 提供默认映射 |
| D9 | 中 | 绑定停用后仍可登录：`FindByIdpAndExternalSub/OpenId/UnionId` 不校验 `state` | `repositoryRam/RamIdpBindingRepository.go`；用例 `idp-login/disabled-binding-still-logs-in` | 绑定查询加 `state=1`，或按 `state_bind` 判定 |
| D10 | 低 | 元数据缓存无法通过管理接口写入 `metadataRaw`（`CreateUpdateCt` 无该字段）；SAML 又把 `sp_cert` 当作 IdP 证书读取，凭证类型语义混用 | `modRamIdpMetadataCache/createUpdateCt.go`；`SamlLoginService.buildSpFromSource` | 提供元数据拉取/写入任务；区分 `idp_cert` 与 `sp_cert` 凭证类型 |
| D11 | 中（风险） | 绑定/凭证/元数据的 `Query` 使用 `Where(...).Or(...)` 未分组，多租户 scope 的 `AND` 可能被 `OR` 破坏，存在跨租户数据泄漏风险（本次单租户未实测触发） | 三个 `Ram*Service.Query` 实现 | 使用 `Where("(a LIKE ? OR b LIKE ?)")` 或 `db.Group` 包裹 |
| D12 | 低 | `Create`（区别于 `CreateUpdate`）不校验 `code` 格式与唯一性，provider/source 可创建重复 code | 各 `Ram*Service.Create` | 统一走 `CreateUpdate` 或补 code 校验 |
| D13 | 低 | CAS `format` 参数被读取但未实现；无 PGT/代理票据能力；`service` 为空时被 `binding:"required"` 先拦下，提示通用“操作失败”而非服务层文案 | `AuthCas.go`、`CasService.go`、`modRamCas.CasLoginCt` | 实现 format 或移除参数；调整 binding tag/提示 |
| D14 | 信息 | `existName` 成功时“可以使用”放在 `message`，`data` 为空字符串（`rg.OkMessage`）；前端需按 message 判断 | `Ram*Service.ExistName`；用例 `existName` | 统一成功响应契约（data 返回布尔或文案） |

## 6. 复现

```bash
export GOROOT=/Users/mac/develop/go
export PATH=$PATH:$GOROOT/bin
export GOCACHE=/tmp/dsh-go-build
export GOTMPDIR=/tmp
export MOCKEY_CHECK_GCFLAGS=false

go test ./test/idp -run TestIDP -v -count=1
IDP_TEST_SOFT_DELETE=true go test ./test/idp -run TestIDP -v -count=1

# 实机冒烟
go run . &
IDP_LIVE_BASE=http://127.0.0.1:9981 go test ./test/idp -run TestIDP/live-server-smoke -v -count=1
```

## 7. 产物

| 文件 | 说明 |
| --- | --- |
| `test/idp/reports/idp-e2e-physical-delete.log` | 默认（物理删除）模式的用例级过程日志 |
| `test/idp/reports/idp-e2e-soft-delete.log` | 逻辑删除模式的用例级过程日志 |
| `test/idp/reports/idp-e2e-live-smoke.log` | 实机冒烟日志 |
| `test/idp/` | 测试程序 |
| `data/templates/test/idp/authorize.html` | 模拟 IdP 授权页模板 |
| `pkg/middleware/runnerPg/table/a-init-table.go` | D1 修复：补充 IdP 及审计/MFA 建表实体 |
