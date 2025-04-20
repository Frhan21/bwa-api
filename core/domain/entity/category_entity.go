package entity

type CategoryEntity struct {
	ID    int
	Title string
	Slug  string
	User  UserEntity
}

// Error implements error.
func (c CategoryEntity) Error() string {
	panic("unimplemented")
}
