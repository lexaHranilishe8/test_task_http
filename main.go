package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"http_service/handlers"
	"http_service/middlewares"
	"http_service/models"
	"http_service/repository"
	"os"
)

func main() {
	// Получаем URL базы данных из переменной окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "file:myapp.db?cache=shared&_foreign_keys=1" // Значение по умолчанию
	}

	// Подключение к SQLite
	sqldb, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		panic(err)
	}
	defer sqldb.Close()

	// Создаем экземпляр bun.DB с использованием sqlitedialect
	db := bun.NewDB(sqldb, sqlitedialect.New())
	userRepo := repository.NewUserRepository(db)

	userHandler := &handlers.UserHandler{
		Repo: userRepo,
	}

	// Создание таблицы User, если её нет
	if err := models.CreateSchema(db); err != nil { // Исправлено на CreateSchema
		fmt.Println("Ошибка создания схемы:", err)
	} else {
		fmt.Println("Схема создана успешно.")
	}

	e := echo.New()

	// Открытые маршруты
	e.POST("/signin", userHandler.SignIn) // Авторизация
	e.POST("/signup", userHandler.SignUp) // Регистрация

	// Защищённые маршруты с использованием JWTMiddleware
	authGroup := e.Group("/users")
	authGroup.Use(middlewares.JWTMiddleware())
	authGroup.GET("", userHandler.GetUsers)
	authGroup.POST("", userHandler.CreateUser)
	authGroup.PUT("/:id", userHandler.UpdateUser)
	authGroup.DELETE("/:id", userHandler.DeleteUser)

	e.Logger.Fatal(e.Start(":8080"))
}
