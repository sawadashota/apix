package apix

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type JSONReadWriter struct {
	logger *logrus.Logger
}

func NewJSONReadWriter(logger *logrus.Logger) *JSONReadWriter {
	return &JSONReadWriter{logger: logger}
}

func (rw *JSONReadWriter) ReadAll(r *http.Request, v interface{}) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

func (rw *JSONReadWriter) Write(w http.ResponseWriter, r *http.Request, v interface{}) {
	rw.WriteCode(w, r, http.StatusOK, v)
}

func (rw *JSONReadWriter) WriteCode(w http.ResponseWriter, r *http.Request, code int, v interface{}) {
	js, err := rw.encode(r, code, v)
	if err != nil {
		rw.WriteError(w, r, NewWithError(err, KindInternalServerError, SeverityCritical))
		return
	}

	if code == 0 {
		code = http.StatusOK
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(js)

	rw.log(r, code, v)
}

func (rw *JSONReadWriter) encode(r *http.Request, code int, v interface{}) (js []byte, err error) {
	switch i := v.(type) {
	case *Error:
		js, err = json.Marshal(i.Response())
	default:
		js, err = json.Marshal(v)
	}
	return
}

func (rw *JSONReadWriter) log(r *http.Request, code int, v interface{}) {
	logger := rw.logger.
		WithField("path", r.RequestURI).
		WithField("method", r.Method).
		WithField("code", code).
		WithField("user_agent", r.UserAgent()).
		WithField("remote_addr", r.RemoteAddr).
		WithField("referer", r.Referer()).
		WithField("duration", rw.duration(r).String())

	switch i := v.(type) {
	case *Error:
		logger = logger.
			WithError(i.Error()).
			WithField("operation", i.Op).
			WithField("kind", i.Kind).
			WithField("severity", i.Severity.String())
	}

	switch {
	case code >= 500:
		logger.Error()
	default:
		logger.Info()
	}
}

func (rw *JSONReadWriter) duration(r *http.Request) time.Duration {
	reqStart := r.Context().Value("time").(time.Time)
	durationNanoSecond := time.Now().UnixNano() - reqStart.UnixNano()
	return time.Duration(durationNanoSecond) * time.Nanosecond
}

func (rw *JSONReadWriter) WriteError(w http.ResponseWriter, r *http.Request, err *Error) {
	rw.WriteErrorCode(w, r, http.StatusInternalServerError, err)
}

func (rw *JSONReadWriter) WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err *Error) {
	rw.WriteCode(w, r, code, err)
}
