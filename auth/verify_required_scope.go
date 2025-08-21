package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/constants"
)

func VerifyRequiredScope(c *gin.Context, logger *slog.Logger, requiredScopes map[string]string) {
	claimsValue, exists := c.Get(constants.AuthTokenClaimsContextKey)
	if !exists {
		logger.Error("No claims found in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing required claims"})
		return
	}
	claims, ok := claimsValue.(map[string]interface{})
	if !ok {
		logger.Error("Claims are not in the expected format", "claims", claimsValue)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid claims format"})
		return
	}
	missingScopes := map[string]string{}
	for scope, requiredValue := range requiredScopes {
		value, exists := claims[scope]
		if !exists || value != requiredValue {
			missingScopes[scope] = requiredValue
		}
	}
	if len(missingScopes) > 0 {
		logger.Error("Required scope not found or does not match", "scopes", missingScopes)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing or invalid scope"})
		return
	}
	c.Next()
}
