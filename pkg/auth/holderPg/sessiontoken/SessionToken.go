package sessiontoken

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

type ISessionTokenPg interface {
}

type SessionToken struct {
	Id int64
	/**
	 * 版本号
	 */
	Version string

	/**
	 * 创建人
	 */
	Creator string

	/**
	 * 创建人ID
	 */
	CreatorId int64

	/**
	 * 创建时间
	 */
	CreatedTime time.Time

	/**
	 * 登陆人Ip地址
	 */
	LoginIp string

	/**
	 * 登录地址
	 */
	LoginLocation string

	AccessToken AccessToken

	/**
	 * 0 在线 10已刷新 20 离线
	 */
	Status int
	Holder interfaces.IHolderPg //用户Session 会话信息
	/**
	 * 额外的，扩展
	 */
	Extra interface{}
	/**
	 * 其他对象
	 */
	Data interface{}
}
