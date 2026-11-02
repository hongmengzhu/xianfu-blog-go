package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	goldap "github.com/go-ldap/ldap/v3"
)

// LdapConn LDAP 连接封装
type LdapConn struct {
	Conn *goldap.Conn
	IsAD bool
}

// LdapUser LDAP 用户信息
type LdapUser struct {
	UidNumber   string `json:"uidNumber"`
	Uid         string `json:"uid"`
	Cn          string `json:"cn"`
	GidNumber   string `json:"gidNumber"`
	Uuid        string `json:"uuid"`
	DisplayName string `json:"displayName"`
	Mail        string `json:"mail"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Telephone   string `json:"telephoneNumber"`
	Country     string `json:"country"`
	Address     string `json:"address"`
	MemberOf    []string `json:"memberOf"`
}

// LdapConfig LDAP 连接配置
type LdapConfig struct {
	Host                string
	Port                int
	EnableSsl           bool
	AllowSelfSignedCert bool
	Username            string // bind DN
	Password            string
	BaseDn              string
	Filter              string
	FilterFields        []string
}

// GetLdapConn 建立 LDAP 连接并绑定
func GetLdapConn(cfg *LdapConfig) (*LdapConn, error) {
	var conn *goldap.Conn
	var err error

	tlsConfig := tls.Config{
		InsecureSkipVerify: cfg.AllowSelfSignedCert,
	}
	if cfg.EnableSsl {
		conn, err = goldap.DialTLS("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), &tlsConfig)
	} else {
		conn, err = goldap.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
	}
	if err != nil {
		return nil, fmt.Errorf("LDAP 连接失败: %w", err)
	}

	err = conn.Bind(cfg.Username, cfg.Password)
	if err != nil {
		conn.Unbind()
		return nil, fmt.Errorf("LDAP 绑定失败: %w", err)
	}

	isAD, err := isMicrosoftAD(conn)
	if err != nil {
		conn.Unbind()
		return nil, fmt.Errorf("检测 LDAP 服务器类型失败: %w", err)
	}

	return &LdapConn{Conn: conn, IsAD: isAD}, nil
}

// Close 关闭 LDAP 连接
func (l *LdapConn) Close() {
	if l.Conn == nil {
		return
	}
	_ = l.Conn.Unbind()
}

// SearchUsers 搜索 LDAP 用户
func (l *LdapConn) SearchUsers(cfg *LdapConfig) ([]LdapUser, error) {
	searchAttrs := []string{
		"uidNumber", "cn", "sn", "gidNumber", "entryUUID", "displayName",
		"mail", "email", "telephoneNumber", "mobile", "c", "memberOf",
	}
	if l.IsAD {
		searchAttrs = append(searchAttrs, "sAMAccountName")
	} else {
		searchAttrs = append(searchAttrs, "uid")
	}

	filter := strings.TrimSpace(cfg.Filter)
	if filter == "" {
		filter = "(objectClass=*)"
	} else if strings.Contains(filter, "{}") {
		filter = strings.ReplaceAll(filter, "{}", "*")
	}

	searchReq := goldap.NewSearchRequest(cfg.BaseDn, goldap.ScopeWholeSubtree, goldap.NeverDerefAliases,
		0, 0, false,
		filter, searchAttrs, nil)

	searchResult, err := l.Conn.SearchWithPaging(searchReq, 100)
	if err != nil {
		return nil, fmt.Errorf("LDAP 搜索失败: %w", err)
	}

	if len(searchResult.Entries) == 0 {
		return []LdapUser{}, nil
	}

	var users []LdapUser
	for _, entry := range searchResult.Entries {
		var user LdapUser
		for _, attr := range entry.Attributes {
			switch attr.Name {
			case "uidNumber":
				user.UidNumber = attr.Values[0]
			case "uid":
				user.Uid = attr.Values[0]
			case "sAMAccountName":
				user.Uid = attr.Values[0]
			case "cn":
				user.Cn = attr.Values[0]
			case "gidNumber":
				user.GidNumber = attr.Values[0]
			case "entryUUID":
				user.Uuid = attr.Values[0]
			case "displayName":
				user.DisplayName = attr.Values[0]
			case "mail":
				user.Mail = attr.Values[0]
			case "email":
				user.Email = attr.Values[0]
			case "telephoneNumber":
				user.Telephone = attr.Values[0]
			case "mobile":
				user.Mobile = attr.Values[0]
			case "c":
				user.Country = attr.Values[0]
			case "memberOf":
				user.MemberOf = attr.Values
			}
		}
		users = append(users, user)
	}

	return users, nil
}

// CheckUserPassword 验证 LDAP 用户密码
func (l *LdapConn) CheckUserPassword(cfg *LdapConfig, username string, password string) (*LdapUser, error) {
	filter := buildAuthFilterString(cfg, username, l.IsAD)

	searchReq := goldap.NewSearchRequest(cfg.BaseDn, goldap.ScopeWholeSubtree, goldap.NeverDerefAliases,
		0, 0, false,
		filter, []string{"dn", "uid", "cn", "displayName", "mail", "mobile"}, nil)

	searchResult, err := l.Conn.Search(searchReq)
	if err != nil {
		return nil, fmt.Errorf("LDAP 搜索用户失败: %w", err)
	}

	if len(searchResult.Entries) == 0 {
		return nil, fmt.Errorf("用户 %s 不存在", username)
	}
	if len(searchResult.Entries) > 1 {
		return nil, fmt.Errorf("存在多个同名用户 %s，请检查 LDAP 服务器", username)
	}

	entry := searchResult.Entries[0]
	userDn := entry.DN

	// 使用用户 DN 和密码重新绑定以验证密码
	err = l.Conn.Bind(userDn, password)
	if err != nil {
		// 重新绑定回管理员账号
		_ = l.Conn.Bind(cfg.Username, cfg.Password)
		return nil, fmt.Errorf("密码验证失败")
	}

	// 重新绑定回管理员账号
	_ = l.Conn.Bind(cfg.Username, cfg.Password)

	user := &LdapUser{
		Uuid: entry.GetAttributeValue("entryUUID"),
		Cn:   entry.GetAttributeValue("cn"),
		Mail: entry.GetAttributeValue("mail"),
	}
	if l.IsAD {
		user.Uid = entry.GetAttributeValue("sAMAccountName")
	} else {
		user.Uid = entry.GetAttributeValue("uid")
	}
	user.DisplayName = entry.GetAttributeValue("displayName")
	user.Mobile = entry.GetAttributeValue("mobile")

	if user.Uuid == "" {
		user.Uuid = user.Uid
	}
	if user.Uuid == "" {
		user.Uuid = user.Cn
	}

	return user, nil
}

// buildAuthFilterString 构建认证过滤器
func buildAuthFilterString(cfg *LdapConfig, username string, isAD bool) string {
	baseFilter := strings.TrimSpace(cfg.Filter)
	if baseFilter == "" {
		baseFilter = "(objectClass=*)"
	} else if strings.Contains(baseFilter, "{}") {
		baseFilter = strings.ReplaceAll(baseFilter, "{}", "*")
	}

	escapedUsername := goldap.EscapeFilter(username)

	if len(cfg.FilterFields) == 0 {
		field := "uid"
		if isAD {
			field = "sAMAccountName"
		}
		return fmt.Sprintf("(&%s(%s=%s))", baseFilter, field, escapedUsername)
	}

	filter := fmt.Sprintf("(&%s(|", baseFilter)
	for _, field := range cfg.FilterFields {
		filter = fmt.Sprintf("%s(%s=%s)", filter, field, escapedUsername)
	}
	filter = fmt.Sprintf("%s))", filter)

	return filter
}

// isMicrosoftAD 检测是否为 Microsoft Active Directory
func isMicrosoftAD(conn *goldap.Conn) (bool, error) {
	searchFilter := "(objectClass=*)"
	searchAttrs := []string{"vendorname", "vendorversion", "isGlobalCatalogReady", "forestFunctionality"}

	searchReq := goldap.NewSearchRequest("",
		goldap.ScopeBaseObject, goldap.NeverDerefAliases, 0, 0, false,
		searchFilter, searchAttrs, nil)

	searchResult, err := conn.Search(searchReq)
	if err != nil {
		return false, err
	}
	if len(searchResult.Entries) == 0 {
		return false, nil
	}

	var vendorName, vendorVersion, isGlobalCatalogReady, forestFunctionality string
	for _, entry := range searchResult.Entries {
		for _, attr := range entry.Attributes {
			switch attr.Name {
			case "vendorname":
				vendorName = attr.Values[0]
			case "vendorversion":
				vendorVersion = attr.Values[0]
			case "isGlobalCatalogReady":
				isGlobalCatalogReady = attr.Values[0]
			case "forestFunctionality":
				forestFunctionality = attr.Values[0]
			}
		}
	}

	isMicrosoft := vendorName == "" &&
		vendorVersion == "" &&
		isGlobalCatalogReady == "TRUE" &&
		forestFunctionality != ""

	return isMicrosoft, nil
}

// GetLdapUuid 获取用户的唯一标识
func (u *LdapUser) GetLdapUuid() string {
	if u.Uuid != "" {
		return u.Uuid
	}
	if u.Uid != "" {
		return u.Uid
	}
	return u.Cn
}

// GetEmail 获取用户邮箱（优先取非空值）
func (u *LdapUser) GetEmail() string {
	if u.Email != "" {
		return u.Email
	}
	if u.Mail != "" {
		return u.Mail
	}
	return ""
}

// GetDisplayName 获取用户显示名
func (u *LdapUser) GetDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Cn
}
