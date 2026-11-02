package menuTypePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// MenuType 菜单资源关系类型
type MenuType string

const (
	Group    MenuType = "group"    //资源组
	Resource MenuType = "resource" //资源
)

// Name 名称
func (this MenuType) Name() string {
	switch this {
	case "group":
		return "资源组"
	case "resource":
		return "资源"
	default:
		return "未知"
	}
}

// 值
func (this MenuType) String() string {
	return string(this)
}

// 值
func (this MenuType) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this MenuType) IsEqual(id string) bool {
	return string(this) == id
}

var MenuTypeMap = map[string]enumBasePg.EnumString{
	Resource.String(): enumBasePg.EnumString{Resource.String(), Resource.Name()},
	Group.String():    enumBasePg.EnumString{Group.String(), Group.Name()},
}

func IsExistMenuType(id string) (MenuType, bool) {
	_, ok := MenuTypeMap[id]
	return MenuType(id), ok
}
