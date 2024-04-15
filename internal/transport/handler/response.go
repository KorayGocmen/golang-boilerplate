package handler

import (
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	"github.com/koraygocmen/golang-boilerplate/internal/errhandle"
)

type ApiResponse struct {
	Success bool         `json:"success"`
	Body    interface{}  `json:"body,omitempty"`
	Error   errapi.Error `json:"error,omitempty"`
}

func (h *Handler) Success(ctx context.Ctx, body interface{}) ApiResponse {
	return ApiResponse{
		Success: true,
		Body:    body,
	}
}

func (h *Handler) Failure(ctx context.Ctx, err error, aerr errapi.Error) ApiResponse {
	errhandle.Handle(ctx, aerr, err, false)
	return ApiResponse{
		Success: false,
		Error:   aerr.Localize(ctx),
	}
}
