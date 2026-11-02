package modTcAccount

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	TenantNo          string             `json:"tenantNo" label:"租户" `
	Name              string             `json:"name" label:"名称" `    // 名称
	Account           string             `json:"account" label:"账户" ` // 账户
	Code              string             `json:"code" label:"码值" `
	Wd                string             `json:"wd" label:"关键词" `         // 关键词
	Description       string             `json:"description" label:"描述" ` // 描述
	Sex               string             `json:"sex" label:"性别" `
	State             typePg.Int8        `json:"state" label:"状态:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	CreateAt          *time.Time         `json:"createAt" label:"创建时间" `
	DepartmentNo      string             `json:"departmentNo" label:"主部门id" `
	Departments       []string           `json:"departments" label:"部门" `
	Roles             []string           `json:"roles" label:"角色" `
	Levels            []string           `json:"levels" label:"级别" `
	Teams             []string           `json:"teams" label:"团队" `
	Groups            []string           `json:"groups" label:"组" `
	LevelNo           string             `json:"levelNo" label:"级别" `
	GroupNo           string             `json:"groupNo" label:"组" `
	RoleNo            string             `json:"roleNo" label:"角色" `
	Job               string             `json:"job" label:"职位" `
	Position          string             `json:"position" label:"岗位" `
	RegisterTimeRange []*typePg.Time     `json:"registerTimeRange" label:"注册时间" `
	LoginTimeRange    []*typePg.Time     `json:"loginTimeRange" label:"登陆时间" `
	BirthdayRange     []*typePg.DateOnly `json:"birthdayRange" label:"生日" `
}
