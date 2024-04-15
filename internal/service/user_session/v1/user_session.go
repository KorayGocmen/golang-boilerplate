package user_session_v1

import (
	"fmt"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"github.com/koraygocmen/golang-boilerplate/internal/repo"
	UserService "github.com/koraygocmen/golang-boilerplate/internal/service/user"
	"github.com/koraygocmen/null"
)

// Function definitions to make it easier to reference the functions.
type CreateParams struct {
	Email          null.String `json:"email"`
	Password       null.String `json:"password"`
	Purpose        string      `json:"purpose"`
	DeliveryMethod string      `json:"deliveryMethod"`
	ClientIP       string      // internal use only
}
type CreateFn func(ctx context.Ctx, params *CreateParams) (*UserSession.UserSession, errapi.Error, error)
type GetFn func(ctx context.Ctx, id int64) (*UserSession.UserSession, errapi.Error, error)
type DeleteFn func(ctx context.Ctx, userID int64, userSessionActive *UserSession.UserSession) (errapi.Error, error)

// Service definition.
type Service struct {
	Create CreateFn
	Get    GetFn
	Delete DeleteFn
}

func New(tx *repo.Transaction, userService *UserService.Service) *Service {
	return &Service{
		Create: create(tx, userService),
		Get:    get(tx),
		Delete: delete(tx),
	}
}

// Methods.
func create(tx *repo.Transaction, userService *UserService.Service) CreateFn {
	return func(ctx context.Ctx, params *CreateParams) (*UserSession.UserSession, errapi.Error, error) {
		if params == nil {
			return nil, ErrCreate.UserSessionCreateParamsMissing, nil
		}

		var (
			clientIP = params.ClientIP
			email    = params.Email
			password = params.Password
			purpose  = UserSession.ToPurpose(params.Purpose)
		)

		if clientIP == "" {
			err := fmt.Errorf("user session service create error: client ip is missing")
			return nil, nil, err
		}

		if purpose == "" {
			return nil, ErrCreate.UserSessionPurposeMissing, nil
		}

		if !UserSession.Purposes[purpose] {
			return nil, ErrCreate.UserSessionPurposeInvalid, nil
		}

		// Create the user object.
		var user *User.User

		if email.Valid {
			var err error
			user, err = tx.User.GetByEmail(ctx, email.String)
			if err != nil {
				err = fmt.Errorf("user session service create error: %w", err)
				return nil, nil, err
			}
		}

		if user == nil {
			return nil, ErrCreate.UserSessionCredentialsInvalid, nil
		}

		// Check the password only if the user does have a password.
		// Allow users to signup with only a phonenumber/email and verify.
		if purpose == UserSession.PurposeSessionCreate && user.PasswordHash.Valid {
			if !password.Valid {
				return nil, ErrCreate.UserSessionPasswordMissing, nil
			}

			if !user.PasswordHashCompare(password.String) {
				return nil, ErrCreate.UserSessionCredentialsInvalid, nil
			}
		}

		userSession := &UserSession.UserSession{
			UserID:   user.ID,
			ClientIP: clientIP,
			Purpose:  purpose,
		}
		if err := userSession.TokenHashCreate(); err != nil {
			err = fmt.Errorf("user session service create error: %w", err)
			return nil, nil, err
		}

		userSession.ExpireAtCreate()
		if err := tx.UserSession.Create(ctx, userSession); err != nil {
			err = fmt.Errorf("user session service create error: %w", err)
			return nil, nil, err
		}

		return userSession, nil, nil
	}
}

func get(tx *repo.Transaction) GetFn {
	return func(ctx context.Ctx, id int64) (*UserSession.UserSession, errapi.Error, error) {
		userSession, err := tx.UserSession.GetByID(ctx, id)
		if err != nil {
			err = fmt.Errorf("user session service get by id error: %w", err)
			return nil, nil, err
		}

		if userSession == nil {
			return nil, ErrGet.UserSessionNotFound, nil
		}

		return userSession, nil, nil
	}
}

func delete(tx *repo.Transaction) DeleteFn {
	return func(ctx context.Ctx, userID int64, userSessionActive *UserSession.UserSession) (errapi.Error, error) {
		if userID == 0 {
			return ErrDelete.UserIdMissing, nil
		}

		userSessions, err := tx.UserSession.ListActive(ctx, userID)
		if err != nil {
			err = fmt.Errorf("user session service delete error: %w", err)
			return nil, err
		}

		for _, userSession := range userSessions {
			if userSessionActive != nil && userSession.ID == userSessionActive.ID {
				continue
			}

			if err := tx.UserSession.Delete(ctx, userSession); err != nil {
				err = fmt.Errorf("user session service delete error: %w", err)
				return nil, err
			}
		}

		return nil, nil
	}
}
