package endpoints

import (
	"ama/api/interfaces"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIContext struct {
	gc *gin.Context
}

func NewAPIContext(gc *gin.Context) interfaces.APIContext {
	return &APIContext{gc: gc}
}

func (a *APIContext) Get(key string) (value any, exists bool) {
	return a.gc.Get(key)
}

func (a *APIContext) Set(key string, value any) {
	a.gc.Set(key, value)
}

func (a *APIContext) AbortWithStatusJSON(code int, jsonObj any) {
	a.gc.AbortWithStatusJSON(code, jsonObj)
}

func (a *APIContext) Next() {
	a.gc.Next()
}

func (a *APIContext) BindJSON(obj any) error {
	return a.gc.BindJSON(obj)
}

func (a *APIContext) DefaultQuery(key string, defaultValue string) string {
	return a.gc.DefaultQuery(key, defaultValue)
}

func (a *APIContext) GetString(key string) string {
	return a.gc.GetString(key)
}

func (a *APIContext) IndentedJSON(code int, obj any) {
	a.gc.IndentedJSON(code, obj)
}

func (a *APIContext) Param(key string) string {
	return a.gc.Param(key)
}

func (a *APIContext) GetQueryArray(key string) ([]string, bool) {
	return a.gc.GetQueryArray(key)
}

func (a *APIContext) GetHeader(key string) string {
	return a.gc.GetHeader(key)
}

func (a *APIContext) Header(key string, value string) {
	a.gc.Header(key, value)
}

func (a *APIContext) Deadline() (deadline time.Time, ok bool) {
	return a.gc.Deadline()
}

func (a *APIContext) Done() <-chan struct{} {
	return a.gc.Done()
}

func (a *APIContext) Err() error {
	return a.gc.Err()
}

func (a *APIContext) Value(key any) any {
	return a.gc.Value(key)
}

func (a *APIContext) Request() http.Request {
	return *a.gc.Request
}
