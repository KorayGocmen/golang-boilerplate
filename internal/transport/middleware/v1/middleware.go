package middleware_v1

import v1 "github.com/koraygocmen/golang-boilerplate/internal/transport/response/v1"

// Handler.
type Handler struct {
	*v1.Response
}

func New(v1Response *v1.Response) *Handler {
	return &Handler{v1Response}
}
