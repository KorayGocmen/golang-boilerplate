package user_v1

import (
	"testing"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
	User "github.com/koraygocmen/golang-boilerplate/internal/model/user"
	UserSession "github.com/koraygocmen/golang-boilerplate/internal/model/user_session"
	"github.com/koraygocmen/golang-boilerplate/internal/repo"
	UserRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user"
	UserSessionRepo "github.com/koraygocmen/golang-boilerplate/internal/repo/user_session"
	"github.com/koraygocmen/null"
)

func TestCreate(t *testing.T) {
	// Override the transaction function in order
	// to return a transaction with the userRepoTest
	// which we can modify the methods of.
	tx := &repo.Transaction{
		User: &UserRepo.Repo{
			Create: func(ctx context.Ctx, u *User.User) error {
				return nil
			},
			GetByEmail: func(ctx context.Ctx, email string) (*User.User, error) {
				return nil, nil
			},
		},
	}

	userService := New(tx)

	params := &CreateParams{
		Email:    null.String{},
		Password: null.String{},
	}

	// Test missing fields.
	_, aerr, err := userService.Create(context.Background(), nil)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserCreateParamsMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserCreateParamsMissing, aerr)
	}

	// Test missing email and phone number.
	_, aerr, err = userService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserEmailPhoneNumberMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserEmailPhoneNumberMissing, aerr)
	}

	// Test invalid email.
	params.Email = null.StringFrom("invalid")
	_, aerr, err = userService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserEmailInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserEmailInvalid, aerr)
	}
	params.Email = null.StringFrom("koray@test.com")

	// Test email already exists.
	tx.User.GetByEmail = func(ctx context.Ctx, email string) (*User.User, error) {
		return &User.User{ID: 1}, nil
	}
	_, aerr, err = userService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserEmailExists) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserEmailExists, aerr)
	}
	tx.User.GetByEmail = func(ctx context.Ctx, email string) (*User.User, error) {
		return nil, nil
	}

	// Test missing password.
	params.Password = null.StringFrom("")
	_, aerr, err = userService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserPasswordMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserPasswordMissing, aerr)
	}
	params.Password = null.String{}

	// Test invalid password.
	params.Password = null.StringFrom("123")
	_, aerr, err = userService.Create(context.Background(), params)
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserPasswordInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserPasswordInvalid, aerr)
	}
	params.Password = null.String{}

	// Test success with email.
	user, aerr, err := userService.Create(context.Background(), &CreateParams{
		Email:    null.StringFrom("koray@test.com"),
		Password: null.StringFrom("123456"),
	})
	if err != nil {
		t.Fatalf(`want: create err nil; got: err = %v`, err)
	}
	if aerr != nil {
		t.Fatalf(`want: aerr nil; got: aerr = %v`, aerr)
	}
	if user == nil {
		t.Fatalf(`want: user not nil; got: nil`)
	}
}

func TestUpdate(t *testing.T) {
	// Override the transaction function in order
	// to return a transaction with the userRepoTest
	// which we can modify the methods of.
	tx := &repo.Transaction{
		User: &UserRepo.Repo{
			Save: func(ctx context.Ctx, user *User.User) error {
				return nil
			},
		},
		UserSession: &UserSessionRepo.Repo{
			ListActive: func(ctx context.Ctx, userId int64) ([]*UserSession.UserSession, error) {
				return []*UserSession.UserSession{
					{ID: 1},
					{ID: 2},
				}, nil
			},
		},
	}

	userService := New(tx)

	user := &User.User{
		Email:         null.StringFrom("koray@test.com"),
		EmailVerified: null.BoolFrom(true),
		Password:      null.StringFrom("123456"),
		PasswordHash:  null.StringFrom("hash"),
		GivenNames:    null.StringFrom("Koray"),
		Surname:       null.StringFrom("Gocmen"),
	}
	userSession := &UserSession.UserSession{ID: 1}

	params := &UpdateParams{
		Password:   null.String{},
		GivenNames: null.String{},
		Surname:    null.String{},
	}

	// Test missing fields.

	// Test missing user.
	_, aerr, err := userService.Update(context.Background(), nil, nil, nil)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserMissing, aerr)
	}

	// Test missing params.
	_, aerr, err = userService.Update(context.Background(), user, userSession, nil)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserUpdateParamsMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserUpdateParamsMissing, aerr)
	}

	// Test missing email.
	params.Email = null.StringFrom("")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserEmailMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserEmailMissing, aerr)
	}

	// Test invalid email.
	params.Email = null.StringFrom("koray")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserEmailInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserEmailInvalid, aerr)
	}

	// Test email already exists.
	params.Email = null.StringFrom("koray2@test.com")
	tx.User.GetByEmail = func(ctx context.Ctx, email string) (*User.User, error) {
		return &User.User{ID: 1}, nil
	}
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserEmailExists) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserEmailExists, aerr)
	}
	tx.User.GetByEmail = func(ctx context.Ctx, email string) (*User.User, error) {
		return nil, nil
	}

	// Test missing password.
	params.Password = null.StringFrom("")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserPasswordMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserPasswordMissing, aerr)
	}
	params.Password = null.StringFrom("123456")

	// Test invalid password.
	params.Password = null.StringFrom("12345")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserPasswordInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserPasswordInvalid, aerr)
	}
	params.Password = null.StringFrom("123456")

	// Test missing given names.
	params.GivenNames = null.StringFrom("")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserGivenNamesMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserGivenNamesMissing, aerr)
	}

	// Test invalid given names.
	params.GivenNames = null.StringFrom("A")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserGivenNamesInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserGivenNamesInvalid, aerr)
	}
	params.GivenNames = null.StringFrom("koray")

	// Test missing surname.
	params.Surname = null.StringFrom("")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSurnameMissing) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSurnameMissing, aerr)
	}

	// Test invalid surname.
	params.Surname = null.StringFrom("A")
	_, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if !errapi.Is(aerr, errapi.ErrCodeUserSurnameInvalid) {
		t.Fatalf(`want: aerr = %v; got: aerr = %v`, errapi.ErrCodeUserSurnameInvalid, aerr)
	}
	params.Surname = null.StringFrom("göçmen")

	tx.UserSession.Delete = func(ctx context.Ctx, userSession *UserSession.UserSession) error {
		if userSession.ID != 2 {
			t.Fatalf(`want: user session id = 2; got: user session id = %v`, userSession.ID)
		}
		return nil
	}

	user, aerr, err = userService.Update(context.Background(), user, userSession, params)
	if err != nil {
		t.Fatalf(`want: update err nil; got: err = %v`, err)
	}
	if aerr != nil {
		t.Fatalf(`want: aerr nil; got: aerr = %v`, aerr)
	}

	if user == nil {
		t.Fatalf(`want: user not nil; got: nil`)
	}

	if got := user.Email.String; got != "koray2@test.com" {
		t.Fatalf(`want: user email = "koray2@test.com"; got user password = %v`, got)
	}

	if user.EmailVerified.Valid && user.EmailVerified.Bool {
		t.Fatalf(`want: user email not verified; got user email verified`)
	}

	if got := user.Password.String; got != "123456" {
		t.Fatalf(`want: user password = "123456"; got user password = %v`, got)
	}

	if got := user.PasswordHash.String; got == "" {
		t.Fatalf(`want: user password hash not nil; got user password hash nil`)
	}

	if got := user.GivenNames.String; got != "Koray" {
		t.Fatalf(`want: user given names = "Koray"; got user given names = %v`, got)
	}

	if got := user.Surname.String; got != "Göçmen" {
		t.Fatalf(`want: user surname = "Göçmen"; got user surname = %v`, got)
	}
}
