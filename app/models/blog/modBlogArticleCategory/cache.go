package modBlogArticleCategory

type Cache struct {
	ID          int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string `gorm:"column:no;type:varchar(80);default:;comment:编号代号" json:"no" comment:"编号代号" `
	Code        string `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `
	Name        string `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	NameFl      string `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `
	NameFull    string `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `
	Description string `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	Sort        int64  `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	ParentNo    string `gorm:"column:parent_no;type:varchar(80);index;comment:上级" json:"parent_no" comment:"上级" `
	NoLink      string `gorm:"column:no_link;type:text;comment:上级" json:"no_link" comment:"上级链" `
	TypeSys     string `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	CodeLink    string `gorm:"column:code_link;type:text;comment:上级" json:"code_link" comment:"上级链" `
}
