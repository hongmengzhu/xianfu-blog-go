package modPublic

// LoginToken 刷新令牌返回体。
//
// 登录接口返回 camelCase（accessToken/refreshToken），历史刷新接口只返回
// access_token/refreshToken，字段名不一致导致前端刷新时取不到访问令牌。
// 这里同时输出两种命名，保持向下兼容。
type LoginToken struct {
	AccessToken      string `json:"access_token"`
	AccessTokenCamel string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	ExpiresIn        int64  `json:"expiresIn"`
}
