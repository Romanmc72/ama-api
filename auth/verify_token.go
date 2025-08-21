package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"

	"ama/api/constants"
)

func VerifyToken(c *gin.Context, jwtVerifier *auth.Client, logger *slog.Logger) {
	tokenString := c.GetHeader(constants.AuthorizationHeader)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.AuthorizationHeader + " header is required"})
		return
	}
	bearerPrefix := "Bearer "
	tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
	token, err := jwtVerifier.VerifyIDTokenAndCheckRevoked(c.Request.Context(), tokenString)
	if err != nil || token == nil || token.UID == "" {
		logger.Error("Could not verify the token", "error", err, "token", tokenString)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token error"})
		return
	}
	c.Set(constants.AuthTokenUserIdContextKey, token.UID)
	c.Set(constants.AuthTokenClaimsContextKey, token.Claims)
	c.Next()
}
