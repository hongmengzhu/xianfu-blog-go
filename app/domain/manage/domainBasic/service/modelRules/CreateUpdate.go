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

type CreateUpdate struct {
	Sp *Sp `autowire:"?"`

	ct       modBasicModelRules.CreateUpdateCt
	model    *entityBasic.BasicConfigModelEntity
	event    *entityBasic.BasicConfigEventEntity
	fields   []*entityBasic.BasicModelRulesEntity
	isUpdate bool
}

func NewCreateUpdate(sp *Sp,
	ct modBasicModelRules.CreateUpdateCt, isUpdate bool) *CreateUpdate {
	return &CreateUpdate{
		Sp:       sp,
		isUpdate: isUpdate,
		ct:       ct,
		fields:   make([]*entityBasic.BasicModelRulesEntity, 0),
		model:    &entityBasic.BasicConfigModelEntity{},
		event:    &entityBasic.BasicConfigEventEntity{},
	}
}
func (c *CreateUpdate) Process(ctx *gin.Context) (rt rg.Rs[string]) {
	return c.verify(ctx)
}

func (c *CreateUpdate) verify(ctx *gin.Context) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", c.ct)
	ct := c.ct
	if strPg.IsBlank(ct.FieldNo) {
		return rt.ErrorMessage("字段编号不能为空")
	}
	typeModel := "model"
	no := ""
	tenantNo := ""
	event, result := c.Sp.repModelFields.FindByNo(ctx, ct.FieldNo)
	if !result {
		info, r := c.Sp.repEventFields.FindByNo(ctx, ct.FieldNo)
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
	if strPg.IsBlank(ct.Name) {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.RuleMode) {
		return rt.ErrorMessage(ct.Name + " 验证模式类型不能为空")
	}
	//
	obj := entityBasic.BasicModelRulesEntity{}
	err := copier.Copy(&obj, &ct)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	obj.Name = strings.TrimSpace(obj.Name)
	obj.RuleMode = strings.TrimSpace(obj.RuleMode)
	//
	tags := make([]string, 0)
	if nil != ct.RuleTarget && len(ct.RuleTarget) > 0 {
		for _, v := range ct.RuleTarget {
			if strPg.IsNotBlank(v) {
				tags = append(tags, strings.TrimSpace(v))
			}
		}
	}
	obj.RuleTarget = datatypes.NewJSONType[[]string](tags)
	//
	if ct.Id.ToInt64() > 0 {
		info, r := c.Sp.repRules.FindByIdString(ctx, ct.Id.ToString())
		if !r {
			return rt.ErrorMessage("字段规则不存在")
		}
		obj.No = ""
		obj.TenantNo = ""
		obj.ValueNo = ""
		err := c.Sp.repRules.Update(ctx, obj, info.ID)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "save err=%+v", err)
			return rt.ErrorMessage("保存失败：")
		}
	} else {
		obj.ValueNo = no
		obj.TenantNo = tenantNo
		obj.ID = 0
		obj.No = noPg.No()
		obj.State = enumStatePg.ENABLE.IndexInt8()
		obj.Show = yesNoIntPg.Yes.IndexInt8()
		obj.Sort = 0
		obj.TypeModel = typeModel
		err, _ := c.Sp.repRules.Create(ctx, &obj)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "save err=%+v", err)
			return rt.ErrorMessage("保存失败：")
		}
	}
	return rt.Ok()
}
