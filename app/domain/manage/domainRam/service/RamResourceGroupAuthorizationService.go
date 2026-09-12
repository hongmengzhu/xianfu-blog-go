package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/resourceTypeCategoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeAttrPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamResourceGroupAuthorizationService))
}

// RamResourceGroupAuthorizationService 资源组 授权方法动作
// @Description:
type RamResourceGroupAuthorizationService struct {
	sv     *repositoryRam.RamResourceGroupRepository     `autowire:"?"`
	authDb *repositoryRam.RamResourceAuthorityRepository `autowire:"?"`
	ra     *RamResourceAuthorizationService              `autowire:"?"`
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			r.Update(ctx, entityRam.RamResourceGroupEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	//有效
	//if enumStatePg.ENABLE.IsExistInt8(state.IndexInt8()) {
	//	for _, info := range finds {
	//		c.ra.AuthorityValid(ctx, utilsRam.ResourceAuthorityMarkByInt64(iamConstant.GroupResourceTypeCategory, info.ID))
	//	}
	//} else {
	//	//无效
	//	for _, info := range finds {
	//		c.ra.AuthorityInValid(ctx, utilsRam.ResourceAuthorityMarkByInt64(iamConstant.GroupResourceTypeCategory, info.ID))
	//	}
	//}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	return c.State(ctx, ids, state)
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	//物理删除
	if c.sv.Config().Data.Delete {
		for _, info := range finds {
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)

			//删除授权及权限规则
			//c.ra.Delete(ctx, utilsRam.ResourceAuthorityMarkByInt64(iamConstant.GroupResourceTypeCategory, info.ID))
		}
		repository.DeleteByIdsString(ctx, ids)

	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityRam.RamResourceGroupEntity{State: reverse.IndexInt8()}, info.ID)
			}

			//设置无效
			//c.ra.AuthorityInValid(ctx, utilsRam.ResourceAuthorityMarkByInt64(iamConstant.GroupResourceTypeCategory, info.ID))
		}
	}

	return rt.Ok()
}

// LogicalRecovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		enum := enumStatePg.State(info.State)
		//  取消 弃置 批量删除，反转 为对应的 有效 停用 停用
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityRam.RamResourceGroupEntity{State: reverse.IndexInt8()}, info.ID)
		}

		//设置有效
		//c.ra.AuthorityValid(ctx, utilsRam.ResourceAuthorityMarkByInt64(iamConstant.GroupResourceTypeCategory, info.ID))
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupAuthorizationService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	authDb := c.authDb
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	idsNew := make([]int64, 0)
	idsCategory := make([]string, 0)
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		//判断 是否是 分类
		if typeAttrPg.CategoryLast.IsEqual(info.TypeAttr) || typeAttrPg.Category.IsEqual(info.TypeAttr) {
			link, b := r.FindAllByIdLink(ctx, numberPg.Int64ToString(info.ID))
			if b {
				for _, entity := range link {
					idsCategory = append(idsCategory, numberPg.Int64ToString(entity.ID))
				}
				if len(idsCategory) > 0 {
					_, result := authDb.FindAllByTypeCategoryAndGroupIdStringIn(ctx, resourceTypeCategoryPg.Group.String(), idsCategory)
					if result {
						return rt.ErrorMessage("该分类下存在子数据，请先删除子数据")
					}
				}
				idsNew = append(idsNew, info.ID)
			}
		}
		//删除授权及权限规则
		//c.ra.Delete(ctx, utilsRam.ResourceAuthorityMarkByInt64(iamConstant.GroupResourceTypeCategory, info.ID))
	}
	if len(idsNew) > 0 {
		r.DeleteByIds(ctx, idsNew)
	}
	return rt.Ok()
}
