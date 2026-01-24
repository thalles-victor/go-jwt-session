package middlewares

import (
	"database/sql"
	"jwt-session/internal/repositories/postgres"
	"jwt-session/internal/shared/date"
	"jwt-session/internal/shared/jwt"
	"jwt-session/internal/shared/logger"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JwtMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	token := parts[1]

	sub, err := jwt.ParseJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	c.Locals("userId", sub)

	return c.Next()
}

func JwtSessionMiddleware(sessionRepo postgres.SessionRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		token := parts[1]

		sub, err := jwt.ParseJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		logger.Info.Println("check if session exist")

		session, err := sessionRepo.GetByID(c.Context(), sub)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Warn.Printf("session with id %s no found", sub)
				return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"message": "sessão não encontrada",
				})
			}
			logger.Error.Printf("internal server error when get session by id %s. error: %s", sub, err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "erro interno no servidor ao tentar buscar a sessão",
				"error":   err.Error(),
			})
		}

		if !date.IsNotExpired(session.ExpiresAt) {
			logger.Warn.Printf("session already expired at %s", session.ExpiresAt)
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": "sessão inválida ou expirada",
			})
		}

		c.Locals("userId", session.UserID)
		c.Locals("sessionId", session.ID)

		return c.Next()
	}
}
