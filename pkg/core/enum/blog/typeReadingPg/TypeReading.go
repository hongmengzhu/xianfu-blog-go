package typeReadingPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeReading 阅读类型||未看|在看|已看
type TypeReading string

const (
	UNREAD    TypeReading = "unread"
	READING   TypeReading = "reading"
	COMPLETED TypeReading = "completed"
)

func (this TypeReading) Name() string {
	switch this {
	case "unread":
		return "未看"
	case "reading":
		return "在看"
	case "completed":
		return "已看"
	default:
		return "未知"
	}
}
func (this TypeReading) String() string {
	return string(this)
}

func (this TypeReading) Index() string {
	return string(this)
}

var MapTypeReading = map[string]enumBasePg.EnumString{
	UNREAD.String():    enumBasePg.EnumString{UNREAD.String(), UNREAD.Name()},
	READING.String():   enumBasePg.EnumString{READING.String(), READING.Name()},
	COMPLETED.String(): enumBasePg.EnumString{COMPLETED.String(), COMPLETED.Name()},
}

func IsExistTypeReading(id string) bool {
	_, ok := MapTypeReading[id]
	return ok
}
