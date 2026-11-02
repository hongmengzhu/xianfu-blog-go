package attachmentTypePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// 附件集  类型
type Type string

const (
	AttachmentGeneral  Type = "general"  //普通
	AttachmentMain     Type = "main"     //主图
	AttachmentFirst    Type = "first"    //首图
	AttachmentCarousel Type = "carousel" //轮播图
	AttachmentVideo    Type = "video"    //视频
	AttachmentList     Type = "list"     //列表图
	AttachmentBanner   Type = "banner"   //banner
)

// Name 名称
func (this Type) Name() string {
	switch this {
	case "general":
		return "普通"
	case "main":
		return "主图"
	case "first":
		return "首图"
	case "carousel":
		return "轮播图"
	case "list":
		return "列表图"
	case "video":
		return "视频"
	case "banner":
		return "banner"
	default:
		return "未知"
	}
}

// 值
func (this Type) String() string {
	return string(this)
}

// 值
func (this Type) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Type) IsEqual(id string) bool {
	return string(this) == id
}

var TypeMap = map[string]enumBasePg.EnumString{
	AttachmentGeneral.String():  enumBasePg.EnumString{AttachmentGeneral.String(), AttachmentGeneral.Name()},
	AttachmentMain.String():     enumBasePg.EnumString{AttachmentMain.String(), AttachmentMain.Name()},
	AttachmentFirst.String():    enumBasePg.EnumString{AttachmentFirst.String(), AttachmentFirst.Name()},
	AttachmentCarousel.String(): enumBasePg.EnumString{AttachmentCarousel.String(), AttachmentCarousel.Name()},
	AttachmentList.String():     enumBasePg.EnumString{AttachmentList.String(), AttachmentList.Name()},
	AttachmentVideo.String():    enumBasePg.EnumString{AttachmentVideo.String(), AttachmentVideo.Name()},
	AttachmentBanner.String():   enumBasePg.EnumString{AttachmentBanner.String(), AttachmentBanner.Name()},
}

func IsExistType(id string) (Type, bool) {
	_, ok := TypeMap[id]
	return Type(id), ok
}
