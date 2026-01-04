package middlewares

import (
	"strings"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/modules/auth/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService service.IJwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			core.HandleApiError(c, core.UnauthorizedError())
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			core.HandleApiError(c, core.UnauthorizedError())
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwtService.Verify(token)
		if err != nil {
			core.HandleApiError(c, err)
			c.Abort()
			return
		}

		c.Set("user_id", (*claims)["user_id"])
		c.Set("role", (*claims)["role"])
		c.Next()
	}
}
