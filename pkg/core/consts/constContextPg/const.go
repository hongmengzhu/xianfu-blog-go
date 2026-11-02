package constContextPg

const (
	//用于存储对象
	AUTH_LOGIN        = "pangu-go.baseAuth.authPg"
	AUTH_LOGIN_IS     = "pangu-go.baseAuth.authIs"
	AUTH_LOGIN_API    = "pangu-go.baseAuthApi.authPg"
	HOLDER            = "pangu-go.baseAuth.holder"
	HOLDER_IS         = "pangu-go.baseAuth.holderIs"
	CTX_MULITI_TENANT = "pangu-go.baseAuth.multiTenant" //租户规则
	CTX_MULITI        = "pangu-go.baseAuth.multi"       //多租户/商户规则
	CTX_RULE          = "pangu-go.baseAuth.rule"        //规则
	CTX               = "pangu-go.baseAuth.MULTI-CTX"
	CTX_LOG           = "pangu-go.baseAuth.CTX_LOG"
	CTX_CONFIG_STORE  = "pangu-go.config.store" //店铺默认配置
	CTX_DOMAIN        = "pangu-go.domain"       //域模式
)
