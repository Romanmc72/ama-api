package test

import (
	"ama/api/application"
	"ama/api/interfaces"
	"errors"
	"time"
)

// MockUserManager implements UserManager interface for testing
type MockUserManager struct {
	// Store users for simulating database
	users map[string]application.User

	// Track method calls for verification in tests
	CreateUserCalls []interfaces.UserConverter
	UpdateUserCalls []interfaces.UserConverter
	ReadUserCalls   []string
	DeleteUserCalls []string

	// Control mock behavior
	ShouldError bool
	Error       error
}

// NewMockUserManager creates a new instance of MockUserManager
func NewMockUserManager() *MockUserManager {
	return &MockUserManager{
		users: make(map[string]application.User),
	}
}

// CreateUser implements UserCreator interface
func (m *MockUserManager) CreateUser(userData interfaces.UserConverter) application.User {
	m.CreateUserCalls = append(m.CreateUserCalls, userData)
	user := userData.User()
	m.users[user.ID] = user
	return user
}

// ReadUser implements UserReader interface
func (m *MockUserManager) ReadUser(id string) (application.User, error) {
	m.ReadUserCalls = append(m.ReadUserCalls, id)
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return application.User{}, errors.New("mock error")
}

// UpdateUser implements UserWriter interface
func (m *MockUserManager) UpdateUser(userData interfaces.UserConverter) application.User {
	m.UpdateUserCalls = append(m.UpdateUserCalls, userData)
	user := userData.User()
	m.users[user.ID] = user
	return user
}

// DeleteUser implements UserDeleter interface
func (m *MockUserManager) DeleteUser(id string) (time.Time, error) {
	if m.ShouldError {
		return time.Time{}, m.Error
	}

	m.DeleteUserCalls = append(m.DeleteUserCalls, id)
	delete(m.users, id)
	return time.Now(), nil
}

// Helper methods for test setup and verification

// SetUser adds a user to the mock database
func (m *MockUserManager) SetUser(user application.User) {
	m.users[user.ID] = user
}

// GetCreateUserCalls returns the number of times CreateUser was called
func (m *MockUserManager) GetCreateUserCalls() int {
	return len(m.CreateUserCalls)
}

// GetUpdateUserCalls returns the number of times UpdateUser was called
func (m *MockUserManager) GetUpdateUserCalls() int {
	return len(m.UpdateUserCalls)
}

// GetReadUserCalls returns the number of times ReadUser was called
func (m *MockUserManager) GetReadUserCalls() int {
	return len(m.ReadUserCalls)
}

// GetDeleteUserCalls returns the number of times DeleteUser was called
func (m *MockUserManager) GetDeleteUserCalls() int {
	return len(m.DeleteUserCalls)
}

// GetUserByID returns a user from the mock database
func (m *MockUserManager) GetUserByID(id string) (application.User, bool) {
	user, exists := m.users[id]
	return user, exists
}
