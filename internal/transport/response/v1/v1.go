package response_v1

import (
	"github.com/koraygocmen/golang-boilerplate/internal/transport/handler"
)

type Response struct {
	*handler.Handler
}

func New(handler *handler.Handler) *Response {
	return &Response{handler}
}
