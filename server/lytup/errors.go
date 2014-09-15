package lytup

import "errors"

var (
	// ErrBlankFirstName    ErrUserField
	ErrInvalidFirstName ErrUserField
	// ErrBlankLastName     ErrUserField
	ErrInvalidLastName ErrUserField
	// ErrBlankEmail        ErrUserField
	ErrInvalidEmail ErrUserField
	// ErrBlankPassword     ErrUserField
	ErrInvalidPassword   ErrUserField
	ErrEmailIsRegistered error
	ErrVerifyEmail       error
	ErrEmailIsVerified   error
	ErrEmailNotFound     error
	ErrResetPassword     error
	ErrUserNotFound      error
	ErrLogin             error
)

type (
	ErrUserField struct {
		error
	}

	HttpError struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}
)

func NewHttpError(code int, err string) *HttpError {
	return &HttpError{
		code,
		err,
	}
}

func init() {
	// ErrBlankFirstName = ErrUserField{errors.New(M["blankFirstName"])}
	ErrInvalidFirstName = ErrUserField{errors.New(M["invalidFirstName"])}
	// ErrBlankLastName = ErrUserField{errors.New(M["blankLastName"])}
	ErrInvalidLastName = ErrUserField{errors.New(M["invalidLastName"])}
	// ErrBlankEmail = ErrUserField{errors.New(M["blankEmail"])}
	ErrInvalidEmail = ErrUserField{errors.New(M["invalidEmail"])}
	// ErrBlankPassword = ErrUserField{errors.New(M["blankPassword"])}
	ErrInvalidPassword = ErrUserField{errors.New(M["invalidPassword"])}
	ErrEmailIsRegistered = errors.New(M["emailIsRegisteredError"])
	ErrVerifyEmail = errors.New(M["verifyEmailFailed"])
	ErrEmailIsVerified = errors.New(M["emailIsVerifiedError"])
	ErrEmailNotFound = errors.New(M["emailNotFoundError"])
	ErrResetPassword = errors.New(M["resetPasswordFailed"])
	ErrUserNotFound = errors.New(M["userNotFoundError"])
	ErrLogin = errors.New(M["loginFailed"])
}
