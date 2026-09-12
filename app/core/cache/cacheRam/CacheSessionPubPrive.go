package cacheRam

import (
	"context"
	"strings"

	"github.com/goccy/go-json"
	authTokenPg2 "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/authTokenPg"
	cacheAuthPubPrivPg2 "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheAuthPubPrivPg"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/sessionKeyTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typePubPrivePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/cryptoPg"
	"github.com/pangu-2/go-tools/tools/jsonPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/userPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CacheSessionPubPrive))
}

type CacheSessionPubPrive struct {
	sessionAk   *repositoryRam.RamAccountSessionAccessKeyRepository `autowire:"?"`
	pubPriveKey typePubPrivePg.PubPrive
}

func (c *CacheSessionPubPrive) SystemLoginKey(ctx context.Context, client clientPg.Client) model.AuthPubPriveDto {
	return c.LoginPubPriveKey(ctx, typeDomainPg.System, client, typeDomainPg.System.Code(), sessionKeyTypePg.Login)
}
func (c *CacheSessionPubPrive) SystemLoginKeyByBrowser(ctx context.Context) model.AuthPubPriveDto {
	return c.LoginPubPriveKey(ctx, typeDomainPg.System, clientPg.Browser, typeDomainPg.System.Code(), sessionKeyTypePg.Login)
}

func (c *CacheSessionPubPrive) ManageLoginKey(ctx context.Context, client clientPg.Client, tenantNo string) model.AuthPubPriveDto {
	return c.LoginPubPriveKey(ctx, typeDomainPg.Manage, client, tenantNo, sessionKeyTypePg.Login)
}

func (c *CacheSessionPubPrive) LoginPubPriveKey(ctx context.Context, typeDomain typeDomainPg.TypeDomain,
	client clientPg.Client, tenantNo string, keyType sessionKeyTypePg.SessionKeyType) model.AuthPubPriveDto {
	get, b := cacheAuthPubPrivPg2.Get(cacheAuthPubPrivPg2.Key(keyType, tenantNo, typeDomain, client))
	if b {
		return model.AuthPubPriveDto{
			PrivateKey: get.Private,
			PublicKey:  get.Public,
		}
	}
	//是否新生成密钥
	isMakeNewKey := true
	var privatePubKeyEnt entityRam.RamAsaJsonPrivatePublicKey
	info, result := c.sessionAk.FindByTenantNoAndTypeDomainInAndClientAndTypeAndState(ctx, tenantNo, typeDomain.Code(), client.Index(), keyType.Code())
	if result {
		isMakeNewKey = false
		if strPg.IsNotBlank(info.Data) {
			err := json.Unmarshal([]byte(info.Data), &privatePubKeyEnt)
			if err != nil {
				//解析失败，重新生成
				isMakeNewKey = true
			}
		}
	}
	return c.LoginPubPriveKeyByNew(ctx, isMakeNewKey, typeDomain, client, tenantNo, keyType, privatePubKeyEnt)
}

func (c *CacheSessionPubPrive) LoginPubPriveKeyByNew(ctx context.Context,
	isMakeNewKey bool,
	typeDomain typeDomainPg.TypeDomain,
	client clientPg.Client, tenantNo string, keyType sessionKeyTypePg.SessionKeyType,
	json entityRam.RamAsaJsonPrivatePublicKey) model.AuthPubPriveDto {
	//重新生成
	if isMakeNewKey {
		ret := cryptoPg.Sm2GenerateKey()
		if ret.ErrorIs() {
			log.Debugf(ctx, log.TagBizDef, "生成密钥失败")
			return model.AuthPubPriveDto{}
		}
		privatePubKey := ret.Data
		//先删除
		c.sessionAk.DeleteByTenantNoAndTypeDomainAndClientAndTypeAndState(ctx, tenantNo, typeDomain.Code(), client.Index(), keyType.Code())
		//
		json.Private = privatePubKey.PrivateKey
		json.Public = privatePubKey.PublicKey
		toJson, _ := jsonPg.ObjToJson(json)
		//
		save := entityRam.RamAccountSessionAccessKeyEntity{
			Ano:        "",
			Data:       toJson,
			No:         noPg.No(),
			TenantNo:   tenantNo,
			TypeDomain: typeDomain.Code(),
			Client:     client.Index(),
			Type:       keyType.Code(),
			Key:        privatePubKey.PublicKey,
			State:      enumStatePg.ENABLE.Index(),
		}
		str := toJson + save.TenantNo + save.TypeDomain + save.Client + keyType.Code()
		save.KindUnique = userPg.SaltMake(str, save.Key)
		err, _ := c.sessionAk.Create(ctx, &save)
		if err != nil {
			log.Debugf(ctx, log.TagBizDef, "创建新密钥失败.%+v", err)
		}
	}

	cacheAuthPubPrivPg2.Set(cacheAuthPubPrivPg2.Key(keyType, tenantNo, typeDomain, client), json)
	//
	return model.AuthPubPriveDto{
		PrivateKey: json.Private,
		PublicKey:  json.Public,
	}
}

func (c *CacheSessionPubPrive) DecodeByLogin(ctx context.Context, pwd string, typeDomain typeDomainPg.TypeDomain,
	client clientPg.Client, tenantNo string, keyType sessionKeyTypePg.SessionKeyType) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagBizDef, "ct.len=%+v,=%+v", len(pwd), pwd)
	if len(pwd) < 194 {
		return rt.ErrorMessage("密文过短，可能已损坏，实际长度：" + string(rune(len(pwd))))
	}
	key := c.LoginPubPriveKey(ctx, typeDomain, client, tenantNo, keyType)
	if !strings.HasPrefix(pwd, "04") {
		pwd = "04" + pwd
	}
	ret := cryptoPg.NewSm2(key.PrivateKey, key.PublicKey).DecodeHex(pwd)
	if ret.ErrorIs() {
		return rt.ErrorMessage(ret.Message)
	}
	return rt.OkData(ret.Data)
}

func (c *CacheSessionPubPrive) DecodeByLoginSystem(ctx context.Context, pwd string) (rt rg.Rs[string]) {
	return c.DecodeByLogin(ctx, pwd, typeDomainPg.System, clientPg.Browser, typeDomainPg.System.Code(), sessionKeyTypePg.Login)
}
func (c *CacheSessionPubPrive) DecodeByLoginManage(ctx context.Context, pwd string, client clientPg.Client, tenantNo string) (rt rg.Rs[string]) {
	return c.DecodeByLogin(ctx, pwd, typeDomainPg.Manage, client, tenantNo, sessionKeyTypePg.Login)
}

//	PaseKeyByNew 获取密钥对
//
// @Description:
// @receiver *CacheSessionPubPrive
// @param
// @return
func (c *CacheSessionPubPrive) PaseKeyByNew(ctx context.Context, isMakeNewKey bool,
	jsonEntity *entityRam.RamAsaJsonPrivatePublicKey,
	keyType sessionKeyTypePg.SessionKeyType, typeDomain typeDomainPg.TypeDomain,
	client clientPg.Client, tenantNo string) entityRam.RamAsaJsonPrivatePublicKey {
	privatePubKey := authTokenPg2.Result{}
	//
	// 需要重新生成
	if isMakeNewKey {
		jsonEntity = &entityRam.RamAsaJsonPrivatePublicKey{}
		privatePubKey = authTokenPg2.MakePublicPrivateKey()
		jsonEntity.Private = privatePubKey.PrivateKey
		jsonEntity.Public = privatePubKey.PublicKey
		// 删除已存在的
		c.sessionAk.DeleteByTenantNoAndTypeDomainAndClientAndTypeAndState(ctx, tenantNo, typeDomain.Code(), client.Index(), keyType.Code())
		//
		toJson, _ := jsonPg.ObjToJson(jsonEntity)
		save := entityRam.RamAccountSessionAccessKeyEntity{
			Ano:        "",
			Data:       toJson,
			No:         noPg.No(),
			TenantNo:   tenantNo,
			TypeDomain: typeDomain.Code(),
			Client:     client.Index(),
			Type:       keyType.Code(),
			Key:        privatePubKey.PublicKey,
			State:      enumStatePg.ENABLE.Index(),
		}
		str := toJson + save.TenantNo + save.TypeDomain + save.Client + keyType.Code()
		save.KindUnique = userPg.SaltMake(str, save.Key)
		err, _ := c.sessionAk.Create(ctx, &save)
		if err != nil {
			log.Debugf(ctx, log.TagBizDef, "创建新密钥失败.%+v", err)
		}
	} else {
		//toJson, _ := jsonPg.ObjToJson(jsonEntity)
	}
	entity := *jsonEntity
	//缓存
	cacheAuthPubPrivPg2.Set(cacheAuthPubPrivPg2.Key(keyType, tenantNo, typeDomain, client), entity)
	return entity
}

//	PaseKey 获取密钥对
//
// @Description:
// @receiver *CacheSessionPubPrive
// @param
// @return
func (c *CacheSessionPubPrive) PaseKey(ctx context.Context,
	keyType sessionKeyTypePg.SessionKeyType,
	typeDomain typeDomainPg.TypeDomain, client clientPg.Client, tenantNo string, jsonEntity *entityRam.RamAsaJsonPrivatePublicKey) entityRam.RamAsaJsonPrivatePublicKey {
	get, b := cacheAuthPubPrivPg2.Get(cacheAuthPubPrivPg2.Key(keyType, tenantNo, typeDomain, client))
	if b {
		return get
	}
	isMakeNewKey := true
	info, result := c.sessionAk.FindByTenantNoAndTypeDomainInAndClientAndTypeAndState(ctx, tenantNo, typeDomain.Code(), client.Index(), keyType.Code())
	if result {
		if strPg.IsNotBlank(info.Data) {
			isMakeNewKey = false
			jsonEntity = &entityRam.RamAsaJsonPrivatePublicKey{}
			err := json.Unmarshal([]byte(info.Data), jsonEntity)
			if err != nil {
				//解析失败，重新生成
				isMakeNewKey = true
			}
		}
	}
	return c.PaseKeyByNew(ctx, isMakeNewKey, jsonEntity, keyType, typeDomain, client, tenantNo)
}

//	PaseKeyAccessToken 获取密钥对
//
// @Description:
// @param
// @return
func (c *CacheSessionPubPrive) PaseKeyAccessToken(ctx context.Context,
	typeDomain typeDomainPg.TypeDomain, client clientPg.Client, tenantNo string, jsonEntity *entityRam.RamAsaJsonPrivatePublicKey) entityRam.RamAsaJsonPrivatePublicKey {
	return c.PaseKey(ctx, sessionKeyTypePg.AccessToken, typeDomain, client, tenantNo, jsonEntity)
}

//	PaseKeyRefreshToken 获取密钥对
//
// @Description:
// @param
// @return
func (c *CacheSessionPubPrive) PaseKeyRefreshToken(ctx context.Context,
	typeDomain typeDomainPg.TypeDomain, client clientPg.Client, tenantNo string, jsonEntity *entityRam.RamAsaJsonPrivatePublicKey) entityRam.RamAsaJsonPrivatePublicKey {
	return c.PaseKey(ctx, sessionKeyTypePg.RefreshToken, typeDomain, client, tenantNo, jsonEntity)
}
