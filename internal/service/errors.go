package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/errapi"
)

var (
	ErrInternalServer = errapi.New(
		fiber.StatusInternalServerError,
		errapi.ErrCodeInternalServerError,
		"Beklenmedik bir hata oluştu. Lütfen daha sonra tekrar dene.",
		"An unexpected error occurred. Please try again later.",
	)
	ErrUrlParamInvalid = errapi.New(
		fiber.StatusBadRequest,
		errapi.ErrCodeUrlParamInvalid,
		"Beklenmedik bir hata oluştu. Lütfen daha sonra tekrar dene.",
		"An unexpected error occurred. Please try again later.",
	)
	ErrNotFound = errapi.New(
		fiber.StatusNotFound,
		errapi.ErrCodeNotFound,
		"Yanlış bir url adresi girdin.",
		"You entered an incorrect url address.",
	)
	ErrMethodNotAllowed = errapi.New(
		fiber.StatusMethodNotAllowed,
		errapi.ErrCodeMethodNotAllowed,
		"Yapmaya çalıştığın işlemi gerçekleştirmek için gerekli izinlere sahip değilsin.",
		"You do not have the necessary permissions to perform the operation you are trying to perform.",
	)
	ErrRequestTimeout = errapi.New(
		fiber.StatusRequestTimeout,
		errapi.ErrCodeRequestTimeout,
		"İşlem yapmak için gerekli olan süre aşıldı. Lütfen daha sonra tekrar dene.",
		"The time required to perform the operation has expired. Please try again later.",
	)
	ErrTooManyRequests = errapi.New(
		fiber.StatusTooManyRequests,
		errapi.ErrCodeTooManyRequests,
		"Çok fazla istek gönderdin. Lütfen daha sonra tekrar dene.",
		"You have sent too many requests. Please try again later.",
	)
	ErrRequestEntityTooLarge = errapi.New(
		fiber.StatusRequestEntityTooLarge,
		errapi.ErrCodeRequestEntityTooLarge,
		"Gönderilen istek boyutu fazla büyük.",
		"The request entity sent is too large.",
	)

	ErrUserAuth = struct {
		AuthorizationMissing errapi.Error
		AuthorizationInvalid errapi.Error
		AuthorizationWrong   errapi.Error
		AuthorizationExpired errapi.Error
	}{
		AuthorizationMissing: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationMissing,
			"Kullanıcı oturum kodu eksik.",
			"Authorization is missing.",
		),
		AuthorizationInvalid: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationInvalid,
			"Kullanıcı oturum kodu geçersiz.",
			"Authorization is invalid.",
		),
		AuthorizationWrong: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationWrong,
			"Kullanıcı oturum kodu hatalı.",
			"Authorization is wrong.",
		),
		AuthorizationExpired: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationExpired,
			"Kullanıcı oturum kodu süresi geçmiş.",
			"Authorization is expired.",
		),
	}

	ErrUserVerify = struct {
		AuthorizationNotVerified errapi.Error
	}{
		AuthorizationNotVerified: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationNotVerified,
			"Kullanıcı oturumu doğrulanmamış.",
			"Authorization is not verified.",
		),
	}

	ErrAgentAuth = struct {
		AuthorizationInvalidIP          errapi.Error
		AuthorizationInvalidAccessLevel errapi.Error
	}{
		AuthorizationInvalidIP: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationInvalidIP,
			"Temsilci oturum kodunun IP adresi geçersiz.",
			"Agent has invalid IP address.",
		),
		AuthorizationInvalidAccessLevel: errapi.New(
			fiber.StatusUnauthorized,
			errapi.ErrCodeAuthorizationInsufficientAccessLevel,
			"Temsilci erişim seviyesi yetersiz.",
			"Agent has invalid access level.",
		),
	}
)
