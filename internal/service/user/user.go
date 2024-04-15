package user

import (
	"github.com/koraygocmen/golang-boilerplate/internal/repo"
	UserServiceV1 "github.com/koraygocmen/golang-boilerplate/internal/service/user/v1"
)

// Service definition.
type Service struct {
	V1 *UserServiceV1.Service
}

func New(tx *repo.Transaction) *Service {
	return &Service{
		V1: UserServiceV1.New(tx),
	}
}
