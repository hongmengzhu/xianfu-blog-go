package modTcAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type UpdateCt struct {
	ID           typePg.Uint64String `json:"id" form:"id" binding:"required" validate:"required" label:"id" `
	Name         string              `json:"name" label:"名称" `                                                                        // 名称
	RealName     string              `json:"realName" label:"真实姓名" `                                                                  // 名称
	Account      string              `json:"account"  binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"账户" ` // 账户
	Cc           string              `json:"cc" label:"国际区号"`                                                                         // 国际区号
	Phone        string              `json:"phone" label:"手机号" `                                                                      // 手机号
	Mail         string              `json:"mail" label:"邮箱" `                                                                        // 邮箱
	Code         string              `json:"code" label:"编码" `                                                                        // 编码
	RegisterTime *typePg.Time        `json:"registerTime" label:"注册时间" `                                                              // 注册时间
	LevelNo      string              `json:"levelNo" label:"级别" `
	GroupNo      string              `json:"groupNo" label:"组" `
	DepartmentNo string              `json:"departmentNo" label:"主部门id" `
	Departments  []string            `json:"departments" label:"部门" `
	Roles        []string            `json:"roles" label:"角色" `
	Teams        []string            `json:"teams" label:"团队"`
	Levels       []string            `json:"levels" label:"级别"`
	Groups       []string            `json:"groups" label:"组" `
	TypeDomain   string              `json:"typeDomain" label:"域类型" `
	TypeIdentity string              `json:"typeIdentity" label:"身份类型;普通;经理;副经理" ` // 身份类型;普通;经理;副经理
	Description  string              `json:"description" label:"描述" `              // 描述
	Position     string              `json:"position" label:"岗位" `
	Job          string              `json:"job" label:"职位" `
	JobTitle     string              `json:"jobTitle" label:"职衔" `
	JobRank      string              `json:"jobRank" label:"职级" `
	Avatar       string              `json:"avatar" label:"头像" `
	Birthday     *typePg.DateOnly    `json:"birthday" label:"生日" `
	Sex          string              `json:"sex" label:"性别:男;女;未知" `
	IdentityCode string              `json:"identityCode" label:"身份编号 "`

	//MailVerify       jsonPg.Int8        `json:"mailVerify" label:"邮箱验证1是2否" `                           // 邮箱验证1是2否
	//PhoneVerify      jsonPg.Int8        `json:"phoneVerify" label:"手机验证1是2否" `                          // 手机验证1是2否
	//AccountVerify    jsonPg.Int8        `json:"accountVerify" label:"账户验证1是2否" `                        // 账户验证1是2否
}
