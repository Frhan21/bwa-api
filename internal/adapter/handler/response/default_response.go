package response

type ErrorResponseDefault struct {
	Meta
}

type Meta struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type DefaultSuccessResponse struct {
	Meta   Meta                `json:"meta"`
	Data   any                 `json:"data,omitempty"`
	Pagina *PaginationResponse `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	TotalRecords int64 `json:"total_records"`
	TotalPages   int64 `json:"total_pages"`
	Page         int64 `json:"page"`
	PerPage      int64 `json:"per_page"`
}
