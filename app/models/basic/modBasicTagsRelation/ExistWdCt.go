package modBasicTagsRelation

type ExistWdCt struct {
	Id       string `json:"id"`
	Wd       string `json:"wd" label:"关键词" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Category string `json:"category" label:"分类" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
}
