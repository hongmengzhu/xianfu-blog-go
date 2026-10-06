package service

import (
	"context"

	"time"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/ram/model/modRamAccount"
	"go-spring.org/log"
)

type AccountLoginLog struct {
	acc      *repositoryRam.RamAccountRepository         `autowire:"?"`
	loginLog *repositoryRam.RamAccountLoginLogRepository `autowire:"?"`
	session  *repositoryRam.RamAccountSessionRepository  `autowire:"?"`
}

func NewAccountLoginLog(
	acc *repositoryRam.RamAccountRepository,
	loginLog *repositoryRam.RamAccountLoginLogRepository,
	session *repositoryRam.RamAccountSessionRepository,
) *AccountLoginLog {
	return &AccountLoginLog{
		acc:      acc,
		loginLog: loginLog,
		session:  session,
	}
}

// Processor
//
//	@Description: 处理
//	@receiver c
//	@param data
func (c *AccountLoginLog) Processor(ctx context.Context, data modRamAccount.LoginLogDto) error {
	acc := data.Account
	ua := ""
	if data.ExtraData != nil {
		if get, ok := data.ExtraData[constHeaderPg.HeaderUserAgent]; ok {
			ua = get.(string)
		}
	}
	now := time.Now()
	//登录时间
	{
		save := entityRam.RamAccountEntity{LoginTime: &now}
		c.acc.Update(ctx, save, acc.ID)
	}
	//会话信息
	{
		session := entityRam.RamAccountSessionEntity{
			Ano:         acc.No,
			AppNo:       data.AppNo,
			Client:      data.Client,
			LoginSource: data.LoginSource,
			Os:          acc.Os,
			TenantNo:    acc.TenantNo,
			UserAgent:   ua,
			Ip:          data.Ip,
		}
		err, _ := c.session.Create(ctx, &session)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, err, "")
		}
	}
	//登录日志
	save := entityRam.RamAccountLoginLogEntity{
		AppNo:       data.AppNo,
		Ano:         acc.No,
		Client:      data.Client,
		LoginSource: data.LoginSource,
		Os:          acc.Os,
		TenantNo:    acc.TenantNo,
		UserAgent:   ua,
		Ip:          data.Ip,
	}
	err, _ := c.loginLog.Create(ctx, &save)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, err, "")
	}
	return err
}
