package modBlogArticle

type StatisticsVo struct {
	Comment        int64  `json:"comment" label:"评论" `
	Read           int64  `json:"read" label:"阅读" `
	SeoKeywords    string `json:"seoKeywords" label:"seo关键词" `
	SeoDescription string `json:"seoDescription" label:"seo描述" `
	PageTitle      string `json:"pageTitle" label:"网页标题" `
}
