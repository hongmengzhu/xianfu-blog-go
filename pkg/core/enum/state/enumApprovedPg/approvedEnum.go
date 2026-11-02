package enumApprovedPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// ApprovedState 审批状态
type ApprovedState string

const (
	ApprovedStateUnApproved    ApprovedState = "unApproved"    //未审批
	ApprovedStateApproved      ApprovedState = "approved"      //审批通过
	ApprovedStateRejected      ApprovedState = "rejected"      //审批不通过,驳回
	ApprovedStateSubmit        ApprovedState = "submit"        //提交
	ApprovedStateUnderApproval ApprovedState = "underApproval" //审批中
)

// Name 名称
func (this ApprovedState) Name() string {
	switch this {
	case "notApproved":
		return "未审批"
	case "approved":
		return "审批通过"
	case "rejected":
		return "驳回"
	case "submit":
		return "提交"
	case "underApproval":
		return "提交"
	default:
		return "未知"
	}
}

// 值
func (this ApprovedState) String() string {
	return string(this)
}

// 值
func (this ApprovedState) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this ApprovedState) IsEqual(id string) bool {
	return string(this) == id
}

var ApprovedStateMap = map[string]enumBasePg.EnumString{
	ApprovedStateUnApproved.String():    enumBasePg.EnumString{ApprovedStateUnApproved.String(), ApprovedStateUnApproved.Name()},
	ApprovedStateApproved.String():      enumBasePg.EnumString{ApprovedStateApproved.String(), ApprovedStateApproved.Name()},
	ApprovedStateRejected.String():      enumBasePg.EnumString{ApprovedStateRejected.String(), ApprovedStateRejected.Name()},
	ApprovedStateSubmit.String():        enumBasePg.EnumString{ApprovedStateSubmit.String(), ApprovedStateSubmit.Name()},
	ApprovedStateUnderApproval.String(): enumBasePg.EnumString{ApprovedStateUnderApproval.String(), ApprovedStateUnderApproval.Name()},
}

func IsExistApprovedState(id string) (ApprovedState, bool) {
	_, ok := ApprovedStateMap[id]
	return ApprovedState(id), ok
}

func IsExistApprovedStateMerchant(id string) (ApprovedState, bool) {
	_, ok := ApprovedStateMap[id]
	return ApprovedState(id), ok
}
