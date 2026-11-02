package bmKeyPg

// BookmarkCategoryTenantNo 文章分类租户编号
func BookmarkCategoryTenantNo(tenantNo, code string) string {
	return "bm:bmCate:" + tenantNo + ":" + code
}

// BookmarkCategoryTenantNoKeys 文章分类租户编号,所有键
func BookmarkCategoryTenantNoKeys(tenantNo string) string {
	return "bm:bmCateKeys:" + tenantNo
}

// BookmarkCategoryTenantNoAndNoByCode 文章分类租户编号和编号
func BookmarkCategoryTenantNoAndNoByCode(tenantNo, code string) string {
	return "bm:bmCateNoCode:" + tenantNo + ":" + code
}
