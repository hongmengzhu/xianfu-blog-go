package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// LoadingLocation 加载位置
type LoadingLocation string

const (
	LoadingStartLoading            LoadingLocation = "startLoading"            //启动加载
	LoadingLoadAfterStartup        LoadingLocation = "loadAfterStartup"        //启动后加载
	LoadingLazyLoadingAfterStartup LoadingLocation = "lazyLoadingAfterStartup" //启动后延迟加载
	LoadingAppStartsLoading        LoadingLocation = "appStartsLoading"        //应用启动加载
)

// Name 名称
func (this LoadingLocation) Name() string {
	switch this {
	case "startLoading":
		return "启动加载"
	case "loadAfterStartup":
		return "启动后加载"
	case "lazyLoadingAfterStartup":
		return "启动后延迟加载"
	case "appStartsLoading":
		return "应用启动加载"
	default:
		return "未知"
	}
}

// 值
func (this LoadingLocation) String() string {
	return string(this)
}

// 值
func (this LoadingLocation) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this LoadingLocation) IsEqual(id string) bool {
	return string(this) == id
}

var LoadingLocationMap = map[string]enumBasePg.EnumString{
	LoadingStartLoading.String():            enumBasePg.EnumString{Id: LoadingStartLoading.String(), Name: LoadingStartLoading.Name()},
	LoadingLoadAfterStartup.String():        enumBasePg.EnumString{Id: LoadingLoadAfterStartup.String(), Name: LoadingLoadAfterStartup.Name()},
	LoadingLazyLoadingAfterStartup.String(): enumBasePg.EnumString{Id: LoadingLazyLoadingAfterStartup.String(), Name: LoadingLazyLoadingAfterStartup.Name()},
	LoadingAppStartsLoading.String():        enumBasePg.EnumString{Id: LoadingAppStartsLoading.String(), Name: LoadingAppStartsLoading.Name()},
}

func IsExistLoadingLocation(id string) (LoadingLocation, bool) {
	_, ok := LoadingLocationMap[id]
	return LoadingLocation(id), ok
}
