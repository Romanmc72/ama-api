package interfaces

import (
	"ama/api/application"
	"ama/api/application/list"
)

type ListReader interface {
	UserReader
	ReadList(userId string, listId string, limit int, finalId string, tags []string) (list.List, []application.Question, error)
}

type ListCreator interface {
	CreateList(userId string, list list.List) (list.List, error)
}

type ListUpdater interface {
	QuestionReader
	AddQuestionToList(userId string, listId string, question application.Question) error
	RemoveQuestionFromList(userId string, listId string, questionId string) error
	UpdateList(userId string, updatedList list.List) error
}

type ListDeleter interface {
	UserReader
	DeleteList(userId string, listId string) error
}

type ListManager interface {
	ListReader
	ListCreator
	ListUpdater
	ListDeleter
}
