package requests

// The request used to change the list name.
type PutListRequest struct {
	// The name for the list.
	Name string `json:"name" binding:"required"`
}
