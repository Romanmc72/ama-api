package test

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

// MockAPIContext implements APIContext for testing purposes
type MockAPIContext struct {
	// Store bound JSON data
	InputJSON []byte

	// Store params and query values
	Params      map[string]string
	QueryValues map[string][]string

	// Store the response code and data
	ResponseCode int
	ResponseData interface{}

	validate *validator.Validate
	headers  map[string]string
	values   map[string]any
	req      http.Request
}

type MockAPIContextConfig struct {
	InputJSON   []byte
	Params      map[string]string
	QueryValues map[string][]string
	Headers     map[string]string
	Values      map[string]any
	Request     http.Request
}

// NewMockAPIContext creates a new instance of MockAPIContext
func NewMockAPIContext(cfg MockAPIContextConfig) *MockAPIContext {
	var params map[string]string
	if cfg.Params != nil {
		params = cfg.Params
	} else {
		params = make(map[string]string)
	}
	var queryValues map[string][]string
	if cfg.QueryValues != nil {
		queryValues = cfg.QueryValues
	} else {
		queryValues = make(map[string][]string)
	}
	var headers map[string]string
	if cfg.Headers != nil {
		headers = cfg.Headers
	} else {
		headers = make(map[string]string)
	}
	var values map[string]any
	if cfg.Values != nil {
		values = cfg.Values
	} else {
		values = make(map[string]any)
	}
	return &MockAPIContext{
		InputJSON:   cfg.InputJSON,
		Params:      params,
		QueryValues: queryValues,
		validate:    validator.New(),
		headers:     headers,
		values:      values,
		req:         cfg.Request,
	}
}

// BindJSON implements the BindJSON method of APIContext
func (m *MockAPIContext) BindJSON(obj any) error {
	if err := json.Unmarshal(m.InputJSON, obj); err != nil {
		return err
	}
	if err := m.validate.Struct(obj); err != nil {
		return err
	}
	return nil
}

// DefaultQuery implements the DefaultQuery method of APIContext
func (m *MockAPIContext) DefaultQuery(key string, defaultValue string) string {
	if values, exists := m.QueryValues[key]; exists && len(values) > 0 {
		return values[0]
	}
	return defaultValue
}

// GetString implements the GetString method of APIContext
func (m *MockAPIContext) GetString(key string) string {
	if v, exists := m.values[key]; exists && v != "" {
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}

// IndentedJSON implements the IndentedJSON method of APIContext
func (m *MockAPIContext) IndentedJSON(code int, obj any) {
	m.ResponseCode = code
	m.ResponseData = obj
}

// Param implements the Param method of APIContext
func (m *MockAPIContext) Param(key string) string {
	return m.Params[key]
}

// GetQueryArray implements the GetQueryArray method of APIContext
func (m *MockAPIContext) GetQueryArray(key string) ([]string, bool) {
	values, exists := m.QueryValues[key]
	return values, exists
}

// Helper methods for setting up test scenarios

// SetParam sets a URL parameter value
func (m *MockAPIContext) SetParam(key, value string) {
	m.Params[key] = value
}

// SetQueryValue sets a query parameter value
func (m *MockAPIContext) SetQueryValue(key string, values []string) {
	m.QueryValues[key] = values
}

// SetInputJSON sets the JSON data to be bound
func (m *MockAPIContext) SetInputJSON(data []byte) {
	m.InputJSON = data
}

func (m *MockAPIContext) GetHeader(key string) string {
	return m.headers[key]
}

func (m *MockAPIContext) Header(key string, value string) {
	m.headers[key] = value
}

func (m *MockAPIContext) Get(key string) (value any, exists bool) {
	v, ok := m.values[key]
	return v, ok
}

func (m *MockAPIContext) Set(key string, value any) {
	m.values[key] = value
}

func (m *MockAPIContext) AbortWithStatusJSON(code int, jsonObj any) {
	m.ResponseCode = code
	m.ResponseData = jsonObj
}

func (m *MockAPIContext) Next() {}

func (m *MockAPIContext) Request() http.Request {
	return m.req
}

func (m *MockAPIContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (m *MockAPIContext) Done() <-chan struct{} {
	return nil
}

func (m *MockAPIContext) Err() error {
	return nil
}

func (m *MockAPIContext) Value(key any) any {
	return m.values[key.(string)]
}

// Just for retrieving and checking the set headers from the test context
func (m *MockAPIContext) GetTestHeaders() map[string]string {
	return m.headers
}
