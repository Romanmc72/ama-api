package fixtures

import (
	"ama/api/application/list"
	"ama/api/application/requests"
)

const ListName = "List 1"

var ValidPostListRequest = requests.PostUserListRequest{
	Name: ListName,
}

var ValidPutUserListRequest = requests.PutUserListRequest{
	Name: ListName,
}

var ValidLists = []list.List{
	{
		ID:   ListId,
		Name: ListName,
	},
	{
		ID:   NewId,
		Name: list.LikedQuestionsListName,
	},
}
