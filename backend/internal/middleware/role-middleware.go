package middlewares

import (
	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole enums.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			core.HandleApiError(c, core.ForbiddenError())
			c.Abort()
			return
		}

		if role != string(requiredRole) {
			core.HandleApiError(c, core.ForbiddenError())
			c.Abort()
			return
		}

		c.Next()
	}
}
