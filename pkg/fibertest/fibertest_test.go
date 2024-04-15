package fibertest

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

var (
	testapp *fiber.App
)

func TestMain(m *testing.M) {
	testapp = fiber.New()
	testapp.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
		})
	})

	testapp.Post("/", func(c *fiber.Ctx) error {
		var req fiber.Map
		if err := c.BodyParser(&req); err != nil {
			panic(err)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"payload": req,
		})
	})

	m.Run()
}

func TestNewRequestGet(t *testing.T) {
	statusCode, resraw, resbody, err := Request(testapp, fiber.MethodGet, "/", nil, nil)
	if err != nil {
		t.Fatalf("GET / error: %v", err)
	}

	if statusCode := statusCode; statusCode != fiber.StatusOK {
		t.Fatalf("GET / status = %d; want %d", statusCode, fiber.StatusOK)
	}

	if len(resraw) == 0 {
		t.Fatalf("GET / response = empty; want not empty")
	}

	if resbody["success"] != true {
		t.Fatalf("GET / success = %t; want %t", resbody["success"], true)
	}
}

func TestNewRequestPost(t *testing.T) {
	statusCode, resraw, resbody, err := Request(testapp, fiber.MethodPost, "/", nil, fiber.Map{
		"key": "val",
	})
	if err != nil {
		t.Fatalf("POST / error: %v", err)
	}

	if statusCode := statusCode; statusCode != fiber.StatusCreated {
		t.Fatalf("POST / status = %d; want %d", statusCode, fiber.StatusCreated)
	}

	if len(resraw) == 0 {
		t.Fatalf("POST / response = empty; want not empty")
	}

	if resbody["success"].(bool) != true {
		t.Fatalf("POST / success = %t; want %t", resbody["success"], true)
	}

	payload := resbody["payload"].(map[string]interface{})
	if payload["key"] != "val" {
		t.Fatalf("POST / payload = %s; want %s", payload["key"], "val")
	}
}
