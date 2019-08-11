package apix

import (
	"net/http"
)

type ReadWriter interface {
	Reader
	Writer
}

type Reader interface {
	ReadAll(r *http.Request, v interface{}) error
}

type Writer interface {
	// Write a response object to the ResponseWriter with status code 200.
	Write(w http.ResponseWriter, r *http.Request, v interface{})

	// WriteCode writes a response object to the ResponseWriter and sets a response code.
	WriteCode(w http.ResponseWriter, r *http.Request, code int, v interface{})

	// WriteError writes an error to ResponseWriter and tries to extract the error's status code by
	// asserting statusCodeCarrier. If the error does not implement statusCodeCarrier, the status code
	// is set to 500.
	WriteError(w http.ResponseWriter, r *http.Request, err *Error)

	// WriteErrorCode writes an error to ResponseWriter and forces an error code.
	WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err *Error)
}
