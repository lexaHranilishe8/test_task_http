package repository

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
	"http_service/models"
	"net/http"
)

// Интерфейс репозитория пользователей, определяющий методы работы с пользователями.
type UserRepositoryInterface interface {
	// Получение всех пользователей.
	GetAll(ctx context.Context) ([]models.User, error)
	// Создание нового пользователя.
	Create(ctx context.Context, user *models.User) error
	// Обновление данных пользователя.
	Update(ctx context.Context, user *models.User) error
	// Удаление пользователя по ID.
	Delete(ctx context.Context, id int) error
	// Вход в систему (проверка учетных данных).
	SignIn(ctx context.Context, username, password string) (*models.User, error)
	// Регистрация нового пользователя.
	SignUp(ctx context.Context, user *models.User) error
}

// Структура UserRepository, реализующая интерфейс UserRepositoryInterface.
type UserRepository struct {
	db *bun.DB // Подключение к базе данных.
}

// Конструктор для создания нового экземпляра UserRepository.
func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Метод для получения всех пользователей из базы данных.
func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.NewSelect().Model(&users).Scan(ctx) // Выполнение запроса на выборку всех пользователей.
	return users, err
}

// Метод для создания нового пользователя в базе данных.
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx) // Выполнение запроса на вставку нового пользователя.
	return err
}

// Метод для обновления данных существующего пользователя.
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.NewUpdate().Model(user).WherePK().Exec(ctx) // Выполнение запроса на обновление пользователя по первичному ключу.
	return err
}

// Метод для удаления пользователя из базы данных по его ID.
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.NewDelete().Model((*models.User)(nil)).Where("id = ?", id).Exec(ctx) // Выполнение запроса на удаление пользователя по ID.
	return err
}

// Метод для регистрации нового пользователя.
func (r *UserRepository) SignUp(ctx context.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) // Хеширование пароля.
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)          // Обновление пароля в модели пользователя хешированным значением.
	_, err = r.db.NewInsert().Model(user).Exec(ctx) // Выполнение запроса на вставку нового пользователя с хешированным паролем.
	return err
}

// Метод для входа в систему. Проверяет учетные данные пользователя.
func (r *UserRepository) SignIn(ctx context.Context, username, password string) (*models.User, error) {
	var user models.User
	err := r.db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx) // Поиск пользователя по имени пользователя.
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil { // Сравнение хешированного пароля из базы данных с введенным паролем.
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials") // Возвращение ошибки, если учетные данные недействительны.
	}
	return &user, nil // Возвращение пользователя, если учетные данные действительны.
}
