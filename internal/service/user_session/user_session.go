package user_session

import (
	"github.com/koraygocmen/golang-boilerplate/internal/repo"
	UserService "github.com/koraygocmen/golang-boilerplate/internal/service/user"
	UserSessionServiceV1 "github.com/koraygocmen/golang-boilerplate/internal/service/user_session/v1"
)

// Service definition.
type Service struct {
	V1 *UserSessionServiceV1.Service
}

func New(tx *repo.Transaction, userService *UserService.Service) *Service {
	return &Service{
		V1: UserSessionServiceV1.New(tx, userService),
	}
}
