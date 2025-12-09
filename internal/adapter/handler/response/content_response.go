package response

type ContentResponse struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Excerpt     string   `json:"excerpt"`
	Description string   `json:"description" omitempty:"description"`
	Image       string   `json:"image"`
	PublicId    string   `json:"public_id" omitempty:"public_id"`
	Tags        []string `json:"tags" omitempty:"tags"`
	Status      string   `json:"status" omitempty:"status"`
	CategoryId  int64    `json:"category_id" omitempty:"category_id"`
	UserID      int64    `json:"user_id" omitempty:"user_id"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
	User        string   `json:"user"`
	Category    string   `json:"category"`
}
