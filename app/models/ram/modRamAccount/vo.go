package modRamAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID                   typePg.Uint64String `json:"id" label:"id" `
	No                   string              `json:"no" label:"编号" `
	Name                 string              `json:"name" label:"名称" `
	RealName             string              `json:"realName" label:"真实姓名" `
	Account              string              `json:"account" label:"账户" `
	Cc                   string              `json:"cc" label:"国际区号"`
	Phone                string              `json:"phone" label:"手机号" `
	Mail                 string              `json:"mail" label:"邮箱" `
	Code                 string              `json:"code" label:"编码" `
	MailVerify           typePg.Int8         `json:"mailVerify" label:"邮箱验证1是2否" `
	PhoneVerify          typePg.Int8         `json:"phoneVerify" label:"手机验证1是2否" `
	AccountVerify        typePg.Int8         `json:"accountVerify" label:"账户验证1是2否" `
	State                typePg.Int8         `json:"state" label:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	RegisterTime         *typePg.Time        `json:"registerTime" label:"注册时间" `
	RegisterIP           string              `json:"registerIp" label:"注册ip" `
	LevelNo              string              `json:"levelNo" label:"级别id"`
	GroupNo              string              `json:"groupNo" label:"组id"`
	DepartmentNo         string              `json:"departmentNo" label:"主部门id"`
	Departments          []string            `json:"departments" label:"部门" `
	Roles                []string            `json:"roles" label:"角色" `
	Teams                []string            `json:"teams" label:"团队"`
	Levels               []string            `json:"levels" label:"级别"`
	Groups               []string            `json:"groups" label:"组" `
	TypeDomain           string              `json:"typeDomain" label:"域类型" `
	TypeIdentity         string              `json:"typeIdentity" label:"身份类型;普通;经理;副经理" `
	Description          string              `json:"description" label:"描述" `
	JobFunction          string              `json:"jobFunction" label:"职务" `
	Post                 string              `json:"post" label:"岗位" `
	Position             string              `json:"position" label:"职位" `
	JobTitle             string              `json:"jobTitle" label:"职衔" `
	JobRank              string              `json:"jobRank" label:"职级" `
	Avatar               string              `json:"avatar" label:"头像" `
	Birthday             *typePg.DateOnly    `json:"birthday" label:"生日" `
	Sex                  string              `json:"sex" label:"性别" `
	IdentityCode         string              `json:"identityCode" label:"身份编号 "`
	Os                   OsVo                `json:"os"`
	LoginTime            *typePg.Time        `json:"loginTime" comment:"登陆时间" `
	PrincipalTypeService typePg.Int8         `json:"principalTypeService" label:"主账号类型" `
	CreateAt             *typePg.Time        `json:"createAt" label:"创建时间" `
	//
	DepartmentNoName string `json:"departmentNoName" label:"主部门"`
	LevelNoName      string `json:"levelNoName" label:"级别"`
	GroupNoName      string `json:"groupNoName" label:"组"`
	PositionName     string `json:"positionName" label:"职位"`
	PostName         string `json:"postName" label:"岗位"`
	JobFunctionName  string `json:"jobFunctionName" label:"职务"`
}
