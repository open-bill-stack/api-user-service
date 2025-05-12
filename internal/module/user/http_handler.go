package user

import (
	"api-user-service/internal/module/user/structure"
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
	app.Post("/user", h.CreateUser)

	group := app.Group("/user")

	group.Get("", h.GetUser)
	group.Patch("", h.UpdateUser)
	group.Delete("", h.DeleteUser)
}

func (h *HttpHandle) GetUser(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "OK"})
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
	var req structure.UserRequest
	// Парсимо JSON з тіла запиту
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	status, err := h.service.ExistsByUUID(c.Context(), req.UUID)
	if err != nil {
		return err
	}
	if !status {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User with this UUID does not exist",
		})
	}

	status, err = h.service.DeleteByUUID(c.Context(), req.UUID)
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
