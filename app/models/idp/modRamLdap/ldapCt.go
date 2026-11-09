package modRamLdap

// LdapLoginCt LDAP 登录请求
type LdapLoginCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号" binding:"required"`
	Username string `json:"username" form:"username" label:"用户名" binding:"required"`
	Password string `json:"password" form:"password" label:"密码" binding:"required"`
}

// LdapTestCt LDAP 连接测试请求
type LdapTestCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号" binding:"required"`
}

// LdapSyncCt LDAP 用户同步请求
type LdapSyncCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号" binding:"required"`
}

// LdapSearchCt LDAP 用户搜索请求
type LdapSearchCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号" binding:"required"`
	Filter   string `json:"filter" form:"filter" label:"过滤条件"`
}
