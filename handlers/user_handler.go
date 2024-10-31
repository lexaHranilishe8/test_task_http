package handlers

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"http_service/models"
	"http_service/repository"
	"net/http"
	"strconv"
	"time"
)

// Ключ для подписи JWT.
var JwtKey = []byte("your_secret_key")

// Структура для хранения пользовательских требований JWT.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Структура обработчика пользователей.
type UserHandler struct {
	Repo repository.UserRepositoryInterface // Репозиторий для работы с пользователями.
}

// Обработчик для получения списка пользователей.
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.Repo.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось получить пользователей"})
	}
	return c.JSON(http.StatusOK, users)
}

// Обработчик для создания нового пользователя.
func (h *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверные входные данные"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Repo.Create(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось создать пользователя"})
	}

	return c.JSON(http.StatusCreated, user)
}

// Обработчик для обновления данных пользователя.
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверные входные данные"})
	}
	user.ID = id

	if err := h.Repo.Update(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось обновить пользователя"})
	}

	return c.JSON(http.StatusOK, user)
}

// Обработчик для удаления пользователя.
func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.Repo.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось удалить пользователя"})
	}
	// Возвращаем сообщение об успешном удалении
	return c.JSON(http.StatusOK, map[string]string{"message": "Пользователь успешно удален"})
}

// Обработчик для регистрации нового пользователя.
func (h *UserHandler) SignUp(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверные входные данные"})
	}
	if err := h.Repo.SignUp(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось создать пользователя"})
	}
	return c.JSON(http.StatusCreated, map[string]string{"message": "Пользователь создан успешно"})
}

// Обработчик для входа пользователя.
func (h *UserHandler) SignIn(c echo.Context) error {
	credentials := new(models.User)
	if err := c.Bind(credentials); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Неверные входные данные"})
	}

	user, err := h.Repo.SignIn(c.Request().Context(), credentials.Username, credentials.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Неверные учетные данные"})
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Время истечения действия токена.
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Не удалось создать токен"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}
