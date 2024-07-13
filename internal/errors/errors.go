package errors

import "fmt"

type APIErr struct {
	code    string
	message string
}

func (l *APIErr) Error() string {
	return fmt.Sprintf("code: %s | message: %s", l.code, l.message)
}

func NewAPIErr(code, message string) *APIErr {
	return &APIErr{
		code:    code,
		message: message,
	}
}
