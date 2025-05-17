package user

import (
	"api-user-service/internal/module/user/structure"
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HttpHandle struct {
	log     *zap.Logger
	service *Service
}

func NewHttpHandler(p Params) (HttpResult, error) {
	return HttpResult{
		Router: &HttpHandle{
			log:     p.Log,
			service: p.Service,
		},
	}, nil
}

func (h *HttpHandle) Register(app *fiber.App) {

	group := app.Group("/user")

	group.Post("", h.CreateUser)
	group.Get("", h.ListUsers)
	group.Get(":id", h.GetUser)
	group.Patch(":id", h.UpdateUser)
	group.Delete(":id", h.DeleteUser)
}

func (h *HttpHandle) ListUsers(c *fiber.Ctx) error {
	users, err := h.service.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var structUsers = make([]structure.ResponseUser, 0)
	for _, user := range users {
		structUsers = append(structUsers, structure.ResponseUser{
			ID:              user.ID.String(),
			Email:           user.Email,
			Phone:           user.Phone,
			IsEmailVerified: user.IsEmailVerified,
			IsPhoneVerified: user.IsPhoneVerified,
			CreatedAt:       user.CreatedAt.String(),
			UpdatedAt:       user.UpdatedAt.String(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(structUsers)

}
func (h *HttpHandle) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var tariffID, errParse = uuid.Parse(id)
	if errParse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}
	user, err := h.service.GetUserByEmail(c.Context(), tariffID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{})
	}
	return c.Status(fiber.StatusOK).JSON(structure.ResponseUser{
		ID:              user.ID.String(),
		Email:           user.Email,
		Phone:           user.Phone,
		IsEmailVerified: user.IsEmailVerified,
		IsPhoneVerified: user.IsPhoneVerified,
		CreatedAt:       user.CreatedAt.String(),
		UpdatedAt:       user.UpdatedAt.String(),
	})
}
func (h *HttpHandle) UpdateUser(c *fiber.Ctx) error {
	return nil
}
func (h *HttpHandle) CreateUser(c *fiber.Ctx) error {
	var req structure.CreateUserRequest
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Парсимо JSON з тіла запиту
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Перевіряємо валідацію структури
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// написати хешування пароля і передау в структуру
	hash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	existsEmail, err := h.service.ExistsByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if existsEmail {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email already exists",
		})
	}

	var userInput = CreateUser{
		PasswordHash: hash,
		Email:        req.Email,
	}

	if _, err := h.service.Create(context.Background(), userInput); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

func (h *HttpHandle) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var userID, errParse = uuid.Parse(id)
	// Парсимо JSON з тіла запиту
	if errParse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	status, err := h.service.ExistsByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !status {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User with this UUID does not exist",
		})
	}

	status, err = h.service.DeleteByID(c.Context(), userID)
	if err != nil {
		return err
	}
	if !status {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting user",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
