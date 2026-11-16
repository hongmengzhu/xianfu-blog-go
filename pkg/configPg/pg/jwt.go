package pg

// jwt
type JwtConfig struct {
	//
	Secret string `json:"secret" value:"${secret:=ssa1ssssaaa234324324394dsfdfs}"`
	// 过期时间，单位 分钟，默认 720（12 小时）；1440:1天
	Expire int `json:"expire" value:"${expire:=720}"`
	// cookie 是否启用
	CookieOn bool `json:"cookieOn" value:"${cookieOn:=false}"`
	//
	Issuer   string `json:"issuer" value:"${issuer:=foxwho.com}" label:"jwt签发者"`
	Audience string `json:"audience" value:"${audience:=pcWeb}" label:"接收jwt的一方"`
}

// 项目
type Jwt struct {
	// 后台
	Admin JwtConfig `json:"admin" value:"${admin}"`
	// 管理
	Manage JwtConfig `json:"manage" value:"${manage}"`
	//系统
	System JwtConfig `json:"system" value:"${system}"`
	// 用户
	User JwtConfig `json:"user" value:"${user}"`
	// 前台
	Web JwtConfig `json:"web" value:"${web}"`
	// app
	App JwtConfig `json:"app" value:"${app}"`
}
