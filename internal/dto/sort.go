package dto

type Sort struct {
	Order string `form:"order,default=desc" binding:"omitempty,oneof=asc desc"`
	By    string `form:"by"`
}
