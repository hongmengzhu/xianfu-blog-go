package entityRam

import (
	"time"

	"gorm.io/datatypes"
)

// RamAccountEntity 账户
type RamAccountEntity struct {
	ID            int64                                       `gorm:"column:id;primaryKey;autoIncrement:true;comment:" json:"id" comment:"" `
	CreateAt      *time.Time                                  `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt      *time.Time                                  `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间;comment:更新时间" json:"update_at" comment:"更新时间" `                    // 更新时间
	CreateBy      string                                      `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy      string                                      `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	No            string                                      `gorm:"column:no;type:varchar(70);index;default:;comment:编码" json:"no" comment:"编码" `
	Code          string                                      `gorm:"column:code;type:varchar(255);index;default:;comment:标志" json:"code" comment:"标志" `
	TenantNo      string                                      `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	StoreNo       string                                      `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	OrgNo         string                                      `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	OwnerNo       string                                      `gorm:"column:owner_no;type:varchar(80);index;default:;comment:所有者编号" json:"owner_no" comment:"所有者编号"`
	TypeDomain    string                                      `gorm:"column:type_domain;type:varchar(80);index;default:'general';comment:域类型" json:"type_domain" comment:"域类型系统-商户" `
	TypeIdentity  string                                      `gorm:"column:type_identity;type:varchar(80);index;default:'general';comment:身份类型|普通|经理|副经理" json:"type_identity" comment:"身份类型;普通;经理;副经理" `
	Account       string                                      `gorm:"column:account;type:varchar(255);index;comment:账户" json:"account" comment:"账户" `
	AccountMd5    string                                      `gorm:"column:account_md5;type:varchar(64);comment:账户md5" json:"account_md5" comment:"账户md5" `
	Cc            string                                      `gorm:"column:cc;type:varchar(20);index;default:;comment:国际区号" json:"cc" comment:"国际区号" `
	Phone         string                                      `gorm:"column:phone;type:varchar(50);index;default:;comment:手机号" json:"phone" comment:"手机号" `
	PhoneMd5      string                                      `gorm:"column:phone_md5;type:varchar(64);comment:手机号" json:"phone_md5" comment:"手机号" `
	Mail          string                                      `gorm:"column:mail;type:varchar(255);index;comment:邮箱" json:"mail" comment:"邮箱" `
	MailMd5       string                                      `gorm:"column:mail_md5;type:varchar(64);comment:" json:"mail_md5" comment:"" `
	MailVerify    int8                                        `gorm:"column:mail_verify;type:int8;not null;default:2;comment:邮箱验证1是2否" json:"mail_verify" comment:"邮箱验证1是2否" `                                  // 邮箱验证1是2否
	PhoneVerify   int8                                        `gorm:"column:phone_verify;type:int8;not null;default:2;comment:手机验证1是2否" json:"phone_verify" comment:"手机验证1是2否" `                                // 手机验证1是2否
	AccountVerify int8                                        `gorm:"column:account_verify;type:int8;not null;default:2;comment:账户验证1是2否" json:"account_verify" comment:"账户验证1是2否" `                            // 账户验证1是2否
	State         int8                                        `gorm:"column:state;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	RegisterTime  *time.Time                                  `gorm:"column:register_time;type:datetime;not null;default:0001-01-01 00:00:00;comment:注册时间" json:"register_time" comment:"注册时间" `                // 注册时间
	RegisterIP    string                                      `gorm:"column:register_ip;type:varchar(100);comment:注册ip" json:"register_ip" comment:"注册ip" `                                                     // 注册ip
	LoginTime     *time.Time                                  `gorm:"column:login_time;type:datetime;default:;comment:登陆时间" json:"login_time" comment:"登陆时间" `                                                  // 登陆时间
	Description   string                                      `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	Os            datatypes.JSONType[RamAccountJsonOs]        `gorm:"column:os;type:jsonb;index;default:'{}';comment:组织架构" json:"os" comment:"组织架构" `
	RoleNo        string                                      `gorm:"column:role_no;type:varchar(80);index;default:;comment:角色编号" json:"role_no" comment:"角色编号" `
	LevelNo       string                                      `gorm:"column:level_no;type:varchar(80);index;default:;comment:级别编号" json:"level_no" comment:"级别编号" `
	GroupNo       string                                      `gorm:"column:group_no;type:varchar(80);index;default:;comment:组编号" json:"group_no" comment:"组编号" `
	DepartmentNo  string                                      `gorm:"column:department_no;type:varchar(80);index;default:;comment:部门编号" json:"department_no" comment:"部门编号" `
	JobFunction   string                                      `gorm:"column:job_function;type:varchar(100);;comment:职务" json:"job_function" comment:"职务" `
	JobTitle      string                                      `gorm:"column:job_title;type:varchar(100);;comment:职衔" json:"job_title" comment:"职衔" `
	JobRank       string                                      `gorm:"column:job_rank;type:varchar(100);;comment:职级" json:"job_rank" comment:"职级" `
	Position      string                                      `gorm:"column:position;type:varchar(100);;comment:职位" json:"position" comment:"职位" `
	Post          string                                      `gorm:"column:post;type:varchar(100);;comment:岗位" json:"post" comment:"岗位" `
	Name          string                                      `gorm:"column:name;type:varchar(255);comment:名称" json:"realName" comment:"名称" `
	RealName      string                                      `gorm:"column:real_name;type:varchar(255);comment:真实姓名" json:"real_name" comment:"真实姓名" `
	IdentityCode  string                                      `gorm:"column:identity_code;type:varchar(255);comment:身份编号" json:"identity_code" comment:"身份编号" `
	Avatar        string                                      `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar" comment:"头像" `
	Birthday      *time.Time                                  `gorm:"column:birthday;type:date;comment:生日" json:"birthday" comment:"生日" `
	Sex           string                                      `gorm:"column:sex;type:varchar(80);index;default:;comment:性别男女未知" json:"sex" comment:"性别男女未知" `
	Founder       int8                                        `gorm:"column:founder;type:int8;index;default:2;comment:创始人1是2否" json:"founder" comment:"创始人1是2否" `
	ExtraData     datatypes.JSON                              `gorm:"column:extra_data;type:jsonb;comment:额外数据" json:"extraData" comment:"额外数据" `
	ExtraCond     datatypes.JSONType[RamAccountJsonExtraCond] `gorm:"column:extra_cond;type:jsonb;index;default:'{}';comment:扩展搜索条件" json:"extraCond" comment:"扩展搜索条件" `
	AppNo         string                                      `gorm:"column:app_no;type:varchar(80);index;default:;comment:应用编号" json:"appNo" comment:"应用编号" `
	SystemNo      string                                      `gorm:"column:system_no;type:varchar(80);index;default:;comment:接入业务系统编号(如:purchase_platform)" json:"systemNo" comment:"接入业务系统编号"`
}

// TableName RamAccount's table name
func (*RamAccountEntity) TableName() string {
	return "ram_account"
}

func (*RamAccountEntity) TableComment() string {
	return "账户"
}
