package auth

import (
	"net/http"
	"strings"

	"quiz-builder/constants"
	"quiz-builder/logger"
	"quiz-builder/service/auth"

	"github.com/gin-gonic/gin"
)

// AuthFilter filter for authentication
type AuthFilter struct {
	service auth.AuthService
	logger  logger.Logger
}

// NewAuthFilter creates new auth filter
func NewAuthFilter(
	logger logger.Logger,
	service auth.AuthService,
) AuthFilter {
	return AuthFilter{
		service: service,
		logger:  logger,
	}
}

// Setup sets up auth filter
func (m AuthFilter) Setup() {}

// Handler handles filter functionality
func (m AuthFilter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, uuidClaim, err := m.service.Authorize(authToken)
			if authorized {
				c.Request.Header.Set(constants.UUIDHeader, uuidClaim)
				c.Next()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			m.logger.Error(err)
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorized",
		})
		c.Abort()
	}
}
