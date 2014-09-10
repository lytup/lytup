package lytup

import "errors"

var (
	ErrEmailIsRegistered error
	ErrEmailNotFound     error
	ErrUserNotFound      error
	ErrLogin             error
)

func init() {
	var msg = Config.Message
	ErrEmailIsRegistered = errors.New(msg.EmailIsRegisteredError)
	ErrEmailNotFound = errors.New(msg.EmailNotFoundError)
	ErrUserNotFound = errors.New(msg.UserNotFoundError)
	ErrLogin = errors.New(msg.LoginFailed)
}
