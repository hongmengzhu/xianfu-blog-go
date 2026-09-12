package ramAccount

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccount"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/appModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/slicePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

type Detail struct {
	sp *Sp `autowire:"?"`
	tp appModulePg.AppModule
}

func NewDetail(sp *Sp, tp appModulePg.AppModule) *Detail {
	return &Detail{
		sp: sp,
		tp: tp,
	}
}

// Process
//
//	@Description: 处理
//	@receiver c
//	@param ctx
//	@param ct
//	@param tp
//	@return rt
func (c *Detail) Process(ctx *gin.Context, id string) (rt rg.Rs[modRamAccount.DetailVo]) {
	if strPg.IsBlank(id) {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sp.accDb.FindByIdAndTypeDomain(ctx, numberPg.StrToInt64(id), c.tp.ToTypeDomain().String(), optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modRamAccount.DetailVo
	copier.Copy(&info, &find)
	//
	mapDep := make(map[string]*entityRam.RamDepartmentEntity)
	mapRole := make(map[string]*entityRam.RamRoleEntity)
	mapLevel := make(map[string]*entityRam.RamLevelEntity)
	mapGroup := make(map[string]*entityRam.RamGroupEntity)
	mapTeam := make(map[string]*entityRam.RamTeamEntity)
	idsDep := make([]string, 0)
	idsRole := make([]string, 0)
	idsLevel := make([]string, 0)
	idsGroup := make([]string, 0)
	idsTeam := make([]string, 0)
	//部门
	if nil != find.Os.Data().Departments && len(find.Os.Data().Departments) > 0 {
		for _, obj := range find.Os.Data().Departments {
			idsDep = append(idsDep, obj)
		}
	}
	if strPg.IsNotBlank(find.DepartmentNo) {
		idsDep = append(idsDep, find.DepartmentNo)
	}
	//角色
	if nil != find.Os.Data().Roles && len(find.Os.Data().Roles) > 0 {
		for _, obj := range find.Os.Data().Roles {
			idsRole = append(idsRole, obj)
		}
	}
	if strPg.IsNotBlank(find.RoleNo) {
		idsRole = append(idsRole, find.RoleNo)
	}
	//级别
	if nil != find.Os.Data().Levels && len(find.Os.Data().Levels) > 0 {
		for _, obj := range find.Os.Data().Levels {
			idsLevel = append(idsLevel, obj)
		}
	}
	if strPg.IsNotBlank(find.LevelNo) {
		idsLevel = append(idsLevel, find.LevelNo)
	}
	//分组
	if nil != find.Os.Data().Groups && len(find.Os.Data().Groups) > 0 {
		for _, obj := range find.Os.Data().Groups {
			idsGroup = append(idsGroup, obj)
		}
	}
	if strPg.IsNotBlank(find.GroupNo) {
		idsGroup = append(idsGroup, find.GroupNo)
	}
	//团队
	if nil != find.Os.Data().Teams && len(find.Os.Data().Teams) > 0 {
		for _, obj := range find.Os.Data().Teams {
			idsTeam = append(idsTeam, obj)
		}
	}
	//部门
	{
		if len(idsDep) > 0 {
			infos, result := c.sp.depDb.FindAllByNoIn(ctx, idsDep)
			if result {
				mapDep = slicePg.ToMap(infos, func(t *entityRam.RamDepartmentEntity) (string, *entityRam.RamDepartmentEntity) {
					return t.No, t
				})
			}
		}
	}
	//角色
	{
		if len(idsRole) > 0 {
			infos, result := c.sp.roleDb.FindAllByNoIn(ctx, idsRole)
			if result {
				mapRole = slicePg.ToMap(infos, func(t *entityRam.RamRoleEntity) (string, *entityRam.RamRoleEntity) {
					return t.No, t
				})
			}
		}
	}
	//级别
	{
		if len(idsLevel) > 0 {
			infos, result := c.sp.levelDb.FindAllByNoIn(ctx, idsLevel)
			if result {
				mapLevel = slicePg.ToMap(infos, func(t *entityRam.RamLevelEntity) (string, *entityRam.RamLevelEntity) {
					return t.No, t
				})
			}
		}
	}
	//分组
	{
		if len(idsGroup) > 0 {
			infos, result := c.sp.groupDb.FindAllByNoIn(ctx, idsGroup)
			if result {
				mapGroup = slicePg.ToMap(infos, func(t *entityRam.RamGroupEntity) (string, *entityRam.RamGroupEntity) {
					return t.No, t
				})
			}
		}
	}
	//分组
	{
		if len(idsTeam) > 0 {
			infos, result := c.sp.teamDb.FindAllByNoIn(ctx, idsTeam)
			if result {
				mapTeam = slicePg.ToMap(infos, func(t *entityRam.RamTeamEntity) (string, *entityRam.RamTeamEntity) {
					return t.No, t
				})
			}
		}
	}
	//
	//部门
	if nil != find.Os.Data().Departments && len(find.Os.Data().Departments) > 0 {
		info.Departments = find.Os.Data().Departments
		info.Os.No.Departments = find.Os.Data().Departments
		info.Os.NoName.Departments = make([]string, 0)
		for _, obj := range find.Os.Data().Departments {
			if get, ok := mapDep[obj]; ok {
				info.Os.NoName.Departments = append(info.Os.NoName.Departments, get.Name)
			}
		}
	}
	if strPg.IsNotBlank(find.DepartmentNo) {
		if obj, ok := mapDep[find.DepartmentNo]; ok {
			info.DepartmentNoName = obj.Name
		}
	}
	//角色
	if nil != find.Os.Data().Roles && len(find.Os.Data().Roles) > 0 {
		info.Roles = find.Os.Data().Roles
		info.Os.No.Roles = find.Os.Data().Roles
		info.Os.NoName.Roles = make([]string, 0)
		for _, obj := range find.Os.Data().Roles {
			if get, ok := mapRole[obj]; ok {
				info.Os.NoName.Roles = append(info.Os.NoName.Roles, get.Name)
			}
		}
	}
	//级别
	if nil != find.Os.Data().Levels && len(find.Os.Data().Levels) > 0 {
		info.Levels = find.Os.Data().Levels
		info.Os.No.Levels = find.Os.Data().Levels
		info.Os.NoName.Levels = make([]string, 0)
		for _, obj := range find.Os.Data().Levels {
			if get, ok := mapLevel[obj]; ok {
				info.Os.NoName.Levels = append(info.Os.NoName.Levels, get.Name)
			}
		}
	}
	if strPg.IsNotBlank(find.LevelNo) {
		if obj, ok := mapLevel[find.LevelNo]; ok {
			info.LevelNoName = obj.Name
		}
	}
	//分组
	if nil != find.Os.Data().Groups && len(find.Os.Data().Groups) > 0 {
		info.Groups = find.Os.Data().Groups
		info.Os.No.Groups = find.Os.Data().Groups
		info.Os.NoName.Groups = make([]string, 0)
		for _, obj := range find.Os.Data().Groups {
			if get, ok := mapGroup[obj]; ok {
				info.Os.NoName.Groups = append(info.Os.NoName.Levels, get.Name)
			}
		}
	}
	if strPg.IsNotBlank(find.GroupNo) {
		if obj, ok := mapLevel[find.GroupNo]; ok {
			info.GroupNoName = obj.Name
		}
	}
	//团队
	if nil != find.Os.Data().Teams && len(find.Os.Data().Teams) > 0 {
		info.Teams = find.Os.Data().Teams
		info.Os.No.Teams = find.Os.Data().Teams
		info.Os.NoName.Teams = make([]string, 0)
		for _, obj := range find.Os.Data().Teams {
			if get, ok := mapTeam[obj]; ok {
				info.Os.NoName.Teams = append(info.Os.NoName.Teams, get.Name)
			}
		}
	}
	return rt.OkData(info)
}
