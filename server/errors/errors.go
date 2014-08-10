package errors

import (
	"fmt"
)

type Enum uint8

const (
	DbNotFoundError uint8 = 1 << iota
	MongoDbError
	JwtError
)

type LytupError struct {
	msg  string
	Code uint8
}

func NewError(msg string, code uint8) *LytupError {
	return &LytupError{msg, code}
}

func (self *LytupError) Error() string {
	fmt.Println(fmt.Sprintf("[%d] %s", self.Code, self.msg))
	return fmt.Sprintf("[%d] %s", self.msg, self.Code)
}
