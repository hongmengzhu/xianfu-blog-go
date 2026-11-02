package clientPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Client 客户端
type Client string

const (
	Desktop               Client = "desktop"
	Browser               Client = "browser"
	ClientDesktopApp      Client = "desktopApp"            //桌面app 指的是运行在桌面操作系统（如windows、macos、ubuntu）上的客户端应用
	ClientOfWeb           Client = "clientOfWeb"           //桌面web 指的是运行在桌面浏览器（如chrome、safari、firefox）中的超文本内容（网页）
	ClientOfMobileApp     Client = "clientOfMobileApp"     //移动app 指的是运行在移动设备操作系统（如iOS、windows mobile、android）上的客户端应用；
	ClientOfMobileWeb     Client = "clientOfMobileWeb"     //移动web 指的是运行在移动设备浏览器（如Chrome for Android, iOS Safari）上的超文本内容（移动网页）
	ClientOfMobileAppMini Client = "clientOfMobileAppMini" //移动app内 小程序
	ClientDefault         Client = "default"               //默认
	ClientAll             Client = "all"                   //默认
)

// Name 名称
func (this Client) Name() string {
	switch this {
	case "browser":
		return "浏览器"
	case "desktop":
		return "桌面"
	case "desktopApp":
		return "桌面app"
	case "clientOfWeb":
		return "桌面web"
	case "clientOfMobileApp":
		return "移动app"
	case "clientOfMobileWeb":
		return "移动web"
	case "default":
		return "默认"
	case "all":
		return "全部"
	default:
		return "未知"
	}
}

// 值
func (this Client) String() string {
	return string(this)
}

// 值
func (this Client) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Client) IsEqual(id string) bool {
	return string(this) == id
}

var ClientMap = map[string]enumBasePg.EnumString{
	Browser.String():               enumBasePg.EnumString{Browser.String(), Browser.Name()},
	Desktop.String():               enumBasePg.EnumString{Desktop.String(), Desktop.Name()},
	ClientDesktopApp.String():      enumBasePg.EnumString{ClientDesktopApp.String(), ClientDesktopApp.Name()},
	ClientOfWeb.String():           enumBasePg.EnumString{ClientOfWeb.String(), ClientOfWeb.Name()},
	ClientOfMobileApp.String():     enumBasePg.EnumString{ClientOfMobileApp.String(), ClientOfMobileApp.Name()},
	ClientOfMobileWeb.String():     enumBasePg.EnumString{ClientOfMobileWeb.String(), ClientOfMobileWeb.Name()},
	ClientOfMobileAppMini.String(): enumBasePg.EnumString{ClientOfMobileAppMini.String(), ClientOfMobileAppMini.Name()},
	ClientDefault.String():         enumBasePg.EnumString{ClientDefault.String(), ClientDefault.Name()},
	ClientAll.String():             enumBasePg.EnumString{ClientAll.String(), ClientAll.Name()},
}

func IsExistClient(id string) (Client, bool) {
	_, ok := ClientMap[id]
	return Client(id), ok
}
