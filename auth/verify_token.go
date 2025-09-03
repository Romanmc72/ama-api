package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ama/api/constants"
	"ama/api/interfaces"
)

func VerifyToken(c interfaces.APIContext, jwtVerifier interfaces.JWTVerifier, logger *slog.Logger) {
	tokenString := c.GetHeader(constants.AuthorizationHeader)
	if tokenString == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": constants.AuthorizationHeader + " header is required"})
		return
	}
	bearerPrefix := "Bearer "
	tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
	token, err := jwtVerifier.VerifyIDTokenAndCheckRevoked(c, tokenString)
	if err != nil || token == nil || token.UID == "" {
		logger.Error("Could not verify the token", "error", err, "token", tokenString)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token error"})
		return
	}
	c.Set(constants.AuthTokenUserIdContextKey, token.UID)
	c.Set(constants.AuthTokenClaimsContextKey, token.Claims)
	c.Next()
}
