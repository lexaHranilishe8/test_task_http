package models

import (
	"context"
	"github.com/uptrace/bun"
)

// Функция для создания схемы таблицы пользователей в базе данных.
func CreateSchema(db *bun.DB) error {
	// Создание таблицы пользователей, если она еще не существует.
	_, err := db.NewCreateTable().
		Model((*User)(nil)).
		IfNotExists().
		Exec(context.Background())
	return err
}

// Структура, представляющая модель пользователя.
type User struct {
	// Уникальный идентификатор пользователя, автоматически увеличивается.
	ID int `bun:",pk,autoincrement" json:"id"`

	Username string `bun:",unique" json:"username"`

	Password string `json:"password"`
}
