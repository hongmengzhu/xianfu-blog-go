package model

// BaseIdsApprovedCt 基础 审批
type BaseIdsApprovedCt[ID any] struct {
	Ids      []ID   `json:"ids" validate:"required"`
	Approved string `json:"approved" validate:"required"`
	Content  string `json:"content" `
}
