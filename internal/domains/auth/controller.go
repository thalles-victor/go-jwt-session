package auth

import (
	"errors"
	"jwt-session/internal/shared/logger"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) SignUp(ctx *fiber.Ctx) error {
	var dto SignUpDto

	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  "INVALID_PAYLOAD",
			"error": "invalid request payload",
		})
	}

	result, err := c.service.SignUp(ctx.Context(), dto)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "usuário cadastrado com sucesso",
		"user":        result.User,
		"access_token": result.AccessToken,
	})
}

func (c *Controller) SignIn(ctx *fiber.Ctx) error {
	var dto SignInDto

	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  "INVALID_PAYLOAD",
			"error": "invalid request payload",
		})
	}

	result, err := c.service.SignIn(ctx.Context(), dto)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "usuário logado com sucesso",
		"user":        result.User,
		"access_token": result.AccessToken,
	})
}

func (c *Controller) SignUpWithSession(ctx *fiber.Ctx) error {
	var dto SignUpDto

	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  "INVALID_PAYLOAD",
			"error": "invalid request payload",
		})
	}

	browser := ctx.Get("User-Agent")
	ip := ctx.IP()

	result, err := c.service.SignUpWithSession(ctx.Context(), dto, browser, ip)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "usuário cadastrado com sucesso",
		"user":        result.User,
		"access_token": result.AccessToken,
		"session":     result.Session,
	})
}

func (c *Controller) SignInWithSession(ctx *fiber.Ctx) error {
	var dto SignInDto

	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  "INVALID_PAYLOAD",
			"error": "invalid request payload",
		})
	}

	browser := ctx.Get("User-Agent")
	ip := ctx.IP()

	result, err := c.service.SignInWithSession(ctx.Context(), dto, browser, ip)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "usuário logado com sucesso",
		"user":        result.User,
		"access_token": result.AccessToken,
		"session":     result.Session,
	})
}

func (c *Controller) GetAllSessionsFromUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(string)

	sessions, err := c.service.GetAllSessionsFromUser(ctx.Context(), userID)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "sessões encontradas com sucesso",
		"sessions": sessions,
	})
}

func (c *Controller) DeleteSessionFromUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userId").(string)
	sessionID := ctx.Params("sessionId")

	err := c.service.DeleteSessionFromUser(ctx.Context(), userID, sessionID)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "sessão deletada com sucesso",
	})
}

func (c *Controller) RequestRecovery(ctx *fiber.Ctx) error {
	email := ctx.Params("email")

	err := c.service.RequestRecovery(ctx.Context(), email)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "código de recuperação enviado com sucesso",
	})
}

func (c *Controller) ChangePasswordRequestRecovery(ctx *fiber.Ctx) error {
	var dto ChangePasswordRequestRecoveryDto

	if err := ctx.BodyParser(&dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":  "INVALID_PAYLOAD",
			"error": "invalid request payload",
		})
	}

	err := c.service.ChangePasswordRequestRecovery(ctx.Context(), dto)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "senha alterada com sucesso",
	})
}

// handleError trata DomainErrors e retorna respostas apropriadas
func handleError(ctx *fiber.Ctx, err error) error {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		logger.Error.Printf("[%s] %s", domainErr.Code, domainErr.Message)
		return ctx.Status(domainErr.HTTPStatus).JSON(fiber.Map{
			"code":  domainErr.Code,
			"error": domainErr.Message,
		})
	}

	// Fallback para erros não-domínio
	logger.Error.Printf("internal error: %v", err)
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":  "INTERNAL_ERROR",
		"error": "internal server error",
	})
}
