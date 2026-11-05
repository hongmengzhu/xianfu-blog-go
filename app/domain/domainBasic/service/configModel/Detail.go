package configModel

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigModel"
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

func (c *Detail) Process(ctx *gin.Context, id string) (rt rg.Rs[modBasicConfigModel.CreateUpdateCt]) {
	var vo modBasicConfigModel.CreateUpdateCt
	vo.BodyDelIds = make([]string, 0)
	if strPg.IsBlank(id) {
		return rt.ErrorMessage("模型ID不能为空")
	}
	info, result := c.Sp.repModel.FindByIdString(ctx, id)
	if !result {
		return rt.ErrorMessage("模型不存在")
	}
	err := copier.Copy(&vo.Header, info)
	if err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	// 模型级规则
	if rules, b := c.Sp.repRule.FindAllByModelNo(ctx, info.No); b {
		for _, item := range rules {
			var obj modBasicConfigModel.RuleItemVo
			if err := copier.Copy(&obj, item); err != nil {
				log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
			}
			obj.RuleTarget = item.RuleTarget.Data()
			vo.Header.RuleItems = append(vo.Header.RuleItems, obj)
			vo.Header.Rules = append(vo.Header.Rules, item.RuleMode)
		}
	}
	// 字段规则
	fieldRuleMap := make(map[string][]modBasicConfigModel.RuleItemVo)
	if rules, b := c.Sp.repFieldRules.FindAllByValueNo(ctx, info.No); b {
		for _, item := range rules {
			var obj modBasicConfigModel.RuleItemVo
			if err := copier.Copy(&obj, item); err != nil {
				log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
			}
			obj.RuleTarget = item.RuleTarget.Data()
			fieldRuleMap[item.FieldNo] = append(fieldRuleMap[item.FieldNo], obj)
		}
	}
	//
	vo.Body = make([]modBasicConfigModel.ItemVo, 0)
	//
	data, r := c.Sp.repField.FindAllByModelNo(ctx, info.No)
	if r {
		for _, item := range data {
			var obj modBasicConfigModel.ItemVo
			copier.Copy(&obj, item)
			obj.Rules = item.Rules.Data()
			obj.SharedScope = item.SharedScope.Data()
			obj.RuleItems = fieldRuleMap[item.No]
			//
			vo.Body = append(vo.Body, obj)
		}
	}
	rt.Data = vo
	return rt.Ok()
}
