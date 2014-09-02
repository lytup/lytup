package lytup

import (
	"fmt"
)

// http://golang.org/src/pkg/net/http/status.go

const (
	MongoDbError uint8 = 1 << iota
	MongoDbNotFoundError
	JwtError
	LoginError
)

var codeText = []string{
	MongoDbError:         "MongoDbError",
	MongoDbNotFoundError: "MongoDbNotFoundError",
	JwtError:             "JwtError",
	LoginError:           "LoginError",
}

type Error struct {
	msg  string
	Code uint8
}

func NewError(msg string, code uint8) *Error {
	return &Error{msg, code}
}

func (le *Error) Error() string {
	return fmt.Sprintf("[%s] %s", codeText[le.Code], le.msg)
}
