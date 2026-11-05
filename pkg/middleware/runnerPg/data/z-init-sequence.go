package data

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbMakePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbManagerPg"
	"go-spring.org/log"
	_ "go-spring.org/spring/gs"
)

// ZInitSequence
// @Description: 初始化完成后校正自增序列
// 基础数据初始化中存在显式指定主键（ID=1、2...）的插入，PostgreSQL 不会因显式
// 插入推进序列，序列会落后于 max(id)，之后正常插入会取到重复主键，报
// duplicate key value violates unique constraint "..._pkey"（SQLSTATE 23505）。
// 该 Runner 注册在数据初始化之后，统一把序列推进到 max(id)+1。
type ZInitSequence struct {
	DbManger *dbManagerPg.Manager `autowire:"?"`
}

func (b *ZInitSequence) Run(ctx context.Context) error {
	log.Info(ctx, log.TagAppDef, "[init].[自增序列校正]===================")
	if nil == b.DbManger {
		return nil
	}
	if err := dbMakePg.SyncSequences(b.DbManger.Db()); nil != err {
		log.Error(ctx, log.TagAppDef, err, "自增序列校正失败")
	}
	return nil
}
