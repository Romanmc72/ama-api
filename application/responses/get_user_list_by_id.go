package responses

import (
	"ama/api/application"
	"ama/api/application/list"
)

type GetUserListByIdResponse struct {
	List      list.List              `json:"list"`
	Questions []application.Question `json:"questions"`
}

func NewGetUserListByIdResponse(l list.List, q []application.Question) GetUserListByIdResponse {
	return GetUserListByIdResponse{
		List:      l,
		Questions: q,
	}
}
