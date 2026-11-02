package modApiDiplAccessKey

type MakeNewCt struct {
	No            string `json:"no" label:"编号代号" validate:"required,min=0,max=255"`
	PrincipalType string `json:"principalType" label:"主体类型" validate:"required,min=0,max=255"`
}
