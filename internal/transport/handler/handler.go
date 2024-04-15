package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/service"
)

type Config struct {
	SHASUM string
}

type Handler struct {
	SHASUM string
}

// Create the handler object with the logger to use the
// error handler. gofiber expects the error handler to be
// available before creating the fiber app.
func New(c Config) *Handler {
	return &Handler{
		SHASUM: c.SHASUM,
	}
}

// Shared handlers.
func (h *Handler) Health(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	if err := srv.Ping(); err != nil {
		err = fmt.Errorf("ping error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(h.Failure(ctx, err, service.ErrInternalServer))
	}

	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(h.Failure(ctx, err, service.ErrInternalServer))
	}

	var shasum string
	if h.SHASUM != "" {
		shasum = h.SHASUM[:7]
	}

	return c.Status(fiber.StatusOK).
		JSON(h.Success(ctx, fiber.Map{
			"shasum": shasum,
		}))
}

// Base error handler for generic errors.
func (h *Handler) Error(c *fiber.Ctx, err error) error {
	ctx := context.FromFiberCtx(c)

	switch err {
	case fiber.ErrNotFound:
		return c.Status(fiber.StatusNotFound).
			JSON(h.Failure(ctx, nil, service.ErrNotFound))
	case fiber.ErrMethodNotAllowed:
		return c.Status(fiber.StatusMethodNotAllowed).
			JSON(h.Failure(ctx, nil, service.ErrMethodNotAllowed))
	case fiber.ErrRequestTimeout:
		return c.Status(fiber.StatusRequestTimeout).
			JSON(h.Failure(ctx, nil, service.ErrRequestTimeout))
	case fiber.ErrTooManyRequests:
		return c.Status(fiber.StatusTooManyRequests).
			JSON(h.Failure(ctx, nil, service.ErrTooManyRequests))
	case fiber.ErrRequestEntityTooLarge:
		return c.Status(fiber.StatusRequestEntityTooLarge).
			JSON(h.Failure(ctx, nil, service.ErrRequestEntityTooLarge))
	}

	return c.Status(fiber.StatusInternalServerError).
		JSON(h.Failure(ctx, err, service.ErrInternalServer))
}
