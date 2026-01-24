package routes

import (
	"jwt-session/internal/domains/auth"
	"jwt-session/internal/repositories/postgres"
	"jwt-session/internal/shared/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type DomainControllers struct {
	Auth *auth.Controller
}

func SetupRoutes(app *fiber.App, db *sqlx.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Inicializar repositórios
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	recoveryRepo := postgres.NewRecoveryRepository(db)

	// Inicializar adapters
	jwtAdapter := auth.NewJWTAdapter()
	mailAdapter := auth.NewMailAdapter()
	codeAdapter := auth.NewCodeAdapter()
	dateAdapter := auth.NewDateAdapter()

	// Inicializar service
	authService := auth.NewService(
		auth.Repositories{
			User:     userRepo,
			Session:  sessionRepo,
			Recovery: recoveryRepo,
		},
		jwtAdapter,
		mailAdapter,
		codeAdapter,
		dateAdapter,
	)

	// Inicializar controllers
	authController := auth.NewController(authService)

	controllers := &DomainControllers{
		Auth: authController,
	}

	v1 := app.Group("v1")

	//====================================================================================
	// (v1) Auth
	//====================================================================================
	{
		v1Auth := v1.Group("auth")
		v1Auth.Post("/sign-in", controllers.Auth.SignIn)
		v1Auth.Post("/sign-up", controllers.Auth.SignUp)

		//====================================================================================
		// (v1) Session Auth
		//====================================================================================
		v1SessionAuth := v1Auth.Group("/session")
		v1SessionAuth.Post("/sign-in", controllers.Auth.SignInWithSession)
		v1SessionAuth.Post("/sign-up", controllers.Auth.SignUpWithSession)
		v1SessionAuth.Get("/", middlewares.JwtSessionMiddleware(sessionRepo), func(c *fiber.Ctx) error {
			userId := c.Locals("userId").(string)

			return c.JSON(fiber.Map{
				"message": "sessão válida",
				"user_id": userId,
			})
		})
		v1SessionAuth.Get("/all", middlewares.JwtSessionMiddleware(sessionRepo), controllers.Auth.GetAllSessionsFromUser)
		v1SessionAuth.Delete("/:sessionId", middlewares.JwtSessionMiddleware(sessionRepo), controllers.Auth.DeleteSessionFromUser)

		//====================================================================================
		// (v1) Recovery Auth
		//====================================================================================
		v1AuthRecovery := v1Auth.Group("/recovery")
		v1AuthRecovery.Post("/request/:email", controllers.Auth.RequestRecovery)
		v1AuthRecovery.Post("/change-password", controllers.Auth.ChangePasswordRequestRecovery)
	}
}
