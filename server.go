package apix

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// TODO: make all these configurable
const (
	port       = 8080
	driverName = "mysql"
	dsn        = "docker:docker@tcp(mysql:3306)/apix"
	logLevel   = logrus.DebugLevel
)

type Server struct {
	registry *Registry
	router   *Router
	rw       ReadWriter
	port     uint
}

func NewDefaultServer() (*Server, error) {
	registry, err := newRegistry()
	if err != nil {
		return nil, err
	}

	s := &Server{
		registry: registry,
		rw:       NewJSONReadWriter(registry.Logger()),
		port:     port,
	}

	var notFoundHandler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		err := New("404 Not Found", KindNotFound, SeverityInfo)
		err.SetOp("notFoundHandler")
		s.Writer().WriteErrorCode(w, req, http.StatusNotFound, err)
	}

	var panicHandler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
		err := New("Internal Server Error", KindInternalServerError, SeveritySevere)
		err.SetOp("panicHandler")
		s.Writer().WriteError(w, req, err)
	}

	s.router = NewRouter(WithNotFound(notFoundHandler), WithPanicHandler(panicHandler))

	return s, nil
}

func (s *Server) Writer() Writer {
	return s.rw
}

func (s *Server) Reader() Reader {
	return s.rw
}

func (s *Server) Port() uint {
	return port
}

func (s *Server) Registry() *Registry {
	return s.registry
}

func (s *Server) Router() *Router {
	return s.router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router().ServeHTTP(w, req)
}

func (s *Server) Serve() error {
	s.Registry().Logger().Infof("Starting admin server on 127.0.0.1:%d", s.Port())
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port()), s)
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logLevel)
	return logger
}
