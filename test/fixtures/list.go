package fixtures

import (
	"ama/api/application/list"
	"ama/api/application/requests"
)

const ListName = "List 1"

var ValidListRequest = requests.PutListRequest{
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
