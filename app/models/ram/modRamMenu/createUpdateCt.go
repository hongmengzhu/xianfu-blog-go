package modRamMenu

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID                     typePg.Int64String `json:"id" form:"id" label:"id" `
	OrgId                  string             `json:"orgId" label:"组织id" `                                                                             // 组织id
	Name                   string             `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl                 string             `json:"nameFl" label:"名称外文" `                                                                            // 名称外文
	Code                   string             `json:"code" form:"code"  label:"编号代号" `
	NameFull               string             `json:"nameFull" label:"全称" `    // 全称
	Description            string             `json:"description" label:"描述" ` // 描述
	ParentId               string             `json:"parentId" label:"上级" `
	ParentNo               string             `json:"parentNo" label:"上级编号" `
	TypeSys                string             `json:"typeSys" binding:"required" validate:"required" comment:"类型;普通;系统;api" `
	TypeAttr               string             `json:"typeAttr" binding:"required" validate:"required" label:"属性" `
	Path                   string             `json:"path" comment:"路由路径" `
	Method                 string             `json:"method" comment:"方法" `
	Show                   typePg.Int8        `json:"show" comment:"列表显示" `
	Component              string             `json:"component" comment:"对应前端文件路径" `
	ActiveName             string             `json:"activeName" comment:"高亮菜单" `
	KeepAlive              typePg.Int8        `json:"keepAlive" comment:"缓存" `
	Icon                   string             `json:"icon" comment:"菜单图标" `
	CloseTab               typePg.Int8        `json:"closeTab" comment:"关闭tab" `
	TypeMenu               string             `json:"typeMenu" comment:"菜单类型" `                // 菜单类型:目录|菜单|按钮|内嵌|外链
	AuthCode               string             `json:"authCode" comment:"权限标识" `                // 权限标识
	ActivePath             string             `json:"activePath" comment:"激活路径" `              // 激活路径
	MetaActiveIcon         string             `json:"metaActiveIcon" comment:"激活图标" `          // 激活图标
	MetaBadgeType          string             `json:"metaBadgeType" comment:"徽标类型" `           // 徽标类型
	MetaBadge              string             `json:"metaBadge" comment:"徽章内容" `               // 徽章内容
	MetaBadgeVariants      string             `json:"metaBadgeVariants" comment:"徽标样式" `       // 徽标样式
	MetaHideInMenu         string             `json:"metaHideInMenu" comment:"隐藏菜单" `          // 隐藏菜单
	MetaHideChildrenInMenu string             `json:"metaHideChildrenInMenu" comment:"隐藏子菜单" ` // 隐藏子菜单
	MetaHideInBreadcrumb   string             `json:"metaHideInBreadcrumb" comment:"在面包屑中隐藏" ` // 在面包屑中隐藏
	MetaHideInTab          string             `json:"metaHideInTab" comment:"在标签栏中隐藏" `        // 在标签栏中隐藏
	MetaTitle              string             `json:"metaTitle" comment:"标题" `                 // 标题
	MetaKeepAlive          string             `json:"metaKeepAlive" comment:"缓存标签页" `          // 缓存标签页
	MetaAffixTab           string             `json:"metaAffixTab" comment:"固定在标签" `           // 固定在标签
	LinkSrc                string             `json:"linkSrc" comment:"链接地址" `                 // 链接地址
	Api                    string             `json:"api" comment:"后端接口" `                     // 后端接口
	ApiMd5                 string             `json:"apiMd5" comment:"接口MD5" `                 // 接口MD5
}
