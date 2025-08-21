package interfaces

// Defines the methods used by the various API Endpoints and such that can
// be implemented by test suites in order to test various pieces of API endpoint code.
type APIContext interface {
	BindJSON(obj any) error
	DefaultQuery(key string, defaultValue string) string
	GetString(key string) string
	IndentedJSON(code int, obj any)
	Param(key string) string
	GetQueryArray(key string) ([]string, bool)
	GetHeader(key string) (string)
}
