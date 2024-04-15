package response_v1

import (
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	"github.com/koraygocmen/golang-boilerplate/internal/errhandle"
)

type ApiResponse struct {
	Status int
	Body   ApiResponseBody
}

type ApiResponseBody struct {
	Version string       `json:"version"`
	Success bool         `json:"success"`
	Body    interface{}  `json:"body,omitempty"`
	Error   errapi.Error `json:"error,omitempty"`
}

func (r *Response) Success(ctx context.Ctx, body interface{}) ApiResponseBody {
	return ApiResponseBody{
		Version: "1",
		Success: true,
		Body:    body,
	}
}

func (r *Response) Failure(ctx context.Ctx, err error, aerr errapi.Error) ApiResponseBody {
	errhandle.Handle(ctx, aerr, err, false)
	return ApiResponseBody{
		Version: "1",
		Success: false,
		Error:   aerr.Localize(ctx),
	}
}
