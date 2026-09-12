package service

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamCasbin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamResourceCasbinService))
}

// RamResourceCasbinService Casbin中间件
// @Description:
type RamResourceCasbinService struct {
	sv *repositoryRam.RamCasbinRepository `autowire:"?"`
}

// UpdateCasbin 更新Casbin 内权限规则
func (c *RamResourceCasbinService) UpdateCasbin(authorityId string, casbinInfos []*entityRam.RamResourceRelationEntity) error {
	c.ClearCasbin(0, authorityId)
	if nil != casbinInfos && len(casbinInfos) > 0 {
		maps := make(map[string]bool)
		rules := [][]string{}
		for _, v := range casbinInfos {
			if "" == v.Path {
				continue
			}
			if "" == v.Method {
				continue
			}
			md5 := cryptPg.Md5(v.Path + v.Method)
			//重复的直接跳过
			if _, ok := maps[md5]; ok {
				continue
			}
			maps[md5] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
		e := c.Casbin()
		success, _ := e.AddPolicies(rules)
		if !success {
			return errors.New("存在相同api,添加失败,请联系管理员")
		}
	}
	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateCasbinApi
//@description: API更新随动
//@param: oldPath string, newPath string, oldMethod string, newMethod string
//@return: error

func (c *RamResourceCasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := c.sv.Db().Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	e := c.Casbin()
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetPolicyPathByAuthorityId
//@description: 获取权限列表
//@param: authorityId string
//@return: pathMaps []request.CasbinInfo

func (c *RamResourceCasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []modRamCasbin.CasbinInfo) {
	e := c.Casbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	list, _ := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, modRamCasbin.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ClearCasbin
//@description: 清除匹配的权限
//@param: v int, p ...string
//@return: bool

func (c *RamResourceCasbinService) ClearCasbin(v int, p ...string) bool {
	e := c.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Casbin
//@description: 持久化到数据库  引入自定义规则
//@return: *casbin.Enforcer

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

func (c *RamResourceCasbinService) Casbin() *casbin.SyncedCachedEnforcer {
	once.Do(func() {
		a, err := gormadapter.NewAdapterByDB(c.sv.Db())
		if err != nil {
			log.Errorf(context.Background(), log.TagAppDef, "适配数据库失败请检查casbin表是否为InnoDB引擎! %+v", err)
			return
		}
		text := `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		m, err := model.NewModelFromString(text)
		if err != nil {
			log.Errorf(context.Background(), log.TagAppDef, "字符串加载模型失败! %+v", err)
			return
		}
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, a)
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}
