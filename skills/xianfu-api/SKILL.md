---
name: xianfu-api
description: 当用户需要添加新的后端 API 接口时触发。提供完整的步骤引导，涵盖 Entity、Repository、Command、Query、Service、Controller 的创建，触发词包括"添加接口"、"新增API"、"add api"、"新建接口"、"添加增删改查"、"加一个表对应的接口"。适用于golang语言 + gorm ORM + go-spring 技术栈。
---

# 博客项目后端 API 添加指南

## 核心原则

项目采用 **多领域分层架构**，按权限和用途分为三层：

| 层 | 路径 | 职责 | 中间件 |
|---|------|------|--------|
| `../../app/domain/system/` | 系统内部API | 系统级管理接口 | `authPg.GroupSystemMiddleware` |
| `../../app/domain/manage/` | 管理后台API | 运营管理接口 | `authPg.GroupManageMiddleware` |
| `../../app/domain/web/` | 前台Web API | 面向公众/前端页面 | `authPg.GroupWebMiddleware` 或无 |

每层内部按模块组织子目录（如 `basic/`、`ram/`、`blog/`、`tc/`、`api/`），每个模块统一包含三个子目录：

```
app/{layer}/{module}/
├── controller/    # Controller 层 — 路由注册 + HTTP 处理
├── service/       # Service 层 — 业务逻辑
└── model/         # Model 层 — 请求/响应 DTO，按实体分子目录
    └── modXxx/
        ├── createUpdateCt.go # 创建更新合并请求（create/update）
        ├── queryCt.go        # 列表查询请求
        └── vo.go             # 视图对象（响应）
```

**参数绑定和陆游：**

参数绑定：`routerPg.BindJson(ctx, &ct)` 
路由前缀：`/xianfu/{面向对象}/{module}/{entity}`

**技术栈：** Go + Gin（Web框架）+ GORM（ORM）+ go-spring（依赖注入）+ copier（对象拷贝）+ pagePg（分页）+ rg（响应封装）

## 文件生成顺序

按依赖关系依次创建以下文件。每创建一个文件后确认编译通过再继续。

### 检查清单

```
□ 1. Entity          infrastructure/entityXxx/XxxEntity.go         — 数据库表映射
□ 2. Repository      infrastructure/repositoryXxx/XxxRepository.go — 数据访问
□ 3. Model (Vo)      app/{layer}/{module}/model/modXxx/vo.go       — 响应DTO
□ 4. Model (CreateUpdateCt) app/{layer}/{module}/model/modXxx/createUpdateCt.go — 创建请求/更新请求
□ 5. Model (QueryCt)  app/{layer}/{module}/model/modXxx/queryCt.go  — 查询请求
□ 6. Service         app/{layer}/{module}/service/XxxService.go    — 业务逻辑
□ 7. Controller      app/{layer}/{module}/controller/Xxx.go        — 路由+控制器
```

---

## 步骤详解

### 步骤 1：创建 Entity（数据库实体）

路径：`infrastructure/entityXxx/XxxEntity.go`

```go
package entityXxx

import "time"

// XxxEntity 描述
type XxxEntity struct {
    ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
    No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号"`
    Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称"`
    State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消12弃置13批量删除" json:"state" comment:"状态"`
    Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述"`
    CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at"`
    UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at"`
    CreateBy    string     `gorm:"column:create_by;type:varchar(80);default:;comment:创建人" json:"create_by"`
    UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by"`
    TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no"`
}

func (*XxxEntity) TableName() string {
    return "table_name"
}

func (*XxxEntity) TableComment() string {
    return "表描述"
}
```

**Entity 关键规则：**
- 表名用单数蛇形（如 `basic_tags`）
- 主键 `ID` 类型 `int64`，使用 `bigserial` 自增
- 状态字段 `State` 类型 `int8`，值域：1有效 2停用 11取消 12弃置 13批量删除
- 时间字段用 `*time.Time`，创建时间用 `autoCreateTime`，更新时间用 `autoUpdateTime`
- 多租户字段 `TenantNo` 几乎每个表都有
- GORM tag 格式：`gorm:"column:xxx;type:xxx;comment:xxx"`

### 步骤 2：创建 Repository（数据访问层）

路径：`infrastructure/repositoryXxx/XxxRepository.go`

```go
package repositoryXxx

import (
    "context"
    "reflect"

    "github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityXxx"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
    "go-spring.org/log"
    "go-spring.org/spring/gs"
)

func init() {
    gs.Provide(new(XxxRepository))

    gs.Provide(new(support.BaseService[XxxRepository]))
}

type XxxRepository struct {
    repositoryPg.BaseRepository[entityXxx.XxxEntity, int64]
}
```

**Repository 关键规则：**
- 嵌入泛型 `BaseRepository[Entity, ID类型]`，自动获得 CRUD 方法
- `init()` 中用 `gs.Provide` 注册到 go-spring 容器
- 同时注册 `support.BaseService[XxxRepository]`
- BaseRepository 已提供：`Create`, `Update`, `FindById`, `FindAll`, `FindAllPage`, `DeleteByIds`, `DeleteByIdsString` 等方法

### 步骤 3：创建 Model（请求/响应 DTO）

路径：`app/{layer}/{module}/model/modXxx/`

#### Vo（响应视图）

```go
package modXxx

import (
    "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
    "time"
)

type Vo struct {
    ID          typePg.Uint64String `json:"id" label:"id"`
    Name        string              `json:"name" label:"名称"`
    State       typePg.Int8         `json:"state" label:"状态"`
    Description string              `json:"description" label:"描述"`
    CreateAt    *time.Time          `json:"createAt" label:"创建时间"`
    UpdateAt    *time.Time          `json:"updateAt" label:"更新时间"`
    CreateBy    string              `json:"createBy" label:"创建人"`
    UpdateBy    string              `json:"updateBy" label:"更新人"`
}
```

#### CreateUpdateCt（创建更新合并请求）

```go
package modXxx

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type CreateUpdateCt struct {
    ID   typePg.Uint64String `json:"id" form:"id" label:"id"`
    Name string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称"`
    // ID > 0 时走更新，否则走创建
}
```

#### QueryCt（列表查询请求）

```go
package modXxx

import (
    "github.com/hongmengzhu/xianfu-blog-go/pkg/model"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
    "time"
)

type QueryCt struct {
    model.BaseQueryCt // 嵌入获得 PageNum, PageSize, Wd 字段
    Name     string              `json:"name" label:"名称"`
    State    typePg.Int8         `json:"state" label:"状态"`
    // 其他筛选字段
}
```

**Model 关键规则：**
- `typePg.Uint64String` — ID 字段类型，JSON 序列化为字符串，防止前端精度丢失
- `typePg.Int8` — 状态等小整数字段类型
- `model.BaseQueryCt` — 嵌入后自动拥有 `PageNum`、`PageSize`、`Wd`（关键词搜索）字段
- `validate` tag — 用于参数校验（`required`, `min`, `max` 等）
- `label` tag — 用于校验错误的中文提示
- manage 层分开 `CreateCt` 和 `UpdateCt`；system 层常合并为 `CreateUpdateCt`

### 步骤 4：编写 Service（业务逻辑层）

路径：`app/{layer}/{module}/service/XxxService.go`

```go
package service

import (
    "context"
    "reflect"

    "github.com/gin-gonic/gin"
    "github.com/hongmengzhu/xianfu-blog-go/app/{layer}/{module}/model/modXxx"
    "github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityXxx"
    "github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryXxx"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/holderPg"
    "go-spring.org/log"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/model"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
    "github.com/jinzhu/copier"
    "github.com/pangu-2/go-tools/tools/dbPg/pagePg"
    "github.com/pangu-2/go-tools/tools/noPg"
    "github.com/pangu-2/go-tools/tools/numberPg"
    "github.com/pangu-2/go-tools/tools/strPg"
    "github.com/pangu-2/go-tools/tools/wrapperPg/rg"
    "go-spring.org/log"
    "go-spring.org/spring/gs"
)

func init() {
    gs.Provide(new(XxxService))
}

type XxxService struct {
    sv  *repositoryXxx.XxxRepository `autowire:"?"`
    // 其他依赖的 Repository 或 Service
}
```

#### Create（新增）

```go
func (c *XxxService) Create(ctx *gin.Context, ct modXxx.CreateCt) (rt rg.Rs[string]) {
	log.Debugf(ctx, log.TagAppDef, "ct=%+v", ct)
    var info entityXxx.XxxEntity
    copier.Copy(&info, &ct)
    // 1. 参数校验
    if "" == ct.Name {
        return rt.ErrorMessage("名称不能为空")
    }
    // 2. 业务校验（唯一性等）
    // 3. 设置系统字段
    holder := holderPg.GetContextAccount(ctx)
    info.TenantNo = holder.GetTenantNo()
    info.No = noPg.No()
    // 4. 保存
    err, _ := c.sv.Create(ctx, &info)
    if err != nil {
        return rt.ErrorMessage("保存失败 " + err.Error())
    }
    return rg.OkData(numberPg.Int64ToString(info.ID))
}
```

#### Update（更新）

```go
func (c *XxxService) Update(ctx *gin.Context, ct modXxx.UpdateCt) (rt rg.Rs[string]) {
    log.Debugf(ctx, log.TagAppDef, "ct=%+v", ct)
	var info entityXxx.XxxEntity
    copier.Copy(&info, &ct)
    if ct.ID < 1 {
        return rt.ErrorMessage("id错误")
    }
    // 1. 查找原记录
    find, b := c.sv.FindById(ctx, ct.ID.ToInt64())
    if !b {
        return rt.ErrorMessage("数据不存在")
    }
    // 2. 业务校验
    // 3. 清除不可更新字段
    info.ID = 0
    info.No = ""
    err := c.sv.Update(ctx, info, find.ID)
    if err != nil {
        return rt.ErrorMessage(err.Error())
    }
    return rt.Ok()
}
```

#### Detail（详情）

```go
func (c *XxxService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modXxx.Vo]) {
    if id < 1 {
        return rt.ErrorMessage("id错误")
    }
    find, b := c.sv.FindById(ctx, id)
    if !b {
        return rt.ErrorMessage("数据不存在")
    }
    var info modXxx.Vo
    copier.Copy(&info, &find)
    return rt.OkData(info)
}
```

#### Query（分页查询）

```go
func (c *XxxService) Query(ctx *gin.Context, ct modXxx.QueryCt) (rt rg.Rs[pagePg.Paginator[modXxx.Vo]]) {
    log.Debugf(ctx, log.TagAppDef, "ct=%+v", ct)
    var query entityXxx.XxxEntity
    copier.Copy(&query, &ct)
    slice := make([]modXxx.Vo, 0)
    rt.Data.Data = slice
    page, err := c.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
        if ct.PageSize < 1 {
            ct.PageSize = 20
        }
        arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
        arg.Db = arg.Db.Order("create_at desc")
        // 自定义筛选
        if strPg.IsNotBlank(ct.Wd) {
            arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
        }
    }), optionsPg.WithCtx(ctx))
    if nil != err {
        return rt.Ok()
    }
    if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
        pg := pagePg.NewPaginatorByPageable[modXxx.Vo](page.Pageable)
        for _, item := range page.Data {
            var vo modXxx.Vo
            copier.Copy(&vo, &item)
            slice = append(slice, vo)
        }
        pg.Data = slice
        pg.Pageable = page.Pageable
        rt.Data = pg
        return rt.Ok()
    }
    return rt.Ok()
}
```

#### Enable / Disable（启用/禁用）

```go
func (c *XxxService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
    log.Debugf(ctx, log.TagAppDef, "ct=%+v", ct)
    return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

func (c *XxxService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
    return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}
```

#### LogicalDeletion（逻辑删除）

```go
func (c *XxxService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
    log.Debugf(ctx, log.TagAppDef, "ct=%+v", ids)
    if len(ids) < 1 {
        return rt.ErrorMessage("id错误")
    }
    finds, b := c.sv.FindAllByIdStringIn(ctx, ids)
    if !b {
        return rt.ErrorMessage("数据不存在")
    }
    if c.sv.Config().Data.Delete {
        c.sv.DeleteByIdsString(ctx, ids)
    } else {
        for _, info := range finds {
            enum := enumStatePg.State(info.State)
            if ok, reverse := enum.ReverseEnableDisable(); ok {
                c.sv.Update(ctx, entityXxx.XxxEntity{State: reverse.IndexInt8()}, info.ID)
            }
        }
    }
    return rt.Ok()
}
```

#### ExistName（名称查重）

```go
func (c *XxxService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
    log.Debugf(ctx, log.TagAppDef, "ct=%+v", ct)
    if "" == ct.Wd {
        return rt.ErrorMessage("关键词不能为空")
    }
    id := "0"
    if strPg.IsNotBlank(ct.Id) {
        id = ct.Id
    }
    _, result := c.sv.FindByNameAndIdNot(ctx, ct.Wd, id)
    if result {
        return rt.ErrorMessage("重复，已存在")
    }
    return rt.OkMessage("可以使用")
}
```

### 步骤 5：编写 Controller（控制器层）

路径：`app/{layer}/{module}/controller/Xxx.go`

```go
package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/hongmengzhu/xianfu-blog-go/app/{layer}/{module}/model/modXxx"
    "github.com/hongmengzhu/xianfu-blog-go/app/{layer}/{module}/service"
    "github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"
    "github.com/hongmengzhu/xianfu-blog-go/middleware/validatorPg"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/model"
    "github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
    "github.com/pangu-2/go-tools/tools/strPg"
    "github.com/pangu-2/go-tools/tools/wrapperPg/rg"
    "go-spring.org/spring/gs"
)

func init() {
    gs.Provide(new(XxxController)).Name("XxxController").Export(gs.As[routerPg.RouteRegistrar]())
}

type XxxController struct {
    routerPg.RouteRegistrar
    Sp *authPg.GroupSystemMiddlewareSp `autowire:""`   // system层
    // Sp *authPg.GroupManageMiddlewareSp `autowire:""` // manage层
    sv *service.XxxService             `autowire:"?"`
}

func (c *XxxController) RegisterRoutes(e *gin.Engine) {
    // system层路径前缀
    group := e.Group("/xianfu/sys/module/xxx", authPg.GroupSystemMiddleware(c.Sp))
    // manage层路径前缀
    // group := e.Group("/xianfu/manage/module/xxx", authPg.GroupManageMiddleware(c.Sp))

    group.POST("/createUpdate", c.CreateUpdate)  // 合并创建更新
    group.GET("/detail/:id", c.Detail)
    group.POST("/enable", c.Enable)
    group.POST("/disable", c.Disable)
    group.POST("/delete", c.Delete)
    group.POST("/recovery", c.Recovery)
    group.POST("/physicalDeletion", c.PhysicalDeletion)
    group.POST("/query", c.Query)
    group.POST("/selectNodeAll", c.SelectNodeAll)
    group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
    group.POST("/existName", c.ExistName)
    group.POST("/existCode", c.ExistCode)
}
```

#### Controller 标准方法

**CreateUpdate（合并写法）：**

```go
func (c *XxxController) CreateUpdate(ctx *gin.Context) {
    var ct modXxx.CreateUpdateCt
    if !routerPg.BindJson(ctx, &ct) {
        return
    }
    if ct.ID.ToInt64() > 0 {
        ctx.JSON(200, c.sv.Update(ctx, ct))
    } else {
        ctx.JSON(200, c.sv.Create(ctx, ct))
    }
}
```

**Detail（路由参数获取 ID）：**

```go
func (c *XxxController) Detail(ctx *gin.Context) {
    param := ctx.Param("id")
    if "" == param {
        ctx.JSON(200, rg.Error[string]("id不能为空"))
        return
    }
    ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}
```

**Query（列表查询）：**

```go
func (c *XxxController) Query(ctx *gin.Context) {
    var ct modXxx.QueryCt
    if !routerPg.BindJson(ctx, &ct) {
        return
    }
    ctx.JSON(200, c.sv.Query(ctx, ct))
}
```

**Delete（批量逻辑删除）：**

```go
func (c *XxxController) Delete(ctx *gin.Context) {
    var ct model.BaseIdsCt[string]
    if !routerPg.BindJson(ctx, &ct) {
        return
    }
    ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}
```

**Enable / Disable（批量启用/禁用）：**

```go
func (c *XxxController) Enable(ctx *gin.Context) {
    var ct model.BaseIdsCt[string]
    if !routerPg.BindJson(ctx, &ct) {
        return
    }
    ctx.JSON(200, c.sv.Enable(ctx, ct))
}
```

---

## 常见模式速查

### 唯一性校验

```go
_, result := c.sv.FindByNameAndIdNot(ctx, ct.Wd, id)
if result {
    return rt.ErrorMessage("重复，已存在")
}
```

### 对象拷贝（Entity ↔ DTO）

```go
var info entityXxx.XxxEntity
copier.Copy(&info, &ct)  // DTO → Entity
var vo modXxx.Vo
copier.Copy(&vo, &item)  // Entity → Vo
```

### 分页查询自定义条件

```go
page, err := c.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
    arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
    arg.Db = arg.Db.Order("create_at desc")
    if strPg.IsNotBlank(ct.Wd) {
        arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
    }
}), optionsPg.WithCtx(ctx))
```

### 多租户

```go
holder := holderPg.GetContextAccount(ctx)
info.TenantNo = holder.GetTenantNo()
```

- 通过 `holderPg.GetContextAccount(ctx)` 获取当前登录信息
- 创建时设置 `info.TenantNo = holder.GetTenantNo()`
- 查询时 Repository 自动附加租户过滤

### 状态枚举

使用 `enumStatePg.State`：

| 值 | 含义 |
|---|------|
| 1 (`enumStatePg.ENABLE`) | 有效 |
| 2 (`enumStatePg.DISABLE`) | 停用 |

- `ReverseEnableDisable()` — 有效↔停用 翻转

---

## 依赖注入规则

所有需要 go-spring 管理的组件必须在 `init()` 中注册：

```go
func init() {
    gs.Provide(new(XxxService))
}
```

**依赖注入 tag：**

| Tag | 说明 |
|-----|------|
| `autowire:"?"` | 可选注入，找不到时为 nil（常用） |
| `autowire:""` | 必须注入，找不到时报错（用于中间件等） |
| `value:"${xxx}"` | 从配置文件读取值 |

**Controller 注册为路由：**

```go
gs.Provide(new(XxxController)).Name("XxxController").Export(gs.As[routerPg.RouteRegistrar]())
```

- 必须设置 `.Name("唯一名称")` 避免同名 Controller 冲突
- 必须 `.Export(gs.As[routerPg.RouteRegistrar]())` 才能被自动发现并注册路由

## 参数绑定方式

| 方式 | 用法 | 场景 |
|------|------|------|
| `routerPg.BindJson(ctx, &ct)` |自动校验+中文错误提示 |
| `ctx.Param("id")` | 路由参数 | GET `/detail/:id` |

## 统一响应格式

所有接口返回 `rg.Rs[T]` 泛型结构：

```go
rg.OkData(data)              // 成功，带数据
rt.Ok()                      // 成功，无数据
rt.OkMessage("msg")          // 成功，带消息
rt.ErrorMessage("msg")       // 失败，带错误消息
rg.Error[T]("msg")           // 错误响应
rg.ErrorDefault[T]()         // 默认错误响应
rg.ErrorMessageData[T](msgs) // 带校验错误数据的响应
```

## 常用公共 Model

| Model | 路径 | 用途 |
|-------|------|------|
| `model.BaseQueryCt` | `pkg/model/BaseQueryCt.go` | 分页查询基类（PageNum, PageSize, Wd） |
| `model.BaseIdsCt[T]` | `pkg/model/BaseIdsCt.go` | 批量操作（Ids []T） |
| `model.BaseExistWdCt[T]` | `pkg/model/BaseExistWdCt.go` | 查重（Id排除, Wd关键词） |
| `model.BaseStateIdsCt[T]` | `pkg/model/BaseStateIdsCt.go` | 批量状态变更 |
| `model.BaseNode` | `pkg/model/BaseNode.go` | 树形节点（Key, Label, ParentId, Extend） |
| `typePg.Uint64String` | `pkg/tools/typePg/` | ID类型，JSON序列化为字符串 |
| `typePg.Int8` | `pkg/tools/typePg/` | 小整数类型 |

## 路由注册规则

- 实现 `routerPg.RouteRegistrar` 接口的 `RegisterRoutes(e *gin.Engine)` 方法
- 使用 `e.Group(路径前缀, 中间件)` 创建路由组
- 路径规范：
  - {layer}层：`/xianfu/{layer}/{module}/{entity}

**标准路由动作：**

| HTTP方法 | 路径                | 动作 | 说明            |
|---------|---------------------|------|-----------------|
| POST | `/createUpdate`     | 创建 | 合并            |
| GET | `/detail/:id`       | 详情 | 路由参数传ID    |
| POST | `/enable`           | 启用 | 批量，Body传ids |
| POST | `/disable`          | 禁用 | 批量，Body传ids |
| POST | `/delete`           | 逻辑删除 | 批量            |
| POST | `/recovery`         | 恢复删除 | 批量            |
| POST | `/physicalDeletion` | 物理删除 | 批量            |
| POST | `/query`            | 分页查询 | Body传筛选条件  |
| POST | `/selectNodePublic` | 节点列表 | 用于下拉选择    |
| POST | `/selectNodeAll`    | 全量节点 | 含所有状态      |
| POST | `/existName`        | 名称查重 |                 |
| POST | `/existCode`        | 标志查重 |                 |

---

## 基础设施层（共享）

```
infrastructure/
├── entityBasic/       # 数据库实体（GORM模型）
├── entityApi/
├── entityRam/
├── repositoryBasic/   # 数据访问层（泛型BaseRepository）
├── repositoryApi/
└── repositoryRam/
```

---

## 快速检查

每添加完一个接口后确认：

1. Entity 的 `TableName()` 返回正确的表名
2. Repository 的 `init()` 中注册了自身和 `support.BaseService`
3. Controller 的 `init()` 中设置了唯一的 `.Name()` 和 `.Export(gs.As[routerPg.RouteRegistrar]())`
4. Controller 的 `RegisterRoutes` 使用了正确的路由前缀和中间件
5. Service 中 `autowire` 注入了所需的 Repository
6. Model 的 `json` tag 与前端约定一致
7. `validate` tag 设置了必要的校验规则
8. 创建时设置了 `TenantNo` 和 `No`
