package user_v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
)

var (
	ErrCreate = struct {
		UserCreateParamsMissing     errapi.Error
		UserEmailPhoneNumberMissing errapi.Error
		UserEmailExists             errapi.Error
	}{
		UserCreateParamsMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserCreateParamsMissing,
			"Lütfen tüm eksik alanları doldur.",
			"Please fill in all missing fields.",
		),
		UserEmailPhoneNumberMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserEmailPhoneNumberMissing,
			"E-posta adresi ve telefon numarası eksik.",
			"Email address and phone number are missing.",
		),
		UserEmailExists: errapi.New(
			fiber.StatusConflict,
			errapi.ErrCodeUserEmailExists,
			"Bu e-posta adresi zaten kayıtlı.",
			"This email address is already registered.",
		),
	}

	ErrUpdate = struct {
		UserMissing             errapi.Error
		UserUpdateParamsMissing errapi.Error
		UserEmailExists         errapi.Error
	}{
		UserMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserMissing,
			"Kullanıcı bulunamadı.",
			"User not found.",
		),
		UserUpdateParamsMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserUpdateParamsMissing,
			"Lütfen tüm eksik alanları doldur.",
			"Please fill in all missing fields.",
		),
		UserEmailExists: errapi.New(
			fiber.StatusConflict,
			errapi.ErrCodeUserEmailExists,
			"Bu e-posta adresi zaten kayıtlı.",
			"This email address is already registered.",
		),
	}

	ErrGet = struct {
		UserIdMissing errapi.Error
		UserNotFound  errapi.Error
	}{
		UserIdMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserIdMissing,
			"Kullanıcı kimligi eksik.",
			"User id missing.",
		),
		UserNotFound: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserNotFound,
			"Kullanıcı bulunamadı.",
			"User not found.",
		),
	}

	ErrDelete = struct {
		UserMissing errapi.Error
	}{
		UserMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserMissing,
			"Kullanıcı bulunamadı.",
			"User not found.",
		),
	}

	ErrGetByEmailOrPhoneNumber = struct {
		UserEmailPhoneNumberMissing errapi.Error
		UserNotFound                errapi.Error
	}{
		UserEmailPhoneNumberMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserEmailPhoneNumberMissing,
			"Lütfen e-posta adresinizi veya telefon numaranızı girin.",
			"Please enter your email address or phone number.",
		),
		UserNotFound: errapi.New(
			fiber.StatusNotFound,
			errapi.ErrCodeUserNotFound,
			"Kullanıcı bulunamadı.",
			"User not found.",
		),
	}

	ErrValidateParams = struct {
		UserEmailMissing      errapi.Error
		UserEmailInvalid      errapi.Error
		UserPasswordMissing   errapi.Error
		UserPasswordInvalid   errapi.Error
		UserGivenNamesMissing errapi.Error
		UserGivenNamesInvalid errapi.Error
		UserSurnameMissing    errapi.Error
		UserSurnameInvalid    errapi.Error
	}{
		UserEmailMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserEmailMissing,
			"Kullanıcı e-posta adresi eksik.",
			"User email address is missing.",
		),
		UserEmailInvalid: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserEmailInvalid,
			"Kullanıcı e-posta adresi geçersiz.",
			"User email address is invalid.",
		),
		UserPasswordMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserPasswordMissing,
			"Kullanıcı şifresi eksik.",
			"User password is missing.",
		),
		UserPasswordInvalid: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserPasswordInvalid,
			"Kullanıcı şifresi 6 karakter uzunluğunda olmalı.",
			"User password must be 6 characters long.",
		),
		UserGivenNamesMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserGivenNamesMissing,
			"Kullanıcı ismi eksik.",
			"User given names is missing.",
		),
		UserGivenNamesInvalid: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserGivenNamesInvalid,
			"Kullanıcı ismi geçersiz.",
			"User given names is invalid.",
		),
		UserSurnameMissing: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserSurnameMissing,
			"Kullanıcı soyismi eksik.",
			"User surname is missing.",
		),
		UserSurnameInvalid: errapi.New(
			fiber.StatusBadRequest,
			errapi.ErrCodeUserSurnameInvalid,
			"Kullanıcı soyismi geçersiz.",
			"User surname is invalid.",
		),
	}
)
