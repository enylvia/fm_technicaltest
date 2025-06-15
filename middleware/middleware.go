package middleware

import (
	"FM_techincaltest/helpers"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

func AuthenticateMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header is missing"})
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header format must be Bearer <token>"})
			}

			tokenString := parts[1]

			claims, err := helpers.VerifyToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token", "error": err.Error()})
			}

			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}
