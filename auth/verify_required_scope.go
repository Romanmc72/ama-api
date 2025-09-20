package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/constants"
	"ama/api/interfaces"
)

func VerifyRequiredScope(c interfaces.APIContext, logger *slog.Logger, requiredScopes map[string]any) {
	if len(requiredScopes) == 0 {
		// If there are no required scopes, allow access
		c.Next()
		return
	}
	claimsValue, exists := c.Get(constants.AuthTokenClaimsContextKey)
	if !exists && len(requiredScopes) > 0 {
		logger.Error("No claims found in context")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Missing required claims"})
		return
	}
	claims, ok := claimsValue.(map[string]any)
	if !ok {
		logger.Error("Claims are not in the expected format", "claims", claimsValue)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid claims format"})
		return
	}
	missingScopes := map[string]any{}
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
