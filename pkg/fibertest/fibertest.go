package fibertest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
)

func Request(app *fiber.App, method, path string, headers map[string]string, body interface{}) (int, []byte, fiber.Map, error) {
	reqbody, err := json.Marshal(body)
	if err != nil {
		err = fmt.Errorf("fibertest request error: marshal body error: %w", err)
		return 0, nil, nil, err
	}

	req := httptest.NewRequest(method, path, bytes.NewReader(reqbody))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	for headerKey, headerVal := range headers {
		req.Header.Set(headerKey, headerVal)
	}

	res, err := app.Test(req)
	if err != nil {
		err = fmt.Errorf("fibertest request error: fiber app test error: %w", err)
		return 0, nil, nil, err
	}

	resraw, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("fibertest request error: io read all res body error: %w", err)
		return 0, nil, nil, err
	}

	var resbody fiber.Map
	if len(resraw) > 0 {
		if err := json.Unmarshal(resraw, &resbody); err != nil {
			err = fmt.Errorf(`fibertest request error: unmarshal res body error: %w, res body raw: "%s"`, err, string(resraw))
			return 0, resraw, nil, err
		}
	}

	return res.StatusCode, resraw, resbody, nil
}
