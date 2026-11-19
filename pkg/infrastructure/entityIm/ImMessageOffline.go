package entityIm

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"
)

type ImMessageOffline struct {
	entityPg.BaseCrUpAtEntity
	Ano        string     `gorm:"column:ano;type:varchar(80);index;default:;comment:接收用户编号" json:"ano" comment:"接收用户编号" `
	MsgNo      string     `gorm:"column:im_message_no;type:varchar(80);index;default:;comment:消息编号" json:"im_message_no" comment:"消息编号" `
	Pulled     int8       `gorm:"column:pulled;type:smallint;index;default:0;comment:是否拉取2否1是" json:"pulled" comment:"是否拉取" `
	PulledTime *time.Time `gorm:"column:pulled_time;default:;comment:拉取时间" json:"pulled_time" comment:"拉取时间" `
}

func (*ImMessageOffline) TableName() string {
	return "im_message_offline"
}
func (*ImMessageOffline) TableComment() string {
	return "消息离线表"
}
