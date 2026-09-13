package configPg

// Database 数据库 配置
type Database struct {
	URL     string `value:"${url:=}"`
	Enabled bool   `value:"${enabled:=false}"`
	Driver  string `value:"${driver:='postgres'}"`
}
