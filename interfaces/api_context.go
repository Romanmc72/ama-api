package interfaces

import (
	"net/http"
	"time"
)

// Defines the methods used by the various API Endpoints and such that can
// be implemented by test suites in order to test various pieces of API endpoint code.
type APIContext interface {
	Get(key string) (value any, exists bool)
	Set(key string, value any)
	AbortWithStatusJSON(code int, jsonObj any)
	Next()
	BindJSON(obj any) error
	DefaultQuery(key string, defaultValue string) string
	GetString(key string) string
	IndentedJSON(code int, obj any)
	Param(key string) string
	GetQueryArray(key string) ([]string, bool)
	GetHeader(key string) string
	Header(key string, value string)

	/**
	 * These are required to implement the context.Context interface.
	 */
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any

	/**
	 * Returns the request method, annoying that this is not a part of the gin interface itself so we must implement a wrapper to be able to test this thing.
	 */
	Request() http.Request
}
