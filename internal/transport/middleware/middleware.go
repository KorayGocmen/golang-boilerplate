package middleware

import (
	"runtime"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/internal/transport/handler"
)

// Setup the shared middleware to the fiber app.
func Setup(app *fiber.App, handler *handler.Handler) {
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New())
	app.Use(etag.New())
	app.Use(requestid.New(requestid.Config{
		ContextKey: string(context.KeyRequestID),
	}))
	app.Use(cors.New(cors.Config{
		AllowMethods:  strings.Join(allowedMethods, ","),
		ExposeHeaders: strings.Join(exposedHeaders, ","),
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	// Context middleware.
	app.Use(func(c *fiber.Ctx) error {
		// Context created here is used by each handler to log the request.
		// Defered cancel is called when each handler returns.
		ctx, cancel := context.NewFiberCtx(c)
		c.Locals("ctx", ctx)
		c.Locals("cancel", cancel)

		defer func() {
			ctx := c.Locals("ctx").(context.Ctx)

			// Response status and body are available when returning from handler.
			ctx = context.WithValue(ctx, context.KeyStatus, c.Context().Response.StatusCode())
			if !strings.Contains(string(c.Context().Path()), "/health") {
				logger.Logger.Infof(ctx, "")
			}
			cancel()
		}()

		return c.Next()
	})

	// Recover middleware with stack trace.
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			ctx := context.FromFiberCtx(c)

			buf := make([]byte, 1<<16)
			stackSize := runtime.Stack(buf, true)

			ctx = context.WithValue(ctx, context.KeyError, e)
			ctx = context.WithValue(ctx, context.KeyErrorStack, string(buf[0:stackSize]))
			c.Locals("ctx", ctx)
		},
	}))

	// Rate limiter middleware.
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			isGetRequest := c.Method() == fiber.MethodGet
			isHealthCheck := strings.Contains(string(c.Context().Path()), "/health")
			return (env.IsDev() || isHealthCheck || isGetRequest)
		},
		Max:               20,
		Expiration:        60 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			ctx := context.FromFiberCtx(c)
			ip := context.RemoteIP(ctx)
			return ip
		},
		LimitReached: func(c *fiber.Ctx) error {
			// Errors returned from middleware are handled by the handler.Error
			// which was setup when initializing the fiber app.
			return fiber.ErrTooManyRequests
		},
	}))

}
