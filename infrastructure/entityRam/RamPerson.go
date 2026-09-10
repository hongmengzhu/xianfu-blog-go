package entityRam

import (
	"time"
)

// RamPersonEntity 用户
type RamPersonEntity struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:id" json:"id" comment:"" `
	CreateAt         *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt         *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                 // 更新时间
	CreateBy         int64      `gorm:"column:create_by;type:bigint;not null;index;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy         int64      `gorm:"column:update_by;type:bigint;not null;comment:更新人" json:"update_by" comment:"更新人" `
	No               string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	TenantNo         string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo            string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo          string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	CountryCode      string     `gorm:"column:country_code;type:varchar(20);comment:国际区号" json:"country_code" comment:"国际区号" `                                     // 国际区号
	Phone            string     `gorm:"column:phone;type:varchar(50);index;comment:手机号" json:"phone" comment:"手机号" `                                               // 手机号
	Mail             string     `gorm:"column:mail;type:varchar(255);index;comment:邮箱" json:"mail" comment:"邮箱" `                                                  // 邮箱
	Code             string     `gorm:"column:code;type:varchar(70);index;comment:编码" json:"code" comment:"编码" `                                                   // 编码
	MailVerify       int64      `gorm:"column:mail_verify;type:int8;not null;default:2;comment:邮箱验证1是2否" json:"mail_verify" comment:"邮箱验证1是2否" `                   // 邮箱验证1是2否
	PhoneVerify      int64      `gorm:"column:phone_verify;type:int8;not null;default:2;comment:手机验证1是2否" json:"phone_verify" comment:"手机验证1是2否" `                 // 手机验证1是2否
	State            int8       `gorm:"column:state;index;default:1;comment:启用1是2否" json:"state" comment:"启用1是2否" `                                                // 启用1是2否
	RegisterTime     *time.Time `gorm:"column:register_time;type:datetime;not null;default:0001-01-01 00:00:00;comment:注册时间" json:"register_time" comment:"注册时间" ` // 注册时间
	RegisterIP       string     `gorm:"column:register_ip;type:varchar(100)" json:"register_ip;comment:注册ip" comment:"注册ip" `                                      // 注册ip
	LoginTime        *time.Time `gorm:"column:login_time;type:datetime;default:0001-01-01 00:00:00;comment:登陆时间" json:"login_time" comment:"登陆时间" `                // 登陆时间
	RoleId           int64      `gorm:"column:role_id;type:bigint;not null;index;comment:角色id" json:"role_id" comment:"角色id" `                                     // 角色id
	LevelID          int64      `gorm:"column:level_id;type:bigint;not null;index;comment:级别id" json:"level_id" comment:"级别id" `                                   // 级别id
	GroupID          int64      `gorm:"column:group_id;type:bigint;not null;index;comment:组id" json:"group_id" comment:"组id" `                                     // 组id
	DepartmentIDMain int64      `gorm:"column:department_id_main;type:bigint;not null;index;comment:主部门id" json:"department_id_main" comment:"主部门id" `             // 主部门id
	DepartmentIds    string     `gorm:"column:department_ids;type:text;comment:部门" json:"department_ids" comment:"部门" `                                            // 部门
	TeamIds          string     `gorm:"column:team_ids;type:text;comment:团队" json:"team_ids" comment:"团队" `                                                        // 团队
	TypeDomain       int64      `gorm:"column:type_domain;type:int4;not null;index;comment:域类型" json:"type_domain" comment:"域类型" `                                 // 域类型
	TypeIdentity     int64      `gorm:"column:type_identity;type:int4;not null;index;comment:身份类型\\;普通\\;经理\\;副经理" json:"type_identity" comment:"身份类型;普通;经理;副经理" ` // 域类型
	Position         string     `gorm:"column:position;type:varchar(100);;comment:职务" json:"position" comment:"职务" `                                               // 域类型
	Description      string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	RealName         string     `gorm:"column:real_name;type:varchar(255);comment:真实姓名" json:"realName" comment:"真实姓名" `
	Avatar           string     `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar" comment:"头像" `
	Roles            string     `gorm:"column:roles;type:text;comment:岗位" json:"roles" comment:"岗位" `
	Birthday         *time.Time `gorm:"column:birthday;type:date;comment:生日" json:"birthday" comment:"生日" `
	Sex              int64      `gorm:"column:sex;type:int8;index;default:1;comment:性别:男;女;未知" json:"sex" comment:"性别:男;女;未知" `
}

// TableName RamPersonEntity's table name
func (*RamPersonEntity) TableName() string {
	return "ram_person"
}

func (*RamPersonEntity) TableComment() string {
	return "用户"
}
