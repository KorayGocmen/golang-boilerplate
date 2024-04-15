package middleware

import "github.com/gofiber/fiber/v2"

var (
	allowedMethods = []string{
		fiber.MethodGet,
		fiber.MethodPost,
		fiber.MethodPut,
		fiber.MethodDelete,
		fiber.MethodOptions,
	}

	exposedHeaders = []string{
		fiber.HeaderContentType,
		fiber.HeaderAccept,
		fiber.HeaderAuthorization,
	}
)
