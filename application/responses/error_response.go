package responses

// A typed response object that can be used to create an error message for the
// json responses as well as for the logger
type ErrorResponse struct {
	// The text of the error itself
	ErrorMessage string `json:"error" binding:"required"`
}

// Satisfies the Error interface
func (e *ErrorResponse) Error() string {
	return e.ErrorMessage
}

func (e *ErrorResponse) String() string {
	return e.ErrorMessage
}

// Creates a new error response object
func NewError(msg string) ErrorResponse {
	return ErrorResponse{ErrorMessage: msg}
}
