package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceRelation"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeAttrPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"time"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamResourceAuthorizationService))
}

// RamResourceAuthorizationService 资源授权动作 公共调用
// @Description:
type RamResourceAuthorizationService struct {
	ra            *RamResourceAuthorityService                  `autowire:"?"`
	resourceAuth  *repositoryRam.RamResourceAuthorityRepository `autowire:"?"`
	resourceGroup *repositoryRam.RamResourceGroupRepository     `autowire:"?"`
	resourceRel   *repositoryRam.RamResourceRelationRepository  `autowire:"?"`
	resource      *repositoryRam.RamResourceRepository          `autowire:"?"`
	casbin        *RamResourceCasbinService                     `autowire:"?"`
}

// UpdateByResourceGroup 更新权限
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *RamResourceAuthorizationService) UpdateByResourceGroup(ctx *gin.Context, ct modRamResourceRelation.UpdateByResourceGroupCt) (rt rg.Rs[string]) {
	if ct.SourceId <= 0 {
		return rt.ErrorMessage("资源组id不能为空")
	}
	source, b := c.resourceGroup.FindByIdString(ctx, ct.SourceId.ToString())
	if !b {
		return rt.ErrorMessage("资源组 不存在")
	}
	auth, b := c.ra.CreateByResourceGroup(ctx, source)
	if !b {
		return rt.ErrorMessage("创建资源权限失败")
	}
	log.Infof(ctx, log.TagAppDef, "authPg=%+v", auth)
	holder := holderPg.GetContextAccount(ctx)
	c.resourceRel.DeleteByAuthorityId(ctx, auth.ID)
	if len(ct.Ids) > 0 {
		c.resourceRel.DeleteByAuthorityId(ctx, auth.ID)
		now := time.Now()
		//增加权限
		infos, b := c.resource.FindAllByIdStringIn(ctx, ct.Ids)
		if !b {
			slice := make([]*entityRam.RamResourceRelationEntity, 0)
			//保存数据库
			for _, item := range infos {
				if !typeAttrPg.Resource.IsEqual(item.TypeAttr) {
					continue
				}
				if "" == item.Path {
					continue
				}
				if "" == item.Method {
					continue
				}
				var info entityRam.RamResourceRelationEntity
				copier.Copy(&info, &item)
				info.ID = 0
				info.CreateAt = &now
				info.TypeCategory = auth.TypeCategory
				info.Mark = auth.Mark
				info.AuthorityId = auth.ID
				info.TypeValue = ct.SourceId.ToString()
				info.ResourceId = item.ID
				log.Infof(ctx, log.TagAppDef, "info%+v", info)
				info.TenantNo = holder.GetTenantNo()
				c.resourceRel.Create(ctx, &info)
				log.Infof(ctx, log.TagAppDef, "save=%+v", info)

				slice = append(slice, &info)
			}
			//保存到 casbin
			c.casbin.UpdateCasbin(numberPg.Int64ToString(auth.ID), slice)
		} else {
			return rt.ErrorMessage("权限 不存在")
		}
	}

	return rt.Ok()
}

// Delete 删除
func (c *RamResourceAuthorizationService) Delete(ctx *gin.Context, code string) {
	//删除授权id
	info, result := c.resourceAuth.FindByMark(ctx, code)
	if result {
		c.resourceAuth.DeleteByMark(ctx, code)
		//删除 资源与授权id关系
		c.resourceRel.DeleteByAuthorityId(ctx, info.ID)
		//删除 casbin 内权限规则
		c.casbin.ClearCasbin(int(info.ID))
	}
}

// AuthorityValid code string 标记
//
// AuthorityValid 设置授权id 有效
func (c *RamResourceAuthorizationService) AuthorityValid(ctx *gin.Context, code string) {
	//删除授权id
	info, result := c.resourceAuth.FindByMark(ctx, code)
	if result {
		//设置有效
		c.resourceAuth.Update(ctx, entityRam.RamResourceAuthorityEntity{State: enumStatePg.ENABLE.Index()}, info.ID)
		//删除 资源与授权id关系
		infos := c.resourceRel.FindAll(ctx, entityRam.RamResourceRelationEntity{Mark: code})
		//保存到 casbin
		c.casbin.UpdateCasbin(numberPg.Int64ToString(info.ID), infos)
	}
}

// AuthorityInValid code string 标记
//
// AuthorityInValid 设置授权id 无效
func (c *RamResourceAuthorizationService) AuthorityInValid(ctx *gin.Context, code string) {
	//删除授权id
	info, result := c.resourceAuth.FindByMark(ctx, code)
	if result {
		//设置无效
		c.resourceAuth.Update(ctx, entityRam.RamResourceAuthorityEntity{State: enumStatePg.DISABLE.Index()}, info.ID)
		//保存到 casbin ，直接清空规则
		c.casbin.ClearCasbin(int(info.ID))
	}
}
