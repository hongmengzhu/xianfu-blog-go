package modelRules

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicModelRules"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoIntPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"gorm.io/datatypes"
)

type CreateUpdateData struct {
	Sp *Sp `autowire:"?"`

	ct       modBasicModelRules.CreateUpdateDataCt
	model    *entityBasic.BasicConfigModelEntity
	event    *entityBasic.BasicConfigEventEntity
	fields   []*entityBasic.BasicModelRulesEntity
	isUpdate bool
}

func NewCreateUpdateData(sp *Sp,
	ct modBasicModelRules.CreateUpdateDataCt, isUpdate bool) *CreateUpdateData {
	return &CreateUpdateData{
		Sp:       sp,
		isUpdate: isUpdate,
		ct:       ct,
		fields:   make([]*entityBasic.BasicModelRulesEntity, 0),
		model:    &entityBasic.BasicConfigModelEntity{},
		event:    &entityBasic.BasicConfigEventEntity{},
	}
}

func (c *CreateUpdateData) Process(ctx *gin.Context) (rt rg.Rs[string]) {
	return c.verify(ctx)
}

func (c *CreateUpdateData) verify(ctx *gin.Context) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", c.ct)
	ct := c.ct
	if nil == ct.Body || len(ct.Body) < 1 {
		return rt.ErrorMessage("规则不能为空")
	}
	if strPg.IsBlank(ct.FieldNo) {
		return rt.ErrorMessage("字段编号不能为空")
	}
	typeModel := "model"
	no := ""
	tenantNo := ""
	event, result := c.Sp.repEvent.FindByNo(ctx, ct.FieldNo)
	if !result {
		info, r := c.Sp.repModel.FindByNo(ctx, ct.FieldNo)
		if !r {
			return rt.ErrorMessage("模型/事件不存在")
		}
		no = info.No
		tenantNo = info.TenantNo
		typeModel = "event"
	} else {
		no = event.No
		tenantNo = event.TenantNo
	}

	if nil != ct.BodyDelIds || len(ct.BodyDelIds) > 0 {
		delIds := make([]string, 0)
		for _, id := range ct.BodyDelIds {
			if strPg.IsNotBlank(id) {
				delIds = append(delIds, strings.TrimSpace(id))
			}
		}
		if len(delIds) > 0 {
			c.Sp.repRules.DeleteAllByValueNoAndIds(ctx, no, delIds)
		}
	}
	dataAdd := make([]*entityBasic.BasicModelRulesEntity, 0)
	dataUpdate := make([]*entityBasic.BasicModelRulesEntity, 0)
	for _, item := range ct.Body {
		//
		if strPg.IsBlank(item.Name) {
			return rt.ErrorMessage("名称不能为空")
		}
		if strPg.IsBlank(item.RuleMode) {
			return rt.ErrorMessage(item.Name + " 验证模式类型不能为空")
		}
	}
	for _, item := range ct.Body {
		obj := entityBasic.BasicModelRulesEntity{}
		err := copier.Copy(&obj, &item)
		if err != nil {
			log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
		}
		obj.Name = strings.TrimSpace(obj.Name)
		obj.RuleMode = strings.TrimSpace(obj.RuleMode)
		//
		tags := make([]string, 0)
		if nil != item.RuleTarget && len(item.RuleTarget) > 0 {
			for _, v := range item.RuleTarget {
				if strPg.IsNotBlank(v) {
					tags = append(tags, strings.TrimSpace(v))
				}
			}
		}
		obj.RuleTarget = datatypes.NewJSONType[[]string](tags)
		//
		if item.Id.ToInt64() > 0 {
			obj.No = ""
			obj.TenantNo = ""
			obj.ValueNo = ""
			dataUpdate = append(dataUpdate, &obj)
		} else {
			obj.ValueNo = no
			obj.TenantNo = tenantNo
			obj.ID = 0
			obj.No = noPg.No()
			obj.State = enumStatePg.ENABLE.IndexInt8()
			obj.Show = yesNoIntPg.Yes.IndexInt8()
			obj.Sort = 0
			obj.TypeModel = typeModel
			dataAdd = append(dataAdd, &obj)
		}
	}
	if len(dataAdd) > 0 {
		tx := c.Sp.repRules.DbModel().CreateInBatches(dataAdd, 1000000)
		if tx.Error != nil {
			log.Errorf(ctx, log.TagAppDef, "save err=%+v", tx.Error)
			return rt.ErrorMessage("保存失败：")
		}
	}
	if len(dataUpdate) > 0 {
		for _, entity := range dataUpdate {
			c.Sp.repRules.Update(ctx, *entity, entity.ID)
		}
	}
	return rt.Ok()
}
