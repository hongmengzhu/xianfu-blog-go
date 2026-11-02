package constPagePg

import "regexp"

// FormatVerify 格式验证
//
//	@Description: 字母数字,:,#,@,-
//	@param str
//	@return bool
func FormatVerify(str string) bool {
	return regexp.MustCompile("^[a-zA-Z0-9:@_\\-]+$").MatchString(str)
}
