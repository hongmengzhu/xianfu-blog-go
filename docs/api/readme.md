# API 接口文档

本目录存放 `app/domain/web/api/controller` 下各控制器对外暴露的接口文档，内容依据控制器、服务层与请求/响应模型源码整理。

最后同步：2026-10-27

| 模块 | 路由前缀 | 接口数 | 文档 |
| --- | --- | --- | --- |
| 书签 | `/api/bookmark` | 6 | [bookmark.md](./bookmark.md) |
| 书签分类 | `/api/bookmark-category` | 3 | [bookmark-category.md](./bookmark-category.md) |

## 通用约定

### 基础地址

开发环境配置（`data/config/app.toml`）：

- `[spring.gin.server] addr = ":9981"`
- `[server] domain = "http://localhost:9981"`

即 `http://localhost:9981` + 各接口的绝对路径，无额外路由前缀（见 `main.go` 的 `NewRouterRegister`，由 starter-gin 在根 engine 上注册）。

### 请求

- 方法：本文档所有接口均为 `POST`
- 请求头：`Content-Type: application/json`
- 请求体：JSON 对象；`ShouldBind` 会按 Content-Type 选择绑定器，未使用 `application/json` 时字段不会被正确解析，随后触发参数校验错误
- 对于绑定了请求体的接口，请求体为空或 JSON 语法错误时返回 `code=500`、`message=操作失败`

### 鉴权

两个控制器的路由组都挂载了 `authPg.GroupApiMiddleware`（`app/middleware/authPg/authApi.go`），必须携带：

```
Authorization: <key>:<hash>
```

- `key` / `secret` 来自接口密钥表 `api_dipl_access_key`（`state=1` 且 `expiry_date >= 当前时间`），服务启动时载入缓存
- `hash` 的算法（`cacheDiplPg.HashSha`）：

```
hash = SHA256( SHA512(key + "::" + secret) + "::" + secret )
```

输出为十六进制小写字符串，示例（非真实密钥）：

```
Authorization: 85227155478337338126:3f1c9a0e...e7
```

- 校验流程：按 `:` 拆分为 `key` 与 `hash`，用缓存中的 `key/secret` 重算并比对
- 缺失请求头、格式错误或校验失败，统一返回：

```json
{"code": "500", "message": "Authorization 参数不能为空", "data": ""}
```

- 鉴权通过后建立调用者上下文：`tenantNo`（租户编号）与 `ano`（操作员账号）均取自该密钥行

### 查询过滤

服务层使用 GORM 结构体条件查询（`Where(结构体)`），**零值字段不参与 WHERE 条件**（空字符串、`0` 会被忽略）。因此：

- 密钥行的 `tenant_no` 为空时，租户过滤条件被省略
- 密钥行的 `ano` 为空时，账号过滤条件被省略（如 `getMy` 会返回该租户下所有 `typeOwner=my` 的记录）

### 响应结构

所有接口统一返回 HTTP `200`，业务结果由响应体的 `code` 表示。外层结构为 `rg.Rs[T]`：

```json
{
  "code": "200",
  "message": "操作成功",
  "messageData": null,
  "data": {}
}
```

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `code` | string | `200` 成功；`415` 参数非法；`500` 失败 |
| `message` | string | 提示信息，成功固定为 `操作成功` |
| `messageData` | any | 当前不会出现：模型使用 `validate` 标签，而 gin 默认校验器识别 `binding` 标签，字段级校验不生效（见下文） |
| `data` | any | 业务数据；失败时为返回类型的零值（`""` 或 `null`） |
| `encryption` / `type` / `timestamp` / `extend` | any | 预留字段，当前不使用，正常响应中省略 |

### 错误形态

字段级校验当前不会触发：模型结构体使用 `validate:"..."` 标签，而 gin 默认校验器（`binding.Validator`）识别的是 `binding` 标签，两者不一致，规则不会执行，`messageData` 也就不会出现。

因此所有错误统一由服务层通过 `message` 返回，例如：

```json
{"code": "500", "message": "名称不能为空", "data": ""}
```

> 若后续在 `app/middleware/validatorPg` 初始化时补一行 `validate.SetTagName("validate")`，字段级提示才会以 `messageData` 形式出现；届时本文档与单元测试需同步更新。

### 类型约定

| 约定 | 说明 |
| --- | --- |
| 64 位整数 | 请求中的 `id` 可传 `"123"` 或 `123`；响应中的 `id` 一律为字符串 |
| 时间 | RFC3339 格式，如 `"2026-10-27T16:07:47+08:00"`；无值时为 `null` |
| `state` | `1` 有效、`2` 停用、`11` 取消（对应有效）、`12` 弃置（对应停用）、`13` 批量删除 |
| `typeOwner` | `my` 我的、`team` 团队、`org` 组织、`general` 普通（见 `pkg/enum/blog/bookmarkTypeOwnerPg`） |
| 空集合 | 服务层已显式初始化的集合返回 `[]`；未初始化的集合为 `null`。当前 9 个接口的输出集合均已显式初始化，实际均返回 `[]` |

### 单元测试

接口的单测位于 `app/domain/web/api/controller/api_test.go`，用 go-spring 容器装配控制器/服务/仓库，数据源指向临时 sqlite 库，并用 `httptest` 驱动真实路由与鉴权中间件，覆盖 9 个接口的成功与错误分支（含空 `ids`、不存在的 id、分类父子链校验等）。

```bash
go test ./app/domain/web/api/controller/ -v
```

> 测试不依赖本地 PostgreSQL 与运行中的服务：sqlite 文件库建在 `t.TempDir()` 下，Authorization 密钥在测试内写入 `cacheDiplPg`。若 `gin` 默认校验器改为识别 `validate` 标签，`create-name-required` 等子测试的断言需要同步更新（见上文说明）。
