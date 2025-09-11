package gorouter

import (
	"fmt"
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp"

	"github.com/gorilla/mux"
	"github.com/sunshineOfficial/golib/golog"
)

type Router struct {
	router      *mux.Router
	log         golog.Logger
	middlewares []Middleware
	basePath    string
}

func NewRouter(log golog.Logger) *Router {
	return FromMux(log, nil)
}

func FromMux(log golog.Logger, router *mux.Router, middlewares ...Middleware) *Router {
	if router == nil {
		router = mux.NewRouter()
	}

	copiedMiddlewares := make([]Middleware, len(middlewares))
	copy(copiedMiddlewares, middlewares)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = gohttp.WriteResponseJson(w, http.StatusNotFound,
			NewErrorResponse("", fmt.Sprintf(`Роут "%s" не найден`, r.RequestURI)))
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = gohttp.WriteResponseJson(w, http.StatusMethodNotAllowed,
			NewErrorResponse("", fmt.Sprintf(`Роут "%s" не поддерживает метод "%s"`, r.RequestURI, r.Method)))
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }).Methods(http.MethodGet)

	return &Router{
		router:      router,
		log:         log,
		middlewares: copiedMiddlewares,
	}
}

// Use Устанавливает Middleware для всех роутов и под-роутеров. Имеющиеся Middleware будут затерты
// ВАЖНО! Должен быть вызван до добавления роутов и под-роутеров
func (r *Router) Use(middlewares ...Middleware) *Router {
	r.middlewares = middlewares
	return r
}

// Add Добавляет Middleware для всех роутов и под-роутеров к имеющимся в конец списка
// ВАЖНО! Должен быть вызван до добавления роутов и под-роутеров
func (r *Router) Add(middlewares ...Middleware) *Router {
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

func (r *Router) Install(plugins ...Plugin) *Router {
	for _, plugin := range plugins {
		subRouter := r.SubRouter(plugin.BasePath())

		plugin.Register(subRouter)
	}

	return r
}

func (r *Router) Handle(path string, handler Handler, methods ...string) *Route {
	route := NewRoute(r.log, r.router.Handle(r.basePath+path, nil), handler, r.middlewares)
	if len(methods) > 0 {
		route.Methods(methods...)
	}

	return route
}

func (r *Router) HandleGet(path string, handler Handler) *Route {
	return r.Handle(path, handler, http.MethodGet)
}

func (r *Router) HandlePost(path string, handler Handler) *Route {
	return r.Handle(path, handler, http.MethodPost)
}

func (r *Router) HandlePut(path string, handler Handler) *Route {
	return r.Handle(path, handler, http.MethodPut)
}

func (r *Router) HandlePatch(path string, handler Handler) *Route {
	return r.Handle(path, handler, http.MethodPatch)
}

func (r *Router) HandleDelete(path string, handler Handler) *Route {
	return r.Handle(path, handler, http.MethodDelete)
}

func (r *Router) Path(path string) *Route {
	return NewRoute(r.log, r.router.Path(r.basePath+path), nil, r.middlewares)
}

func (r *Router) PathPrefix(path string) *Route {
	return NewRoute(r.log, r.router.PathPrefix(r.basePath+path), nil, r.middlewares)
}

func (r *Router) SubRouter(path string) *Router {
	subRouter := FromMux(r.log, r.router, r.middlewares...)
	subRouter.setBasePath(r.basePath + path)
	return subRouter
}

func (r *Router) HandleHTTP(path string, handler http.Handler, methods ...string) {
	route := r.router.Handle(r.basePath+path, handler)
	if len(methods) > 0 {
		route.Methods(methods...)
	}
}

func (r *Router) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	r.router.ServeHTTP(rs, rq)
}

func (r *Router) setBasePath(basePath string) {
	r.basePath = basePath
}
