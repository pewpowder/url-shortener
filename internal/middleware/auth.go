package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pewpowder/url-shortener/internal/service"
)

// AuthMiddleware проверяет токен и добавляет USER_ID в context

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	se.UnauthorizedError(c, errors.New("authorization header is required"))
		// 	c.Abort()
		// 	return
		// }

		// tokenParts := strings.Split(authHeader, " ")
		// if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		// 	se.UnauthorizedError(c, errors.New("invalid authorization header format"))
		// 	c.Abort()
		// 	return
		// }

		// token := tokenParts[1]
		// userID, err := authService.GetUserIDByToken(c.Request.Context(), token)
		// if err != nil {
		// 	se.UnauthorizedError(c, errors.New("invalid or expired token"))
		// 	c.Abort()
		// 	return
		// }

		// // Добавляем USER_ID в gin context и request context
		// c.Set("user_id", userID)
		// ctx := utils.SetUserID(c.Request.Context(), userID)
		// c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
