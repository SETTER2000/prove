package er

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrNotFound ошибка в случаи отсутствия данных
// ErrAlreadyExists ошибка в случаи если данные уже существуют
// ErrBadRequest ошибка в случаи не корректного запроса
// ErrAccessDenied ошибка в случаи отсутствия права доступа
// ErrForbidden ошибка доступа к запрошенному ресурсу в случаи когда доступ запрещен
// ErrUserIDRequired -.
// ErrStatusConflict -.
var (
	ErrNotFound                 = errors.New("not found")
	ErrAlreadyExists            = errors.New("already exists")
	ErrBadRequest               = errors.New("bad request")
	ErrAccessDenied             = errors.New(`access denied`)
	ErrForbidden                = errors.New(`forbidden`)
	ErrUserIDRequired           = errors.New("user id required")
	ErrStatusConflict           = errors.New("status conflict")
	ErrStatusGone               = errors.New("status gone")
	ErrValidCard                = errors.New("not a valid card")
	ErrInsufficientFundsAccount = errors.New("insufficient funds in the account")
	ErrAlreadyBeenUploaded      = errors.New("the order number has already been uploaded by another user")
	ErrNotDataAnswer            = errors.New("no data to answer")
	ErrBadFormat                = errors.New("invalid request format")
	ErrBadFormatOrder           = errors.New("invalid order number format")
	ErrIncorrectLoginOrPass     = errors.New("incorrect login or password")
	ErrIncorrectToken           = errors.New("incorrect token, the token is rotten")
)

type response struct {
	Error string `json:"error" example:"message"`
}

// TimeError предназначен для ошибок с фиксацией времени возникновения.
type TimeError struct {
	Time time.Time
	Err  error
}

// Error добавляет поддержку интерфейса error для типа TimeError.
func (te *TimeError) Error() string {
	return fmt.Sprintf("%v %v", te.Time.Format(`2006/01/02 15:04:05`), te.Err)
}

// NewTimeError упаковывает ошибку err в тип TimeError c текущим временем.
func NewTimeError(err error) error {
	return &TimeError{
		Time: time.Now(),
		Err:  err,
	}
}
func (te *TimeError) Unwrap() error {
	return te.Err
}

func (te *TimeError) Is(err error) bool {
	return te.Err == err
}

// LabelError описывает ошибку с дополнительной меткой.
type LabelError struct {
	Err   error
	Label string
}

// NewLabelError упаковывает ошибку err в тип LabelError.
func NewLabelError(label string, err error) error {
	return &LabelError{
		Label: strings.ToUpper(label),
		Err:   err,
	}
}

// Error добавляет поддержку интерфейса error для типа LabelError.
func (le *LabelError) Error() string {
	return fmt.Sprintf("[%s] %v", le.Label, le.Err)
	//return fmt.Sprintf("[%s] %v", le.Label, le.Err)
}

func (le *LabelError) Unwrap() error {
	return le.Err
}

func (le *LabelError) Is(err error) bool {
	return le.Err == err
}

// ConflictError описывает ошибку с дополнительной меткой и значением.
type ConflictError struct {
	Err   error
	Label string
	URL   string
}

// NewConflictError упаковывает ошибку err в тип LabelError.
func NewConflictError(label string, url string, err error) error {
	return &ConflictError{
		Label: strings.ToUpper(label),
		URL:   url,
		Err:   err,
	}
}

// Error добавляет поддержку интерфейса error для типа LabelError.
func (ce *ConflictError) Error() string {
	return fmt.Sprintf("%v", ce.Err)
	//return fmt.Sprintf("[%s] %s %v", ce.Label, ce.URL, ce.Err)
}

func (ce *ConflictError) Unwrap() error {
	return ce.Err
}

func (ce *ConflictError) Is(err error) bool {
	return ce.Err == err
}

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func NewAppError(message, code, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}
