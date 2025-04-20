package request

// CreateCategoryRequest struct
type CreateCategoryRequest struct {
	Title string `json:"title" validate:"required"`
}
