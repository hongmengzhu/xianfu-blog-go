package configEvent

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigEvent"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model/modelBasePg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/validatePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
)

type CreateUpdate struct {
	Sp *Sp `autowire:"?"`

	ct       modBasicConfigEvent.CreateUpdateCt
	model    *entityBasic.BasicConfigModelEntity
	event    *entityBasic.BasicConfigEventEntity
	module   *entityBasic.BasicModuleEntity
	fields   []*entityBasic.BasicConfigEventFieldsEntity
	isUpdate bool
}

func NewCreateUpdate(sp *Sp,
	ct modBasicConfigEvent.CreateUpdateCt, isUpdate bool) *CreateUpdate {
	return &CreateUpdate{
		Sp:       sp,
		isUpdate: isUpdate,
		ct:       ct,
		fields:   make([]*entityBasic.BasicConfigEventFieldsEntity, 0),
		model:    &entityBasic.BasicConfigModelEntity{},
		event:    &entityBasic.BasicConfigEventEntity{},
		module:   &entityBasic.BasicModuleEntity{},
	}
}

func (c *CreateUpdate) Process(ctx *gin.Context) (rt rg.Rs[string]) {
	return c.verify(ctx)
}

func (c *CreateUpdate) verify(ctx *gin.Context) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", c.ct)
	header := c.ct.Header
	//
	err := copier.Copy(c.event, &header)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	c.event.Name = strings.TrimSpace(c.event.Name)
	c.event.Model = strings.TrimSpace(c.event.Model)
	c.event.ModelNo = strings.TrimSpace(c.event.ModelNo)
	c.event.ModuleSub = strings.TrimSpace(c.event.ModuleSub)
	//
	if strPg.IsBlank(c.event.Name) {
		return rt.ErrorMessage("模型中文名称不能为空")
	}
	if strPg.IsBlank(c.event.ModelNo) {
		return rt.ErrorMessage("模型不能为空")
	}

	if strPg.IsBlank(header.Field) {
		return rt.ErrorMessage("字段名称不能为空")
	}
	if !validatePg.ValidateString(header.Field) {
		return rt.ErrorMessage("字段名称格式错误")
	}
	model, b := c.Sp.repModel.FindByNo(ctx, c.event.ModelNo)
	if !b {
		return rt.ErrorMessage("模型不存在")
	}

	rt.Extend = make(map[string]interface{})
	errs := make([]modelBasePg.ItemResult, 0)
	//
	//
	ids := make([]string, 0)
	if nil != c.ct.Body && len(c.ct.Body) > 0 {
		i := int64(0)
		//
		uniqueField := make(map[string]bool)
		//
		for _, item := range c.ct.Body {
			i++

			var field entityBasic.BasicConfigEventFieldsEntity
			err := copier.Copy(&field, &item)
			if err != nil {
				log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
				errs = append(errs, modelBasePg.ItemResult{
					Col:   i,
					Field: "col",
					Msg:   "整行格式转换错误",
				})
				continue
			}
			field.Name = strings.TrimSpace(field.Name)
			field.Field = strings.TrimSpace(field.Field)
			if strPg.IsBlank(field.Name) {
				errs = append(errs, modelBasePg.ItemResult{
					Col:   i,
					Field: "field",
					Msg:   "字段中文名称不能为空",
				})
				continue
			}
			if strPg.IsBlank(field.Field) {
				errs = append(errs, modelBasePg.ItemResult{
					Col:   i,
					Field: "field",
					Msg:   "字段英文标识不能为空",
				})
				continue
			}
			if _, ok := uniqueField[field.Field]; ok {
				errs = append(errs, modelBasePg.ItemResult{
					Col:   i,
					Field: "field",
					Msg:   "字段英文标识重复",
				})
				continue
			}
			if !validatePg.ValidateString(field.Field) {
				errs = append(errs, modelBasePg.ItemResult{
					Col:   i,
					Field: "field",
					Msg:   "字段英文标识格式错误",
				})
				continue
			}
			uniqueField[field.Field] = true
			field.Show = yesNoIntPg.No.IndexInt8()
			if yesNoIntPg.Yes.IsEqual(item.Show.ToInt8()) {
				field.Show = yesNoIntPg.Yes.IndexInt8()
			}
			field.Binary = yesNoIntPg.No.IndexInt8()
			if yesNoIntPg.Yes.IsEqual(item.Binary.ToInt8()) {
				field.Binary = yesNoIntPg.Yes.IndexInt8()
			}
			c.fields = append(c.fields, &field)
			//
			if item.Id.ToInt64() > 0 {
				ids = append(ids, item.Id.ToString())
			}
		}
		if len(errs) > 0 {
			rt.Extend["errors"] = errs
			return rt.ErrorMessage("保存失败")
		}
	}
	//
	//
	if header.Id.ToInt64() <= 0 {
		c.event.Model = model.Model
		c.event.ModuleSub = model.ModuleSub
		c.event.Module = model.Module
		c.event.TenantNo = model.TenantNo
		c.event.State = enumStatePg.ENABLE.Index()
		c.event.No = noPg.No()
		c.event.Sort = 0
		c.event.KindUnique = cryptPg.Md5(c.event.Field)
		err, _ := c.Sp.repEvent.Create(ctx, c.event)
		if err != nil {
			return rt.ErrorMessage("保存失败 " + err.Error())
		}
		//
		if len(c.fields) > 0 {
			for _, item := range c.fields {
				item.ID = 0
				item.EventNo = c.event.No
				item.ModelNo = c.event.ModelNo
				item.Model = c.event.Model
				item.ModuleSub = c.event.ModuleSub
				item.Module = c.event.Module
				item.TenantNo = c.event.TenantNo
				item.State = enumStatePg.ENABLE.Index()
				item.No = noPg.No()
				item.Sort = 0
				item.KindUnique = cryptPg.Md5(c.event.No + item.Field)
			}
			//
			{
				tx := c.Sp.repEventField.DbModel().CreateInBatches(c.fields, 1000000)
				if tx.Error != nil {
					log.Errorf(ctx, log.TagAppDef, "save err=%+v", tx.Error)
					return rt.ErrorMessage("保存失败：")
				}
				//if 0 == tx.RowsAffected {
				//	return rt.ErrorMessage("保存失败，没有更新任何数据")
				//}
			}
		}
	} else {
		//
		info, b := c.Sp.repEvent.FindByIdString(ctx, header.Id.ToString())
		if !b {
			return rt.ErrorMessage("事件不存在")
		}
		//
		var save entityBasic.BasicConfigEventEntity
		{
			err := copier.Copy(&save, info)
			if err != nil {
				log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
			}
		}
		{
			err := copier.Copy(&save, &header)
			if err != nil {
				log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
			}
		}
		save.KindUnique = cryptPg.Md5(save.Field)
		//
		delIds := make([]string, 0)
		dataInsert := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
		dataUpdate := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
		mapField := make(map[int64]*entityBasic.BasicConfigEventFieldsEntity)
		if len(ids) > 0 {
			infos, result := c.Sp.repEventField.FindAllByIdStringIn(ctx, ids)
			if result {
				for _, item := range infos {
					mapField[item.ID] = item
				}
			}
		}
		if len(c.fields) > 0 {
			for _, item := range c.fields {
				if _, ok := mapField[item.ID]; ok {
					item.Model = info.Model
					item.ModuleSub = info.ModuleSub
					item.Module = info.Module
					item.ModelNo = info.ModelNo
					item.TenantNo = info.TenantNo
					item.KindUnique = cryptPg.Md5(info.No + item.Field)
					dataUpdate = append(dataUpdate, item)
				} else {
					item.ID = 0
					item.EventNo = info.No
					item.ModelNo = info.ModelNo
					item.Model = info.Model
					item.ModuleSub = info.ModuleSub
					item.Module = info.Module
					item.TenantNo = info.TenantNo
					item.State = enumStatePg.ENABLE.Index()
					item.No = noPg.No()
					item.Sort = 0
					item.KindUnique = cryptPg.Md5(info.No + item.Field)
					dataInsert = append(dataInsert, item)
				}
			}
		}
		//
		err := c.Sp.repEvent.Update(ctx, save, info.ID)
		if err != nil {
			log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
		}
		//
		if nil != c.ct.BodyDelIds && len(c.ct.BodyDelIds) > 0 {
			for _, id := range c.ct.BodyDelIds {
				if strPg.IsNotBlank(id) {
					delIds = append(delIds, strings.TrimSpace(id))
				}
			}
			if len(delIds) > 0 {
				c.Sp.repEventField.DeleteAllByEventNoAndIds(ctx, info.No, delIds)
			}
		}
		//
		if len(dataInsert) > 0 {
			{
				tx := c.Sp.repEventField.DbModel().CreateInBatches(dataInsert, 1000000)
				if tx.Error != nil {
					log.Errorf(ctx, log.TagAppDef, "save err=%+v", tx.Error)
					return rt.ErrorMessage("保存失败：")
				}
				//if 0 == tx.RowsAffected {
				//	return rt.ErrorMessage("保存失败，没有更新任何数据")
				//}
			}
		}
		if len(dataUpdate) > 0 {
			for _, entity := range dataUpdate {
				c.Sp.repEventField.Update(ctx, *entity, entity.ID)
			}
		}

	}
	if len(errs) > 0 {
		rt.Extend["errors"] = errs
	}
	return rt.Ok()
}
