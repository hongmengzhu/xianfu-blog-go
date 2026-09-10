package holderPg

import (
	"encoding/json"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// HolderSimple 用户Session 会话信息 登录人信息
type HolderSimple struct {
	MultiTenant   interfaces.IMultiTenantPg `json:"multiTenant"`   //多租户
	LoginId       int64                     `json:"loginId"`       //登录用户ID
	No            string                    `json:"no"`            //系统编号
	LoginUserName string                    `json:"loginUserName"` //登录用户名
	Name          string                    `json:"name"`          //显示名称
	LevelId       int64                     `json:"levelId"`       //级别ID
	GroupId       int64                     `json:"groupId"`       //组织ID
	GroupName     string                    `json:"groupName"`     //组织名称
	Type          int64                     `json:"type"`          //类型
	Token         string                    `json:"token"`         //令牌
	Other         string                    `json:"other"`         //其他信息
	Version       string                    `json:"version"`       //版本
	Extra         interface{}               `json:"extra"`         //额外的，扩展
}

func (s HolderSimple) Get() HolderSimple {
	return s
}

func (s HolderSimple) ToJsonString() (string, error) {
	marshal, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
