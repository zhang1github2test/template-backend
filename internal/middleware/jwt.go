package middleware

import (
	"errors"
	"net/http"
	"strings"
	"template-backend/config"
	"template-backend/pkg/utils"

	"go.uber.org/zap"

	"template-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret123")

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取不需要进行验证的 url
		skipAuthUrls := config.GetConfig().JWT.SkipAuthUrls
		for _, url := range skipAuthUrls {
			if strings.Contains(c.Request.URL.Path, url) {
				logger.Logger().Info("url should skip auth ", zap.String("url", url))
				c.Next()
				return
			}
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.JSON(c, utils.Error("missing Authorization header", http.StatusUnauthorized))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.JSON(c, utils.Error("invalid Authorization header", http.StatusUnauthorized))
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := parseToken(tokenStr)
		if err != nil {
			utils.JSON(c, utils.Error(err.Error(), http.StatusUnauthorized))
			c.Abort()
			return
		}

		// 将 username 放到 context，供 handler 使用
		if username, ok := claims["username"].(string); ok {
			c.Set("username", username)
		}
		// 将 username 放到 context，供 handler 使用
		if userId, ok := claims["userId"].(string); ok {
			c.Set("userId", userId)
		}

		c.Next()
	}
}

func parseToken(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		m := map[string]interface{}{}
		for k, v := range claims {
			m[k] = v
		}
		return m, nil
	}
	return nil, errors.New("cannot parse claims")
}
