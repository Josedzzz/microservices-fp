package api

// PaginationQuery represents common pagination query parameters
// It can be used with Gin's ShouldBindQuery.
type PaginationQuery struct {
	Page       int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize   int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	Department string `form:"department" json:"department"`
	Status     string `form:"status" json:"status" binding:"omitempty,oneof=ACTIVE ON_VACATION RETIRED"`
	Position   string `form:"position" json:"position"`
}

// PaginatedResponse is a generic structure for paginated results
type PaginatedResponse struct {
	Data       any            `json:"data"` // Can hold any slice ([]models.Employee to be concrete). Maybe "any" can be replaced by interface{}?
	Pagination PaginationMeta `json:"pagination"`
}

// PaginationMeta contains metadata about the pagination
type PaginationMeta struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}
