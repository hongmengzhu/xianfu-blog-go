package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicAttachmentRepository))

	gs.Provide(new(support.BaseService[BasicAttachmentRepository]))
}

type BasicAttachmentRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicAttachmentEntity, int64]
}

func (c *BasicAttachmentRepository) FindAllByModuleValue(ctx context.Context, module, value string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("module=?", module).Where("value=?", value).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) FindAllByModuleValueIn(ctx context.Context, module string, value []string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("module=?", module).Where("value in ?", value).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) FindByModuleTypeValue(ctx context.Context, module, typ, value string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("module=?", module).Where("type=?", typ).Where("value=?", value).Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) DeleteByModuleTypeValue(ctx context.Context, module, typ, value string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("module=?", module).Where("type=?", typ).Where("value=?", value).Delete(&c.Entity)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) FindByMark(ctx context.Context, mark string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", mark).Order("sort ASC,id desc").Find(&info)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) DeleteByMark(ctx context.Context, mark string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", mark).Delete(&c.Entity)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) DeleteByNoAndFileOwner(ctx context.Context, no []string, fileOwner string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("no in ?", no).Where("file_owner=?", fileOwner).Delete(&c.Entity)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) DeleteByIdAndFileOwner(ctx context.Context, no []string, fileOwner string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", no).Where("file_owner=?", fileOwner).Delete(&c.Entity)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) UpdateByNoAndFileOwnerSetState13(ctx context.Context, no []string, fileOwner string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("no in ?", no).Where("file_owner=?", fileOwner).Update("state", 13)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) UpdateByIdAndFileOwnerSetState13(ctx context.Context, no []string, fileOwner string) (info []*entityBasic.BasicAttachmentEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", no).Where("file_owner=?", fileOwner).Update("state", 13)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicAttachmentRepository) UpdateByIdSetFileOwner(ctx context.Context, no []string, fileOwner string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", no).Update("file_owner", fileOwner)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return false
	}
	if 0 == tx.RowsAffected {
		return false
	}
	return true
}
