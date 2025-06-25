package middleware

import "github.com/gin-gonic/gin"

func JWTAuthMiddleware(handler gin.HandlerFunc, requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement JWT auth
		handler(c)
	}
}

// Middleware cho group: chỉ kiểm tra JWT, không cần handler
func JWTAuthMiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement JWT auth cho group
		// Nếu không hợp lệ: c.AbortWithStatus(401); return
		c.Next()
	}
}
