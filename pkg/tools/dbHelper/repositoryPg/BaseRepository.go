package repositoryPg

import (
	"context"
	"errors"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constContextPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbManagerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPgI"
	"github.com/pangu-2/go-tools/tools/dbPg/genericPg"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"go-spring.org/log"
	"gorm.io/gorm"
)

type RepositoryBase struct {
	//从内部
	db *gorm.DB `autowire:"?"`
	//
	ctx *gin.Context
}

func (b *RepositoryBase) SetCtx(ctx *gin.Context, arg ...any) {
	b.setDb(b.db.WithContext(holderPg.SetContextValue(ctx)))
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case *gorm.DB:
				b.setDb(result)
			}
		}
	}
}

func (b *RepositoryBase) setDb(db *gorm.DB) {
	b.db = db
}

type IRepository[T any, ID genericPg.ID] interface {
	repositoryPgI.IRepositoryBase
	Save(T)
	SaveAll([]T)
	FindById(ID)
	DeleteById(ID)
	DeleteByIds([]ID)
	FindAll(T)
	FindAllPage(T)
	FindAllByIdIn([]ID)
	Db()
	DbScopes()
}

type BaseRepository[T any, ID genericPg.ID] struct {
	// Entity 已废弃：旧实现会对该字段做惰性写入，共享 Bean 并发复用时构成数据竞争；
	// 现各方法改用局部 new(T) 建模，这里仅为兼容旧代码保留声明。
	Entity *T
	//从内部
	DbManger *dbManagerPg.Manager `autowire:"?"`
	Pg       configPg.Pg          `value:"${pg}"`
}

func (b *BaseRepository[T, ID]) DbScopes() *gorm.DB {
	return b.DbManger.Db()
}

func (b *BaseRepository[T, ID]) Db() *gorm.DB {
	return b.DbScopes()
}

func (b *BaseRepository[T, ID]) DbSource() *gorm.DB {
	return b.DbManger.Db()
}
func (b *BaseRepository[T, ID]) DbModel() *gorm.DB {
	// 用局部零值实体建模，不再对共享单例的 Entity 字段做惰性写入：
	// 同一个 repository Bean 被多个请求并发复用时，旧写法的首次赋值会构成数据竞争。
	return b.DbScopes().Model(new(T))
}

func (b *BaseRepository[T, ID]) SetOptionScopes(db *gorm.DB, opts ...optionsPg.Option) *gorm.DB {
	arg := optionsPg.OptionParams{
		Db: db,
	}
	for _, opt := range opts {
		opt(&arg)
	}
	//FIXME 待优化
	if arg.Condition != nil {
		arg.Db = arg.Condition(arg.Db)
	}
	// 调用方通常只传 WithContext(ctx)，这里自动从 GORM Statement 取回 gin.Context，
	// 避免每个 repository 方法都必须显式传 optionsPg.WithCtx(ctx) 才能启用多租户 scope。
	if nil == arg.Ctx && nil != arg.Db && nil != arg.Db.Statement {
		if gc, ok := arg.Db.Statement.Context.(*gin.Context); ok {
			arg.Ctx = gc
		}
	}
	if nil != arg.Ctx {
		_, exists := arg.Ctx.Get(constContextPg.CTX_MULITI_TENANT)
		//log.Error(ctx, log.TagAppDef,"exists=xxxxxxx=%+v", exists)
		if exists {
			//解析表名称
			if err := arg.Db.Statement.Parse(new(T)); err != nil {
				log.Error(arg.Ctx, log.TagAppDef, err, "failed to parse table name")
				// 解析失败时 Statement.Schema 为空，继续访问 Schema.Table 会触发空指针；
				// 把错误挂到会话上，让后续查询显式失败，而不是崩溃或绕过租户过滤。
				_ = arg.Db.AddError(err)
				return arg.Db
			}
			return arg.Db.Scopes(multiTenantPg.ScopeRulePgWhere(arg.Ctx, arg.Db.Statement.Schema.Table))
		}
	}
	return arg.Db
}

func (b *BaseRepository[T, ID]) SetOptionPgScopes(db *gorm.DB, opts ...optionsPg.Option) (*gorm.DB, optionsPg.OptionParams) {
	arg := optionsPg.OptionParams{
		Db: db,
	}
	for _, opt := range opts {
		opt(&arg)
	}
	if nil == arg.Ctx && nil != arg.Db && nil != arg.Db.Statement {
		if gc, ok := arg.Db.Statement.Context.(*gin.Context); ok {
			arg.Ctx = gc
		}
	}
	if nil != arg.Ctx {
		_, exists := arg.Ctx.Get(constContextPg.CTX_MULITI_TENANT)
		//log.Error(ctx, log.TagAppDef,"exists=xxxxxxx=%+v", exists)
		if exists {
			//解析表名称
			if err := arg.Db.Statement.Parse(new(T)); err != nil {
				log.Error(arg.Ctx, log.TagAppDef, err, "failed to parse table name")
				_ = arg.Db.AddError(err)
				return arg.Db, arg
			}
			return arg.Db.Scopes(multiTenantPg.ScopeRulePgWhere(arg.Ctx, arg.Db.Statement.Schema.Table)), arg
		}
	}
	return arg.Db, arg
}

func (b *BaseRepository[T, ID]) Config() configPg.Pg {
	return b.Pg
}

// Create 创建
func (b *BaseRepository[T, ID]) Create(ctx context.Context, v *T) (error, int64) {
	tx := b.Db().WithContext(ctx).Create(&v)
	return tx.Error, tx.RowsAffected
}

// CreateBatch 创建
func (b *BaseRepository[T, ID]) CreateBatch(ctx context.Context, v []T) (error, int64) {
	tx := b.Db().WithContext(ctx).Create(&v)
	return tx.Error, tx.RowsAffected
}

// CreateBatchPointer 创建
func (b *BaseRepository[T, ID]) CreateBatchPointer(ctx context.Context, v []*T) (error, int64) {
	tx := b.Db().WithContext(ctx).Create(&v)
	return tx.Error, tx.RowsAffected
}

// Save 保存 会保存所有的字段，即使字段是零值
func (b *BaseRepository[T, ID]) Save(ctx context.Context, v *T, opts ...optionsPg.Option) error {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Save(&v)
	return tx.Error
}

// SaveAll 保存
// 批量写入：PostgreSQL 下生成一条 INSERT ... ON CONFLICT (主键) DO UPDATE SET ...，
// 由数据库一次性完成插入/更新，避免循环里逐条 Save。
func (b *BaseRepository[T, ID]) SaveAll(ctx context.Context, ts []*T, opts ...optionsPg.Option) error {
	if len(ts) == 0 {
		return nil
	}
	tx := b.SetOptionScopes(b.Db().WithContext(ctx), opts...).Save(&ts)
	if tx.Error != nil {
		log.Error(ctx, log.TagAppDef, tx.Error, "SaveAll failed")
	}
	return tx.Error
}

// Update 更新 更新属性，只会更新非零值的字段
func (b *BaseRepository[T, ID]) Update(ctx context.Context, info T, id ID, opts ...optionsPg.Option) error {
	result := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id=?", id).Updates(&info)
	if result.Error != nil {
		//log.Fatal("failed to connect database")
		// record not found 跳过日志
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, result.Error, "err")
		}
		return result.Error
	}
	return nil
}

// Update 更新 更新属性，只会更新非零值的字段
func (b *BaseRepository[T, ID]) UpdatePointer(ctx context.Context, info T, id ID) error {
	result := b.SetOptionScopes(b.DbModel().WithContext(ctx)).Where("id=?", id).Updates(info)
	if result.Error != nil {
		//log.Fatal("failed to connect database")
		// record not found 跳过日志
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, result.Error, "")
		}
		return result.Error
	}
	return nil
}

// UpdatePointerObject 更新  更新属性，只会更新非零值的字段
func (b *BaseRepository[T, ID]) UpdatePointerObject(ctx context.Context, info any, id ID) error {
	result := b.SetOptionScopes(b.DbModel().WithContext(ctx)).Where("id=?", id).Updates(info)
	if result.Error != nil {
		//log.Fatal("failed to connect database")
		// record not found 跳过日志
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, result.Error, "")
		}
		return result.Error
	}
	return nil
}

// UpdateMap 更新, map里所有属性都会更新
func (b *BaseRepository[T, ID]) UpdateMap(ctx context.Context, info map[string]any, id ID, opts ...optionsPg.Option) error {
	result := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id=?", id).Updates(info)
	if result.Error != nil {
		//log.Fatal("failed to connect database")
		// record not found 跳过日志
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, result.Error, "")
		}
		return result.Error
	}
	return nil
}

// UpdateByIdIn 按主键集合批量更新（一条 UPDATE ... WHERE id IN (...)，只更新非零值字段）
func (b *BaseRepository[T, ID]) UpdateByIdIn(ctx context.Context, info T, ids []ID, opts ...optionsPg.Option) error {
	if len(ids) == 0 {
		return nil
	}
	result := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id IN ?", ids).Updates(&info)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, result.Error, "UpdateByIdIn failed")
		}
		return result.Error
	}
	return nil
}

// UpdateStructMap 更新, 结构体转换为map，map里所有属性都会更新
func (b *BaseRepository[T, ID]) UpdateStructMap(ctx context.Context, info any, id ID, opts ...optionsPg.Option) error {
	toMap, err := convertor.StructToMap(info)
	if nil != err {
		return err
	}
	result := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id=?", id).Updates(toMap)
	if result.Error != nil {
		// record not found 跳过日志
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, result.Error, "")
		}
		return result.Error
	}
	return nil
}

// UpdateAll 更新
// 批量 upsert：按主键冲突一次性写回，分片控制单条 SQL 的参数个数。
func (b *BaseRepository[T, ID]) UpdateAll(ctx context.Context, ts []*T, opts ...optionsPg.Option) error {
	const batchSize = 100
	for i := 0; i < len(ts); i += batchSize {
		end := i + batchSize
		if end > len(ts) {
			end = len(ts)
		}
		if err := b.SaveAll(ctx, ts[i:end], opts...); err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseRepository[T, ID]) DeleteById(ctx context.Context, id ID, opts ...optionsPg.Option) error {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Delete(new(T), id)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return tx.Error
	}
	return nil
}

func (b *BaseRepository[T, ID]) DeleteByIds(ctx context.Context, id []ID, opts ...optionsPg.Option) error {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Delete(new(T), id)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return tx.Error
	}
	return nil
}

func (b *BaseRepository[T, ID]) DeleteByIdsString(ctx context.Context, id []string, opts ...optionsPg.Option) error {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Delete(new(T), id)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return tx.Error
	}
	return nil
}

func (b *BaseRepository[T, ID]) DeleteAllByTenantNoAndIdsString(ctx context.Context, tenantNo string, id []string, opts ...optionsPg.Option) error {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("tenant_no = ?", tenantNo).Delete(new(T), id)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return tx.Error
	}
	return nil
}

func (b *BaseRepository[T, ID]) DeleteByNo(ctx context.Context, no string, opts ...optionsPg.Option) error {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("no = ?", no).Delete(new(T))
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return tx.Error
	}
	return nil
}
func (b *BaseRepository[T, ID]) DeleteByCondition(ctx context.Context, t *T, opts ...optionsPg.Condition) error {
	where := b.SetOptionScopes(b.Db().WithContext(ctx))
	if nil != t {
		where = where.Where(t)
	}
	if nil != opts {
		for _, item := range opts {
			where = item(where)
		}
	}
	tx := where.Delete(new(T))
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return tx.Error
	}
	return nil
}

// FindById 根据主键查询
//
//	@Description:
//	@receiver b
//	@param id
//	@return info
//	@return result 是否查询到值
//	@return err
func (b *BaseRepository[T, ID]) FindById(ctx context.Context, id ID, opts ...optionsPg.Option) (info *T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id=?", id).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindByIdString 根据主键查询
//
//	@Description:
//	@receiver b
//	@param id
//	@return info
//	@return result 是否查询到值
//	@return err
func (b *BaseRepository[T, ID]) FindByIdString(ctx context.Context, id string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id=?", id).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindAll 查询所有
func (b *BaseRepository[T, ID]) FindAll(ctx context.Context, t T, arg ...any) (infos []*T) {
	where := b.SetOptionScopes(b.Db().WithContext(ctx)).Where(t)
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil
	}
	return
}

// FindAllOption 查询所有
func (b *BaseRepository[T, ID]) FindAllOption(ctx context.Context, arg ...any) (infos []*T) {
	where := b.SetOptionScopes(b.Db().WithContext(ctx))
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return
	}
	return
}

// FindAllLimit 查询所有-但是限制条数
func (b *BaseRepository[T, ID]) FindAllLimit(ctx context.Context, t T, limit int, arg ...any) (infos []*T, result bool) {
	where := b.SetOptionScopes(b.Db().WithContext(ctx)).Where(t).Limit(limit)
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return infos, true
}

// FindAllData 查询所有
func (b *BaseRepository[T, ID]) FindAllData(ctx context.Context, arg ...any) (infos []*T, result bool) {
	where := b.SetOptionScopes(b.Db().WithContext(ctx))
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return infos, true
}

// FindAllPage 分页
func (b *BaseRepository[T, ID]) FindAllPage(ctx context.Context, t T, opts ...optionsPg.Option) (pagePg.Paginator[*T], error) {
	var total int64
	//
	condition, arg := b.SetOptionPgScopes(b.DbModel().WithContext(ctx), opts...)
	//
	pageable := arg.Pageable
	if nil == pageable {
		pageable = &pagePg.Pageable{Total: total, PageNum: 0, PageSize: 10}
	}
	pg := pagePg.NewPaginator[*T]()
	//
	countTx := condition.Where(t).Count(&total)
	if nil != countTx.Error {
		return pg, countTx.Error
	}
	var infos []*T
	tx := countTx.Scopes(optionsPg.WithScopes(pageable)).Find(&infos)
	//log.Info(ctx, log.TagAppDef,"sql=%+v", tx.Statement.SQL.String())
	if tx.Error != nil {
		return pg, tx.Error
	}
	pg.Data = infos
	pg.Total = total
	pg.Pageable = pagePg.NewPageable(total, pageable.PageNum, pageable.PageSize)
	pg.TotalPage = pg.Pageable.TotalPage
	//if total >0 {
	//	t2 := infos[len(infos)-1]
	//	pg.OffsetId=t2
	//}
	return pg, nil
}

// FindAllByIdIn
//
//	@Description:
//	@receiver b
//	@param ids
//	@return infos
//	@return result true: 有值;    false: 错误或 没查询到
func (b *BaseRepository[T, ID]) FindAllByIdIn(ctx context.Context, ids []ID, opts ...optionsPg.Option) (infos []*T, result bool) {
	infos = make([]*T, 0)
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id in (?)", ids).Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	result = true
	return
}

// FindAllByIdStringIn
//
//	@Description:
//	@receiver b
//	@param ids
//	@return infos
//	@return result true: 有值;    false: 错误或 没查询到
func (b *BaseRepository[T, ID]) FindAllByIdStringIn(ctx context.Context, ids []string, opts ...optionsPg.Option) (infos []*T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("id in (?)", ids).Find(&infos)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	result = true
	return
}

// Count
//
//	@Description:   统计
//	@receiver b
//	@param arg
//	@return total
//	@return result
func (b *BaseRepository[T, ID]) Count(ctx context.Context, arg ...any) (total int64, result bool) {
	where := b.SetOptionScopes(b.Db().WithContext(ctx))
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				// 条件
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Count(&total)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return 0, false
	}
	if tx.RowsAffected == 0 {
		return 0, false
	}
	return total, true
}

// FindByNo 根据主键查询
//
//	@Description:
//	@receiver b
//	@param no
//	@return info
//	@return result 是否查询到值
//	@return err
func (b *BaseRepository[T, ID]) FindByNo(ctx context.Context, no string, arg ...any) (info *T, result bool) {
	where := b.SetOptionScopes(b.Db().WithContext(ctx))
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				// 条件
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Where("no=?", no).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		log.Debug(ctx, log.TagAppDef, "RowsAffected=0", func() []log.Field {
			return []log.Field{
				log.Any("error", tx.Error),
			}
		})
		return nil, false
	}
	return info, true
}

// FindAllByNoIn 根据主键查询
//
//	@Description:
//	@receiver b
//	@param no
//	@return info
//	@return result 是否查询到值
//	@return err
func (b *BaseRepository[T, ID]) FindAllByNoIn(ctx context.Context, no []string, opts ...optionsPg.Option) (info []*T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("no in ?", no).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindAllByNameIn 根据主键查询
//
//	@Description:
//	@receiver b
//	@param no
//	@return info
//	@return result 是否查询到值
//	@return err
func (b *BaseRepository[T, ID]) FindAllByNameIn(ctx context.Context, no []string, opts ...optionsPg.Option) (info []*T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("name in ?", no).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindByName 根据主键查询
//
//	@Description:
//	@receiver b
//	@param no
//	@return info
//	@return result 是否查询到值
//	@return err
func (b *BaseRepository[T, ID]) FindByName(ctx context.Context, no string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("name=?", no).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindByNameAndIdNot
//
//	@Description:
//	@receiver c
//	@param name
//	@param id
//	@return info
//	@return result
func (b *BaseRepository[T, ID]) FindByNameAndIdNot(ctx context.Context, name string, id string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("name=?", name).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// FindByNoAndIdNot
//
//	@Description:
//	@receiver c
//	@param name
//	@param id
//	@return info
//	@return result
func (c *BaseRepository[T, ID]) FindByNoAndIdNot(ctx context.Context, name string, id string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("no=?", name).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// FindAllByNoLink
//
//	@Description:
//	@receiver c
//	@param name
//	@return info
//	@return result
func (c *BaseRepository[T, ID]) FindAllByNoLink(ctx context.Context, code string, opts ...optionsPg.Option) (info []*T, result bool) {
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("no_link like ?", "%|"+code+"|%").Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// FindByCode 根据主键查询
//
//	@Description:
//	@receiver b
//	@param no
//	@return info
//	@return result 是否查询到值
//	@return err
func (c *BaseRepository[T, ID]) FindByCode(ctx context.Context, no string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := c.SetOptionScopes(c.DbModel().WithContext(ctx), opts...).Where("code=?", no).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}

// FindByCodeAndIdNot
//
//	@Description:
//	@receiver c
//	@param name
//	@param id
//	@return info
//	@return result
func (b *BaseRepository[T, ID]) FindByCodeAndIdNot(ctx context.Context, name string, id string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("code=?", name).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// FindByCodeAndNoNot
//
//	@Description:
//	@receiver c
//	@param name
//	@param no
//	@return info
//	@return result
func (b *BaseRepository[T, ID]) FindByCodeAndNoNot(ctx context.Context, name string, no string, opts ...optionsPg.Option) (info *T, result bool) {
	tx := b.SetOptionScopes(b.DbModel().WithContext(ctx), opts...).Where("code=?", name).Where("no != ?", no).First(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// Exist 是否存在
func (b *BaseRepository[T, ID]) Exist(ctx context.Context, t T, arg ...any) bool {
	var cnt int64
	where := b.SetOptionScopes(b.DbModel().WithContext(ctx)).Where(t)
	if nil != arg {
		for _, item := range arg {
			switch result := item.(type) {
			case optionsPg.Condition:
				where = result(where)
			case optionsPg.Option:
				where = b.SetOptionScopes(where, item.(optionsPg.Option))
			}
		}
	}
	tx := where.Count(&cnt)
	if tx.Error != nil {
		// 这里返回 false 表示"查不到"，但底层可能是连接/超时等故障，
		// 至少要把错误记进日志，避免故障被当成"数据不存在"静默吞掉。
		log.Error(ctx, log.TagAppDef, tx.Error, "Exist count failed")
		return false
	}
	return cnt > 0
}
