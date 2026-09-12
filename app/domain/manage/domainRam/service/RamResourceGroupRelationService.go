package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceGroupRelation"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceRelation"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamResourceGroupRelationService))
}

// RamResourceGroupRelationService 资源关系
// @Description:
type RamResourceGroupRelationService struct {
	sv *repositoryRam.RamResourceGroupRelationRepository `autowire:"?"`
	//

}

// SelectNodePublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupRelationService) SelectNodePublic(ctx *gin.Context, ct modRamResourceGroupRelation.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID),
				Label: item.Name,
				Id:    numberPg.Int64ToString(item.ID),
				Code:  numberPg.Int64ToString(item.ID),
			})
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectNodeAllPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupRelationService) SelectNodeAllPublic(ctx *gin.Context, ct modRamResourceGroupRelation.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamResourceRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID),
				Label:  item.Name,
				Id:     numberPg.Int64ToString(item.ID),
				Code:   numberPg.Int64ToString(item.ID),
				Extend: vo,
			})
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupRelationService) SelectPublic(ctx *gin.Context, ct modRamResourceGroupRelation.QueryCt) (rt rg.Rs[[]modRamResourceRelation.Vo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupRelationEntity
	copier.Copy(&query, &ct)
	rt.Data = []modRamResourceRelation.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modRamResourceRelation.Vo, 0)
		for _, item := range infos {
			var vo modRamResourceRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// Selected 查询已选中的
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupRelationService) Selected(ctx *gin.Context, ct modRamResourceGroupRelation.QueryByTypeValueCt) (rt rg.Rs[[]string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if strPg.IsBlank(ct.TypeValue) {
		return rt.ErrorMessage("资源组类型id 不能为空")
	}
	if strPg.IsBlank(ct.TypeCategory) {
		return rt.ErrorMessage("类型 不能为空")
	}
	var query entityRam.RamResourceGroupRelationEntity
	query.TypeValue = ct.TypeValue
	query.TypeCategory = ct.TypeCategory
	slice := make([]string, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			slice = append(slice, item.GroupNo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}
