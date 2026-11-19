package entityIm

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"
)

type ImPushTask struct {
	entityPg.BaseCrUpAtEntity
	MsgNo      string `gorm:"column:im_message_no;type:varchar(80);index;default:;comment:消息编号" json:"im_message_no" comment:"消息编号" `
	Ano        string `gorm:"column:ano;type:varchar(80);index;default:;comment:接收用户编号" json:"ano" comment:"接收用户编号" `
	Title      string `gorm:"column:title;type:varchar(80);default:;comment:标题" json:"title" comment:"标题" `
	Content    string `gorm:"column:content;type:varchar(255);default:;comment:内容" json:"content" comment:"内容" `
	PushChanel string `gorm:"column:push_chanel;type:varchar(80);default:;comment:推送渠道apns/huawei/xiaomi/oppo/vivo" json:"push_chanel" comment:"推送渠道" `
	State      int8   `gorm:"column:state;not null;index;default:1;comment:状态|1成功|2失败" json:"state" comment:"状态:1启用;2禁用" `
	RetryCount int64  `gorm:"column:retry_count;type:int;index;default:0;comment:重试次数" json:"retry_count" comment:"重试次数" `
}
