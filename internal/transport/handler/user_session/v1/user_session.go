package user_session_v1

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"github.com/koraygocmen/golang-boilerplate/internal/service"
	UserSessionServiceV1 "github.com/koraygocmen/golang-boilerplate/internal/service/user_session/v1"
	v1 "github.com/koraygocmen/golang-boilerplate/internal/transport/response/v1"
)

// Handler.
type Handler struct {
	*v1.Response
}

func New(v1Response *v1.Response) *Handler {
	return &Handler{v1Response}
}

// POST /v1/users/sessions
func (v1 *Handler) Create(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	var userSessionCreateParams UserSessionServiceV1.CreateParams
	if err := c.BodyParser(&userSessionCreateParams); err != nil {
		err = fmt.Errorf("user session handle create error: body parser error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}
	userSessionCreateParams.ClientIP = context.RemoteIP(ctx)

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	userSession, aerr, err := srv.UserSession.V1.Create(ctx, &userSessionCreateParams)
	if err != nil {
		err = fmt.Errorf("user session handle create error: %w", err)
		srv.Rollback(err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("user handle create error: commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	return c.Status(fiber.StatusCreated).
		JSON(v1.Handler.Success(ctx, fiber.Map{
			"userSession": userSession,
		}))
}

// DELETE /v1/users/sessions
func (v1 *Handler) Delete(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	userSession, ok := c.Locals("userSession").(*UserSession.UserSession)
	if !ok {
		err := fmt.Errorf("user sessions handle delete error: user session not found in local ctx")
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	aerr, err := srv.UserSession.V1.Delete(ctx, userSession.UserID, nil)
	if err != nil {
		err = fmt.Errorf("user sessions handle delete error: %w", err)
		srv.Rollback(err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("user sessions handle delete error: commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	return c.Status(fiber.StatusOK).
		JSON(v1.Handler.Success(ctx, fiber.Map{
			"userSession": userSession,
		}))
}
