package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"http_service/handlers"
)

// Функция для создания middleware для проверки JWT.
func JWTMiddleware() echo.MiddlewareFunc {
	// Возвращение middleware для проверки JWT с использованием ключа из обработчиков.
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: handlers.JwtKey,
	})
}
