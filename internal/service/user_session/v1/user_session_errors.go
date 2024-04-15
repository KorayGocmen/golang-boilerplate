package user_session_v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
)

var (
	ErrCreate = struct {
		UserSessionCreateParamsMissing errapi.Error
		UserSessionPurposeMissing      errapi.Error
		UserSessionPurposeInvalid      errapi.Error
		UserSessionCredentialsInvalid  errapi.Error
		UserSessionPasswordMissing     errapi.Error
		UserNotAgent                   errapi.Error
	}{
		UserSessionCreateParamsMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserSessionCreateParamsMissing,
			"Lütfen tüm eksik alanları doldur.",
			"Please fill in all missing fields.",
		),
		UserSessionPurposeMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserSessionPurposeMissing,
			"Lütfen oturum açma amacını belirtin.",
			"Please specify the purpose of the session.",
		),
		UserSessionPurposeInvalid: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserSessionPurposeInvalid,
			"Oturum açma amacı geçersiz.",
			"Session purpose is invalid.",
		),
		UserSessionCredentialsInvalid: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeUserSessionCredentialsInvalid,
			"Kullanıcı telefon numarası veya şifre geçersiz.",
			"User phone number or password is invalid.",
		),
		UserSessionPasswordMissing: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeUserSessionPasswordMissing,
			"Şifre eksik.",
			"Password is missing.",
		),
		UserNotAgent: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationInsufficientAccessLevel,
			"Kullanıcı bir temsilci değil.",
			"User is not an agent.",
		),
	}

	ErrGet = struct {
		UserSessionNotFound errapi.Error
	}{
		UserSessionNotFound: errapi.New(
			fiber.StatusNotFound,
			errapi.ErrCodeUserSessionNotFound,
			"Kullanıcı oturumu bulunamadı.",
			"Session not found.",
		),
	}

	ErrDelete = struct {
		UserIdMissing errapi.Error
	}{
		UserIdMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserIdMissing,
			"Kullanıcı kimliği eksik.",
			"User ID is missing.",
		),
	}

	ErrImpersonate = struct {
		UserIdMissing errapi.Error
		UserNotFound  errapi.Error
	}{
		UserIdMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserIdMissing,
			"Kullanıcı kimliği eksik.",
			"User ID is missing.",
		),
		UserNotFound: errapi.New(
			fiber.StatusNotFound,
			errapi.ErrCodeUserNotFound,
			"Kullanıcı bulunamadı.",
			"User not found.",
		),
	}
)
