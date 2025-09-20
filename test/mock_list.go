package test

import (
	"ama/api/application"
	"ama/api/application/list"
)

type MockListManager struct {
	MockReadUser               func(id string) (application.User, error)
	MockReadQuestion           func(id string) (application.Question, error)
	MockReadQuestions          func(limit int, finalId string, tags []string) ([]application.Question, error)
	MockReadList               func(userId string, listId string, limit int, finalId string, tags []string) (list.List, []application.Question, error)
	MockCreateList             func(userId string, list list.List) (list.List, error)
	MockAddQuestionToList      func(userId string, listId string, question application.Question) error
	MockRemoveQuestionFromList func(userId string, listId string, questionId string) error
	MockUpdateList             func(userId string, updatedList list.List) error
	MockDeleteList             func(userId string, listId string) error
}

type MockListManagerConfig struct {
	ReadUser               func(id string) (application.User, error)
	ReadQuestion           func(id string) (application.Question, error)
	ReadQuestions          func(limit int, finalId string, tags []string) ([]application.Question, error)
	ReadList               func(userId string, listId string, limit int, finalId string, tags []string) (list.List, []application.Question, error)
	CreateList             func(userId string, list list.List) (list.List, error)
	AddQuestionToList      func(userId string, listId string, question application.Question) error
	RemoveQuestionFromList func(userId string, listId string, questionId string) error
	UpdateList             func(userId string, updatedList list.List) error
	DeleteList             func(userId string, listId string) error
}

func NewMockListManager(cfg MockListManagerConfig) *MockListManager {
	return &MockListManager{
		MockReadUser:               cfg.ReadUser,
		MockReadQuestion:           cfg.ReadQuestion,
		MockReadQuestions:          cfg.ReadQuestions,
		MockReadList:               cfg.ReadList,
		MockCreateList:             cfg.CreateList,
		MockAddQuestionToList:      cfg.AddQuestionToList,
		MockRemoveQuestionFromList: cfg.RemoveQuestionFromList,
		MockUpdateList:             cfg.UpdateList,
		MockDeleteList:             cfg.DeleteList,
	}
}

func (m *MockListManager) ReadUser(id string) (application.User, error) {
	if m.MockReadUser != nil {
		return m.MockReadUser(id)
	}
	return application.User{}, nil
}

func (m *MockListManager) ReadQuestion(id string) (application.Question, error) {
	if m.MockReadQuestion != nil {
		return m.MockReadQuestion(id)
	}
	return application.Question{}, nil
}

func (m *MockListManager) ReadQuestions(limit int, finalId string, tags []string) ([]application.Question, error) {
	if m.MockReadQuestions != nil {
		return m.MockReadQuestions(limit, finalId, tags)
	}
	return []application.Question{}, nil
}

func (m *MockListManager) ReadList(userId string, listId string, limit int, finalId string, tags []string) (list.List, []application.Question, error) {
	if m.MockReadList != nil {
		return m.MockReadList(userId, listId, limit, finalId, tags)
	}
	return list.List{}, []application.Question{}, nil
}

func (m *MockListManager) CreateList(userId string, l list.List) (list.List, error) {
	if m.MockCreateList != nil {
		return m.MockCreateList(userId, l)
	}
	return list.List{}, nil
}

func (m *MockListManager) AddQuestionToList(userId string, listId string, question application.Question) error {
	if m.MockAddQuestionToList != nil {
		return m.MockAddQuestionToList(userId, listId, question)
	}
	return nil
}

func (m *MockListManager) RemoveQuestionFromList(userId string, listId string, questionId string) error {
	if m.MockRemoveQuestionFromList != nil {
		return m.MockRemoveQuestionFromList(userId, listId, questionId)
	}
	return nil
}

func (m *MockListManager) UpdateList(userId string, updatedList list.List) error {
	if m.MockUpdateList != nil {
		return m.MockUpdateList(userId, updatedList)
	}
	return nil
}

func (m *MockListManager) DeleteList(userId string, listId string) error {
	if m.MockDeleteList != nil {
		return m.MockDeleteList(userId, listId)
	}
	return nil
}
