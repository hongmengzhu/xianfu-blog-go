package tagsBasicEvent

import (
	"context"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/model/modEventBasicTags"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
)

// SaveByCategory
// @Description: 标签处理
type SaveByCategory struct {
	sp *Sp

	dto modEventBasicTags.TagsRelation
}

func NewSaveByCategory(sp *Sp, dto modEventBasicTags.TagsRelation) *SaveByCategory {
	return &SaveByCategory{
		sp:  sp,
		dto: dto,
	}
}

// Processor 处理
//
//	@Description:
//	@receiver c
//	@param ctx
//	@return error
func (c *SaveByCategory) Processor() error {
	if strPg.IsBlank(c.dto.Category) {
		return nil
	}
	if nil == c.dto.Tags {
		return nil
	}
	if len(c.dto.Tags) < 1 {
		return nil
	}
	//类别是否存在
	info, result := c.sp.TagCate.FindByNo(context.Background(), c.dto.Category)
	if !result {
		return nil
	}
	categoryRoot := make([]string, 0)
	categoryRoot = append(categoryRoot, c.dto.Category)
	mapTagsOnly := make(map[string]string)
	tagsEntity := make(map[string]*entityBasic.BasicTagsEntity)
	tagsRelation := make(map[string]*entityBasic.BasicTagsRelationEntity)
	tagsNo := make([]string, 0)
	in, r := c.sp.TagsDb.FindAllByNameIn(context.Background(), c.dto.Tags)
	if r {
		for _, item := range in {
			mapTagsOnly[item.Name] = item.No
			tagsEntity[item.No] = item
			tagsNo = append(tagsNo, item.No)
		}
	}
	//如果不存在设置 false
	for _, item := range c.dto.Tags {
		if _, ok := mapTagsOnly[item]; !ok {
			mapTagsOnly[item] = ""
		}
	}
	//获取标签关系
	infos, b := c.sp.TagRela.FindAllByTagNoInAndCategoryRoot(context.Background(), tagsNo, c.dto.Category)
	if b {
		for _, item := range infos {
			tagsRelation[item.TagNo] = item
		}
	}
	for _, v := range c.dto.Tags {
		//标签是否存在
		if get, ok := mapTagsOnly[v]; ok {
			if strPg.IsNotBlank(get) {
				//关系存在，则跳过
				if _, ok2 := tagsRelation[get]; ok2 {
					continue
				}
			}
			save := &entityBasic.BasicTagsEntity{}
			//标签 不存在，创建标签
			if strPg.IsBlank(get) {
				save.Name = strings.TrimSpace(v)
				save.NameFull = save.Name
				save.State = enumStatePg.ENABLE.Index()
				save.TypeSys = typeSysPg.General.Index()
				save.No = noPg.No()
				save.Code = save.No
				save.CategoryNo = info.No
				save.TenantNo = c.dto.Holder.GetTenantNo()
				save.CreateBy = c.dto.Holder.GetAccountNo()
				c.sp.TagsDb.Create(context.Background(), save)
				//
				mapTagsOnly[v] = save.No
				tagsEntity[save.No] = save
			} else {
				save = tagsEntity[get]
			}
			//关系不存在
			if _, ok2 := tagsRelation[save.No]; !ok2 {
				save = tagsEntity[save.No]
				rel := entityBasic.BasicTagsRelationEntity{}
				rel.TagNo = save.No
				rel.CategoryNo = save.No
				rel.State = enumStatePg.ENABLE.Index()
				rel.Name = save.Name
				rel.NameFl = save.NameFl
				rel.NameFull = save.NameFull
				rel.Code = save.Code
				rel.TenantNo = save.TenantNo
				rel.TypeSys = save.TypeSys
				rel.Module = save.Module
				//rel.Attribute = save.Attribute
				rel.CategoryRoot = c.dto.Category
				c.sp.TagRela.Create(context.Background(), &rel)
				//
				tagsRelation[rel.TagNo] = &rel
			}
		}
	}
	c.cache(categoryRoot)
	return nil
}

// 缓存更新
func (t *SaveByCategory) cache(categoryRoot []string) {
	dto := modEventBasicTags.TagsCacheDto{CategoryRoot: categoryRoot}
	//缓存更新
	err := NewCachePush(t.sp, dto).Processor(context.Background())
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, err, "tags.push.error:=", err)
	}
}
