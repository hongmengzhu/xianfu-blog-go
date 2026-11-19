package entityIm

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"
)

type ImMessage struct {
	entityPg.BaseCrUpAtEntity
	SenderAno  string     `gorm:"column:sender_ano;type:varchar(80);index;default:'';comment:发送人编号" json:"sender_ano" comment:"发送人编号" `
	MsgType    string     `gorm:"column:msg_type;type:varchar(80);index;default:'';comment:消息类型文本|图片|文件|语音|系统通知|业务消息|指令消息" json:"msg_type" comment:"消息类型" `
	Content    string     `gorm:"column:content;type:varchar(255);default:;comment:内容" json:"content" comment:"内容" `
	ExtraInfo  string     `gorm:"column:extra_info;type:varchar(255);default:;comment:额外信息" json:"extra_info" comment:"额外信息" `
	ReCallIs   int8       `gorm:"column:re_call_is;not null;index;default:2;comment:是否撤回2否1是" json:"re_call_is" comment:"是否撤回" `
	RecallTime *time.Time `gorm:"column:recall_time;default:;comment:撤回时间" json:"recall_time" comment:"撤回时间" `
}

func (*ImMessage) TableName() string {
	return "im_message"
}

func (*ImMessage) TableComment() string {
	return "消息"
}
