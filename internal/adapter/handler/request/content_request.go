package request

type CreateContent struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Excerpt     string `json:"excerpt" validate:"required"`
	CategoryID  int    `json:"category_id" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Tags        string `json:"tags" validate:"required"`
	Status      string `json:"status" validate:"required"`
}
