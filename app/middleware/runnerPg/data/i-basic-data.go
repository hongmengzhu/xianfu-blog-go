package data

import (
	"context"
	"slices"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constNodePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constTags"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
	"gorm.io/datatypes"
)

// IBasicData
// @Description: 初始化基础数据
type IBasicData struct {
	Bus         event.Bus                                         `autowire:"?"`
	country     *repositoryBasic.BasicCountryRepository           `autowire:"?"`
	tagsCat     *repositoryBasic.BasicTagsCategoryRepository      `autowire:"?"`
	conList     *repositoryBasic.BasicConfigListRepository        `autowire:"?"`
	config      *repositoryBasic.BasicConfigRepository            `autowire:"?"`
	model       *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	modelFields *repositoryBasic.BasicConfigModelFieldsRepository `autowire:"?"`
	event       *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	eventFields *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
}

func (b *IBasicData) Run(ctx context.Context) error {
	log.Infof(context.Background(), log.TagAppDef, "[init].[基础数据：国家:中国]===================")
	{
		save := entityBasic.BasicCountryEntity{
			ID:           100000,
			No:           "100000",
			Code:         "86",
			Name:         "中国",
			NameFl:       "CHN",
			NameFull:     "中国",
			State:        enumStatePg.ENABLE.Index(),
			Iso3:         "CHN",
			CountryCode:  "86",
			PhoneUse:     enumStatePg.ENABLE.Index(),
			DomainSuffix: "cn",
		}
		//不存在时创建
		_, result := b.country.FindByNo(ctx, save.No)
		if !result {
			save.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(save.ID))
			save.NoLink = constNodePg.NoLinkDefault(save.No)
			b.country.Create(ctx, &save)
		}
	}
	//
	log.Infof(context.Background(), log.TagAppDef, "[init].[基础数据：标签:分类]===================")
	{
		save := entityBasic.BasicTagsCategoryEntity{
			ID:       100000,
			No:       constTags.ArticleInfo.Index(),
			Code:     constTags.ArticleInfo.Index(),
			Name:     "文章",
			NameFl:   "文章",
			NameFull: "文章",
			State:    enumStatePg.ENABLE.Index(),
			TenantNo: constsPg.ACCOUNT_MANAGE_No,
			TypeSys:  typeSysPg.System.Index(),
			Tags:     datatypes.NewJSONType(make([]string, 0)),
		}
		//不存在时创建
		_, result := b.tagsCat.FindByNo(ctx, save.No)
		if !result {
			save.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(save.ID))
			save.NoLink = constNodePg.NoLinkDefault(save.No)
			b.tagsCat.Create(ctx, &save)
		}
	}
	//
	b.initSysConfig()
	//
	b.initDefaultTenantConfig()
	//
	b.initSystemModel()
	//
	b.initDefaultTenantModel()
	//
	b.initSystemEvent()
	//
	b.initDefaultTenantEvent()
	return nil
}

// initSysConfig 初始化系统配置
func (b *IBasicData) initSysConfig() {
	ctx := context.Background()
	log.Infof(context.Background(), log.TagAppDef, "[init].[基础数据：配置.系统配置]===================")
	save := entityBasic.BasicConfigListEntity{
		State:       enumStatePg.ENABLE.Index(),
		Show:        yesNoIntPg.Yes.Index(),
		TypeDomain:  "system",
		ID:          1,
		No:          "systemConfig",
		Name:        "系统配置",
		EventNo:     "systemConfig",
		ModelNo:     "systemConfig",
		Field:       "systemConfig",
		FieldPath:   "systemConfig",
		Description: "系统配置",
	}
	save.KindUnique = cryptPg.Md5(save.No)
	//不存在时创建
	_, result := b.conList.FindByNo(ctx, save.No)
	if !result {
		b.conList.Create(ctx, &save)
	}
	fields := make([]string, 0)
	dataInset := make([]*entityBasic.BasicConfigEntity, 0)
	{
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				ID:           1,
				EventNo:      save.EventNo,
				ModelNo:      save.ModelNo,
				Field:        "name",
				FieldPath:    "name",
				Name:         "系统名称",
				Description:  "",
				DefaultValue: "例如：仙府系统",
				Value:        "仙府系统",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				ID:           2,
				EventNo:      save.EventNo,
				ModelNo:      save.ModelNo,
				Field:        "logo",
				FieldPath:    save.FieldPath,
				Name:         "logo图",
				Description:  "",
				DefaultValue: "",
				Value:        "",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//

		//不存在时创建
		info, r := b.config.FindByEventNoAndFieldIn(ctx, save.EventNo, fields)
		if !r {
			b.config.DbModel().Create(dataInset)
		} else {
			find := make([]string, 0)
			for _, v := range info {
				find = append(find, v.Field)
			}
			dataInsetNew := make([]*entityBasic.BasicConfigEntity, 0)
			//
			for _, v := range dataInset {
				if !slices.Contains(find, v.Field) {
					dataInsetNew = append(dataInsetNew, v)
				}
			}
			//
			if len(dataInsetNew) > 0 {
				b.config.DbModel().Create(dataInsetNew)
			}
		}
	}
}

// 默认租户配置
func (b *IBasicData) initDefaultTenantConfig() {
	ctx := context.Background()
	log.Infof(context.Background(), log.TagAppDef, "[init].[基础数据：配置.默认租户配置]===================")
	save := entityBasic.BasicConfigListEntity{
		State:       enumStatePg.ENABLE.Index(),
		Show:        yesNoIntPg.Yes.Index(),
		TypeDomain:  "tenant",
		ID:          2,
		No:          constsPg.ACCOUNT_MANAGE_No,
		Name:        "默认租户配置",
		EventNo:     constsPg.ACCOUNT_MANAGE_No,
		ModelNo:     constsPg.ACCOUNT_MANAGE_No,
		TenantNo:    constsPg.ACCOUNT_MANAGE_No,
		Field:       "tenantConfig",
		FieldPath:   "tenantConfig",
		Description: "默认租户配置",
	}
	save.KindUnique = cryptPg.Md5(save.No)
	//不存在时创建
	_, result := b.conList.FindByNo(ctx, save.No)
	if !result {
		b.conList.Create(ctx, &save)
	}
	fields := make([]string, 0)
	dataInset := make([]*entityBasic.BasicConfigEntity, 0)
	//默认租户
	{
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "name",
				FieldPath:    "name",
				Name:         "站点名称",
				Description:  "例如：博客",
				DefaultValue: "",
				Value:        "仙府博客",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepage",
				FieldPath:    "homepage",
				Name:         "首页名称",
				Description:  "例如：首页",
				DefaultValue: "",
				Value:        "首页",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageSubtitle",
				FieldPath:    "homepageSubtitle",
				Name:         "首页副标题",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageDescription",
				FieldPath:    "homepageDescription",
				Name:         "首页SEO描述",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageKeyword",
				FieldPath:    "homepageKeyword",
				Name:         "首页SEO关键词",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//不存在时创建
		info, r := b.config.FindByEventNoAndFieldIn(ctx, save.EventNo, fields)
		if !r {
			b.config.DbModel().Create(dataInset)
		} else {
			find := make([]string, 0)
			for _, v := range info {
				find = append(find, v.Field)
			}
			dataInsetNew := make([]*entityBasic.BasicConfigEntity, 0)
			//
			for _, v := range dataInset {
				if !slices.Contains(find, v.Field) {
					dataInsetNew = append(dataInsetNew, v)
				}
			}
			//
			if len(dataInsetNew) > 0 {
				b.config.DbModel().Create(dataInsetNew)
			}
		}
	}
}

// 系统模型
func (b *IBasicData) initSystemModel() {
	ctx := context.Background()
	save := entityBasic.BasicConfigModelEntity{
		State:         enumStatePg.ENABLE.Index(),
		Show:          yesNoIntPg.Yes.Index(),
		ID:            1,
		No:            "systemConfig",
		Model:         "systemConfig",
		Table:         "systemConfig",
		Name:          "系统配置",
		ModelCategory: configModelPg.ModelCategoryTable.Index(),
	}
	save.Tags = datatypes.NewJSONType[[]string](make([]string, 0))
	save.KindUnique = cryptPg.Md5(save.Model)
	//不存在时创建
	_, result := b.model.FindByNo(ctx, save.No)
	if !result {
		b.model.Create(ctx, &save)
	}

	fields := make([]string, 0)
	dataInset := make([]*entityBasic.BasicConfigModelFieldsEntity, 0)
	{
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				ID:           1,
				ModelNo:      save.No,
				Field:        "name",
				FieldPath:    "name",
				Name:         "系统名称",
				Description:  "",
				DefaultValue: "例如：仙府系统",
				Value:        "仙府系统",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				ID:           2,
				ModelNo:      save.No,
				Field:        "logo",
				FieldPath:    "logo",
				Name:         "logo图",
				Description:  "",
				DefaultValue: "",
				Value:        "",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//

		//不存在时创建
		info, r := b.modelFields.FindByEventNoAndFieldIn(ctx, save.No, fields)
		if !r {
			b.modelFields.DbModel().Create(dataInset)
		} else {
			find := make([]string, 0)
			for _, v := range info {
				find = append(find, v.Field)
			}
			dataInsetNew := make([]*entityBasic.BasicConfigModelFieldsEntity, 0)
			//
			for _, v := range dataInset {
				if !slices.Contains(find, v.Field) {
					dataInsetNew = append(dataInsetNew, v)
				}
			}
			//
			if len(dataInsetNew) > 0 {
				b.modelFields.DbModel().Create(dataInsetNew)
			}
		}
	}
}
func (b *IBasicData) initDefaultTenantModel() {
	ctx := context.Background()
	save := entityBasic.BasicConfigModelEntity{
		State:         enumStatePg.ENABLE.Index(),
		Show:          yesNoIntPg.Yes.Index(),
		ID:            2,
		No:            constsPg.ACCOUNT_MANAGE_No,
		Model:         "tenantConfig",
		Table:         "tenantConfig",
		Name:          "默认租户配置",
		ModelCategory: configModelPg.ModelCategoryTable.Index(),
	}
	save.Tags = datatypes.NewJSONType[[]string](make([]string, 0))
	save.KindUnique = cryptPg.Md5(save.Model)
	//不存在时创建
	_, result := b.model.FindByNo(ctx, save.No)
	if !result {
		b.model.Create(ctx, &save)
	}
	fields := make([]string, 0)
	dataInset := make([]*entityBasic.BasicConfigModelFieldsEntity, 0)
	//默认租户
	{
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "name",
				FieldPath:    "name",
				Name:         "站点名称",
				Description:  "例如：博客",
				DefaultValue: "",
				Value:        "仙府博客",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepage",
				FieldPath:    "homepage",
				Name:         "首页名称",
				Description:  "例如：首页",
				DefaultValue: "",
				Value:        "首页",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageSubtitle",
				FieldPath:    "homepageSubtitle",
				Name:         "首页副标题",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageDescription",
				FieldPath:    "homepageDescription",
				Name:         "首页SEO描述",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigModelFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageKeyword",
				FieldPath:    "homepageKeyword",
				Name:         "首页SEO关键词",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//不存在时创建
		info, r := b.modelFields.FindByEventNoAndFieldIn(ctx, save.No, fields)
		if !r {
			b.modelFields.DbModel().Create(dataInset)
		} else {
			find := make([]string, 0)
			for _, v := range info {
				find = append(find, v.Field)
			}
			dataInsetNew := make([]*entityBasic.BasicConfigModelFieldsEntity, 0)
			//
			for _, v := range dataInset {
				if !slices.Contains(find, v.Field) {
					dataInsetNew = append(dataInsetNew, v)
				}
			}
			//
			if len(dataInsetNew) > 0 {
				b.modelFields.DbModel().Create(dataInsetNew)
			}
		}
	}
}

func (b *IBasicData) initSystemEvent() {
	ctx := context.Background()
	save := entityBasic.BasicConfigEventEntity{
		State:   enumStatePg.ENABLE.Index(),
		Show:    yesNoIntPg.Yes.Index(),
		ID:      1,
		No:      "systemConfig",
		Model:   "systemConfig",
		ModelNo: "systemConfig",
		Name:    "系统配置",
	}
	save.Tags = datatypes.NewJSONType[[]string](make([]string, 0))
	save.KindUnique = cryptPg.Md5(save.Model)
	//不存在时创建
	_, result := b.event.FindByNo(ctx, save.No)
	if !result {
		b.event.Create(ctx, &save)
	}

	fields := make([]string, 0)
	dataInset := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
	{
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ID:           1,
				ModelNo:      save.ModelNo,
				EventNo:      save.No,
				Field:        "name",
				FieldPath:    "name",
				Name:         "系统名称",
				Description:  "",
				DefaultValue: "例如：仙府系统",
				Value:        "仙府系统",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ID:           2,
				ModelNo:      save.ModelNo,
				EventNo:      save.No,
				Field:        "logo",
				FieldPath:    "logo",
				Name:         "logo图",
				Description:  "",
				DefaultValue: "",
				Value:        "",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//

		//不存在时创建
		info, r := b.eventFields.FindByEventNoAndFieldIn(ctx, save.No, fields)
		if !r {
			b.eventFields.DbModel().Create(dataInset)
		} else {
			find := make([]string, 0)
			for _, v := range info {
				find = append(find, v.Field)
			}
			dataInsetNew := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
			//
			for _, v := range dataInset {
				if !slices.Contains(find, v.Field) {
					dataInsetNew = append(dataInsetNew, v)
				}
			}
			//
			if len(dataInsetNew) > 0 {
				b.eventFields.DbModel().Create(dataInsetNew)
			}
		}
	}
}
func (b *IBasicData) initDefaultTenantEvent() {
	ctx := context.Background()
	save := entityBasic.BasicConfigEventEntity{
		State:    enumStatePg.ENABLE.Index(),
		Show:     yesNoIntPg.Yes.Index(),
		ID:       2,
		No:       constsPg.ACCOUNT_MANAGE_No,
		ModelNo:  constsPg.ACCOUNT_MANAGE_No,
		TenantNo: constsPg.ACCOUNT_MANAGE_No,
		Model:    "tenantConfig",
		Name:     "默认租户配置",
	}
	save.Tags = datatypes.NewJSONType[[]string](make([]string, 0))
	save.KindUnique = cryptPg.Md5(save.Model)
	//不存在时创建
	_, result := b.event.FindByNo(ctx, save.No)
	if !result {
		b.event.Create(ctx, &save)
	}
	fields := make([]string, 0)
	dataInset := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
	//默认租户
	{
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				Field:        "name",
				FieldPath:    "name",
				Name:         "站点名称",
				Description:  "例如：博客",
				DefaultValue: "",
				Value:        "仙府博客",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepage",
				FieldPath:    "homepage",
				Name:         "首页名称",
				Description:  "例如：首页",
				DefaultValue: "",
				Value:        "首页",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageSubtitle",
				FieldPath:    "homepageSubtitle",
				Name:         "首页副标题",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageDescription",
				FieldPath:    "homepageDescription",
				Name:         "首页SEO描述",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//
		{
			item := entityBasic.BasicConfigEventFieldsEntity{
				State:        enumStatePg.ENABLE.Index(),
				Show:         yesNoIntPg.Yes.Index(),
				Binary:       yesNoIntPg.No.Index(),
				ModelNo:      constsPg.ACCOUNT_MANAGE_No,
				TenantNo:     constsPg.ACCOUNT_MANAGE_No,
				EventNo:      constsPg.ACCOUNT_MANAGE_No,
				Field:        "homepageKeyword",
				FieldPath:    "homepageKeyword",
				Name:         "首页SEO关键词",
				Description:  "例如：创造未来",
				DefaultValue: "",
				Value:        "创造未来",
				Rules:        datatypes.NewJSONType(make([]string, 0)),
			}
			item.No = noPg.No()
			item.KindUnique = cryptPg.Md5(item.Field)
			fields = append(fields, item.Field)
			dataInset = append(dataInset, &item)
		}
		//不存在时创建
		info, r := b.eventFields.FindByEventNoAndFieldIn(ctx, save.No, fields)
		if !r {
			b.eventFields.DbModel().Create(dataInset)
		} else {
			find := make([]string, 0)
			for _, v := range info {
				find = append(find, v.Field)
			}
			dataInsetNew := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
			//
			for _, v := range dataInset {
				if !slices.Contains(find, v.Field) {
					dataInsetNew = append(dataInsetNew, v)
				}
			}
			//
			if len(dataInsetNew) > 0 {
				b.eventFields.DbModel().Create(dataInsetNew)
			}
		}
	}
}
