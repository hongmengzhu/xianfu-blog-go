package ramResource

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/utilsRam"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceAuthority"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/resourceTypeCategoryPg"
	iamConstant2 "github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeAttrPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"go-spring.org/log"

	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

// UpdateByTypeValue
// @Description: 根据资源类型和资源组id修改资源组权限
type UpdateByTypeValue struct {
	roleDb          *repositoryRam.RamRoleRepository                  `autowire:"?"`
	authDb          *repositoryRam.RamResourceAuthorityRepository     `autowire:"?"`
	resDb           *repositoryRam.RamResourceRepository              `autowire:"?"`
	groupDb         *repositoryRam.RamResourceGroupRepository         `autowire:"?"`
	groupRelationDb *repositoryRam.RamResourceGroupRelationRepository `autowire:"?"`
	relationDb      *repositoryRam.RamResourceRelationRepository      `autowire:"?"`
	typeCategory    resourceTypeCategoryPg.ResourceTypeCategory
	ct              modRamResourceAuthority.UpdateByTypeValueCt
	ctx             *gin.Context
	//
	groupData     []*entityRam.RamResourceGroupEntity
	groupAuthData []*entityRam.RamResourceAuthorityEntity
	role          *entityRam.RamRoleEntity
}

// NewUpdateByTypeValue
//
//	@Description: 根据资源类型和资源组id修改资源组权限
//	@param log
//	@param roleDb
//	@param authDb
//	@param resDb
//	@param groupDb
//	@param typeCategory
//	@param ct
//	@param ctx
//	@return *UpdateByTypeValue
func NewUpdateByTypeValue(
	roleDb *repositoryRam.RamRoleRepository,
	authDb *repositoryRam.RamResourceAuthorityRepository,
	resDb *repositoryRam.RamResourceRepository,
	groupDb *repositoryRam.RamResourceGroupRepository,
	groupRelationDb *repositoryRam.RamResourceGroupRelationRepository,
	relationDb *repositoryRam.RamResourceRelationRepository,
	typeCategory resourceTypeCategoryPg.ResourceTypeCategory,
	ct modRamResourceAuthority.UpdateByTypeValueCt,
	ctx *gin.Context) *UpdateByTypeValue {
	return &UpdateByTypeValue{
		groupDb:         groupDb,
		authDb:          authDb,
		roleDb:          roleDb,
		groupRelationDb: groupRelationDb,
		relationDb:      relationDb,
		ctx:             ctx,
		typeCategory:    typeCategory,
		ct:              ct,
		resDb:           resDb,
	}
}

// Process
//
//	@Description:  处理
//	@receiver c
//	@return rt
func (c *UpdateByTypeValue) Process() (rt rg.Rs[string]) {
	if c.ct.Ids == nil || 0 == len(c.ct.Ids) {
		return rt.ErrorMessage("ids不能为空")
	}
	if strPg.IsBlank(c.ct.TypeValue) {
		return rt.ErrorMessage("资源组id不能为空")
	}
	ids := c.getIds()
	if ids.ErrorIs() {
		return ids
	}
	//角色
	if resourceTypeCategoryPg.Role.IsEqual(c.typeCategory.Index()) {
		process := c.roleProcess()
		if process.ErrorIs() {
			return rt.ErrorMessage(process.Message)
		}

		return c.saveResourceRelationByRole()
	}
	return rt.Ok()
}

// getIds
//
//	@Description: 数据
//	@receiver c
//	@return rt
func (c *UpdateByTypeValue) getIds() (rt rg.Rs[string]) {
	ids := make([]string, 0)
	for _, id := range c.ct.Ids {
		if strPg.IsNotBlank(id) {
			ids = append(ids, id)
		}
	}
	if len(ids) > 0 {
		result := false
		c.groupData, result = c.groupDb.FindAllByIdStringIn(c.ctx, ids)
		if !result {
			return rt.ErrorMessage("部分资源组匹配失败")
		}
		ids = make([]string, 0)
		for _, item := range c.groupData {
			ids = append(ids, numberPg.Int64ToString(item.ID))
		}
		//查询所有资	源组 权限
		c.groupAuthData, result = c.authDb.FindAllByTypeCategoryAndGroupIdStringIn(c.ctx, resourceTypeCategoryPg.Group.String(), ids)
		if !result {
			return rt.ErrorMessage("没有获取到任何资源组权限")
		}
		return rt.Ok()
	}
	return rt.ErrorMessage("资源组id不能为空")
}

// role
//
//	@Description: 角色权限
//	@receiver c
//	@return rt
func (c *UpdateByTypeValue) roleProcess() (rt rg.Rs[string]) {
	result := false
	c.role, result = c.roleDb.FindByIdString(c.ctx, c.ct.TypeValue)
	if !result {
		return rt.ErrorMessage("角色不存在")
	}

	//删除 角色 对应的资源组
	c.groupRelationDb.DeleteByTypeCategoryAndTypeValue(c.ctx, c.typeCategory.Index(), numberPg.Int64ToString(c.role.ID))
	//
	log.Infof(c.ctx, log.TagAppDef, "c.groupData=%+v", len(c.groupData))
	//插入数据
	for _, item := range c.groupData {
		//不是资源属性，跳过
		if !iamConstant2.Resource.IsEqual(item.TypeAttr) {
			continue
		}
		var info entityRam.RamResourceGroupRelationEntity
		info.Name = item.Name
		info.NameFl = item.NameFl
		info.NameFull = item.NameFull
		info.Code = item.Code
		info.Description = item.Description
		info.State = enumStatePg.ENABLE.Index()
		info.TypeCategory = c.typeCategory.String()
		info.TypeSys = typeSysPg.General.String()
		info.TypeDomain = item.TypeDomain
		info.GroupNo = item.No
		info.TypeValue = numberPg.Int64ToString(c.role.ID)
		info.Mark = utilsRam.ResourceAuthorityMark(c.typeCategory, info.TypeValue, numberPg.Int64ToString(item.ID))
		//已存在 则跳过
		_, b := c.groupRelationDb.FindByMark(c.ctx, info.Mark)
		if b {
			continue
		}
		//
		log.Infof(c.ctx, log.TagAppDef, "info%+v", info)
		err, _ := c.groupRelationDb.Create(c.ctx, &info)
		if err != nil {
			return rt.ErrorMessage("保存失败")
		}
		log.Infof(c.ctx, log.TagAppDef, "save=%+v", info)
	}
	return rt.Ok()
}

// saveResourceRelationByRole
//
//	@Description: 角色权限
//	@receiver c
//	@return rt
func (c *UpdateByTypeValue) saveResourceRelationByRole() (rt rg.Rs[string]) {
	log.Infof(c.ctx, log.TagAppDef, "c.groupAuthData=%+v", len(c.groupAuthData))
	if len(c.groupAuthData) > 0 {
		ids := make([]string, 0)
		for _, item := range c.groupAuthData {
			ids = append(ids, item.ResourceNo)
		}
		result := false
		data := make([]*entityRam.RamResourceEntity, 0)
		data, result = c.resDb.FindAllByNoIn(c.ctx, ids)
		if !result {
			return rt.ErrorMessage("没有获取到任何资源组权限")
		}
		//插入数据
		for _, item := range data {
			//不是资源属性，跳过
			if !iamConstant2.Resource.IsEqual(item.TypeAttr) {
				continue
			}
			var info entityRam.RamResourceRelationEntity
			info.Code = item.Code
			info.TypeAttr = item.TypeAttr
			info.TypeCategory = c.typeCategory.String()
			info.TypeSys = typeSysPg.General.String()
			info.TypeDomain = item.TypeDomain
			info.ResourceId = item.ID
			info.TypeValue = numberPg.Int64ToString(c.role.ID)
			info.Path = item.Path
			info.Method = item.Method
			info.Mark = utilsRam.ResourceRelationMark(c.typeCategory, info.TypeValue, numberPg.Int64ToString(info.ResourceId))
			//已存在 则跳过
			_, b := c.relationDb.FindByMark(c.ctx, info.Mark)
			if b {
				continue
			}
			//
			log.Debugf(c.ctx, log.TagAppDef, "info%+v", info)
			err, _ := c.relationDb.Create(c.ctx, &info)
			if err != nil {
				return rt.ErrorMessage("保存失败")
			}
			log.Infof(c.ctx, log.TagAppDef, "save=%+v", info)
		}
	}

	return rt.Ok()
}
