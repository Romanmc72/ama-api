package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"ama/api/constants"
)

func VerifyUserID(c *gin.Context, logger *slog.Logger) {
	tokenId := c.GetString(constants.AuthTokenUserIdContextKey)
	userId := c.Param(constants.UserIdPathIdentifier)
	if tokenId != userId {
		logger.Error("User ID does not match token ID", "tokenId", tokenId, constants.UserIdPathIdentifier, userId)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User ID does not match token ID"})
		return
	}
	c.Next()
}
