# RAM 身份提供商（IdP）端到端测试说明

## 1. 测试目标

基于库表、控制器与服务层，对 RAM 身份提供商（IdP）相关能力做一轮端到端验证：

- 6 张 IdP 库表对应的管理端接口：增、删、改、查、启用/禁用、逻辑删除/恢复/物理删除
- 身份提供商与认证源的关联、身份绑定记录
- OAuth2 / OpenID 授权码登录（首次自动建号、二次绑定登录、刷新 token 及失败分支）
- CAS 1.0 / 2.0 / 3.0 的 Service Ticket 签发与校验
- SAML2 登录请求生成与回调异常分支
- LDAP 连接测试、用户搜索、同步、登录

## 2. 被测对象

| 控制器 | 路由前缀 | 库表 | 服务 |
| --- | --- | --- | --- |
| `app/handler/manage/ram/IdentityProvider.go` | `/xianfu/manage/ram/identity-provider` | `ram_identity_provider` | `RamIdentityProviderService` |
| `app/handler/manage/ram/IdentitySource.go` | `/xianfu/manage/ram/identity-source` | `ram_identity_source` | `RamIdentitySourceService` |
| `app/handler/manage/ram/IdentitySourceCallback.go` | `/xianfu/manage/ram/identity-source-callback` | `ram_identity_source_callback` | `RamIdentitySourceCallbackService` |
| `app/handler/manage/ram/IdpBinding.go` | `/xianfu/manage/ram/idp-binding` | `ram_idp_binding` | `RamIdpBindingService` |
| `app/handler/manage/ram/IdpCredential.go` | `/xianfu/manage/ram/idp-credential` | `ram_idp_credential` | `RamIdpCredentialService` |
| `app/handler/manage/ram/IdpMetadataCache.go` | `/xianfu/manage/ram/idp-metadata-cache` | `ram_idp_metadata_cache` | `RamIdpMetadataCacheService` |
| `app/handler/manage/ram/Ldap.go` | `/xianfu/manage/ram/ldap` | 复用认证源/账号/绑定 | `LdapService` |
| `app/handler/manage/ram/Cas.go` | `/xianfu/auth/manage/cas` | 复用账号 | `CasService` |
| `app/handler/manage/ram/authLoginIdp.go` | `/xianfu/auth/manage/idp` | 复用认证源/绑定/账号 | `AccountLoginIdpService`、`SamlLoginService`、`CasService` |

协议 SDK：`pkg/sdk/auth/idpPg/idp`（Custom/OIDC 提供方）、`pkg/sdk/auth/idpPg/cas`、`pkg/sdk/auth/idpPg/saml`、`pkg/sdk/auth/idpPg/ldap`。

## 3. 测试环境

| 项 | 值 |
| --- | --- |
| Go | 1.27.3 darwin/amd64 |
| 数据库 | PostgreSQL `127.0.0.1:5432`，库 `test-make-01`，用户 `postgres` |
| 配置 | `data/config/app.toml`（`pg.data.delete=true`，gin `:9981`，echo `:9985`） |
| 租户 | `1000` |
| 管理账号 | `1000` / `manage`（测试通过数据库内 AccessToken 密钥签发 PASETO，不依赖登录密码） |
| 测试数据前缀 | `E2E-<runID>-*`、`e2e_*` 账号、`e2e-*` 外部身份标识 |

测试不依赖正在运行的服务进程：go-spring 容器在测试进程内启动，注册真实控制器与中间件，`httptest` 驱动真实路由，数据库连接与生产配置一致。

## 4. 测试程序

目录：`test/idp/`

| 文件 | 作用 |
| --- | --- |
| `support_test.go` | 容器装配、建表兜底、PASETO 签发、HTTP 断言辅助、数据清理 |
| `crud_test.go` | 6 个控制器 × 增删改查/状态/删除全链路；路由与库表断言 |
| `idplogin_test.go` | 模拟 OAuth2/OIDC 提供方 + OAuth2/OpenID 登录与刷新 token |
| `cas_test.go` | CAS 票据签发、1.0/2.0/3.0 校验、票据复用与 service 不匹配 |
| `saml_test.go` | 自签证书、SAML AuthnRequest 重定向、回调异常分支 |
| `ldap_test.go` | LDAP 连接测试/搜索/同步/登录的失败路径与参数校验 |
| `reports/` | 三次执行的过程结果日志 |

测试模板：`data/templates/test/idp/authorize.html`（模拟 IdP 授权页，由 `idplogin_test.go` 的 httptest 服务渲染）。

只启动一次容器：顶层用例 `TestIDP` 内通过 `t.Run` 组织全部子场景，避免重复装配与重复建连。

## 5. 运行方式

```bash
export GOROOT=/Users/mac/develop/go
export PATH=$PATH:$GOROOT/bin
export GOCACHE=/tmp/dsh-go-build
export GOTMPDIR=/tmp
export MOCKEY_CHECK_GCFLAGS=false

# 默认模式：沿用 app.toml，pg.data.delete=true（delete 即物理删除）
go test ./test/idp -run TestIDP -v -count=1

# 逻辑删除模式：基于 app.toml 生成 delete=false 的临时配置目录
IDP_TEST_SOFT_DELETE=true go test ./test/idp -run TestIDP -v -count=1

# 实机冒烟：需先用 go run . 启动服务（gin :9981）
IDP_LIVE_BASE=http://127.0.0.1:9981 \
  go test ./test/idp -run TestIDP/live-server-smoke -v -count=1
```

可选环境变量：

| 变量 | 默认 | 说明 |
| --- | --- | --- |
| `IDP_TENANT_NO` | `1000` | 签发 token 与数据过滤使用的租户 |
| `IDP_TEST_SOFT_DELETE` | `false` | `true` 时走逻辑删除（state 1→11）与恢复 |
| `IDP_LIVE_BASE` | 空 | 设置后执行真实进程冒烟，空则跳过 |

## 6. 数据准备与清理

- 每次运行生成唯一 `runID`（`MMDD-HHMMSS`），名称前缀 `E2E-`，避免与既有数据冲突。
- 创建记录时登记 `(table, id)`，运行结束在容器关闭前统一 `DELETE`；自动登录创建的账号按 `bind_ano` 清理，审计日志按 `ano` 清理。
- IdP 登录会真实创建账号与绑定（用户要求允许在当前程序数据上增删），但只创建测试前缀账号，并在结束时删除。
- 清理已实测：provider/source/callback/binding/account/sessionlog/credential/metadata 残留均为 0。

## 7. 覆盖矩阵

### 7.1 管理端接口

| 能力 | 身份提供商 | 认证源 | 回调白名单 | 身份绑定 | 凭证 | 元数据缓存 |
| --- | --- | --- | --- | --- | --- | --- |
| create | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| update（含 id 缺失/不存在） | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| detail（含 0/不存在） | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| query（含 `wd` 过滤） | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| selectPublic（仅启用项） | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| existName（重复/排除自身/空值） | ✓ | ✓ | ✓ | — | — | — |
| enable / disable / state（非法状态、空 ids） | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| delete / recovery / physicalDeletion | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| 鉴权缺失 | ✓ | — | — | — | — | — |

删除语义说明：`pg.data.delete=true` 时 `delete` 直接物理删除，`recovery` 必然返回“数据不存在”；`IDP_TEST_SOFT_DELETE=true` 时验证 `delete→state=11`、`recovery→state=1`、`physicalDeletion` 再物理删除。

### 7.2 协议与登录

| 场景 | 覆盖 |
| --- | --- |
| OAuth2 授权码换 token | ✓（mock token endpoint） |
| OIDC/OpenID id_token 解析 | ✓（`sub/name/preferred_username/email`） |
| 首次登录自动创建账号与绑定 | ✓（校验 `isSignup=true`、`state_bind=2`、token 落库） |
| 二次登录命中绑定 | ✓（`isSignup=false`） |
| 未绑定 + `method=link` | ✓（拒绝） |
| 认证源/提供商停用 | ✓（拒绝登录，恢复后可用） |
| 授权码为空/认证源不存在/非法 JSON | ✓ |
| token 端点错误 / userinfo 映射缺失 / id_token 缺 sub | ✓（分别拒绝） |
| 刷新 token / 缺少刷新头 | ✓ |
| 绑定停用后登录 | ✓（当前仍可登录，见缺陷 D9） |
| CAS 签发 ST、1.0 校验通过/复用/service 不匹配/缺参 | ✓ |
| CAS 2.0、3.0 XML 成功与 INVALID_TICKET_SPEC/INVALID_TICKET/INVALID_SERVICE/INVALID_REQUEST | ✓ |
| CAS 登录路由无登录上下文 | ✓（返回“用户未登录”，见缺陷 D5） |
| SAML AuthnRequest 302 重定向 | ✓（自签 IdP 证书经 `sp_cert` 凭证注入） |
| SAML 未知源/协议不符/停用/无证书 | ✓ |
| SAML 回调缺响应/坏 RelayState/非法断言 | ✓ |
| LDAP 连接测试（连接被拒） | ✓ |
| LDAP 搜索/同步/登录连接失败、停用源、缺参 | ✓ |
| LDAP 登录不校验协议 | ✓（非 ldap 源仍走 LDAP 连接，见缺陷 D6） |

## 8. 关键流程说明

### OAuth2 / OIDC

1. 测试进程启动 `httptest` 模拟 IdP：`/oauth/authorize` 渲染模板，`/oauth/token` 返回 `access_token`、`refresh_token`、`id_token`。
2. 通过管理端接口创建 `Code=Custom` 的身份提供商与认证源，`baseConfig` 写入 `clientId/clientSecret/tokenURL/authURL`，`userInfoURL` 置空。
3. 调用 `/xianfu/auth/manage/idp/login`（该路由未挂载鉴权中间件，登录请求本身公开），服务用授权码换 token，再从 `id_token` claims 解析用户并查找/创建绑定。
4. `CustomIdProvider` 在 `UserInfoURL` 为空时回退到 id_token，因此无需 `userMapping`。

### CAS

`CasController` 的 `/xianfu/auth/manage/cas/login` 与 `AuthLoginIdpController` 的 `/xianfu/auth/manage/idp/cas/login` 均未挂载鉴权中间件，`CasService.Login` 依赖 `ctx.Get("ano")`，因此 HTTP 直连必然“用户未登录”。测试通过直接调用控制器方法并注入 `ano=1000` 验证票据签发，再用真实 GET 路由验证 1.0/2.0/3.0 校验。

### SAML

测试用自签 X.509 证书写入 `ram_idp_credential`（`credType=sp_cert`，`source_no` 指向认证源）；`SamlLoginService` 会从该凭证读取 IdP 证书。未覆盖签名断言的成功回调。

### LDAP

使用 `127.0.0.1:1` 保证快速连接失败，验证错误处理与参数校验；未覆盖真实 LDAP 目录的成功登录/同步。

## 9. 已知缺陷与限制

缺陷清单与证据见 [idp-e2e-test-result.md](./idp-e2e-test-result.md#缺陷与风险清单)。

未覆盖项：

- 真实 LDAP 目录的成功搜索/同步/登录
- 带签名断言的成功 SAML 回调（需要完整 IdP 私钥签名链路）
- MFA 二次验证（仅覆盖无 MFA 记录时的直接登录）
- 多租户隔离下的 `Or` 条件越权（仅单租户实测）
- 微信/支付宝/钉钉等需要真实第三方凭据的 provider
- 前端页面与浏览器跳转链路

## 10. 测试结果

- 结果报告：[idp-e2e-test-result.md](./idp-e2e-test-result.md)
- 过程日志：`test/idp/reports/idp-e2e-physical-delete.log`、`test/idp/reports/idp-e2e-soft-delete.log`、`test/idp/reports/idp-e2e-live-smoke.log`
