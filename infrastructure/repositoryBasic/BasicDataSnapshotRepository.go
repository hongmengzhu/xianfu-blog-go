package repositoryBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"github.com/pangu-2/go-tools/tools/jsonPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicDataSnapshotRepository))

	gs.Provide(new(support.BaseService[BasicDataSnapshotRepository]))
}

type BasicDataSnapshotRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicDataSnapshotEntity, int64]
}

func (c *BasicDataSnapshotRepository) FindByNameAndIdNot(ctx context.Context, name string, id int64) (info *entityBasic.BasicDataSnapshotEntity, result bool) {
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

func (c *BasicDataSnapshotRepository) SnapshotVersion(ctx context.Context, obj interface{}, module, tenantNo, value, version, extend string) {
	toJson, err := jsonPg.ObjToJson(obj)
	if nil != err {
		log.Errorf(ctx, log.TagAppDef, "snapshot version error: %+v", err)
	} else {
		mark := value + "|" + value
		c.Create(ctx, &entityBasic.BasicDataSnapshotEntity{
			Module:   module,
			Data:     toJson,
			Mark:     mark,
			Name:     "",
			TenantNo: tenantNo,
			Value:    value,
			Version:  version,
			Extend:   extend,
		})
	}
}

func (c *BasicDataSnapshotRepository) SnapshotVersionAll(ctx context.Context, obj interface{}, module, tenantNo, value, version, name, extend string) {
	toJson, err := jsonPg.ObjToJson(obj)
	if nil != err {
		log.Errorf(ctx, log.TagAppDef, "snapshot version error: %+v", err)
	} else {
		mark := value + "|" + value
		c.Create(ctx, &entityBasic.BasicDataSnapshotEntity{
			Module:   module,
			Data:     toJson,
			Mark:     mark,
			Name:     name,
			TenantNo: tenantNo,
			Value:    value,
			Version:  version,
			Extend:   extend,
		})
	}
}
