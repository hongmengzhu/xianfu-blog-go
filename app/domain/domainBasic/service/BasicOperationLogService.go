package service

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicOperationLog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/excelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicOperationLogService))
}

// BasicOperationLogService 操作日志(审计日志)查询。
type BasicOperationLogService struct {
	sv *repositoryBasic.BasicOperationLogRepository `autowire:"?"`
}

const (
	defaultQueryDays = 7
	maxQueryDays     = 90
	defaultStatDays  = 30
)

// Query 分页查询。
func (s *BasicOperationLogService) Query(ctx *gin.Context, ct modBasicOperationLog.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicOperationLog.Vo]]) {
	log.Info(ctx, log.TagAppDef, "audit query", log.Any("ct", ct))
	query, f, err := s.buildQuery(ctx, ct)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	page, err := s.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		ct.PageNum, ct.PageSize = optionsPg.ClampPage(ct.PageNum, ct.PageSize)
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = f.Apply(arg.Db)
		arg.Db = arg.Db.Order("operate_at desc").Order("id desc")
	}), optionsPg.WithCtx(ctx))
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	if page.Total == 0 || len(page.Data) == 0 {
		return rt.Ok()
	}
	pg := pagePg.NewPaginatorByPageable[modBasicOperationLog.Vo](page.Pageable)
	slice := make([]modBasicOperationLog.Vo, 0, len(page.Data))
	for _, item := range page.Data {
		slice = append(slice, toVo(item))
	}
	pg.Data = slice
	return rt.OkData(pg)
}

// Detail 日志详情。
func (s *BasicOperationLogService) Detail(ctx *gin.Context, id typePg.Int64String) (rt rg.Rs[modBasicOperationLog.DetailVo]) {
	if id.ToInt64() <= 0 {
		return rt.ErrorMessage("id 不能为空")
	}
	info, ok := s.sv.FindById(ctx, id.ToInt64(), optionsPg.WithCtx(ctx))
	if !ok || info == nil {
		return rt.ErrorMessage("日志不存在")
	}
	if err := s.checkTenant(ctx, info.TenantNo); err != nil {
		return rt.ErrorMessage(err.Error())
	}
	var vo modBasicOperationLog.DetailVo
	copier.Copy(&vo, info)
	vo.Vo = toVo(info)
	vo.OperatorTypeName = entityBasic.OperationLogOperatorName(int8(vo.OperatorType))
	return rt.OkData(vo)
}

// Statistics 统计趋势与分布。
func (s *BasicOperationLogService) Statistics(ctx *gin.Context, ct modBasicOperationLog.StatisticsCt) (rt rg.Rs[modBasicOperationLog.StatisticsVo]) {
	f, err := s.buildFilter(ctx, ct.TenantNo, ct.TypeDomain, ct.Module, ct.Action, ct.OperatorNo, "", "", 0, ct.StartAt, ct.EndAt, defaultStatDays)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	days, err := s.sv.StatByDay(ctx, f)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	out := modBasicOperationLog.StatisticsVo{
		Trend:   make([]modBasicOperationLog.Point, 0, len(days)),
		Modules: make([]modBasicOperationLog.GroupItem, 0),
		Actions: make([]modBasicOperationLog.GroupItem, 0),
		Risks:   make([]modBasicOperationLog.GroupItem, 0),
		Results: make([]modBasicOperationLog.GroupItem, 0),
	}
	for _, d := range days {
		out.Total += d.Total
		out.Success += d.Success
		out.Fail += d.Fail
		out.Partial += d.Partial
		out.Unknown += d.Unknown
		out.Trend = append(out.Trend, modBasicOperationLog.Point{
			Day: d.Day, Total: d.Total, Success: d.Success, Fail: d.Fail, Partial: d.Partial, Unknown: d.Unknown,
		})
	}
	for _, field := range []string{"module", "action", "risk", "result"} {
		rows, err := s.sv.GroupBy(ctx, f, field, 20)
		if err != nil {
			return rt.ErrorMessage(err.Error())
		}
		items := make([]modBasicOperationLog.GroupItem, 0, len(rows))
		for _, row := range rows {
			item := modBasicOperationLog.GroupItem{Key: row.Key, Name: row.Name, Count: row.Count}
			applyGroupLabels(&item, field)
			items = append(items, item)
		}
		switch field {
		case "module":
			out.Modules = items
		case "action":
			out.Actions = items
		case "risk":
			out.Risks = items
		case "result":
			out.Results = items
		}
	}
	return rt.OkData(out)
}

// Options 前端字典。
func (s *BasicOperationLogService) Options(ctx *gin.Context) (rt rg.Rs[modBasicOperationLog.OptionsVo]) {
	out := modBasicOperationLog.OptionsVo{
		OperationTypes: []modBasicOperationLog.Option{
			{Value: entityBasic.OperationLogTypeCreate, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeCreate)},
			{Value: entityBasic.OperationLogTypeUpdate, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeUpdate)},
			{Value: entityBasic.OperationLogTypeDelete, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeDelete)},
			{Value: entityBasic.OperationLogTypeQuery, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeQuery)},
			{Value: entityBasic.OperationLogTypeExport, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeExport)},
			{Value: entityBasic.OperationLogTypeImport, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeImport)},
			{Value: entityBasic.OperationLogTypeGrant, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeGrant)},
			{Value: entityBasic.OperationLogTypeStateChange, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeStateChange)},
			{Value: entityBasic.OperationLogTypeOther, Label: entityBasic.OperationLogTypeName(entityBasic.OperationLogTypeOther)},
		},
		RiskLevels: []modBasicOperationLog.Option{
			{Value: entityBasic.OperationLogRiskLow, Label: entityBasic.OperationLogRiskName(entityBasic.OperationLogRiskLow)},
			{Value: entityBasic.OperationLogRiskMedium, Label: entityBasic.OperationLogRiskName(entityBasic.OperationLogRiskMedium)},
			{Value: entityBasic.OperationLogRiskHigh, Label: entityBasic.OperationLogRiskName(entityBasic.OperationLogRiskHigh)},
		},
		Results: []modBasicOperationLog.Option{
			{Value: entityBasic.OperationLogResultSuccess, Label: entityBasic.OperationLogResultName(entityBasic.OperationLogResultSuccess)},
			{Value: entityBasic.OperationLogResultFail, Label: entityBasic.OperationLogResultName(entityBasic.OperationLogResultFail)},
			{Value: entityBasic.OperationLogResultPartial, Label: entityBasic.OperationLogResultName(entityBasic.OperationLogResultPartial)},
			{Value: entityBasic.OperationLogResultUnknown, Label: entityBasic.OperationLogResultName(entityBasic.OperationLogResultUnknown)},
		},
		Sources: []modBasicOperationLog.Option{
			{Value: entityBasic.OperationLogSourceApi, Label: entityBasic.OperationLogSourceName(entityBasic.OperationLogSourceApi)},
			{Value: entityBasic.OperationLogSourceEvent, Label: entityBasic.OperationLogSourceName(entityBasic.OperationLogSourceEvent)},
			{Value: entityBasic.OperationLogSourceCron, Label: entityBasic.OperationLogSourceName(entityBasic.OperationLogSourceCron)},
			{Value: entityBasic.OperationLogSourceManual, Label: entityBasic.OperationLogSourceName(entityBasic.OperationLogSourceManual)},
		},
		States: []modBasicOperationLog.Option{
			{Value: entityBasic.OperationLogStateNormal, Label: entityBasic.OperationLogStateName(entityBasic.OperationLogStateNormal)},
			{Value: entityBasic.OperationLogStateFrozen, Label: entityBasic.OperationLogStateName(entityBasic.OperationLogStateFrozen)},
		},
		OperatorTypes: []modBasicOperationLog.Option{
			{Value: entityBasic.OperationLogOperatorPlatform, Label: entityBasic.OperationLogOperatorName(entityBasic.OperationLogOperatorPlatform)},
			{Value: entityBasic.OperationLogOperatorTenant, Label: entityBasic.OperationLogOperatorName(entityBasic.OperationLogOperatorTenant)},
			{Value: entityBasic.OperationLogOperatorNormal, Label: entityBasic.OperationLogOperatorName(entityBasic.OperationLogOperatorNormal)},
			{Value: entityBasic.OperationLogOperatorAnonymous, Label: entityBasic.OperationLogOperatorName(entityBasic.OperationLogOperatorAnonymous)},
		},
		Modules: make([]modBasicOperationLog.GroupItem, 0),
	}
	tenant := s.currentTenant(ctx)
	rows, err := s.sv.GroupBy(ctx, repositoryBasic.OperationLogFilter{TenantNo: tenant}, "module", 50)
	if err == nil {
		for _, row := range rows {
			out.Modules = append(out.Modules, modBasicOperationLog.GroupItem{Key: row.Key, Name: row.Name, Count: row.Count})
		}
	}
	return rt.OkData(out)
}

// Operators 操作人下拉。
func (s *BasicOperationLogService) Operators(ctx *gin.Context, ct modBasicOperationLog.OperatorsCt) (rt rg.Rs[[]modBasicOperationLog.OperatorVo]) {
	f, err := s.buildFilter(ctx, "", "", "", "", "", "", "", 0, ct.StartAt, ct.EndAt, defaultStatDays)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	rows, err := s.sv.TopOperators(ctx, f, strings.TrimSpace(ct.Keyword), 20)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	out := make([]modBasicOperationLog.OperatorVo, 0, len(rows))
	for _, row := range rows {
		out = append(out, modBasicOperationLog.OperatorVo{
			OperatorNo: row.OperatorNo, OperatorAccount: row.OperatorAccount, OperatorName: row.OperatorName, Count: row.Count,
		})
	}
	return rt.OkData(out)
}

// TargetHistory 对象操作历史。
func (s *BasicOperationLogService) TargetHistory(ctx *gin.Context, ct modBasicOperationLog.TargetHistoryCt) (rt rg.Rs[pagePg.Paginator[modBasicOperationLog.Vo]]) {
	if strings.TrimSpace(ct.TargetType) == "" {
		return rt.ErrorMessage("对象类型不能为空")
	}
	if ct.TargetNo == "" && ct.TargetID.ToInt64() <= 0 {
		return rt.ErrorMessage("对象编号或对象ID不能为空")
	}
	f, err := s.buildFilter(ctx, "", "", "", "", "", ct.TargetType, ct.TargetNo, 0, ct.StartAt, ct.EndAt, defaultQueryDays)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	var query entityBasic.BasicOperationLogEntity
	if ct.TargetID.ToInt64() > 0 {
		query.TargetID = strconv.FormatInt(ct.TargetID.ToInt64(), 10)
	}
	page, err := s.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		ct.PageNum, ct.PageSize = optionsPg.ClampPage(ct.PageNum, ct.PageSize)
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = f.Apply(arg.Db).Order("operate_at desc").Order("id desc")
	}), optionsPg.WithCtx(ctx))
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	pg := pagePg.NewPaginatorByPageable[modBasicOperationLog.Vo](page.Pageable)
	slice := make([]modBasicOperationLog.Vo, 0, len(page.Data))
	for _, item := range page.Data {
		slice = append(slice, toVo(item))
	}
	pg.Data = slice
	return rt.OkData(pg)
}

// Export 导出日志；limit 由调用端控制（管理端 10000、系统端 50000）。
func (s *BasicOperationLogService) Export(ctx *gin.Context, ct modBasicOperationLog.QueryCt, limit int) {
	f, err := s.buildFilter(ctx, ct.TenantNo, ct.TypeDomain, ct.Module, ct.Action, ct.OperatorNo, ct.TargetType, ct.TargetNo, int8(ct.Result), ct.StartAt, ct.EndAt, defaultQueryDays)
	if err != nil {
		ctx.JSON(200, rg.Rs[string]{}.ErrorMessage(err.Error()))
		return
	}
	rows, err := s.sv.FindByPageForExport(ctx, f, limit)
	if err != nil {
		ctx.JSON(200, rg.Rs[string]{}.ErrorMessage(err.Error()))
		return
	}
	if len(rows) == 0 {
		ctx.JSON(200, rg.Rs[string]{}.ErrorMessage("没有任何数据"))
		return
	}
	slice := make([]interface{}, 0, len(rows))
	for i := range rows {
		vo := toVo(&rows[i])
		slice = append(slice, modBasicOperationLog.ExportVo{
			No: vo.No, OperateAt: vo.OperateAt, TypeDomain: vo.TypeDomain, ModuleName: vo.ModuleName,
			ActionName: vo.ActionName, OperationTypeName: vo.OperationTypeName, RiskLevelName: vo.RiskLevelName,
			Description: vo.Description, OperatorNo: vo.OperatorNo, OperatorAccount: vo.OperatorAccount,
			OperatorName: vo.OperatorName, TargetType: vo.TargetType, TargetNo: vo.TargetNo, TargetName: vo.TargetName,
			ResultName: vo.ResultName, ResultCode: vo.ResultCode, ErrorMessage: vo.ErrorMessage,
			DurationMs: vo.DurationMs, IPAddress: vo.IPAddress, RequestID: vo.RequestID, TraceID: vo.TraceID, TenantNo: vo.TenantNo,
		})
	}
	titles := []string{"日志编号", "操作时间", "域", "模块", "动作", "操作类型", "风险等级", "描述",
		"操作人编号", "操作人账号", "操作人姓名", "对象类型", "对象编号", "对象名称", "结果", "结果码",
		"失败原因", "耗时(毫秒)", "IP", "请求ID", "链路ID", "租户编号"}
	if err := excelPg.ExportExcelByStruct(ctx, titles, slice, "operation-log", "Sheet1"); err != nil {
		ctx.JSON(200, rg.Rs[string]{}.ErrorMessage(err.Error()))
	}
}

// buildQuery 生成实体条件与统计过滤条件。
func (s *BasicOperationLogService) buildQuery(ctx *gin.Context, ct modBasicOperationLog.QueryCt) (entityBasic.BasicOperationLogEntity, repositoryBasic.OperationLogFilter, error) {
	var query entityBasic.BasicOperationLogEntity
	copier.Copy(&query, &ct)
	f, err := s.buildFilter(ctx, ct.TenantNo, ct.TypeDomain, ct.Module, ct.Action, ct.OperatorNo, ct.TargetType, ct.TargetNo, int8(ct.Result), ct.StartAt, ct.EndAt, defaultQueryDays)
	if err != nil {
		return query, f, err
	}
	f.Keyword = strings.TrimSpace(ct.Wd)
	return query, f, nil
}

// buildFilter 统一构造过滤条件：租户收敛 + 时间窗（默认 N 天、最长 90 天）。
func (s *BasicOperationLogService) buildFilter(ctx *gin.Context, tenantNo, typeDomain, module, action, operatorNo, targetType, targetNo string, result int8, startAt, endAt *time.Time, defaultDays int) (repositoryBasic.OperationLogFilter, error) {
	f := repositoryBasic.OperationLogFilter{
		TenantNo: tenantNo, TypeDomain: typeDomain, Module: module, Action: action,
		OperatorNo: operatorNo, TargetType: targetType, TargetNo: targetNo, Result: result,
	}
	if err := s.checkTenant(ctx, tenantNo); err != nil && tenantNo != "" {
		return f, err
	}
	if holder := holderPg.GetContextAccount(ctx); !holder.IsFounder() {
		tenant := normalizeTenantNo(holder.GetTenantNo())
		if tenant == "" {
			return f, errors.New("无法识别当前租户")
		}
		if f.TenantNo != "" && f.TenantNo != tenant {
			return f, errors.New("无权查询其他租户的日志")
		}
		f.TenantNo = tenant
	}
	now := time.Now()
	if endAt == nil {
		t := now
		endAt = &t
	}
	if startAt == nil {
		t := endAt.AddDate(0, 0, -defaultDays)
		startAt = &t
	}
	if endAt.Sub(*startAt) > maxQueryDays*24*time.Hour {
		t := endAt.AddDate(0, 0, -maxQueryDays)
		startAt = &t
	}
	f.StartAt = startAt
	f.EndAt = endAt
	return f, nil
}

func (s *BasicOperationLogService) currentTenant(ctx *gin.Context) string {
	holder := holderPg.GetContextAccount(ctx)
	if holder.IsFounder() {
		return ""
	}
	return normalizeTenantNo(holder.GetTenantNo())
}

// checkTenant 校验日志是否属于当前租户；founder 不限制。
func (s *BasicOperationLogService) checkTenant(ctx *gin.Context, tenantNo string) error {
	holder := holderPg.GetContextAccount(ctx)
	if holder.IsFounder() || tenantNo == "" {
		return nil
	}
	if normalizeTenantNo(holder.GetTenantNo()) != tenantNo {
		return errors.New("无权查看其他租户的日志")
	}
	return nil
}

func normalizeTenantNo(tenantNo string) string {
	if tenantNo == "-1" {
		return ""
	}
	return tenantNo
}

func toVo(info *entityBasic.BasicOperationLogEntity) modBasicOperationLog.Vo {
	var vo modBasicOperationLog.Vo
	copier.Copy(&vo, info)
	vo.OperationTypeName = entityBasic.OperationLogTypeName(int8(vo.OperationType))
	vo.RiskLevelName = entityBasic.OperationLogRiskName(int8(vo.RiskLevel))
	vo.ResultName = entityBasic.OperationLogResultName(int8(vo.Result))
	vo.SourceName = entityBasic.OperationLogSourceName(int8(vo.Source))
	return vo
}

func applyGroupLabels(item *modBasicOperationLog.GroupItem, field string) {
	v, _ := parseSmallInt(item.Key)
	switch field {
	case "module":
		if item.Name == "" {
			item.Name = item.Key
		}
	case "action":
		if item.Name == "" {
			item.Name = item.Key
		}
	case "risk":
		item.Name = entityBasic.OperationLogRiskName(v)
	case "result":
		item.Name = entityBasic.OperationLogResultName(v)
	}
}

func parseSmallInt(s string) (int8, error) {
	if s == "" {
		return 0, errors.New("empty")
	}
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, errors.New("not number")
		}
		n = n*10 + int(r-'0')
	}
	return int8(n), nil
}
