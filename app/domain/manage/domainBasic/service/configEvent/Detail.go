package configEvent

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigEvent"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
)

type Detail struct {
	Sp *Sp `autowire:"?"`
}

func NewDetail(sp *Sp) *Detail {
	return &Detail{
		Sp: sp,
	}
}

func (c *Detail) Process(ctx *gin.Context, id string) (rt rg.Rs[modBasicConfigEvent.CreateUpdateCt]) {
	var vo modBasicConfigEvent.CreateUpdateCt
	vo.BodyDelIds = make([]string, 0)
	if strPg.IsBlank(id) {
		return rt.ErrorMessage("模型ID不能为空")
	}
	info, result := c.Sp.repEvent.FindByIdString(ctx, id)
	if !result {
		return rt.ErrorMessage("模型不存在")
	}
	err := copier.Copy(&vo.Header, info)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	//
	vo.Body = make([]modBasicConfigEvent.ItemVo, 0)
	//
	data, r := c.Sp.repEventField.FindAllByModelNo(ctx, info.No)
	if r {
		for _, item := range data {
			var obj modBasicConfigEvent.ItemVo
			copier.Copy(&obj, item)
			//
			vo.Body = append(vo.Body, obj)
		}
	}
	rt.Data = vo
	return rt.Ok()
}
