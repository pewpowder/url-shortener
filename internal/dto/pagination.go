package dto

type Pagination struct {
	Page int `form:"page,default=1" binding:"omitempty,gte=1"`
	Size int `form:"size,default=10" binding:"omitempty,gte=1"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Size
}
