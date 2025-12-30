package middlewares

import (
	"komando/internal/shared/response"

	"github.com/gin-gonic/gin"
)

func RBAC(allowed ...string) gin.HandlerFunc {
	allowedSet := map[string]bool{}
	for _, r := range allowed {
		allowedSet[r] = true
	}

	return func(c *gin.Context) {
		roleAny, ok := c.Get("role_code")
		if !ok {
			response.Fail(c, 401, "UNAUTHORIZED", "missing role")
			c.Abort()
			return
		}
		role := roleAny.(string)

		if !allowedSet[role] {
			response.Fail(c, 403, "FORBIDDEN", "insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}
