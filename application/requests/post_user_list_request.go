package requests

// The request used to create a list
type PostUserListRequest struct {
	// The name of the new list
	Name string `json:"name" binding:"required" validate:"required"`
}
