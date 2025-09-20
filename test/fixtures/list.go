package fixtures

import (
	"ama/api/application"
	"ama/api/application/list"
	"ama/api/application/requests"
	"ama/api/application/responses"
)

const ListName = "List 1"

var ValidPostListRequest = requests.PostUserListRequest{
	Name: ListName,
}

var ValidPutUserListRequest = requests.PutUserListRequest{
	Name: ListName,
}

var ValidList = list.List{
	ID:   ListId,
	Name: ListName,
}

var ValidLists = []list.List{
	ValidList,
	list.GetLikedQuestionList(),
}

var ValidGetUserListByIdResponse = responses.NewGetUserListByIdResponse(
	ValidList, []application.Question{ValidQuestion},
)
