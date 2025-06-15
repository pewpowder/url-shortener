package errors

import (
	stderrors "errors"
	"runtime"
)

type ErrorType uint

const (
	// ErrorTypeUnknown - неизвестная ошибка (500)
	ErrTypeInternal ErrorType = iota
	// ErrorTypeValidation - ошибка валидации (400)
	ErrTypeValidation
	// ErrorTypeNotFound - ресурс не найден (404)
	ErrTypeNotFound
	// ErrorTypeConflict - конфликт ресурсов (409)
	ErrTypeConflict
	// ErrorTypeUnauthorized - не авторизован (401)
	ErrTypeUnauthorized
	// ErrorTypeForbidden - доступ запрещен (403)
	ErrTypeForbidden
	// ErrorTypeExpired - ресурс истек (410)
	ErrTypeExpired
)

type ServiceError struct {
	Message string
	Type    ErrorType
	File    string
	Line    int
	Err     error
}

func NewServiceError(message string, errorType ErrorType, err error) error {
	_, file, line, _ := runtime.Caller(1)

	return &ServiceError{
		Message: message,
		Type:    errorType,
		Err:     err,
		File:    file,
		Line:    line,
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

// GormExample
// return nil, &ServiceError{
// 	Message: "some message",
// 	Type: getTypeByGormError(err), // map gorm error to my ErrorType
// 	InitialError: err,
// 	File: "file.go",
// 	Line: 10,
// }

// General error
// return nil, &ServiceError{
// 	Message: "some message",
// 	Type: InternalServerError // I exactly know what error occurred
// 	InitialError: err,
// 	File: "file.go",
// 	Line: 10,
// }

// Custom errors flow
// 1. Ошибка происходит на любом из уровней, я создаю новую структуру ServiceError
// и наполняю ее данными в месте создания (тк в этом месте есть все необходимые данные для идентификации ошибки)
// 2. Возвращаю ошибку выше по уровню (через error в объявлении функций)
// 3. На каждом из уровней могу дополнительно лишь обогащать ошибку через fmt.Errorf("%w: add info")
// 4. В контроллере вызываю HandleError метод (ошибка передается не связанная с протоколом HTTP, GRPC и тд)
// 5. В методе проверяю ошибку, что это моя кастомная ошибка и достаю всю необходимую инфу оттуда
// 6. Статус код мапится по Type в теле ошибки
// 7. Отправляется ответ клиенту

// func NewServiceError(errType ErrorType, message string, initialErr error) *ServiceError {
//     _, file, line, _ := runtime.Caller(1)
//     return &ServiceError{
//         Message:      message,
//         Type:         errType,
//         InitialError: initialErr,
//         File:         file,
//         Line:         line,
//     }
// }

// type AppError interface {
//     error
//     ErrorType() ErrorType
//     Unwrap() error
//     File() string
//     Line() int
//     Code() string
// }

// Sentinel errors flow
// 1. Ошибка происходит на любом из уровней, я просто возвращаю одну из предсозданных ошибок
// 2. Если ошибка пришла из вне (к примеру, gorm) мне нужно мапить уже на этом уровне ошибку gorm
// к объявленным внутри приложения ошибкам
// 3. Так же на каждом уровне могу дополнительно обогащать ошибку через fmt.Errorf("%w: add info")
// 4. остальное все так же как и в кастомных ошибках за исключением, что статус я достаю не из тела
// а маплю сервисную ошибку на нужный статус код руками
