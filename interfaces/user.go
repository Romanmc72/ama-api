// Describes the various interfaces for working with user objects
package interfaces

import (
	"ama/api/application"
	"time"
)

// Anything that can be converted into a user
type UserConverter interface {
	// The method that converts or returns a user object
	User() application.User
}

// Anything that can create a new user
type UserCreator interface {
	// Create a new user and return the newly created user
	CreateUser(userData UserConverter) application.User
}

// Anything that can retrieve a user
type UserReader interface {
	// Get a particular user from the database given their user id
	ReadUser(id string) application.User
}

// Anything that can update a user
type UserWriter interface {
	// Update a user using new user data
	UpdateUser(userData UserConverter) application.User
}

// Anything capable of deleting a user
type UserDeleter interface {
	// Delete a user and return when they were deleted along with any errors
	DeleteUser(id string) (time.Time, error)
}

// Anything that can manage the full user lifecycle
type UserManager interface {
	UserCreator
	UserDeleter
	UserReader
	UserWriter
}
