package automatedPg

import "regexp"

// CREATE_CODE 自动创建标志
const CREATE_CODE = "系统自动建立"

// IsCreateCode 是否 自动创建标志
//
//	@Description:
//	@param str
//	@return bool
func IsCreateCode(str string) bool {
	return CREATE_CODE == str
}

// FormatVerify 格式验证
//
//	@Description: 字母数字,:,#,@,-
//	@param str
//	@return bool
func FormatVerify(str string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9:#@\\-]+$").MatchString(str)
}

// VerifyMax100 长度不能大于100
//
//	@Description: 字母数字,:,#,@,-
//	@param str
//	@return bool
func VerifyMax100(str string) bool {
	return len(str) > 100
}
