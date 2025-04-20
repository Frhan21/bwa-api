package response

type SuccessCategoryResponse struct {
	ID    int64        `json:"id"`
	Title string       `json:"title"`
	Slug  string       `json:"slug"`
	User  UserResponse `json:"user"`
}
