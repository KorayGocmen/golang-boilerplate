package user_v1

import (
	"fmt"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	"github.com/koraygocmen/golang-boilerplate/internal/errhandle"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"github.com/koraygocmen/golang-boilerplate/internal/repo"
	"github.com/koraygocmen/golang-boilerplate/internal/slack"
	"github.com/koraygocmen/null"
)

// Function definitions to make it easier to reference the functions.
type CreateParams struct {
	Email      null.String `json:"email"`
	Password   null.String `json:"password"`
	GivenNames null.String `json:"givenNames"`
	Surname    null.String `json:"surname"`
}
type CreateFn func(ctx context.Ctx, params *CreateParams) (*User.User, errapi.Error, error)
type UpdateParams struct {
	Email      null.String `json:"email"`
	Password   null.String `json:"password"`
	GivenNames null.String `json:"givenNames"`
	Surname    null.String `json:"surname"`
}
type UpdateFn func(ctx context.Ctx, user *User.User, userSession *UserSession.UserSession, params *UpdateParams) (*User.User, errapi.Error, error)
type GetFn func(ctx context.Ctx, id int64) (*User.User, errapi.Error, error)
type DeleteFn func(ctx context.Ctx, user *User.User) (*User.User, errapi.Error, error)
type ValidateParams struct {
	Email      null.String `json:"email"`
	Password   null.String `json:"password"`
	GivenNames null.String `json:"givenNames"`
	Surname    null.String `json:"surname"`
}

// Service definition.
type Service struct {
	Create CreateFn
	Update UpdateFn
	Get    GetFn
	Delete DeleteFn
}

func New(tx *repo.Transaction) *Service {
	return &Service{
		Create: create(tx),
		Update: update(tx),
		Get:    get(tx),
		Delete: delete(tx),
	}
}

// Methods.
func create(tx *repo.Transaction) CreateFn {
	return func(ctx context.Ctx, params *CreateParams) (*User.User, errapi.Error, error) {
		if params == nil {
			return nil, ErrCreate.UserCreateParamsMissing, nil
		}

		// Create the user object.
		user := &User.User{}

		paramsValidated, aerr := validateParams(&ValidateParams{
			Email:      params.Email,
			Password:   params.Password,
			GivenNames: params.GivenNames,
			Surname:    params.Surname,
		})
		if aerr != nil {
			return nil, aerr, nil
		}

		if paramsValidated.Email.Valid {
			userFound, err := tx.User.GetByEmail(ctx, paramsValidated.Email.String)
			if err != nil {
				err = fmt.Errorf("user service create error: %w", err)
				return nil, nil, err
			}

			if userFound != nil {
				return nil, ErrCreate.UserEmailExists, nil
			}

			user.Email = paramsValidated.Email
		}

		if paramsValidated.Password.Valid {
			user.Password = paramsValidated.Password
			if err := user.PasswordHashCreate(); err != nil {
				err = fmt.Errorf("user service create error: %w", err)
				return nil, nil, err
			}
		}

		if paramsValidated.GivenNames.Valid {
			user.GivenNames = paramsValidated.GivenNames
		}

		if paramsValidated.Surname.Valid {
			user.Surname = paramsValidated.Surname
		}

		// Create the user.
		if err := tx.User.Create(ctx, user); err != nil {
			err = fmt.Errorf("user service create error: %w", err)
			return nil, nil, err
		}

		// Send slack notification.
		if err := slack.Client.MessageEvent(ctx, "New user", user.Email.String); err != nil {
			err = fmt.Errorf("user service create error: %w", err)
			errhandle.Handle(ctx, nil, err, false)
		}

		return user, nil, nil
	}
}

func update(tx *repo.Transaction) UpdateFn {
	return func(ctx context.Ctx, user *User.User, userSession *UserSession.UserSession, params *UpdateParams) (*User.User, errapi.Error, error) {
		if user == nil {
			return nil, ErrUpdate.UserMissing, nil
		}

		if params == nil {
			return nil, ErrUpdate.UserUpdateParamsMissing, nil
		}

		paramsValidated, aerr := validateParams(&ValidateParams{
			Email:      params.Email,
			Password:   params.Password,
			GivenNames: params.GivenNames,
			Surname:    params.Surname,
		})
		if aerr != nil {
			return nil, aerr, nil
		}

		var (
			userPasswordChanged bool
		)

		if paramsValidated.Email.Valid {
			if paramsValidated.Email.String != user.Email.String {
				userFound, err := tx.User.GetByEmail(ctx, paramsValidated.Email.String)
				if err != nil {
					err = fmt.Errorf("user service update error: %w", err)
					return nil, nil, err
				}

				if userFound != nil {
					return nil, ErrUpdate.UserEmailExists, nil
				}

				user.Email = paramsValidated.Email
				user.EmailVerified = null.BoolFrom(false)
			}
		}

		if paramsValidated.Password.Valid {
			user.Password = paramsValidated.Password
			if err := user.PasswordHashCreate(); err != nil {
				err = fmt.Errorf("user service update error: %w", err)
				return nil, nil, err
			}

			userPasswordChanged = true
		}

		if paramsValidated.GivenNames.Valid {
			user.GivenNames = paramsValidated.GivenNames
		}

		if paramsValidated.Surname.Valid {
			user.Surname = paramsValidated.Surname
		}

		// Update the user.
		if err := tx.User.Save(ctx, user); err != nil {
			err = fmt.Errorf("user service update error: %w", err)
			return nil, nil, err
		}

		if userPasswordChanged {
			// Invalidate all user sessions except the current one.
			userSessions, err := tx.UserSession.ListActive(ctx, user.ID)
			if err != nil {
				err = fmt.Errorf("user service update error: %w", err)
				return nil, nil, err
			}

			for _, us := range userSessions {
				if userSession != nil && us.ID == userSession.ID {
					continue
				}

				if err := tx.UserSession.Delete(ctx, us); err != nil {
					err = fmt.Errorf("user service update error: %w", err)
					return nil, nil, err
				}
			}
		}

		return user, nil, nil
	}
}

func get(tx *repo.Transaction) GetFn {
	return func(ctx context.Ctx, id int64) (*User.User, errapi.Error, error) {
		if id == 0 {
			return nil, ErrGet.UserIdMissing, nil
		}

		user, err := tx.User.GetByID(ctx, id)
		if err != nil {
			err = fmt.Errorf("user service get by id error: %w", err)
			return nil, nil, err
		}

		if user == nil {
			return nil, ErrGet.UserNotFound, nil
		}

		return user, nil, nil
	}
}

func delete(tx *repo.Transaction) DeleteFn {
	return func(ctx context.Ctx, user *User.User) (*User.User, errapi.Error, error) {
		if user == nil {
			return nil, ErrDelete.UserMissing, nil
		}

		// Should delete the user.
		return user, nil, nil
	}
}

// validateParams validates the user params.
func validateParams(params *ValidateParams) (*ValidateParams, errapi.Error) {
	// Validate email input.
	if params.Email.Valid {
		email := strings.ToLower(strings.TrimSpace(params.Email.String))

		if email == "" {
			return nil, ErrValidateParams.UserEmailMissing
		}

		if !govalidator.IsEmail(email) {
			return nil, ErrValidateParams.UserEmailInvalid
		}

		email = strings.ToLower(email)
		params.Email = null.StringFrom(email)
	}

	// Validate password input.
	if params.Password.Valid {
		password := params.Password.String
		if password == "" {
			return nil, ErrValidateParams.UserPasswordMissing
		}

		if len(password) < 6 {
			return nil, ErrValidateParams.UserPasswordInvalid
		}

		params.Password = null.StringFrom(password)
	}

	// Validate given names input.
	if params.GivenNames.Valid {
		givenNames := User.ToGivenNames(params.GivenNames.String)
		if givenNames == "" {
			return nil, ErrValidateParams.UserGivenNamesMissing
		}

		if len(givenNames) < 2 {
			return nil, ErrValidateParams.UserGivenNamesInvalid
		}

		params.GivenNames = null.StringFrom(string(givenNames))
	}

	// Validate surname input.
	if params.Surname.Valid {
		surname := User.ToSurname(params.Surname.String)
		if surname == "" {
			return nil, ErrValidateParams.UserSurnameMissing
		}

		if len(surname) < 2 {
			return nil, ErrValidateParams.UserSurnameInvalid
		}

		params.Surname = null.StringFrom(string(surname))
	}

	return params, nil
}
