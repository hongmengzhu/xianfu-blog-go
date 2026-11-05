package modCacheBasicConfig

// Item 配置缓存项。缓存 key 由 tenant_no + model 组成，值是 []Item 的 JSON。
type Item struct {
	ID           int64  `json:"id"`
	No           string `json:"no"`
	TenantNo     string `json:"tenantNo"`
	ModelNo      string `json:"modelNo"`
	Model        string `json:"model"`
	Module       string `json:"module"`
	ModuleSub    string `json:"moduleSub"`
	Field        string `json:"field"`
	FieldPath    string `json:"fieldPath"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Show         int8   `json:"show"`
	Binary       int8   `json:"binary"`
	ValueType    string `json:"valueType"`
	FormCode     string `json:"formCode"`
	DefaultValue string `json:"defaultValue"`
	Value        string `json:"value"`
	ValueBinary  string `json:"valueBinary"`
	CategoryNo   string `json:"categoryNo"`
	Sort         int64  `json:"sort"`
}

// ValueMap 返回 field -> value 的配置值集合；二进制字段返回 value_binary。
func (c Item) ValueMap() string {
	if c.Binary == 1 && c.ValueBinary != "" {
		return c.ValueBinary
	}
	return c.Value
}
