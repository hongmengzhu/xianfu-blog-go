# 书签接口（/api/bookmark）

> 源码：[Bookmark.go](../../app/domain/web/api/controller/Bookmark.go)、[BookmarkService.go](../../app/domain/web/api/service/BookmarkService.go)、[modApiBlogBookmark](../../app/models/blog/modBlogBookmark/modApiBlogBookmark/)
>
> 通用约定（基础地址、鉴权、响应结构、类型约定）见 [readme.md](./readme.md)。

## 路由一览

| 方法 | 路径 | Go 方法 | 说明 |
| --- | --- | --- | --- |
| POST | `/api/bookmark/getAll` | `GetAll` | 汇总查询我的 / 团队 / 组织的书签与分类 |
| POST | `/api/bookmark/getMy` | `GetMy` | 查询「我的」书签列表 |
| POST | `/api/bookmark/getOrg` | `GetOrg` | 查询「组织」书签列表 |
| POST | `/api/bookmark/create` | `Create` | 创建书签 |
| POST | `/api/bookmark/update` | `Update` | 更新书签 |
| POST | `/api/bookmark/delete` | `Delete` | 删除书签（按配置物理删除或切换状态） |

书签对象 `Vo`（`modApiBlogBookmark.Vo`）字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | string | 主键 |
| `no` | string | 编号代号 |
| `name` | string | 名称 |
| `state` | number | 状态，见类型约定 |
| `createAt` | string\|null | 创建时间 |
| `categoryNo` | string | 分类编号（对应分类的 `no`） |
| `categoryName` | string | 分类名称（当前实现未填充，恒为 `""`） |
| `teamNo` | string | 团队编号 |
| `urlSource` | string | 来源地址 |
| `operationTime` | string\|null | 操作时间 |
| `pageTitle` | string | 网页标题（本接口不写入，由其他模块维护） |

---

## 1. 汇总查询

`POST /api/bookmark/getAll`

返回当前租户下的书签与分类：`my` / `team` / `org` 三类书签，`MyCategory` / `teamCategory` / `orgCategory` 三类分类。无请求参数（不读取请求体）。

### 响应 `data`（`modApiBlogBookmark.VoAll`）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `my` | Vo[] | `typeOwner=my` 且 `ano=当前账号` 的书签 |
| `MyCategory` | Vo[] | 我的分类（注意：键名首字母大写，标签为 `json:"MyCategory"`） |
| `team` | Vo[] | `typeOwner=team` 的书签（不按 `ano` 过滤） |
| `teamCategory` | Vo[] | 团队分类 |
| `org` | Vo[] | `typeOwner=org` 的书签（不按 `ano` 过滤）；已初始化为空切片，无数据时为 `[]` |
| `orgCategory` | Vo[] | 组织分类（取自分类服务 `getAll` 的 `org` 分组）；无数据时为 `[]` |

### 示例

```json
{
  "code": "200",
  "message": "操作成功",
  "data": {
    "my": [
      {
        "id": "1",
        "no": "",
        "name": "Go 官方文档",
        "state": 1,
        "createAt": "2026-10-27T15:57:12+08:00",
        "categoryNo": "",
        "categoryName": "",
        "teamNo": "",
        "urlSource": "https://go.dev/doc/",
        "operationTime": null,
        "pageTitle": ""
      }
    ],
    "MyCategory": [],
    "team": [],
    "teamCategory": [],
    "org": [],
    "orgCategory": []
  }
}
```

---

## 2. 我的书签

`POST /api/bookmark/getMy`

无请求参数。查询条件：`state=1`、`tenantNo=当前租户`、`ano=当前账号`、`typeOwner=my`。

### 响应 `data`

`Vo[]`，无数据时为 `[]`。

### 示例

```json
{
  "code": "200",
  "message": "操作成功",
  "data": [
    {
      "id": "1",
      "no": "",
      "name": "Go 官方文档",
      "state": 1,
      "createAt": "2026-10-27T15:57:12+08:00",
      "categoryNo": "",
      "categoryName": "",
      "teamNo": "",
      "urlSource": "https://go.dev/doc/",
      "operationTime": null,
      "pageTitle": ""
    }
  ]
}
```

---

## 3. 组织书签

`POST /api/bookmark/getOrg`

无请求参数。查询条件：`state=1`、`tenantNo=当前租户`、`ano=当前账号`、`typeOwner=org`。

> 与 `getAll` 中的 `org` 不同，本接口额外按 `ano` 过滤。

### 响应 `data`

`Vo[]`，无数据时为 `[]`。示例同 [2. 我的书签](#2-我的书签)。

---

## 4. 创建书签

`POST /api/bookmark/create`

### 请求体（`modApiBlogBookmark.CreateUpdateCt`）

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | string\|number | 否 | 创建时建议不传；传入会作为主键写入 |
| `name` | string | 是 | 名称，长度 1–255；空值由服务层拦截（字段级校验当前不生效，见「实现备注」第 8 条） |
| `description` | string | 否 | 描述 |
| `content` | string | 否 | 内容（text） |
| `urlSource` | string | 是 | 来源地址，服务层校验非空；数据库上限 2000 字符 |
| `publishTime` | string\|null | 否 | 发布时间（RFC3339）；不传或 `null` 时服务端置为当前时间 |
| `typeReading` | string | 否 | 阅读类型：`未看` / `在看` / `已看`；不传时数据库默认为 `no` |
| `remark` | string | 否 | 备注 |
| `tags` | string[] | 否 | 标签（见「实现备注」第 5 条） |
| `categoryNo` | string | 否 | 分类编号（对应分类的 `no`） |
| `orgNo` | string | 否 | 组织编号 |

字段级校验当前不生效（`validate` 标签与 gin 校验器使用的 `binding` 标签不一致），`name` 为空时由服务层返回 `名称不能为空`。

### 响应 `data`

新建书签的 `id`（字符串）。新建记录不会生成 `no`（保持空），`state` 取数据库默认值 `1`。

```json
{"code": "200", "message": "操作成功", "data": "12"}
```

### 请求示例

```json
{
  "name": "Gin 文档",
  "description": "Go Web 框架",
  "content": "",
  "urlSource": "https://gin-gonic.com/zh-cn/docs/",
  "publishTime": "2026-10-27T10:00:00+08:00",
  "typeReading": "在看",
  "remark": "",
  "tags": ["go", "web"],
  "categoryNo": "",
  "orgNo": ""
}
```

### 业务错误

| `message` | 场景 |
| --- | --- |
| `名称不能为空` | `name` 为空（字段级校验未生效，由服务层返回） |
| `来源地址不能为空` | `urlSource` 为空 |
| `保存失败 <err>` | 数据库写入失败 |

---

## 5. 更新书签

`POST /api/bookmark/update`

### 请求体

同 [4. 创建书签](#4-创建书签) 的请求结构，`id` 必填且 `> 0`。

### 响应 `data`

被更新书签的 `id`（字符串）。

```json
{"code": "200", "message": "操作成功", "data": "12"}
```

### 业务错误

| `message` | 场景 |
| --- | --- |
| `ID不能为空` | `id` 缺失或为 `0` |
| `数据不存在` | `id` 对应的书签不存在 |
| `名称不能为空` | 服务层兜底校验 |
| `来源地址不能为空` | `urlSource` 为空 |
| `保存失败 <err>` | 数据库写入失败 |

> 更新使用 GORM 的 `Updates`（仅更新非零值字段）：传空字符串不会清空已有内容；`publishTime` 省略或传 `null` 会被重置为当前时间。

---

## 6. 删除书签

`POST /api/bookmark/delete`

### 请求体（`model.BaseIdsCt[string]`）

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `ids` | string[] | 是 | 主键集合；缺失、`null` 或 `[]` 均由服务层返回 `id错误`（字段级校验当前不生效） |
| `field` | string | 否 | 保留字段，当前未使用 |
| `extend` | object | 否 | 保留字段，当前未使用 |

```json
{"ids": ["12", "13"]}
```

### 响应 `data`

空字符串 `""`。

```json
{"code": "200", "message": "操作成功", "data": ""}
```

### 行为说明

- 按主键 `id in (?)` 查找，任一命中即视为成功；未命中返回 `数据不存在`
- `ids` 缺失、为 `null` 或 `[]` 时都会进入服务层并返回 `id错误`（`required` 校验未生效）
- 删除方式由 `data/config/app.toml` 的 `[pg.data] delete` 决定：
  - `true`（当前配置）→ 物理删除
  - `false` → 逻辑删除，`1↔11`、`2↔12` 状态互转

### 业务错误

| `message` | 场景 |
| --- | --- |
| `id错误` | `ids` 为空数组 |
| `数据不存在` | `ids` 均未命中 |

---

## 实现备注（按当前代码）

1. `getAll` 的 `MyCategory` 键名首字母大写（`json:"MyCategory"`），其余为小写，前端取值时注意大小写。
2. `getAll` 的 `org`、`orgCategory` 均已显式初始化为空切片，无数据时返回 `[]`；`orgCategory` 的内容取自分类服务 `GetAll` 的 `org` 分组。
3. `Vo.categoryName` 在服务层未赋值（实体中无该字段），恒为空字符串。
4. 创建时不写入 `tenantNo`、`ano`、`typeOwner`，由数据库默认值兜底（`type_owner='general'`、`tenant_no=''`、`ano=''`）；因此新建书签不会出现在 `getMy` / `getOrg` / `getAll` 的 `my`、`team`、`org` 结果中，需由其他模块补全归属。
5. `tags` 仅在请求结构中定义，实体字段为 `datatypes.JSONType[[]string]`，当前 `copier` 拷贝不会转换该类型，标签实际不会落库。
6. `update` 使用非零值更新，因此无法通过本接口将字段清空。
7. `delete` 仅按主键匹配，未附加租户或账号过滤。
8. 模型字段使用 `validate` 标签，而 gin 默认校验器使用 `binding` 标签，字段级校验不会执行；必填/长度限制实际由服务层兜底。单元测试 `app/domain/web/api/controller/api_test.go` 锁定了这一现状。
