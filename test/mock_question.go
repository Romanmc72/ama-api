package test

import (
	"ama/api/application"
	"ama/api/interfaces"
	"errors"
	"time"
)

// MockQuestionManager implements QuestionManager interface for testing
type MockQuestionManager struct {
	// Store questions for simulating database
	questions map[string]application.Question

	// Track method calls for verification in tests
	CreateQuestionCalls []interfaces.QuestionConverter
	UpdateQuestionCalls map[string]interfaces.QuestionConverter
	ReadQuestionCalls   []string
	DeleteQuestionCalls []string

	// Control mock behavior
	ShouldError bool
}

// NewMockQuestionManager creates a new instance of MockQuestionManager
func NewMockQuestionManager() *MockQuestionManager {
	return &MockQuestionManager{
		questions:           make(map[string]application.Question),
		UpdateQuestionCalls: make(map[string]interfaces.QuestionConverter),
	}
}

// CreateQuestion implements QuestionWriter interface
func (m *MockQuestionManager) CreateQuestion(questionData interfaces.QuestionConverter) (application.Question, error) {
	if m.ShouldError {
		return application.Question{}, errors.New("mock error")
	}

	m.CreateQuestionCalls = append(m.CreateQuestionCalls, questionData)
	question := questionData.Question("")
	m.questions[question.ID] = question
	return question, nil
}

// UpdateQuestion implements QuestionWriter interface
func (m *MockQuestionManager) UpdateQuestion(id string, questionData interfaces.QuestionConverter) (application.Question, error) {
	if m.ShouldError {
		return application.Question{}, errors.New("mock error")
	}

	m.UpdateQuestionCalls[id] = questionData
	question := questionData.Question(id)
	m.questions[id] = question
	return question, nil
}

// ReadQuestion implements QuestionReader interface
func (m *MockQuestionManager) ReadQuestion(id string) (application.Question, error) {
	if m.ShouldError {
		return application.Question{}, errors.New("mock error")
	}

	m.ReadQuestionCalls = append(m.ReadQuestionCalls, id)
	question, exists := m.questions[id]
	if !exists {
		return application.Question{}, nil
	}
	return question, nil
}

// ReadQuestions implements QuestionReader interface
func (m *MockQuestionManager) ReadQuestions(limit int, finalId string, tags []string) ([]application.Question, error) {
	if m.ShouldError {
		return nil, errors.New("mock error")
	}

	// Convert map to slice and apply filters
	var result []application.Question
	for _, q := range m.questions {
		// Apply tags filter if specified
		if len(tags) > 0 {
			matched := false
			for _, tag := range tags {
				for _, qTag := range q.Tags {
					if tag == qTag {
						matched = true
						break
					}
				}
				if matched {
					break
				}
			}
			if !matched {
				continue
			}
		}
		result = append(result, q)
	}

	// Apply finalId filter
	if finalId != "" {
		filtered := []application.Question{}
		for _, q := range result {
			if q.ID > finalId {
				filtered = append(filtered, q)
			}
		}
		result = filtered
	}

	// Apply limit
	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// DeleteQuestion implements QuestionDeleter interface
func (m *MockQuestionManager) DeleteQuestion(id string) (time.Time, error) {
	if m.ShouldError {
		return time.Time{}, errors.New("mock error")
	}

	m.DeleteQuestionCalls = append(m.DeleteQuestionCalls, id)
	delete(m.questions, id)
	return time.Now(), nil
}

// Helper methods for test setup and verification

// SetQuestion adds a question to the mock database
func (m *MockQuestionManager) SetQuestion(question application.Question) {
	m.questions[question.ID] = question
}

// GetCreateQuestionCalls returns the number of times CreateQuestion was called
func (m *MockQuestionManager) GetCreateQuestionCalls() int {
	return len(m.CreateQuestionCalls)
}

// GetUpdateQuestionCalls returns the number of times UpdateQuestion was called
func (m *MockQuestionManager) GetUpdateQuestionCalls() int {
	return len(m.UpdateQuestionCalls)
}

// GetReadQuestionCalls returns the number of times ReadQuestion was called
func (m *MockQuestionManager) GetReadQuestionCalls() int {
	return len(m.ReadQuestionCalls)
}

// GetDeleteQuestionCalls returns the number of times DeleteQuestion was called
func (m *MockQuestionManager) GetDeleteQuestionCalls() int {
	return len(m.DeleteQuestionCalls)
}
