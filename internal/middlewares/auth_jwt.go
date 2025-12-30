package middlewares

import (
	"strings"
	"time"

	"komando/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID     int64  `json:"user_id"`
	RoleCode   string `json:"role_code"`
	RegionalID int64  `json:"regional_id"`
	jwt.RegisteredClaims
}

func AuthJWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			response.Fail(c, 401, "UNAUTHORIZED", "missing bearer token")
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(h, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.Fail(c, 401, "UNAUTHORIZED", "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*AuthClaims)
		if !ok {
			response.Fail(c, 401, "UNAUTHORIZED", "invalid claims")
			c.Abort()
			return
		}

		// optional: expired check (RegisteredClaims usually handles it if set)
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			response.Fail(c, 401, "UNAUTHORIZED", "token expired")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role_code", claims.RoleCode)
		c.Set("regional_id", claims.RegionalID)
		c.Next()
	}
}
