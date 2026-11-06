package bmKeyPg

// ToolCategoryTenantNo 文章分类租户编号
func ToolCategoryTenantNo(tenantNo, code string) string {
	return "bm:toCate:" + tenantNo + ":" + code
}

// ToolCategoryTenantNoKeys 文章分类租户编号,所有键
func ToolCategoryTenantNoKeys(tenantNo string) string {
	return "bm:toCateKeys:" + tenantNo
}

// ToolCategoryTenantNoAndNoByCode 文章分类租户编号和编号
func ToolCategoryTenantNoAndNoByCode(tenantNo, code string) string {
	return "bm:toCateNoCode:" + tenantNo + ":" + code
}
