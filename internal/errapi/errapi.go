package errapi

import (
	"encoding/json"

	"github.com/koraygocmen/golang-boilerplate/internal/context"
)

// Error type definitions and methods.
type Error interface {
	Status() int
	Code() string
	Error() string
	Localize(ctx context.Ctx) Error
	MarshalJSON() ([]byte, error)
}

type errapi struct {
	status   int
	code     string
	message  string
	messages map[context.Language]string
}

func New(status int, code, tr, en string) Error {
	return &errapi{
		status: status,
		code:   code,
		messages: map[context.Language]string{
			context.LangTR: tr,
			context.LangEN: en,
		},
	}
}

func (e *errapi) Status() int {
	return e.status
}

func (e *errapi) Error() string {
	return e.code
}

func (e *errapi) Code() string {
	return e.code
}

func (e *errapi) Localize(ctx context.Ctx) Error {
	if e == nil {
		return nil
	}

	if message, ok := e.messages[context.Lang(ctx)]; ok {
		e.message = message
	}
	return e
}

func (e *errapi) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"code":    e.code,
		"message": e.message,
	})
}

func Is(aerr Error, code string) bool {
	if aerr == nil {
		return false
	}

	return aerr.Code() == code
}

const (
	// Common.
	ErrCodeInternalServerError   = "internalServerError"
	ErrCodeUrlParamInvalid       = "urlParamInvalid"
	ErrCodeNotFound              = "notFound"
	ErrCodeMethodNotAllowed      = "methodNotAllowed"
	ErrCodeRequestTimeout        = "requestTimeout"
	ErrCodeTooManyRequests       = "tooManyRequests"
	ErrCodeRequestEntityTooLarge = "requestEntityTooLarge"

	// Auth.
	ErrCodeAuthorizationMissing                 = "authorizationMissing"
	ErrCodeAuthorizationInvalid                 = "authorizationInvalid"
	ErrCodeAuthorizationWrong                   = "authorizationWrong"
	ErrCodeAuthorizationExpired                 = "authorizationExpired"
	ErrCodeAuthorizationNotVerified             = "authorizationNotVerified"
	ErrCodeAuthorizationInvalidIP               = "authorizationInvalidIP"
	ErrCodeAuthorizationInsufficientAccessLevel = "authorizationInsufficientAccessLevel"

	// User Service.
	ErrCodeUserMissing                 = "userMissing"
	ErrCodeUserNotFound                = "userNotFound"
	ErrCodeUserCreateParamsMissing     = "userCreateParamsMissing"
	ErrCodeUserUpdateParamsMissing     = "userUpdateParamsMissing"
	ErrCodeUserIdMissing               = "userIdMissing"
	ErrCodeUserEmailPhoneNumberExists  = "userEmailPhoneNumberExists"
	ErrCodeUserEmailPhoneNumberMissing = "userEmailPhoneNumberMissing"
	ErrCodeUserEmailExists             = "userEmailExists"
	ErrCodeUserEmailMissing            = "userEmailMissing"
	ErrCodeUserEmailInvalid            = "userEmailInvalid"
	ErrCodeUserPasswordMissing         = "userPasswordMissing"
	ErrCodeUserPasswordInvalid         = "userPasswordInvalid"
	ErrCodeUserGivenNamesMissing       = "userGivenNamesMissing"
	ErrCodeUserGivenNamesInvalid       = "userGivenNamesInvalid"
	ErrCodeUserSurnameMissing          = "userSurnameMissing"
	ErrCodeUserSurnameInvalid          = "userSurnameInvalid"

	// User Session Service.
	ErrCodeUserSessionMissing             = "userSessionMissing"
	ErrCodeUserSessionNotFound            = "userSessionNotFound"
	ErrCodeUserSessionCreateParamsMissing = "userSessionCreateParamsMissing"
	ErrCodeUserSessionPurposeMissing      = "userSessionPurposeMissing"
	ErrCodeUserSessionPurposeInvalid      = "userSessionPurposeInvalid"
	ErrCodeUserSessionCredentialsInvalid  = "userSessionCredentialsInvalid"
	ErrCodeUserSessionPasswordMissing     = "userSessionPasswordMissing"
)
