package auth

import (
	"api-user-service/internal/module/auth/structure"
	"api-user-service/internal/module/user"
	jwtMiddleware "api-user-service/internal/service/fiber/middleware/jwt"
	jwtService "api-user-service/internal/service/jwt"
	"context"
	"github.com/alexedwards/argon2id"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

type Handle struct {
	log                *zap.Logger
	jwt                *jwtService.JWT
	checkJWTMiddleware *jwtMiddleware.CheckJWTMiddleware
	userService        *user.Service
}

func NewHandler(p Params) (Result, error) {
	return Result{
		Router: &Handle{
			log:                p.Log,
			jwt:                p.JWT,
			checkJWTMiddleware: p.CheckJWTMiddleware,
			userService:        p.UserService,
		},
	}, nil
}

func (h *Handle) Register(app *fiber.App) {
	app.Post("/auth/jwt/login", h.LoginJWT)
	app.Post("/auth/jwt/verify", h.VerifyJWT)
	app.Post("/auth/jwt/refresh", h.checkJWTMiddleware.GetHandler(), h.RefreshJWT)
}

func (h *Handle) LoginJWT(c *fiber.Ctx) error {
	var req structure.LoginUserRequest
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

	status, err := h.userService.ExistsByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking detailUser existence",
		})
	}

	if !status {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User with this email does not exist",
		})
	}

	credential, err := h.userService.GetCredentialsByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting detailUser credentials",
		})
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, credential.PasswordHash)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error comparing password and hash",
		})
	}

	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is incorrect",
		})
	}

	detailUser, err := h.userService.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating detailUser",
		})
	}
	extendClaims := jwtService.ExtendClaims{
		UserID: detailUser.ID.String(),
	}
	timeNow := time.Now()

	accessClams := jwtService.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(15 * time.Minute)),
		},
		ExtendClaims: extendClaims,
	}
	refreshClams := jwtService.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(7 * 24 * time.Hour)),
		},
		ExtendClaims: extendClaims,
	}
	accessToken, err := h.jwt.Create(accessClams)
	if err != nil {
		return err
	}
	refreshToken, err := h.jwt.Create(refreshClams)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

type Body struct {
	Token string `json:"token"`
}

func (h *Handle) VerifyJWT(c *fiber.Ctx) error {
	p := new(Body)

	if err := c.BodyParser(p); err != nil {
		return err
	}
	token, err := h.jwt.Verify(p.Token)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"result": token})
}

func (h *Handle) RefreshJWT(c *fiber.Ctx) error {
	return nil
}
