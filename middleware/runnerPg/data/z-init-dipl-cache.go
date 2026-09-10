package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityApi"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryApi"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/components/cachePg/cacheDiplPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
	"gorm.io/gorm"
)

// ZInitDiplCache
// @Description: 初始化 dipl 缓存
type ZInitDiplCache struct {
	sv *repositoryApi.ApiDiplAccessKeyRepository `autowire:"?"`
}

func (b *ZInitDiplCache) Run(ctx context.Context) error {
	log.Infof(ctx, log.TagAppDef, "初始化 => 接口密钥")
	var query entityApi.ApiDiplAccessKeyEntity
	query.State = enumStatePg.ENABLE.Index()
	infos := b.sv.FindAll(context.Background(), query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
		db = db.Order("create_at desc")
		db = db.Where("expiry_date >= ?", datetimePg.Now())
		return db
	}))
	if infos != nil && len(infos) > 0 {
		for _, info := range infos {
			// 添加缓存
			//加入缓存
			sha := cacheDiplPg.HashSha(info.Key, info.Secret)
			obj := cacheDiplPg.DiplCo{
				HashSha:  sha,
				No:       info.DiplNo,
				TenantNo: info.TenantNo,
				Key:      info.Key,
				Secret:   info.Secret,
				Ano:      info.Ano,
			}
			cacheDiplPg.Set(info.Key, obj)
		}
	}
	return nil
}
