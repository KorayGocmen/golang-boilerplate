package user_session_v1

import (
	"testing"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	"github.com/koraygocmen/null"

	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"github.com/koraygocmen/golang-boilerplate/internal/repo"
	UserRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user"
	UserSessionRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user_session"
	UserService "github.com/koraygocmen/golang-boilerplate/internal/service/user"
	UserServiceV1 "github.com/koraygocmen/golang-boilerplate/internal/service/user/v1"
)

func TestCreate(t *testing.T) {
	user := &User.User{
		ID:       1,
		Email:    null.StringFrom("koray@test.com"),
		Password: null.StringFrom("123456"),
	}
	user.PasswordHashCreate()

	tx := &repo.Transaction{
		User: &UserRepo.Repo{},
		UserSession: &UserSessionRepo.Repo{
			Create: func(ctx context.Ctx, userSession *UserSession.UserSession) error {
				return nil
			},
		},
	}

	// Create service dependencies for testing.
	userService := &UserService.Service{
		V1: &UserServiceV1.Service{},
	}

	userSessionService := New(tx, userService)

	params := &CreateParams{
		ClientIP: "",
	}

	// Test missing params.
	_, aerr, err := userSessionService.Create(context.Background(), nil)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionCreateParamsMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionCreateParamsMissing, aerr)
	}

	// Test missing client IP. This should not happen in real life.
	_, _, err = userSessionService.Create(context.Background(), params)
	if err == nil {
		t.Fatalf(`want: create err not nil; got: err = nil`)
	}
	params.ClientIP = "0.0.0.0"

	// Test mssing purpose.
	_, aerr, err = userSessionService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionPurposeMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionPurposeMissing, aerr)
	}

	// Test unknown purpose.
	params.Purpose = "invalid"
	_, aerr, err = userSessionService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionPurposeInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionPurposeInvalid, aerr)
	}
	params.Purpose = "session_create"

	// Test email not found.
	params.Email = null.StringFrom("koray@test.com")
	tx.User.GetByEmail = func(ctx context.Ctx, email string) (*User.User, error) {
		return nil, nil
	}
	_, aerr, err = userSessionService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionCredentialsInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionCredentialsInvalid, aerr)
	}

	// Test missing password.
	_, aerr, err = userSessionService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionCredentialsInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionCredentialsInvalid, aerr)
	}
	params.Password = null.StringFrom("123456")

	// Test password compare (wrong password).
	params.Password = null.StringFrom("1234567")
	_, aerr, err = userSessionService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionCredentialsInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionCredentialsInvalid, aerr)
	}
	params.Password = null.StringFrom("123456")

	// Test password compare (correct password).
	_, aerr, err = userSessionService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if aerr != nil {
		t.Fatalf(`want: aerr = nil; got: aerr = %v`, aerr)
	}
}

func TestGetByID(t *testing.T) {
	tx := &repo.Transaction{
		UserSession: &UserSessionRepo.Repo{
			GetByID: func(ctx context.Ctx, id int64) (*UserSession.UserSession, error) {
				return &UserSession.UserSession{ID: 1}, nil
			},
		},
	}

	// Create service dependencies for testing.
	userService := &UserService.Service{
		V1: &UserServiceV1.Service{},
	}

	userSessionService := New(tx, userService)

	// Test user session not found.
	getByID := tx.UserSession.GetByID
	tx.UserSession.GetByID = func(ctx context.Ctx, id int64) (*UserSession.UserSession, error) {
		return nil, nil
	}
	_, aerr, err := userSessionService.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf(`want: get err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSessionNotFound) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSessionNotFound, aerr)
	}
	tx.UserSession.GetByID = getByID

	// Test user session found.
	_, aerr, err = userSessionService.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf(`want: get err nil; got: err = %v`, err)
	}
	if aerr != nil {
		t.Fatalf(`want: aerr = nil; got: aerr = %v`, aerr)
	}
}

func TestDelete(t *testing.T) {
	tx := &repo.Transaction{
		UserSession: &UserSessionRepo.Repo{
			ListActive: func(ctx context.Ctx, userID int64) ([]*UserSession.UserSession, error) {
				return []*UserSession.UserSession{
					{ID: 1},
					{ID: 2},
				}, nil
			},
			Delete: func(ctx context.Ctx, userSession *UserSession.UserSession) error {
				return nil
			},
		},
	}

	// Create service dependencies for testing.
	userService := &UserService.Service{
		V1: &UserServiceV1.Service{},
	}

	userSessionService := New(tx, userService)

	// Test user session service delete success.
	aerr, err := userSessionService.Delete(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf(`want: delete err nil; got: err = %v`, err)
	}
	if aerr != nil {
		t.Fatalf(`want: aerr = nil; got: aerr = %v`, aerr)
	}
}
