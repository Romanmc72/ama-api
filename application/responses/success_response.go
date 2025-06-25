package responses

import "time"

// Generic success response that contains no data
type SuccessResponse struct {
	// This response means that the whole operation was successful, but this
	// can be false if for example the delete request attempted to delete
	// something that does not exist. It is "successful" in the sense
	// that the object intended for deletion does not exist after the request,
	// and unsuccessful in that the request played no part in the deletion.
	Success bool `json:"success"`
	// The approximate time on the server when the operation was completed.
	Time int64 `json:"time"`
}

// Create a new success response with the given success status.
func NewSuccessResponse(success bool) SuccessResponse {
	return SuccessResponse{
		Success: success,
		Time:    time.Now().Unix(),
	}
}
