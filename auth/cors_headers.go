package auth

import (
	"ama/api/interfaces"
	"maps"
	"net/http"
)

func GetCorsHeaders() map[string]string {
	copy := map[string]string{}
	maps.Copy(copy, corsHeaders)
	return copy
}

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":      "*",
	"Access-Control-Allow-Credentials": "true",
	"Access-Control-Allow-Headers":     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
	"Access-Control-Allow-Methods":     "POST, OPTIONS, GET, PUT, DELETE",
}

func CORSHeaders(c interfaces.APIContext) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

	if c.Request().Method == "OPTIONS" {
		c.AbortWithStatusJSON(http.StatusNoContent, nil)
		return
	}

	c.Next()
}
