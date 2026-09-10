package templatePg

// SitePage 站点信息
type SitePage struct {
	Title           string `json:"title" label:"网页标题"`
	Description     string `json:"description" label:"网页描述"`
	Keywords        string `json:"keywords" label:"网页关键字"`
	SiteName        string `json:"siteName" label:"网站名称"`
	SiteDescription string `json:"siteDescription" label:"网站描述"`
	SiteUrl         string `json:"siteUrl" label:"网站地址"`
	SiteLogo        string `json:"siteLogo" label:"网站logo"`
}
