package modRamResourceRelation

type QuerySelectedCt struct {
	Code string `json:"code" label:"码值" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
}
