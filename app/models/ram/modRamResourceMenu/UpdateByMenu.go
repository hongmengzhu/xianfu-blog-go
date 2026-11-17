package modRamResourceMenu

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type UpdateByMenuCt struct {
	MenuId typePg.Int64String   `json:"menuId" binding:"required" validate:"required" label:"菜单id" `
	Data   []ResourceAndGroupVo `json:"data" binding:"required" validate:"required" label:"数据" `
}
