package errors

import (
	stderrors "errors"
	"runtime"
)

type ErrorCode uint

const (
	// ErrorTypeUnknown - unknown error (500)
	ErrCodeInternal ErrorCode = iota
	// ErrCodeBadRequest - request is invalid (400)
	ErrCodeBadRequest
	// ErrorTypeNotFound - resource not found (404)
	ErrCodeNotFound
	// ErrorTypeConflict - resource already exists 	(409)
	ErrCodeConflict
	// ErrCodeUnauthorized - unauthorized (401)
	ErrCodeUnauthorized
	// ErrCodeForbidden - forbidden (403)
	ErrCodeForbidden
	// ErrorTypeExpired - expired (410)
	ErrCodeExpired
)

const (
	ErrInternal       = "INTERNAL_ERROR"
	ErrInvalidRequest = "INVALID_REQUEST"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrGone           = "GONE"
	ErrNotFound       = "NOT_FOUND"
	ErrConflict       = "CONFLICT"
)

type ServiceError struct {
	Message  string
	Code     ErrorCode
	CodeText string
	Err      error
	File     string
	Line     int
}

func NewServiceError(message string, code ErrorCode, codeText string, err error) error {
	_, file, line, _ := runtime.Caller(1)

	return &ServiceError{
		Message:  message,
		Code:     code,
		CodeText: codeText,
		Err:      err,
		File:     file,
		Line:     line,
	}
}

func (se *ServiceError) Error() string {
	return se.Message
}

func (se *ServiceError) Unwrap() error {
	return se.Err
}

func IsAnyOf(err error, errors ...error) bool {
	for _, e := range errors {
		if stderrors.Is(err, e) {
			return true
		}
	}
	return false
}

func GetCodeTextByCode(code ErrorCode) string {
	switch code {
	case ErrCodeBadRequest:
		return ErrInvalidRequest
	case ErrCodeNotFound:
		return ErrNotFound
	case ErrCodeConflict:
		return ErrConflict
	case ErrCodeUnauthorized:
		return ErrUnauthorized
	case ErrCodeForbidden:
		return ErrForbidden
	case ErrCodeExpired:
		return ErrGone
	default:
		return ErrInternal
	}
}
