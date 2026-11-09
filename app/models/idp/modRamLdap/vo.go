package modRamLdap

// LdapUserVo LDAP 用户信息视图
type LdapUserVo struct {
	Uid         string   `json:"uid"`
	Cn          string   `json:"cn"`
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email"`
	Mobile      string   `json:"mobile"`
	Uuid        string   `json:"uuid"`
	MemberOf    []string `json:"memberOf"`
}

// LdapTestVo LDAP 连接测试响应
type LdapTestVo struct {
	Connected bool   `json:"connected"`
	IsAD      bool   `json:"isAd"`
	Message   string `json:"message"`
	UserCount int    `json:"userCount"`
}

// LdapSyncVo LDAP 同步结果响应
type LdapSyncVo struct {
	NewUsers     int      `json:"newUsers"`
	UpdatedUsers int      `json:"updatedUsers"`
	FailedUsers  int      `json:"failedUsers"`
	Errors       []string `json:"errors,omitempty"`
}
