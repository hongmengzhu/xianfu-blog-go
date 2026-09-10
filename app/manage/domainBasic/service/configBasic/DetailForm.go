package configBasic

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/basic/model/modCacheBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigList"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/rdsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/basic/key/basicEventKey"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
)

type DetailForm struct {
	Sp *Sp `autowire:"?"`

	rdt *rdsPg.BatchString `autowire:"?"`
}

func NewDetailForm(sp *Sp) *DetailForm {
	return &DetailForm{
		Sp:  sp,
		rdt: sp.rdt,
	}
}

// Process
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param ct
//	@return rt
func (c *DetailForm) Process(ctx *gin.Context, ct modBasicConfigList.DetailFormCt) (rt rg.Rs[modBasicConfigList.DetailFormVo]) {
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("参数错误")
	}
	holder := holderPg.GetContextAccount(ctx)
	config, b2 := c.Sp.repConfigList.FindByEventNo(ctx, ct.EventNo)
	if !b2 {
		return rt.ErrorMessage("配置不存在")
	}
	if holder.GetTenantNo() != config.TenantNo {
		return rt.ErrorMessage("配置不存在")
	}
	//info, b := c.Sp.repEvent.FindByNo(ctx, ct.EventNo)
	//if !b {
	//	return rt.ErrorMessage("事件不存在")
	//}
	//if holder.GetTenantNo() != info.TenantNo {
	//	return rt.ErrorMessage("事件不存在")
	//}
	data := modBasicConfigList.DetailFormVo{
		Form: modBasicConfigEvent.ModelForm{
			Item: make([]modBasicConfigEvent.ModelFormItem, 0),
		},
		Data: make(map[string]interface{}),
	}
	key2 := basicEventKey.EventTenantNoAllFields(config.TenantNo, config.EventNo)
	all, b2 := c.rdt.HGetAll(ctx, key2)
	if b2 {
		tmp := make([]modCacheBasicEvent.FieldCache, 0)
		for _, v := range all {
			var obj modCacheBasicEvent.FieldCache
			err := json.Unmarshal([]byte(v), &obj)
			if err != nil {
				log.Errorf(ctx, log.TagAppDef, "json.Unmarshal.err:%+v", err)
			} else {
				copier.Copy(&obj, v)
				//
				tmp = append(tmp, obj)
			}
		}
		//排序
		slices.SortFunc(tmp, func(a, b modCacheBasicEvent.FieldCache) int {
			if a.Sort < b.Sort {
				return -1
			} else if a.Sort > b.Sort {
				return 1
			}
			return 0
		})
		//遍历
		for _, v := range tmp {
			var obj2 modBasicConfigEvent.ModelFormItem
			copier.Copy(&obj2, &v)
			//
			data.Form.Item = append(data.Form.Item, obj2)
		}
	}
	info, result := c.Sp.repConfig.FindByEventNo(ctx, config.EventNo)
	if result {
		for _, v := range info {
			if yesNoIntPg.Yes.IsEqual(v.Binary) {
				data.Data[v.Field] = v.ValueBinary
			} else {
				data.Data[v.Field] = v.Value
			}
		}
	}
	rt.Data = data
	return rt.Ok()
}
