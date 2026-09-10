package holderPg

import (
	"time"

	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/interfaces"
)

type AccountHolder struct {
	ID            int64           `json:"id" comment:"" `
	No            string          `json:"no" comment:"编码"`
	Code          string          `json:"code" comment:"编码" `
	CreateAt      time.Time       `json:"createAt" comment:"创建时间" `
	UpdateAt      time.Time       `json:"updateAt" comment:"更新时间" `
	CreateBy      string          `json:"createBy" comment:"创建人" `
	UpdateBy      string          `json:"updateBy" comment:"更新人" `
	TenantNo      string          `json:"tenantNo" comment:"租户编码"`
	OrgNo         string          `json:"orgNo" comment:"组织编码"`
	OwnerNo       string          `json:"ownerNo" comment:"所有者"`
	Name          string          `json:"name" comment:"名称" `
	Account       string          `json:"account" comment:"账户" `
	AccountMd5    string          `json:"accountMd5" comment:"账户md5" `
	Cc            string          `json:"cc" comment:"国际区号" `
	Phone         string          `json:"phone" comment:"手机号" `
	PhoneMd5      string          `json:"phoneMd5" comment:"手机号" `
	Mail          string          `json:"mail" comment:"邮箱" `
	MailMd5       string          `json:"mailMd5" comment:"" `
	MailVerify    int64           `json:"mailVerify" comment:"邮箱验证1是2否" `
	PhoneVerify   int64           `json:"phoneVerify" comment:"手机验证1是2否" `
	AccountVerify int64           `json:"accountVerify" comment:"账户验证1是2否" `
	State         int64           `json:"state" comment:"启用1是2否" `
	RegisterTime  time.Time       `json:"registerTime" comment:"注册时间" `
	RegisterIP    string          `json:"registerIp" comment:"注册ip" `
	LoginTime     time.Time       `json:"loginTime" comment:"登陆时间" `
	RoleNo        string          `json:"roleNo" comment:"角色id"`
	LevelNo       string          `json:"levelNo" comment:"级别id"`
	GroupNo       string          `json:"groupNo" comment:"组id"`
	DepartmentNo  string          `json:"departmentNo" comment:"主部门id"`
	TypeDomain    string          `json:"typeDomain" comment:"域类型" `
	Os            AccountHolderOs `json:"os" comment:"组织架构" `
}

func NewAccountHolder() *AccountHolder {
	return new(AccountHolder)
}
