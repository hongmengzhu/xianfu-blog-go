package resourceTypeCategoryPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// ResourceTypeCategory 资源 权限类别
type ResourceTypeCategory string

const (
	Menu       ResourceTypeCategory = "menu"       //菜单
	Group      ResourceTypeCategory = "group"      //资源组
	Role       ResourceTypeCategory = "role"       //角色
	Department ResourceTypeCategory = "department" //部门
)

// Name 名称
func (this ResourceTypeCategory) Name() string {
	switch this {
	case "menu":
		return "菜单"
	case "group":
		return "资源组"
	case "role":
		return "角色"
	case "department":
		return "部门"
	default:
		return "未知"
	}
}

// 值
func (this ResourceTypeCategory) String() string {
	return string(this)
}

// 值
func (this ResourceTypeCategory) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this ResourceTypeCategory) IsEqual(id string) bool {
	return string(this) == id
}

var ResourceAuthorityTypeMap = map[string]enumBasePg.EnumString{
	Menu.String():       enumBasePg.EnumString{Menu.String(), Menu.Name()},
	Group.String():      enumBasePg.EnumString{Group.String(), Group.Name()},
	Role.String():       enumBasePg.EnumString{Role.String(), Role.Name()},
	Department.String(): enumBasePg.EnumString{Department.String(), Department.Name()},
}

func IsExistResourceAuthorityType(id string) (ResourceTypeCategory, bool) {
	_, ok := ResourceAuthorityTypeMap[id]
	return ResourceTypeCategory(id), ok
}
