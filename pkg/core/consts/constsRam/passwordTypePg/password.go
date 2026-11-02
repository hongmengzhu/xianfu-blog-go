package passwordTypePg

// PasswordType 密码类型
type PasswordType string

const (
	Password PasswordType = "password" //密码
)

// Name 名称
func (this PasswordType) Name() string {
	switch this {
	case "password":
		return "密码"
	default:
		return "未知"
	}
}

// 值
func (this PasswordType) String() string {
	return string(this)
}
