package modAttachment

// Ext 附件上传的扩展元数据，可由调用方传入任意实现
type Ext interface {
}

// OwnerExt 附件归属信息：上传时由业务层填充，落库后用于归属校验，
// 避免附件被其他账号/租户越权读取或改写。
type OwnerExt struct {
	// TenantNo 租户编码
	TenantNo string
	// CreateBy 上传人账号编码
	CreateBy string
	// OwnerNo 上传人主体编码
	OwnerNo string
}

// OwnerFromExt 从扩展参数中提取归属信息
func OwnerFromExt(ext Ext) (OwnerExt, bool) {
	owner, ok := ext.(OwnerExt)
	if !ok {
		ownerPtr, okPtr := ext.(*OwnerExt)
		if nil == ownerPtr || !okPtr {
			return OwnerExt{}, false
		}
		return *ownerPtr, true
	}
	return owner, true
}
