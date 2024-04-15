package middleware_v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	"github.com/koraygocmen/golang-boilerplate/internal/service"
)

func (v1 *Handler) UserAuth(c *fiber.Ctx) error {
	ctx := context.FromFiberCtx(c)

	authorization := c.Get(fiber.HeaderAuthorization)
	if authorization == "" {
		aerr := service.ErrUserAuth.AuthorizationMissing
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	authorizationPieces := strings.Split(authorization, "-")
	if len(authorizationPieces) != 2 {
		aerr := service.ErrUserAuth.AuthorizationInvalid
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}
	sessionIDStr, sessionToken := authorizationPieces[0], authorizationPieces[1]

	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		aerr := service.ErrUserAuth.AuthorizationInvalid
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	// Create the service with a new transaction.
	srv := service.Service.Transaction(ctx, 30*time.Second)

	userSession, aerr, err := srv.UserSession.V1.Get(ctx, sessionID)
	if err != nil {
		err = fmt.Errorf("auth user middleware error: %w", err)
		srv.Rollback(err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	if !userSession.TokenHashCompare(sessionToken) {
		srv.Rollback(nil)
		aerr := service.ErrUserAuth.AuthorizationWrong
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	if userSession.IsExpired() {
		srv.Rollback(nil)
		aerr := service.ErrUserAuth.AuthorizationExpired
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	user, aerr, err := srv.User.V1.Get(ctx, userSession.UserID)
	if err != nil {
		srv.Rollback(nil)
		err = fmt.Errorf("auth user middleware error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	if aerr != nil {
		srv.Rollback(nil)

		if !errapi.Is(aerr, errapi.ErrCodeUserNotFound) {
			err = fmt.Errorf("auth user middleware error: user not found for user session: %d", userSession.ID)
			return c.Status(fiber.StatusInternalServerError).
				JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
		}

		// User service get by id only returns ErrUserNotFound,
		// technically this should never happen.
		return c.Status(aerr.Status()).
			JSON(v1.Handler.Failure(ctx, nil, aerr))
	}

	// Commit to release the transaction.
	if err := srv.Commit(); err != nil {
		err = fmt.Errorf("auth user middleware error: commit error: %w", err)
		return c.Status(fiber.StatusInternalServerError).
			JSON(v1.Handler.Failure(ctx, err, service.ErrInternalServer))
	}

	// Attach user and user session to context.
	ctx = context.WithValue(ctx, context.KeyUserID, user.ID)
	ctx = context.WithValue(ctx, context.KeyUserSessionID, userSession.ID)

	c.Locals("ctx", ctx)
	c.Locals("user", user)
	c.Locals("userSession", userSession)

	return c.Next()
}
