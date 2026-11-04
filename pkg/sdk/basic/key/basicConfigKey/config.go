package basicConfigKey

// ConfigByTenantNoAndModel 配置值缓存键：tenant_no + model。
// 其他服务读取某个租户、某个模型的全部配置时，统一使用该键。
func ConfigByTenantNoAndModel(tenantNo, model string) string {
	return "basicConfig:config:" + tenantNo + ":" + model
}

// ConfigByTenantNoAndModelNo 配置值缓存键（按模型编号），用于模型编号与 model 不一致的内部场景。
func ConfigByTenantNoAndModelNo(tenantNo, modelNo string) string {
	return "basicConfig:configNo:" + tenantNo + ":" + modelNo
}

// ConfigKeysByTenantNo 租户下所有配置缓存键集合，用于缓存重建/清理。
func ConfigKeysByTenantNo(tenantNo string) string {
	return "basicConfig:configKeys:" + tenantNo
}
