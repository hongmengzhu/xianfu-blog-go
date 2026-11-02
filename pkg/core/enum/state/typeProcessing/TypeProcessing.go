package typeProcessing

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeProcessing 处理状态
type TypeProcessing string

const (
	TypeProcessingUnProcessed      TypeProcessing = "unProcessed"      //未处理
	TypeProcessingCompleted        TypeProcessing = "completed"        //处理完成
	TypeProcessingProcessing       TypeProcessing = "processing"       //处理中
	TypeProcessingProcessingFailed TypeProcessing = "processingFailed" //处理失败
	TypeProcessingSubmit           TypeProcessing = "submit"           //提交
)

// Name 名称
func (this TypeProcessing) Name() string {
	switch this {
	case "unProcessed":
		return "未处理"
	case "completed":
		return "处理完成"
	case "processing":
		return "处理中"
	case "processingFailed":
		return "处理失败"
	case "submit":
		return "提交"
	default:
		return "未知"
	}
}

// 值
func (this TypeProcessing) String() string {
	return string(this)
}

// 值
func (this TypeProcessing) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TypeProcessing) IsEqual(id string) bool {
	return string(this) == id
}

var TypeProcessingMap = map[string]enumBasePg.EnumString{
	TypeProcessingUnProcessed.String():      enumBasePg.EnumString{TypeProcessingUnProcessed.String(), TypeProcessingUnProcessed.Name()},
	TypeProcessingCompleted.String():        enumBasePg.EnumString{TypeProcessingCompleted.String(), TypeProcessingCompleted.Name()},
	TypeProcessingProcessingFailed.String(): enumBasePg.EnumString{TypeProcessingProcessingFailed.String(), TypeProcessingProcessingFailed.Name()},
	TypeProcessingProcessing.String():       enumBasePg.EnumString{TypeProcessingProcessing.String(), TypeProcessingProcessing.Name()},
	TypeProcessingSubmit.String():           enumBasePg.EnumString{TypeProcessingSubmit.String(), TypeProcessingSubmit.Name()},
}

func IsExistTypeProcessing(id string) (TypeProcessing, bool) {
	_, ok := TypeProcessingMap[id]
	return TypeProcessing(id), ok
}
