package cacheDiplPg

import "time"

type DiplCo struct {
	Name          string     `json:"name"`
	No            string     `json:"no"`
	Ano           string     `json:"ano"`
	TenantNo      string     `json:"tenantNo"`
	Key           string     `json:"key"`
	Secret        string     `json:"secret"`
	HashSha       string     `json:"hashSha"`
	PrincipalType string     `json:"principalType"`
	State         int8       `json:"state"`
	ExpiryDate    *time.Time `json:"expiryDate"`
}
