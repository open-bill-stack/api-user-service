package user

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type User struct {
	ID              pgtype.UUID `db:"id"`                // UUID користувача
	Email           string      `db:"email"`             // Email користувача
	Phone           *string     `db:"phone"`             // Телефон (може бути NULL)
	IsEmailVerified bool        `db:"is_email_verified"` // Підтвердження пошти
	IsPhoneVerified bool        `db:"is_phone_verified"` // Підтвердження телефона
	CreatedAt       time.Time   `db:"created_at"`        // Коли створено
	UpdatedAt       time.Time   `db:"updated_at"`        // Коли оновлено
}

type Credential struct {
	UserID       pgtype.UUID `db:"user_id"`       // UUID користувача
	PasswordHash string      `db:"password_hash"` // Хеш пароля
	CreatedAt    time.Time   `db:"created_at"`    //
}

type CreateUser struct {
	Email        string
	PasswordHash string
}
