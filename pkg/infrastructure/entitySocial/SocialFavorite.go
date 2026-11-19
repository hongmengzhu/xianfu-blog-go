package entitySocial

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"

type SocialFavorite struct {
	entityPg.BaseCrUpAtEntity
	Ano        string `gorm:"column:ano;type:varchar(80);index;default:'';comment:收藏人编号" json:"ano" comment:"收藏人编号" `
	TargetType string `gorm:"column:target_type;type:varchar(80);index;default:'';comment:目标类型动态|评论|商品" json:"target_type" comment:"目标类型" `
	TargetNo   string `gorm:"column:target_no;type:varchar(80);index;default:'';comment:目标编号" json:"target_no" comment:"目标编号" `
	State      int8   `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用|11审核中" json:"state" comment:"状态:1启用;2禁用" `
	FolderNo   string `gorm:"column:folder_no;type:varchar(80);index;default:'';comment:收藏夹编号" json:"folder_no" comment:"收藏夹编号" `
}
