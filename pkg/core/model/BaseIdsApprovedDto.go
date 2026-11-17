package model

// BaseIdsApprovedDto 基础 审批
type BaseIdsApprovedDto[ID any] struct {
	Ids      []ID   `json:"ids" binding:"required" validate:"required"`
	Approved string `json:"approved" binding:"required" validate:"required"`
}
