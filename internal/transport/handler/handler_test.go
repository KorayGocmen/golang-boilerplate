package handler

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/logger"
	"github.com/koraygocmen/golang-boilerplate/internal/service"
)

var (
	appTest *fiber.App
)

func TestMain(m *testing.M) {
	// Setup logger.
	logger.Logger, _ = logger.New(logger.Config{
		Mode: string(logger.ModeNone),
	})

	// Set up service.
	service.Service.Transaction = func(ctx context.Ctx, timeout time.Duration) *service.Transaction {
		return &service.Transaction{
			Ping: func() error {
				return nil
			},
			Commit: func() error {
				return nil
			},
		}
	}

	// Setup app.
	appTest = fiber.New(fiber.Config{})

	// Create handler to be tested.
	handlerTest := New(Config{
		SHASUM: "1111111",
	})

	// Set up the routes to be tested.
	appTest.Get("/health", handlerTest.Health)

	m.Run()
}

func TestSetup(t *testing.T) {
	if appTest == nil {
		t.Fatalf("app is nil")
	}

	if len(appTest.GetRoutes()) == 0 {
		t.Fatalf("app has no routes")
	}
}

func TestHealth(t *testing.T) {
	resp, err := appTest.Test(httptest.NewRequest(fiber.MethodGet, "/health", nil))

	if err != nil {
		t.Fatalf("/health error: %v", err)
	}

	if statusCode := resp.StatusCode; statusCode != fiber.StatusOK {
		t.Fatalf("/health status = %d; want %d", statusCode, fiber.StatusOK)
	}
}
