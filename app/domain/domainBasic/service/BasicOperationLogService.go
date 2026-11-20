package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicOperationLog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicOperationLogService))
}

// BasicOperationLogService 操作日志查询。
type BasicOperationLogService struct {
	sv *repositoryBasic.BasicOperationLogRepository `autowire:"?"`
}

const (
	defaultQueryDays = 7
	maxQueryDays     = 90
)

// Query 分页查询。
func (s *BasicOperationLogService) Query(ctx *gin.Context, ct modBasicOperationLog.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicOperationLog.Vo]]) {
	f := s.buildFilter(ct)
	page, err := s.sv.FindAllPage(ctx, entityBasic.BasicOperationLogEntity{}, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		ct.PageNum, ct.PageSize = optionsPg.ClampPage(ct.PageNum, ct.PageSize)
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = f.Apply(arg.Db).Order("create_at desc").Order("id desc")
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

// buildFilter 统一构造过滤条件：时间窗默认 7 天、最长 90 天。
func (s *BasicOperationLogService) buildFilter(ct modBasicOperationLog.QueryCt) repositoryBasic.OperationLogFilter {
	f := repositoryBasic.OperationLogFilter{
		TraceID:    ct.TraceID,
		Method:     ct.Method,
		Path:       ct.Path,
		TypeDomain: ct.TypeDomain,
		Module:     ct.Module,
		Action:     ct.Action,
		OperatorNo: ct.OperatorNo,
		IP:         ct.IP,
		Keyword:    ct.Wd,
		StatusCode: ct.StatusCode,
	}
	now := time.Now()
	endAt := now
	if ct.EndAt != nil {
		endAt = *ct.EndAt
	}
	startAt := endAt.AddDate(0, 0, -defaultQueryDays)
	if ct.StartAt != nil {
		startAt = *ct.StartAt
	}
	if endAt.Sub(startAt) > maxQueryDays*24*time.Hour {
		startAt = endAt.AddDate(0, 0, -maxQueryDays)
	}
	f.StartAt = &startAt
	f.EndAt = &endAt
	return f
}

func toVo(info *entityBasic.BasicOperationLogEntity) modBasicOperationLog.Vo {
	var vo modBasicOperationLog.Vo
	copier.Copy(&vo, info)
	return vo
}
