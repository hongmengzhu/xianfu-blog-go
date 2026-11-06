# 基础配置（BasicConfig）设计说明

> 源码：
>
> - 服务：[BasicConfigService.go](../../app/domain/domainBasic/service/BasicConfigService.go)、[BasicConfigCategoryService.go](../../app/domain/domainBasic/service/BasicConfigCategoryService.go)、[BasicConfigModelService.go](../../app/domain/domainBasic/service/BasicConfigModelService.go)、[configModel](../../app/domain/domainBasic/service/configModel/)、[configCache](../../app/domain/domainBasic/service/configCache/)
> - 控制器：[Config.go](../../app/handler/manage/b72/Config.go)、[ConfigCategory.go](../../app/handler/manage/b72/ConfigCategory.go)、[ConfigModel.go](../../app/handler/manage/b72/ConfigModel.go)
> - 实体：[entityBasic](../../pkg/infrastructure/entityBasic/)、模型请求/响应：[modBasicConfig](../../app/models/basic/modBasicConfig/)、[modBasicConfigCategory](../../app/models/basic/modBasicConfigCategory/)、[modBasicConfigModel](../../app/models/basic/modBasicConfigModel/)
> - 缓存 DTO：[modCacheBasicConfig](../../app/models/cache/basic/modCacheBasicConfig/config.go)
> - 初始化建表：[a-init-table.go](../../pkg/middleware/runnerPg/table/a-init-table.go)

## 1. 目标与整体流程

设计器页面负责“拖控件 → 选模型分类（配置/表）→ 填字段与校验规则 → 保存”。后端一次请求完成：

1. 保存模型头 `basic_config_model`；
2. 保存字段 `basic_config_model_fields`；
3. 保存字段规则 `basic_model_rules`、模型规则 `basic_config_model_rules`；
4. 按字段自动生成/同步【配置】行 `basic_config`，并写入模型分类；
5. 模型种类为【配置模型】`config` 时，同步一条【配置列表】行 `basic_config_list`；
6. 按 `tenant_no + model` 刷新 bigcache。

因此设计器保存后，字段已经自动出现在【配置】列表里，不需要再手工新增配置行。配置值可以在【配置】页面或 `updateValues` 接口按模型批量修改。

### 1.0 配置列表同步（模型种类 = 配置模型）

`header.modelCategory = config` 时，保存完成后按模型维度维护 `basic_config_list`：

- 一个模型对应一行：`model_no = 模型 No`，`event_no = field = field_path = 模型标识`；
- `name` / `description` / `show` / `sort` 取模型头，`module` / `module_sub` / `tenant_no` 取模型归属；
- `type_domain` 默认 `tenant`，模型 `type_sys = system` 时为 `system`；
- `kind_unique = md5(tenant_no + model_no + model)`；
- 重复保存只更新上述字段，保留 `state`、`content` 等运行期数据；
- 模型标识改名时，同一 `model_no` 的行会同步更新 `event_no` / `field` / `field_path`；
- `table`（库表模型）不同步，避免影响种子数据里已有的配置列表行；从 `config` 改为 `table` 也不会自动删除已生成的行，需要时在【配置列表】手工删除。

### 1.1 表与实体

| 表 | 实体 | 作用 |
| --- | --- | --- |
| `basic_config` | `BasicConfigEntity` | 配置行：一个模型字段对应一行，保存字段的值、默认值、展示属性 |
| `basic_config_category` | `BasicConfigCategoryEntity` | 配置分类（树形，供【配置】左侧分类/选项卡使用） |
| `basic_config_list` | `BasicConfigListEntity` | 配置列表：配置模型对应一行，供 `config-list` 的 `detailForm` / `configUpdate` 使用 |
| `basic_config_model` | `BasicConfigModelEntity` | 模型头，`model` 在租户内唯一 |
| `basic_config_model_fields` | `BasicConfigModelFieldsEntity` | 模型字段（表单控件、字段路径、控件属性等） |
| `basic_config_model_rules` | `BasicConfigModelRulesEntity` | 模型级规则 |
| `basic_model_rules` | `BasicModelRulesEntity` | 字段级规则（`value_no=模型 No`、`field_no=字段 No`） |

`basic_config` 不存 `value_attr`；该字段只存在于 `basic_config_model_fields`，配置行保留 `value_type` 和值。

### 1.2 唯一性

- 模型：`model` 在同一 `tenant_no` 下唯一。
  - 新增校验 `FindByTenantNoAndModel`；
  - 更新校验 `FindByTenantNoAndModelAndIdNot`（排除自身）；
  - `basic_config_model.kind_unique = md5(tenant_no + model)`。
- 配置行：同一模型下同一字段唯一，服务层按 `tenant_no + model_no + field` 查重，`kind_unique = md5(tenant_no + model_no + field)`。
- 配置列表行：配置模型按 `tenant_no + model_no` 唯一，`kind_unique = md5(tenant_no + model_no + model)`。
- 模型 `No` 由 `noPg.No()` 生成；对外查询/关联优先使用 `model_no`（即 `No`），用户可读标识使用 `model`。

## 2. 设计器保存接口

`POST /xianfu/manage/b72/config-model/createUpdate`

请求体 `modBasicConfigModel.CreateUpdateCt`：

```json
{
  "header": {
    "id": "",
    "name": "系统配置",
    "model": "systemConfig",
    "table": "systemConfig",
    "typeSys": "general",
    "modelCategory": "config",
    "moduleSub": "system",
    "description": "系统配置模型",
    "category": "config",
    "categoryTab": "basic",
    "tags": ["system"],
    "rules": ["required"],
    "ruleItems": [
      {
        "name": "模型必填规则",
        "ruleMode": "required",
        "condition": "",
        "errorMessage": "模型不能为空",
        "structure": "",
        "ruleTarget": ["frontend", "backend"],
        "show": 1,
        "sort": "0",
        "extraData": null
      }
    ]
  },
  "body": [
    {
      "id": "",
      "no": "",
      "name": "站点名称",
      "field": "siteName",
      "fieldPath": "siteName",
      "description": "站点展示名称",
      "show": 1,
      "independent": 2,
      "binary": 2,
      "defaultValue": "",
      "value": "",
      "valueType": "string",
      "valueAttr": "single",
      "formCode": "input",
      "formAttr": {"maxlength": 100},
      "formGroup": "basic",
      "sharedScope": [],
      "sort": "0",
      "rules": ["required", "length"],
      "ruleItems": [
        {
          "name": "站点名称必填",
          "ruleMode": "required",
          "condition": "",
          "errorMessage": "站点名称不能为空",
          "structure": "",
          "ruleTarget": ["frontend", "backend"],
          "show": 1,
          "sort": "0",
          "extraData": null
        },
        {
          "name": "长度 1-100",
          "ruleMode": "length",
          "condition": "1,100",
          "errorMessage": "长度必须在 1-100 之间",
          "structure": "",
          "ruleTarget": ["frontend", "backend"],
          "show": 1,
          "sort": "1",
          "extraData": null
        }
      ]
    }
  ],
  "bodyDelIds": []
}
```

| 字段 | 说明 |
| --- | --- |
| `header.id` | 空表示新增；非空表示更新 |
| `header.model` | 租户内唯一英文标识 |
| `header.modelCategory` | `config`=配置模型，`table`=库表模型；创建时必填且必须存在 |
| `header.moduleSub` | 子模块 `No`，必须存在 |
| `header.table` | 表名，必填并做格式校验 |
| `header.category` | 配置分类编号，写入自动生成的 `basic_config.category_no`；空时默认 `config` |
| `header.rules` / `header.ruleItems` | 模型级规则，`rules` 是 `ruleMode` 简写；完整规则写 `ruleItems` |
| `body[].rules` / `body[].ruleItems` | 字段级规则；完整规则写 `ruleItems` |
| `bodyDelIds` | 更新时要删除的字段 `id`，同时删除对应配置行 |

保存后的副作用：

- 字段级规则写入 `basic_model_rules`：`value_no = 模型 No`、`field_no = 字段 No`、`type_model = "model"`。
- 模型级规则写入 `basic_config_model_rules`：`model_no = 模型 No`、`type_category = "model"`。
- 每个字段生成一条 `basic_config`：
  - 已有配置行保留 `value` / `value_binary`，只同步名称、展示、默认值、排序、分类等；
  - 新字段使用字段默认值作为初始值；
  - 字段被删除时删除对应配置行；
  - `category_no` 优先取 `header.category`，为空时取 `config`。
- 模型种类为 `config` 时同步一条 `basic_config_list`（规则见 1.0）；`table` 不同步。
- 最后刷新 `tenant_no + model` 缓存。

`Detail` 会把这些规则重新组装成 `header.rules/ruleItems` 和 `body[].rules/ruleItems`，设计器可直接回显。

## 3. 模型 CRUD 路由

基础路径：`/xianfu/manage/b72/config-model`，统一使用管理端鉴权。

| 方法 | 路径 | 服务方法 | 说明 |
| --- | --- | --- | --- |
| POST | `/createUpdate` | `BasicConfigModelService.CreateUpdate` | 设计器保存（新增/更新模型、字段、规则、配置） |
| GET | `/detail/:id` | `Detail` | 模型详情（含字段与规则） |
| POST | `/query` | `Query` | 分页查询 |
| POST | `/enable` | `Enable` | 启用 |
| POST | `/disable` | `Disable` | 停用 |
| POST | `/state` | `StateEnableDisable` | 设置有效/停用 |
| POST | `/delete` | `LogicalDeletion` | 逻辑删除；配置采用物理删除时级联删除字段/规则/配置 |
| POST | `/recovery` | `LogicalRecovery` | 逻辑删除恢复 |
| POST | `/physicalDeletion` | `PhysicalDeletion` | 物理删除，并级联删除字段、字段规则、模型规则、配置行、缓存 |
| POST | `/selectNodeAllPublic` | `SelectNodeAllPublic` | 模型选择树 |
| POST | `/existName` | `ExistName` | 名称查重 |
| POST | `/existCode` | `ExistCode` | 编码查重 |
| POST | `/existModel` | `ExistModel` | 模型标识查重 |

## 4. 配置 CRUD 与查询接口

基础路径：`/xianfu/manage/b72/config`。

| 方法 | 路径 | 服务方法 | 说明 |
| --- | --- | --- | --- |
| POST | `/createUpdate` | `CreateUpdate` | 新增/更新单条配置（通常由模型保存自动生成，手工维护时使用） |
| GET | `/detail/:id` | `Detail` | 配置详情 |
| POST | `/query` | `Query` | 按租户分页查询，支持 `wd` 模糊匹配名称/字段/描述 |
| POST | `/allByModel` | `AllByModel` | 按 `tenantNo + model` 读取配置（走缓存） |
| POST | `/updateValues` | `UpdateValues` | 按模型批量保存字段值，并按字段规则校验 |
| POST | `/enable` | `Enable` | 启用 |
| POST | `/disable` | `Disable` | 停用 |
| POST | `/state` | `StateEnableDisable` | 设置有效/停用 |
| POST | `/delete` | `LogicalDeletion` | 逻辑删除 |
| POST | `/recovery` | `LogicalRecovery` | 逻辑删除恢复 |
| POST | `/physicalDeletion` | `PhysicalDeletion` | 物理删除 |
| POST | `/cacheOverride` | `CacheOverride` | 从数据库全量重建配置缓存 |
| POST | `/existName` | `ExistName` | 名称查重 |

### 4.1 单条新增/更新

请求体 `modBasicConfig.CreateUpdateCt`：

```json
{
  "id": "",
  "modelNo": "",
  "model": "systemConfig",
  "field": "siteName",
  "name": "站点名称",
  "fieldPath": "siteName",
  "description": "站点展示名称",
  "show": 1,
  "binary": 2,
  "valueType": "string",
  "formCode": "input",
  "defaultValue": "",
  "value": "我的站点",
  "valueBinary": "",
  "categoryNo": "config",
  "sort": "0"
}
```

`modelNo` 与 `model` 二选一，优先 `modelNo`。字段 `field` 必须已存在于该模型，否则返回“模型字段不存在，请先在设计器保存字段”。更新时若 `value` 为空会清空值；新增时为空则使用 `defaultValue`。

### 4.2 按模型批量保存值

`POST /xianfu/manage/b72/config/updateValues`

```json
{
  "model": "systemConfig",
  "data": {
    "siteName": "我的站点",
    "siteUrl": "https://example.com"
  }
}
```

服务会：

1. 按 `tenant_no + model` 找模型；
2. 读取模型字段、配置行、字段规则；
3. 对每个值执行 `required`、`length`、`pattern`、`greaterThan0` 规则；
4. `binary=1` 的字段写 `value_binary`，其余写 `value`；
5. 刷新该模型缓存。

`length` 的 `condition` 支持：

- `"10"`：精确长度 10；
- `"1,100"`：最小 1、最大 100。

### 4.3 按模型读取

`POST /xianfu/manage/b72/config/allByModel`

```json
{
  "model": "systemConfig"
}
```

返回 `modBasicConfig.Vo[]`，字段包含 `field`、`name`、`value`、`valueType`、`show`、`categoryNo` 等。

## 5. 配置分类接口

基础路径：`/xianfu/manage/b72/config-category`。

| 方法 | 路径 | 服务方法 | 说明 |
| --- | --- | --- | --- |
| POST | `/createUpdate` | `Create` / `Update` | `id` 为空新增，否则更新 |
| GET | `/detail/:id` | `Detail` | 分类详情 |
| POST | `/query` | `Query` | 分页查询 |
| POST | `/enable` | `Enable` | 启用 |
| POST | `/disable` | `Disable` | 停用 |
| POST | `/state` | `StateEnableDisable` | 设置有效/停用 |
| POST | `/delete` | `LogicalDeletion` | 逻辑删除 |
| POST | `/recovery` | `LogicalRecovery` | 逻辑删除恢复 |
| POST | `/physicalDeletion` | `PhysicalDeletion` | 物理删除 |
| POST | `/existName` | `ExistName` | 名称查重 |
| POST | `/existCode` | `ExistCode` | 标志查重 |

请求体 `modBasicConfigCategory.CreateUpdateCt`：

```json
{
  "id": "",
  "name": "配置",
  "nameFl": "Config",
  "nameFull": "基础配置",
  "code": "",
  "description": "系统基础配置分类",
  "parentId": "",
  "parentNo": "",
  "sort": "0"
}
```

- `code` 为空时自动生成；非空时按租户校验唯一。
- 分类是树形结构，使用 `parent_no` / `no_link` 维护路径；不能把上级设置成自己的子节点。
- `state` 使用统一状态枚举：`1` 有效、`2` 停用、`11/12` 逻辑删除态、`13` 批量删除态。

## 6. 缓存设计

缓存组件：[configCache](../../app/domain/domainBasic/service/configCache/)，底层为 bigcache。

- Key：`basicConfig:config:<tenantNo>:<model>`
  - 生成方法：`basicConfigKey.ConfigByTenantNoAndModel(tenantNo, model)`
- Value：`[]modCacheBasicConfig.Item` 的 JSON 字符串。
- 二进制字段在 `ValueMap()` 中返回 `value_binary`（当 `binary=1` 且非空），否则返回 `value`。
- 写入使用 `Rdt.Set(ctx, key, payload, 0)`，即不设置单条 TTL；bigcache 自身由全局 `life-window` 管理。
- 缓存未命中时从 `basic_config` 按 `tenant_no + model + state=1` 回源，并回填缓存；空结果也会缓存空数组。
- 变更后刷新：
  - `BasicConfigService.CreateUpdate/UpdateValues/State/删除/恢复`；
  - 设计器 `syncConfig` 保存后；
  - 模型物理删除时清空该模型缓存。
- 全量重建：`POST /xianfu/manage/b72/config/cacheOverride`。

其他服务读取配置时不要直接查表，统一调用：

```go
items, ok := basicConfigService.GetByTenantNoAndModel(ctx, tenantNo, model)
values, ok := basicConfigService.ValuesByTenantNoAndModel(ctx, tenantNo, model)
```

`GetByTenantNoAndModel` 优先读缓存，未命中才回源数据库，因此调用方不需要关心缓存 key。

## 7. 通用响应

所有接口返回 `rg.Rs`：

```json
{
  "code": "200",
  "message": "操作成功",
  "data": {}
}
```

失败时 `code` 非 `200`，`message` 为具体错误。分页接口 `data` 为：

```json
{
  "pageNum": 1,
  "pageSize": 20,
  "total": 1,
  "data": []
}
```

## 8. 初始化与启动

- `a-init-table.go` 已注册：
  - `BasicConfigEntity`
  - `BasicConfigCategoryEntity`
  - `BasicConfigEventEntity`
  - `BasicConfigEventFieldsEntity`
  - `BasicConfigListEntity`
  - `BasicConfigModelEntity`
  - `BasicConfigModelFieldsEntity`
  - `BasicConfigModelRulesEntity`
  - `BasicModelRulesEntity`
- 启动时只为缺失的表建表；已有表不会自动增加列。若旧库缺少新列，需要单独执行迁移 SQL。
- 种子数据中的 `systemConfig` 模型仍为 `table` 分类，已有 `basic_config` 行会保留。

## 9. 删除级联

模型物理删除（或配置为物理删除时的逻辑删除）会级联清理：

1. `basic_config`（按 `model_no`）；
2. `basic_config_list`（按 `model_no`）；
3. `basic_config_model_fields`（按 `model_no`）；
4. `basic_model_rules`（按 `value_no = 模型 No`）；
5. `basic_config_model_rules`（按 `model_no`）；
6. 刷新 `tenant_no + model` 缓存。

删除入口：

- `BasicConfigModelService.PhysicalDeletion`；
- `BasicConfigModelService.LogicalDeletion` 中 `Data.Delete=true` 的分支。
