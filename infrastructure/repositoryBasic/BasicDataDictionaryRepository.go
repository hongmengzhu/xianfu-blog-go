package repositoryBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicDataDictionaryRepository))

	gs.Provide(new(support.BaseService[BasicDataDictionaryRepository]))
}

type BasicDataDictionaryRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicDataDictionaryEntity, int64]
}

// 查询所有
func (b *BasicDataDictionaryRepository) FindAllOwnerId0(ctx context.Context, t entityBasic.BasicDataDictionaryEntity) (infos []*entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := b.DbModel().WithContext(ctx).Where(t).Where("owner_no='0'").Find(&infos)
	if tx.Error != nil {
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

// FindAllOwners 查询所有
func (b *BasicDataDictionaryRepository) FindAllOwners(ctx context.Context, ids []string) (infos []*entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := b.DbModel().WithContext(ctx).Where("owner_no in ?", ids).Find(&infos)
	if tx.Error != nil {
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

// FindAllOwnersMap 查询所有
func (b *BasicDataDictionaryRepository) FindAllOwnersMap(ctx context.Context, ids []string, t entityBasic.BasicDataDictionaryEntity) (maps map[string][]*entityBasic.BasicDataDictionaryEntity, result bool) {
	if nil == maps {
		maps = make(map[string][]*entityBasic.BasicDataDictionaryEntity)
	}
	var infos []*entityBasic.BasicDataDictionaryEntity
	tx := b.DbModel().WithContext(ctx).Where(t).Where("owner_no in ?", ids).Find(&infos)
	if tx.Error != nil {
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	for _, info := range infos {
		if _, ok := maps[info.OwnerNo]; !ok {
			maps[info.OwnerNo] = make([]*entityBasic.BasicDataDictionaryEntity, 0)
		}
		maps[info.OwnerNo] = append(maps[info.OwnerNo], info)
	}
	return maps, true
}

// FindAllByCodeIn 查询所有
func (b *BasicDataDictionaryRepository) FindAllByCodeIn(ctx context.Context, ids []string) (infos []*entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := b.DbModel().WithContext(ctx).Where("code in ?", ids).Find(&infos)
	if tx.Error != nil {
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return infos, true
}

func (c *BasicDataDictionaryRepository) FindByNameAndIdNot(ctx context.Context, name string, id int64) (info *entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("name=?", name).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicDataDictionaryRepository) FindByValueAndIdNotAndOwnerNo(ctx context.Context, name, id, ownerId string) (info *entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("owner_no=?", ownerId).Where("value=?", name).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
func (c *BasicDataDictionaryRepository) FindByCodeAndIdNotAndOwnerNo(ctx context.Context, code, id, ownerId string) (info *entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("code=?", code).Where("owner_no=?", ownerId).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// FindAllByTypeUniqueMd5AndOwnerNo 根据类型编号查询
func (c *BasicDataDictionaryRepository) FindAllByTypeUniqueMd5AndOwnerNo(ctx context.Context, code, id, ownerId string) (info []*entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("owner_no=?", ownerId).Where("type_unique_md5=?", code).Where("id <> ?", id).Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicDataDictionaryRepository) FindByCodeAndTypeCodeAndIdNotAndOwnerNo(ctx context.Context, code, typeCode, id, ownerId string) (info *entityBasic.BasicDataDictionaryEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("code=?", code).Where("type_code=?", typeCode).Where("owner_no=?", ownerId).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
