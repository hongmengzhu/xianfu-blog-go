package model

type AuthPubPriveDto struct {
	PrivateKey string `json:"privateKey" label:""`
	PublicKey  string `json:"publicKey"`
	Type       string `json:"type" label:"类型"`
	No         string `json:"no" label:"编号"`
}
