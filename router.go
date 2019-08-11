package apix

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Router struct {
	HandleOPTIONS bool
	NotFound      http.HandlerFunc
	PanicHandler  http.HandlerFunc
	trees         *trees
}

type Option func(*Router)

func WithNotFound(handle http.HandlerFunc) Option {
	return func(router *Router) {
		router.NotFound = handle
	}
}

func WithPanicHandler(handle http.HandlerFunc) Option {
	return func(router *Router) {
		router.PanicHandler = handle
	}
}

func NewRouter(opts ...Option) *Router {
	r := &Router{
		HandleOPTIONS: true,
	}

	for _, opt := range opts {
		opt(r)
	}
	return r
}

type trees struct {
	routes map[string]*routes
}

func (t *trees) get(method string) (*routes, error) {
	routes := t.routes[method]
	if routes == nil {
		return nil, errors.Errorf("method %s is not found", method)
	}

	return routes, nil
}

func (t *trees) add(method string, routes *routes) {
	t.routes[method] = routes
}

func (t *trees) allowedMethods(path, reqMethod string) string {
	allowed := make([]string, 0, len(t.routes))
	for method := range t.routes {
		if method == http.MethodOptions {
			continue
		}

		allowed = append(allowed, method)
	}

	if len(allowed) > 0 {
		allowed = append(allowed, http.MethodOptions)
	}

	return strings.Join(allowed, ", ")
}

type routes struct {
	handlers map[string]http.HandlerFunc
}

func (r *routes) get(path string) http.HandlerFunc {
	return r.handlers[path]
}

func (r *Router) GET(path string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, path, handler)
}

func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.trees == nil {
		r.trees = &trees{
			routes: make(map[string]*routes),
		}
	}

	rt, err := r.trees.get(method)
	if err != nil {
		rt = &routes{
			handlers: make(map[string]http.HandlerFunc),
		}
		r.trees.add(method, rt)
	}

	rt.addRoute(path, handler)
}

func (r *routes) addRoute(path string, handler http.HandlerFunc) {
	r.handlers[path] = handler
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if r.PanicHandler != nil {
		defer r.recover(w, req)
	}

	path := req.URL.Path

	ctx := context.WithValue(context.Background(), "time", time.Now())
	req = req.WithContext(ctx)

	if routes, err := r.trees.get(req.Method); err == nil {
		if handle := routes.get(path); handle != nil {
			handle(w, req)
			return
		}
	}

	if req.Method == http.MethodOptions && r.HandleOPTIONS {
		if allowed := r.trees.allowedMethods(path, req.Method); len(allowed) > 0 {
			w.Header().Set("Allow", allowed)
			return
		}
	}

	r.handleNotFound(w, req)
}

func (r *Router) handleNotFound(w http.ResponseWriter, req *http.Request) {
	if r.NotFound != nil {
		r.NotFound.ServeHTTP(w, req)
		return
	}

	http.NotFound(w, req)
}

func (r *Router) recover(w http.ResponseWriter, req *http.Request) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(w, req)
	}
}

type Handler interface {
	SetRoutes(r *Router)
}

func (r *Router) Register(hs ...Handler) {
	for _, h := range hs {
		h.SetRoutes(r)
	}
}
