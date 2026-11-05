package configModel

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigModel"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model/modelBasePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/validatePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"gorm.io/datatypes"
)

type CreateUpdate struct {
	Sp *Sp `autowire:"?"`

	ct       modBasicConfigModel.CreateUpdateCt
	model    *entityB72.B72ConfigModelEntity
	module   *entityBasic.BasicModuleEntity
	fields   []*entityB72.B72ConfigModelFieldsEntity
	isUpdate bool
}

func NewCreateUpdate(sp *Sp,
	ct modBasicConfigModel.CreateUpdateCt, isUpdate bool) *CreateUpdate {
	return &CreateUpdate{
		Sp:       sp,
		isUpdate: isUpdate,
		ct:       ct,
		fields:   make([]*entityB72.B72ConfigModelFieldsEntity, 0),
		model:    &entityB72.B72ConfigModelEntity{},
		module:   &entityBasic.BasicModuleEntity{},
	}
}

func (c *CreateUpdate) Process(ctx *gin.Context) (rt rg.Rs[string]) {
	return c.verify(ctx)
}

func (c *CreateUpdate) verify(ctx *gin.Context) (rt rg.Rs[string]) {
	header := c.ct.Header
	//
	err := copier.Copy(c.model, &header)
	if err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	c.model.Name = strings.TrimSpace(c.model.Name)
	c.model.Model = strings.TrimSpace(c.model.Model)
	c.model.ModuleSub = strings.TrimSpace(c.model.ModuleSub)
	c.model.ModelCategory = strings.TrimSpace(c.model.ModelCategory)
	c.model.Category = strings.TrimSpace(c.model.Category)
	//
	if strPg.IsBlank(c.model.Name) {
		return rt.ErrorMessage("模型中文名称不能为空")
	}
	if strPg.IsBlank(c.model.Model) {
		return rt.ErrorMessage("模型英文标识不能为空")
	}
	if strPg.IsBlank(c.model.ModelCategory) {
		return rt.ErrorMessage("类型种类不能为空")
	}
	if _, ok := configModelPg.IsExistModelCategory(c.model.ModelCategory); !ok {
		return rt.ErrorMessage("类型种类不存在")
	}
	if strPg.IsBlank(header.ModuleSub) {
		return rt.ErrorMessage("子模块不能为空")
	}
	//子模块是否存在
	{
		b := false
		c.module, b = c.Sp.repModule.FindByNo(ctx, header.ModuleSub)
		if !b {
			return rt.ErrorMessage("子模块不存在")
		}
	}
	if strPg.IsBlank(header.Table) {
		return rt.ErrorMessage("表名称不能为空")
	}
	if !validatePg.ValidateString(header.Table) {
		return rt.ErrorMessage("表名称格式错误")
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

			var field entityB72.B72ConfigModelFieldsEntity
			err := copier.Copy(&field, &item)
			if err != nil {
				log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
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
			if strPg.IsBlank(field.FieldPath) {
				field.FieldPath = field.Field
			}
			field.Sort = item.Sort.ToInt64()
			field.Rules = datatypes.NewJSONType(normalizeStringSlice(item.Rules))
			field.SharedScope = datatypes.NewJSONType(normalizeStringSlice(item.SharedScope))
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
	holder := holderPg.GetContextAccount(ctx)
	tenantNo := holder.GetTenantNo()
	//
	//
	if header.Id.ToInt64() <= 0 {
		// 租户下 model 唯一
		if _, exists := c.Sp.repModel.FindByTenantNoAndModel(ctx, tenantNo, c.model.Model); exists {
			return rt.ErrorMessage("模型标识在当前租户下已存在")
		}
		c.model.State = enumStatePg.ENABLE.Index()
		c.model.TypeSys = typeSysPg.General.String()
		c.model.No = noPg.No()
		c.model.TenantNo = tenantNo
		c.model.Sort = 0
		c.model.KindUnique = cryptPg.Md5(c.model.TenantNo + c.model.Model)
		err, _ := c.Sp.repModel.Create(ctx, c.model)
		if err != nil {
			log.Error(ctx, log.TagAppDef, err, "save err", log.Any("model", c.model))
			return rt.ErrorMessage("保存失败 " + err.Error())
		}
		//
		if len(c.fields) > 0 {
			for _, item := range c.fields {
				item.ID = 0
				item.ModelNo = c.model.No
				item.Model = c.model.Model
				item.ModuleSub = c.model.ModuleSub
				item.TenantNo = c.model.TenantNo
				item.State = enumStatePg.ENABLE.Index()
				item.No = noPg.No()
				item.KindUnique = cryptPg.Md5(c.model.No + item.Field)
			}
			//
			{
				tx := c.Sp.repField.DbModel().CreateInBatches(c.fields, 1000000)
				if tx.Error != nil {
					log.Error(ctx, log.TagAppDef, tx.Error, "save err")
					return rt.ErrorMessage("保存失败：")
				}
				//if 0 == tx.RowsAffected {
				//	return rt.ErrorMessage("保存失败，没有更新任何数据")
				//}
			}
		}
		if err := c.saveRules(ctx, tenantNo); err != nil {
			rt.Extend["errors"] = err
			return rt.ErrorMessage("规则保存失败：" + err.Error())
		}
		if err := c.syncConfig(ctx, tenantNo); err != nil {
			return rt.ErrorMessage("配置保存失败：" + err.Error())
		}
		if err := c.syncConfigList(ctx, tenantNo); err != nil {
			return rt.ErrorMessage("配置列表保存失败：" + err.Error())
		}
	} else {
		//
		info, b := c.Sp.repModel.FindByIdString(ctx, header.Id.ToString())
		if !b {
			return rt.ErrorMessage("模型不存在")
		}
		if strPg.IsNotBlank(info.TenantNo) && info.TenantNo != tenantNo {
			return rt.ErrorMessage("无权操作其他租户的模型")
		}
		// 租户下 model 唯一（排除自己）
		if _, exists := c.Sp.repModel.FindByTenantNoAndModelAndIdNot(ctx, tenantNo, c.model.Model, numberPg.Int64ToString(info.ID)); exists {
			return rt.ErrorMessage("模型标识在当前租户下已存在")
		}
		//
		var save entityB72.B72ConfigModelEntity
		{
			save = *info
			// header 中未提交的字段保持数据库原值，避免部分更新把字段清空
			err := copier.CopyWithOption(&save, &header, copier.Option{IgnoreEmpty: true})
			if err != nil {
				log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
			}
		}
		if strPg.IsBlank(save.Model) {
			save.Model = info.Model
		}
		if strPg.IsBlank(save.ModelCategory) {
			save.ModelCategory = info.ModelCategory
		}
		if strPg.IsBlank(save.TenantNo) {
			save.TenantNo = info.TenantNo
		}
		save.No = ""
		save.KindUnique = cryptPg.Md5(save.TenantNo + save.Model)
		//
		delIds := make([]string, 0)
		dataInsert := make([]*entityB72.B72ConfigModelFieldsEntity, 0)
		dataUpdate := make([]*entityB72.B72ConfigModelFieldsEntity, 0)
		mapField := make(map[int64]*entityB72.B72ConfigModelFieldsEntity)
		if len(ids) > 0 {
			infos, result := c.Sp.repField.FindAllByIdStringIn(ctx, ids)
			if result {
				for _, item := range infos {
					mapField[item.ID] = item
				}
			}
		}
		if len(c.fields) > 0 {
			for _, item := range c.fields {
				if _, ok := mapField[item.ID]; ok {
					item.No = mapField[item.ID].No
					item.ModelNo = info.No
					item.Model = save.Model
					item.ModuleSub = save.ModuleSub
					item.TenantNo = save.TenantNo
					item.KindUnique = cryptPg.Md5(info.No + item.Field)
					dataUpdate = append(dataUpdate, item)
				} else {
					item.ID = 0
					item.ModelNo = info.No
					item.Model = save.Model
					item.ModuleSub = save.ModuleSub
					item.TenantNo = save.TenantNo
					item.State = enumStatePg.ENABLE.Index()
					item.No = noPg.No()
					item.KindUnique = cryptPg.Md5(info.No + item.Field)
					dataInsert = append(dataInsert, item)
				}
			}
		}
		//
		err := c.Sp.repModel.Update(ctx, save, info.ID)
		if err != nil {
			log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
		}
		//
		if nil != c.ct.BodyDelIds && len(c.ct.BodyDelIds) > 0 {
			for _, id := range c.ct.BodyDelIds {
				if strPg.IsNotBlank(id) {
					delIds = append(delIds, strings.TrimSpace(id))
				}
			}
			if len(delIds) > 0 {
				c.Sp.repField.DeleteAllByModelNoAndIds(ctx, info.No, delIds)
			}
		}
		//
		if len(dataInsert) > 0 {
			{
				tx := c.Sp.repField.DbModel().CreateInBatches(dataInsert, 1000000)
				if tx.Error != nil {
					log.Error(ctx, log.TagAppDef, tx.Error, "save err")
					return rt.ErrorMessage("保存失败：")
				}
				//if 0 == tx.RowsAffected {
				//	return rt.ErrorMessage("保存失败，没有更新任何数据")
				//}
			}
		}
		if len(dataUpdate) > 0 {
			for _, entity := range dataUpdate {
				c.Sp.repField.Update(ctx, *entity, entity.ID)
			}
		}
		// 规则与配置同步使用的模型主信息
		c.model.No = info.No
		c.model.TenantNo = info.TenantNo
		c.model.Model = save.Model
		c.model.ModelCategory = save.ModelCategory
		c.model.Module = info.Module
		c.model.ModuleSub = info.ModuleSub
		c.model.Category = save.Category
		c.model.Name = save.Name
		c.model.Description = save.Description
		c.model.Show = save.Show
		c.model.Sort = save.Sort
		c.model.TypeSys = save.TypeSys
		if err := c.saveRules(ctx, tenantNo); err != nil {
			rt.Extend["errors"] = err
			return rt.ErrorMessage("规则保存失败：" + err.Error())
		}
		if err := c.syncConfig(ctx, tenantNo); err != nil {
			return rt.ErrorMessage("配置保存失败：" + err.Error())
		}
		if err := c.syncConfigList(ctx, tenantNo); err != nil {
			return rt.ErrorMessage("配置列表保存失败：" + err.Error())
		}
		// 模型标识变更时，旧 tenant_no + model 缓存也要失效
		if strPg.IsNotBlank(info.Model) && info.Model != c.model.Model && c.Sp.cache != nil {
			if err := c.Sp.cache.Refresh(ctx, tenantNo, info.Model); err != nil {
				log.Error(ctx, log.TagAppDef, err, "refresh old model config cache error")
			}
		}
	}
	if len(errs) > 0 {
		rt.Extend["errors"] = errs
	}
	return rt.Ok()
}
