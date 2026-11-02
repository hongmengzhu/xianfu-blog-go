package basicModulePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Module 模块
type Module string

const (
	ModuleArea          Module = "area"           //省市区
	ModuleAccount       Module = "account"        //账号
	ModuleDepartment    Module = "department"     //部门
	ModuleCarousel      Module = "channel"        //渠道
	ModuleGrout         Module = "group"          //用户组
	ModuleLevel         Module = "level"          //级别
	ModuleShopGoods     Module = "shop:goods"     //商品
	ModuleShopSku       Module = "shop:sku"       //商品
	ModuleShopOperate   Module = "shop:operate"   //运营
	ModuleOperateCoupon Module = "operate:coupon" //运营-优惠券
)

// Name 名称
func (this Module) Name() string {
	switch this {
	case "area":
		return "省市区"
	case "account":
		return "账号"
	case "department":
		return "部门"
	case "channel":
		return "渠道"
	case "group":
		return "用户组"
	case "level":
		return "级别"
	default:
		return "未知"
	}
}

// 值
func (this Module) String() string {
	return string(this)
}

// 值
func (this Module) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Module) IsEqual(id string) bool {
	return string(this) == id
}

var ModuleMap = map[string]enumBasePg.EnumString{
	ModuleArea.String():          enumBasePg.EnumString{ModuleArea.String(), ModuleArea.Name()},
	ModuleAccount.String():       enumBasePg.EnumString{ModuleAccount.String(), ModuleAccount.Name()},
	ModuleDepartment.String():    enumBasePg.EnumString{ModuleDepartment.String(), ModuleDepartment.Name()},
	ModuleCarousel.String():      enumBasePg.EnumString{ModuleCarousel.String(), ModuleCarousel.Name()},
	ModuleLevel.String():         enumBasePg.EnumString{ModuleLevel.String(), ModuleLevel.Name()},
	ModuleGrout.String():         enumBasePg.EnumString{ModuleGrout.String(), ModuleGrout.Name()},
	ModuleOperateCoupon.String(): enumBasePg.EnumString{ModuleOperateCoupon.String(), ModuleOperateCoupon.Name()},
}

func IsExistModule(id string) (Module, bool) {
	_, ok := ModuleMap[id]
	return Module(id), ok
}
