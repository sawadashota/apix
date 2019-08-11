package apix

import (
	"errors"
	"net/http"
)

type Kind int

const (
	KindNotFound            = http.StatusNotFound
	KindUnauthorized        = http.StatusUnauthorized
	KindBadRequest          = http.StatusBadRequest
	KindRequestTimeout      = http.StatusRequestTimeout
	KindInternalServerError = http.StatusInternalServerError
)

type Op string

type Severity int

const (
	SeverityDebug = iota
	SeverityInfo
	SeverityWarn
	SeverityError
	SeverityCritical
	SeveritySevere
)

func (s *Severity) String() string {
	switch int(*s) {
	case SeverityDebug:
		return "debug"
	case SeverityInfo:
		return "info"
	case SeverityWarn:
		return "warn"
	case SeverityError:
		return "error"
	case SeverityCritical:
		return "critical"
	case SeveritySevere:
		return "severe"
	default:
		return "unknown"
	}
}

type Error struct {
	Op       Op
	Kind     Kind
	Severity Severity
	Err      error
}

func New(err string, k Kind, s Severity) *Error {
	return &Error{
		Kind:     k,
		Severity: s,
		Err:      errors.New(err),
	}
}

func NewWithError(err error, k Kind, s Severity) *Error {
	return &Error{
		Kind:     k,
		Severity: s,
		Err:      err,
	}
}

func (e *Error) SetOp(op Op) {
	e.Op = op
}

func (e *Error) String() string {
	return e.Err.Error()
}

func (e *Error) Error() error {
	return e.Err
}

type Response struct {
	Error string `json:"error"`
}

func (e *Error) Response() *Response {
	return &Response{
		Error: e.String(),
	}
}
