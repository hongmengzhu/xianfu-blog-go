package authTokenPg

type ResultAccessRefresh struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
