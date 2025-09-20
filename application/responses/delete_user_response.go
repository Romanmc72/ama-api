package responses

type DeleteUserResponse struct {
	Success    bool  `json:"success"`
	DeleteTime int64 `json:"deleteTime"`
}

func NewDeleteUserResponse(success bool, deleteTime int64) DeleteUserResponse {
	return DeleteUserResponse{
		Success:    success,
		DeleteTime: deleteTime,
	}
}
