package constTags

type Article string

const (
	ArticleInfo Article = "blog:article" //文章
	CollectInfo Article = "blog:collect" //收集
)

func (this Article) String() string {
	return string(this)
}

func (this Article) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Article) IsEqual(id string) bool {
	return string(this) == id
}
