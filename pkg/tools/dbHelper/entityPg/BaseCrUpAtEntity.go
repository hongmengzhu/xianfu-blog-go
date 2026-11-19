package entityPg

import "time"

type BaseCrUpAtEntity struct {
	ID       int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	CreateAt *time.Time `gorm:"column:create_at;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt *time.Time `gorm:"column:update_at;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	TenantNo string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	OrgNo    string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
}
