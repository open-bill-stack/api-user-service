package structure

type ResponseUser struct {
	ID              string  `json:"id"`
	Email           string  `json:"email"`             // Email користувача
	Phone           *string `json:"phone"`             // Телефон (може бути NULL)
	IsEmailVerified bool    `json:"is_email_verified"` // Підтвердження пошти
	IsPhoneVerified bool    `json:"is_phone_verified"` // Підтвердження телефона
	CreatedAt       string  `json:"created_at"`        // Коли створено
	UpdatedAt       string  `json:"updated_at"`
}
