package responses

import (
	"ama/api/application"
	"ama/api/application/list"
)

type GetUserListByIdResponse struct {
	List      list.List              `json:"list"`
	Questions []application.Question `json:"questions"`
}

func NewGetUserListByIdResponse(list list.List, questions []application.Question) GetUserListByIdResponse {
	return GetUserListByIdResponse{
		List:      list,
		Questions: questions,
	}
}
