package model

// BaseIdsApprovedCt 基础 审批
type BaseIdsApprovedCt[ID any] struct {
	Ids      []ID   `json:"ids" binding:"required" validate:"required"`
	Approved string `json:"approved" binding:"required" validate:"required"`
	Content  string `json:"content" `
}
