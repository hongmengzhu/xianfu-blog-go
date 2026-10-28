# 书签分类接口（/api/bookmark-category）

> 源码：[BookmarkCategory.go](../../app/domain/web/api/controller/BookmarkCategory.go)、[BookmarkCategoryService.go](../../app/domain/web/api/service/BookmarkCategoryService.go)、[modApiBlogBookmarkCategory](../../app/models/blog/modBlogBookmarkCategory/modApiBlogBookmarkCategory/)
>
> 通用约定（基础地址、鉴权、响应结构、类型约定）见 [readme.md](./readme.md)。

## 路由一览

| 方法 | 路径 | Go 方法 | 说明 |
| --- | --- | --- | --- |
| POST | `/api/bookmark-category/getAll` | `GetAll` | 查询我的 / 团队 / 组织的分类 |
| POST | `/api/bookmark-category/createUpdate` | `CreateUpdate` | 创建或更新分类（按 `id` 区分） |
| POST | `/api/bookmark-category/delete` | `Delete` | 删除分类（按配置物理删除或切换状态） |

分类对象 `Vo`（`modApiBlogBookmarkCategory.Vo`）字段：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | string | 主键 |
| `name` | string | 名称 |
| `nameFl` | string | 名称外文 |
| `no` | string | 编号代号（创建时自动生成） |
| `code` | string | 标志（创建时自动生成 22 位数字串） |
| `nameFull` | string | 全称 |
| `parentNo` | string | 上级分类的 `no` |

> 响应中不包含 `description`、`state`、`createAt` 等字段。

---

## 1. 查询全部分类

`POST /api/bookmark-category/getAll`

无请求参数。按 `state=1` + 当前租户查询，分别得到：

- `my`：`typeOwner=my` 且 `ano=当前账号`
- `team`：`typeOwner=team`（不按 `ano` 过滤）
- `org`：`typeOwner=org`（不按 `ano` 过滤）

### 响应 `data`（`modApiBlogBookmarkCategory.VoAll`）

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `my` | Vo[] | 我的分类 |
| `team` | Vo[] | 团队分类 |
| `org` | Vo[] | 组织分类；已初始化为空切片，无数据时为 `[]` |

### 示例

```json
{
  "code": "200",
  "message": "操作成功",
  "data": {
    "my": [
      {
        "id": "3",
        "name": "后端",
        "nameFl": "",
        "no": "2610271557121234567890",
        "code": "2610271557121230987654",
        "nameFull": "",
        "parentNo": ""
      }
    ],
    "team": [],
    "org": []
  }
}
```

---

## 2. 创建或更新分类

`POST /api/bookmark-category/createUpdate`

### 请求体（`modApiBlogBookmarkCategory.CreateUpdateCt`）

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `id` | string\|number | 否 | `> 0` 表示更新，缺失或 `0` 表示创建 |
| `name` | string | 是 | 名称，长度 1–255 |
| `description` | string | 否 | 描述 |
| `parentNo` | string | 否 | 上级分类的 `no`（不是 id）；空表示根节点 |

> 请求体不包含 `code`、`nameFl`、`nameFull`，这三个字段无法通过本接口写入。

### 创建（`id` 缺失或为 0）

- 名称必填；`parentNo` 非空时校验上级是否存在
- `code` 未提供 → 自动生成 22 位数字串；`no` 自动生成
- `tenantNo` 写入 `"-1"`：此处取的是后台登录上下文（`holderPg` 读取 `AUTH_LOGIN`），API 请求只设置了 `AUTH_LOGIN_API`，取不到调用者租户，`GetTenantNo()` 返回兜底值 `"-1"`；`state` 由数据库默认值为 `1`
- 重算 `idLink` / `noLink`

响应 `data` 为新建分类的 `id`（字符串）：

```json
{"code": "200", "message": "操作成功", "data": "3"}
```

### 更新（`id` > 0）

- 分类必须存在；名称必填
- `parentNo` 非空时：上级必须存在、不能等于自己、且不能把分类挂到自己的子孙节点下
- 新上级非空且与原上级不同时，同步重算子孙节点的 `idLink` / `noLink`（移动到根节点时不重算子孙）
- `no` 保持不变；`code` 不被更新

响应 `data` 为空字符串（与创建不同）：

```json
{"code": "200", "message": "操作成功", "data": ""}
```

### 请求示例

```json
{
  "name": "后端",
  "description": "服务端相关",
  "parentNo": ""
}
```

更新时：

```json
{
  "id": "3",
  "name": "后端开发",
  "description": "服务端相关",
  "parentNo": "2610271557121234567890"
}
```

### 业务错误

| `message` | 场景 |
| --- | --- |
| `名称不能为空` | `name` 为空（字段级校验未生效，由服务层返回） |
| `数据不存在` | 更新时 `id` 对应的分类不存在 |
| `上级不存在` | `parentNo` 无对应记录 |
| `上级不能等于自己` | 更新时 `parentNo` 指向自身 |
| `无法保存，不能设置为自己的子集` | 更新时把上级设为自己的子孙节点 |
| `保存失败 <err>` | 数据库写入失败 |

> 服务层还保留了 `标志已存在`、`标志格式不能为空` 两条校验，但请求结构无 `code` 字段（创建时固定走自动生成），当前不会触发。

---

## 3. 删除分类

`POST /api/bookmark-category/delete`

### 请求体（`model.BaseIdsCt[string]`）

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `ids` | string[] | 是 | 主键集合；缺失、`null` 或 `[]` 均由服务层返回 `id错误`（字段级校验当前不生效） |
| `field` | string | 否 | 保留字段，当前未使用 |
| `extend` | object | 否 | 保留字段，当前未使用 |

```json
{"ids": ["3"]}
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
- 删除父分类不会级联删除子孙，子孙的 `parentNo` / `idLink` 不会被清理

### 业务错误

| `message` | 场景 |
| --- | --- |
| `id错误` | `ids` 为空（缺失、`null` 或 `[]`） |
| `数据不存在` | `ids` 均未命中 |

---

## 实现备注（按当前代码）

1. 通过本接口只能创建/更新 `name`、`description`、`parentNo`；`nameFl`、`nameFull` 需由管理端写入，`code` 在创建时自动生成、更新时不修改。
2. 创建成功返回 `id`，更新成功返回空字符串 `""`，前端需按分支处理。
3. `getAll` 的 `org` 已显式初始化为空切片，无数据时返回 `[]`。
4. 创建时 `tenantNo` 实际写入 `"-1"`（读取的是后台登录上下文 `AUTH_LOGIN`，API 请求未设置，`GetTenantNo()` 返回兜底值），且未写入 `ano` 与 `typeOwner`（数据库默认 `type_owner='general'`、`ano=''`）；因此新建分类不会出现在 `my` / `team` / `org` 任一列表中，需由其他模块补全归属。
5. 更新使用 GORM 的 `Updates`（仅更新非零值字段），空字符串不会清空已有字段；`no` 会被显式置空后跳过更新，保持不变。
6. `parentNo` 指向的是上级分类的 `no` 而非 `id`。
7. 删除仅按主键匹配，未附加租户过滤。
8. 模型字段使用 `validate` 标签，而 gin 默认校验器使用 `binding` 标签，字段级校验不会执行；必填限制由服务层兜底。单元测试 `app/domain/web/api/controller/api_test.go` 覆盖了这一现状。
