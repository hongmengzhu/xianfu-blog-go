package model

// BaseIdApprovedCt 基础 审批
type BaseIdApprovedCt[ID any] struct {
	Id       ID     `json:"id"`
	Approved string `json:"approved"`
}
