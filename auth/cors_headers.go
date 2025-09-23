package auth

import (
	"ama/api/interfaces"
	"maps"
	"net/http"
)

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
	"Access-Control-Allow-Methods":     "POST, OPTIONS, GET, PUT, DELETE",
}

func GetCorsHeaders() map[string]string {
	copy := map[string]string{}
	maps.Copy(copy, corsHeaders)
	return copy
}

func CORSHeaders(c interfaces.APIContext) {
	for k, v := range GetCorsHeaders() {
		c.Header(k, v)
	}

	if c.Request().Method == "OPTIONS" {
		c.AbortWithStatusJSON(http.StatusNoContent, nil)
		return
	}

	c.Next()
}
