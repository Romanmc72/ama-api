package test

import (
	"ama/api/application"
	"ama/api/application/user"
	"ama/api/interfaces"
	"time"
)

// MockUserManager implements UserManager interface for testing
type MockUserManager struct {
	Users          map[string]MockUserConfig
	MockCreateUser func(userData user.BaseUser) (application.User, error)
	MockReadUser   func(id string) (application.User, error)
	MockUpdateUser func(userData interfaces.UserConverter) error
	MockDeleteUser func(id string) (time.Time, error)
}

type MockUserConfig struct {
	Data application.User
}

type MockUserManagerConfig struct {
	Users map[string]MockUserConfig
}

// NewMockUserManager creates a new instance of MockUserManager
func NewMockUserManager(cfg MockUserManagerConfig) *MockUserManager {
	return &MockUserManager{
		Users: cfg.Users,
	}
}

// CreateUser implements UserCreator interface
func (m *MockUserManager) CreateUser(userData user.BaseUser) (application.User, error) {
	if m.MockCreateUser != nil {
		return m.MockCreateUser(userData)
	}
	return application.User{}, nil
}

// ReadUser implements UserReader interface
func (m *MockUserManager) ReadUser(id string) (application.User, error) {
	if m.MockReadUser != nil {
		return m.MockReadUser(id)
	}
	return application.User{}, nil
}

// UpdateUser implements UserWriter interface
func (m *MockUserManager) UpdateUser(userData interfaces.UserConverter) error {
	if m.MockUpdateUser != nil {
		return m.MockUpdateUser(userData)
	}
	return nil
}

// DeleteUser implements UserDeleter interface
func (m *MockUserManager) DeleteUser(id string) (time.Time, error) {
	if m.MockDeleteUser != nil {
		return m.MockDeleteUser(id)
	}
	return time.Now(), nil
}
