package entityIm

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"
)

type ImMessageReceipt struct {
	entityPg.BaseCrUpAtEntity
	MsgNo         string     `gorm:"column:im_message_no;type:varchar(80);index;default:;comment:消息编号" json:"im_message_no" comment:"消息编号" `
	Ano           string     `gorm:"column:ano;type:varchar(80);index;default:;comment:接收用户编号" json:"ano" comment:"接收用户编号" `
	DeliverStatus int8       `gorm:"column:deliver_status;type:smallint;index;default:0;comment:投递状态0待投递11投递中1投递成功2投递失败" json:"deliver_status" comment:"投递状态" `
	ReadStatus    int8       `gorm:"column:read_status;type:smallint;index;default:0;comment:阅读状态0待阅读11阅读中1阅读成功2阅读失败" json:"read_status" comment:"阅读状态" `
	DeliverTime   *time.Time `gorm:"column:deliver_time;default:;comment:投递时间" json:"deliver_time" comment:"投递时间" `
	ReadTime      *time.Time `gorm:"column:read_time;default:;comment:阅读时间" json:"read_time" comment:"阅读时间" `
}
