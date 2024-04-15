package user_v1

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"github.com/koraygocmen/golang-boilerplate/internal/service"
	UserServiceV1 "github.com/koraygocmen/golang-boilerplate/internal/service/user/v1"
	v1 "github.com/koraygocmen/golang-boilerplate/internal/transport/response/v1"
)

// Handler.
type Handler struct {
	*v1.Response
}

func New(v1Response *v1.Response) *Handler {
	return &Handler{v1Response}
}

// POST /v1/users
func (v1 *Handler) Create(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	var userCreateParams UserServiceV1.CreateParams
	if err := c.BodyParser(&userCreateParams); err != nil {
		err = fmt.Errorf("user handle create error: body parser error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	user, aerr, err := srv.User.V1.Create(ctx, &userCreateParams)
	if err != nil {
		err = fmt.Errorf("user handle create error: %w", err)
		srv.Rollback(err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, err, aerr))
	}

	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("user handle create error: commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	return c.Status(fiber.StatusCreated).
		JSON(v1.Handler.Success(ctx, fiber.Map{
			"user": user,
		}))
}

// GET /v1/users
func (v1 *Handler) Get(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	user, ok := c.Locals("user").(*User.User)
	if !ok {
		err := fmt.Errorf("user handle get error: user not found in local ctx")
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	return c.Status(fiber.StatusOK).
		JSON(v1.Handler.Success(ctx, fiber.Map{
			"user": user,
		}))
}

// PUT /v1/users
func (v1 *Handler) Update(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	user, ok := c.Locals("user").(*User.User)
	if !ok {
		err := fmt.Errorf("user handle update error: user not found in local ctx")
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	userSession, ok := c.Locals("userSession").(*UserSession.UserSession)
	if !ok {
		err := fmt.Errorf("user handle update error: user session not found in local ctx")
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	var userUpdateParams UserServiceV1.UpdateParams
	if err := c.BodyParser(&userUpdateParams); err != nil {
		err = fmt.Errorf("user handle update error: body parser error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	user, aerr, err := srv.User.V1.Update(ctx, user, userSession, &userUpdateParams)
	if err != nil {
		err = fmt.Errorf("user handle update error: %w", err)
		srv.Rollback(err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, err, aerr))
	}

	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("user handle update error: commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	return c.Status(fiber.StatusOK).
		JSON(v1.Handler.Success(ctx, fiber.Map{
			"user": user,
		}))
}

// DELETE /v1/users
func (v1 *Handler) Delete(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	user, ok := c.Locals("user").(*User.User)
	if !ok {
		err := fmt.Errorf("user handle delete error: user not found in local ctx")
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	user, aerr, err := srv.User.V1.Delete(ctx, user)
	if err != nil {
		err = fmt.Errorf("user handle delete error: %w", err)
		srv.Rollback(err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, err, aerr))
	}

	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("user handle delete error: commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	return c.Status(fiber.StatusOK).
		JSON(v1.Handler.Success(ctx, fiber.Map{
			"user": user,
		}))
}
