package authTokenPg

type Result struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Token      string `json:"token"`
}
