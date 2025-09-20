package test

import (
	"ama/api/application"
	"ama/api/interfaces"
	"time"
)

// MockQuestionManager implements QuestionManager interface for testing
type MockQuestionManager struct {
	MockCreateQuestion func(questionData interfaces.QuestionConverter) (application.Question, error)
	MockReadQuestion   func(id string) (application.Question, error)
	MockReadQuestions  func(limit int, finalId string, tags []string) ([]application.Question, error)
	MockUpdateQuestion func(id string, questionData interfaces.QuestionConverter) (application.Question, error)
	MockDeleteQuestion func(id string) (time.Time, error)
}

type MockQuestionManagerConfig struct {
	CreateQuestion func(questionData interfaces.QuestionConverter) (application.Question, error)
	ReadQuestion   func(id string) (application.Question, error)
	ReadQuestions  func(limit int, finalId string, tags []string) ([]application.Question, error)
	UpdateQuestion func(id string, questionData interfaces.QuestionConverter) (application.Question, error)
	DeleteQuestion func(id string) (time.Time, error)
}

// NewMockQuestionManager creates a new instance of MockQuestionManager
func NewMockQuestionManager(cfg MockQuestionManagerConfig) *MockQuestionManager {
	return &MockQuestionManager{
		MockCreateQuestion: cfg.CreateQuestion,
		MockReadQuestion:   cfg.ReadQuestion,
		MockReadQuestions:  cfg.ReadQuestions,
		MockUpdateQuestion: cfg.UpdateQuestion,
		MockDeleteQuestion: cfg.DeleteQuestion,
	}
}

// CreateQuestion implements QuestionWriter interface
func (m *MockQuestionManager) CreateQuestion(questionData interfaces.QuestionConverter) (application.Question, error) {
	if m.MockCreateQuestion != nil {
		return m.MockCreateQuestion(questionData)
	}
	return application.Question{}, nil
}

// UpdateQuestion implements QuestionWriter interface
func (m *MockQuestionManager) UpdateQuestion(id string, questionData interfaces.QuestionConverter) (application.Question, error) {
	if m.MockUpdateQuestion != nil {
		return m.MockUpdateQuestion(id, questionData)
	}
	return application.Question{}, nil
}

// ReadQuestion implements QuestionReader interface
func (m *MockQuestionManager) ReadQuestion(id string) (application.Question, error) {
	if m.MockReadQuestion != nil {
		return m.MockReadQuestion(id)
	}
	return application.Question{}, nil
}

// ReadQuestions implements QuestionReader interface
func (m *MockQuestionManager) ReadQuestions(limit int, finalId string, tags []string) ([]application.Question, error) {
	if m.MockReadQuestions != nil {
		return m.MockReadQuestions(limit, finalId, tags)
	}
	return []application.Question{}, nil
}

// DeleteQuestion implements QuestionDeleter interface
func (m *MockQuestionManager) DeleteQuestion(id string) (time.Time, error) {
	if m.MockDeleteQuestion != nil {
		return m.MockDeleteQuestion(id)
	}
	return time.Now(), nil
}
