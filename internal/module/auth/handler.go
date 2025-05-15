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
	app.Post("/auth/jwt/refresh", h.RefreshJWT)
}

func (h *Handle) LoginJWT(c *fiber.Ctx) error {
	var req structure.LoginUserRequest
	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Парсимо JSON з тіла запиту
	if err := c.BodyParser(&req); err != nil {
		h.log.Error("Invalid body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	h.log.Info("Login by JWT", zap.String("email", req.Email))

	// Перевіряємо валідацію структури
	if err := validate.Struct(req); err != nil {
		h.log.Error("validating struct", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	status, err := h.userService.ExistsByEmail(context.Background(), req.Email)
	if err != nil {
		h.log.Error("error checking email existence", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking detailUser existence",
		})
	}

	if !status {
		h.log.Info("User with this email does not exist", zap.String("email", req.Email))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User with this email does not exist",
		})
	}

	credential, err := h.userService.GetCredentialsByEmail(context.Background(), req.Email)
	if err != nil {
		h.log.Error("Error getting detailUser credentials", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting detailUser credentials",
		})
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, credential.PasswordHash)

	if err != nil {
		h.log.Error("Error comparing password and hash", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error comparing password and hash",
		})
	}

	if !match {
		h.log.Info("Password is incorrect")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password is incorrect",
		})
	}

	detailUser, err := h.userService.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		h.log.Error("Error getting detailUser", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating detailUser",
		})
	}

	timeNow := time.Now()
	accessClams := jwtService.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(15 * time.Minute)),
		},
		ExtendClaims: jwtService.ExtendClaims{
			UserID: detailUser.ID.String(),
		},
	}
	refreshClams := jwtService.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(7 * 24 * time.Hour)),
		},
		ExtendClaims: jwtService.ExtendClaims{
			UserID:         detailUser.ID.String(),
			IsRefreshToken: true,
		},
	}
	accessToken, err := h.jwt.Create(accessClams)
	if err != nil {
		h.log.Error("Error generating JWT access token", zap.Error(err))
		return err
	}
	refreshToken, err := h.jwt.Create(refreshClams)
	if err != nil {
		h.log.Error("Error generating JWT refresh token", zap.Error(err))
		return err
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

func (h *Handle) RefreshJWT(c *fiber.Ctx) error {
	p := new(structure.RefreshTokenRequest)
	if err := c.BodyParser(p); err != nil {
		return err
	}
	token, err := h.jwt.Verify(p.Token)
	if err != nil {
		return err
	}

	if !token.IsRefreshToken {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token is not refresh token",
		})
	}

	timeNow := time.Now()
	accessClams := jwtService.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(15 * time.Minute)),
		},
		ExtendClaims: jwtService.ExtendClaims{
			UserID: token.UserID,
		},
	}
	refreshClams := jwtService.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(7 * 24 * time.Hour)),
		},

		ExtendClaims: jwtService.ExtendClaims{
			UserID:         token.UserID,
			IsRefreshToken: true,
		},
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
