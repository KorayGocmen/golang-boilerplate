package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/transport/handler"
	v1User "github.com/koraygocmen/golang-boilerplate/internal/transport/handler/user/v1"
	v1UserSession "github.com/koraygocmen/golang-boilerplate/internal/transport/handler/user_session/v1"
	v1Middleware "github.com/koraygocmen/golang-boilerplate/internal/transport/middleware/v1"
	v1 "github.com/koraygocmen/golang-boilerplate/internal/transport/response/v1"
)

// Setup setups the routes.
func Setup(app *fiber.App, handler *handler.Handler) {
	// V1 Requests.
	v1Response := v1.New(handler)

	// V1 Middlewares.
	v1Middleware := v1Middleware.New(v1Response)

	// V1 Handlers.
	v1UserHandler := v1User.New(v1Response)
	v1UserSessionHandler := v1UserSession.New(v1Response)

	// Static files.
	app.Static("/", "./public")

	// Shared endpoints.
	app.Get("/health", handler.Health)

	// Users.
	app.Post("/v1/users", v1UserHandler.Create) // Create a user.

	// User Sessions.
	app.Post("/v1/users/sessions", v1UserSessionHandler.Create) // Create a user session.

	// Requests that require the user to be authenticated.
	v1AuthApp := app.Use(v1Middleware.UserAuth)
	{
		// Users.
		v1AuthApp.Get("/v1/users", v1UserHandler.Get)       // Get authenticated user.
		v1AuthApp.Put("/v1/users", v1UserHandler.Update)    // Update authenticated user.
		v1AuthApp.Delete("/v1/users", v1UserHandler.Delete) // Delete authenticated user.

		// User Sessions.
		v1AuthApp.Delete("/v1/users/sessions", v1UserSessionHandler.Delete) // Delete authenticated user session.
	}

	app.Use(func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})
}
