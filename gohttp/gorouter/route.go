package gorouter

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sunshineOfficial/golib/golog"
)

type Route struct {
	log   golog.Logger
	route *mux.Route

	middlewares []Middleware

	handler Handler
}

func NewRoute(log golog.Logger, route *mux.Route, handler Handler, middlewares []Middleware) *Route {
	middlewaresCopy := make([]Middleware, len(middlewares))
	copy(middlewaresCopy, middlewares)

	r := &Route{
		log:         log,
		route:       route,
		middlewares: middlewaresCopy,
		handler:     handler,
	}
	route.HandlerFunc(r.ServeHTTP)

	return r
}

func (r *Route) Use(middlewares ...Middleware) *Route {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

func (r *Route) Methods(methods ...string) *Route {
	r.route.Methods(methods...)
	return r
}

func (r *Route) Handler(h Handler) *Route {
	r.handler = h
	return r
}

func (r *Route) Router() *Router {
	return FromMux(r.log, r.route.Subrouter(), r.middlewares...)
}

func (r *Route) ServeHTTP(rs http.ResponseWriter, rq *http.Request) {
	var (
		rsProxy = NewResponseWriter(rs)
		ctx     = NewContext(r.log, rsProxy, rq, WithUserInfo())
	)
	defer func() {
		if err := ctx.Close(); err != nil {
			r.log.Debugf("Не удалось закрыть gorouter.Context: %v", err)
		}
	}()

	targetHandler := r.handler
	for i := len(r.middlewares); i > 0; i-- {
		mw := r.middlewares[i-1]
		targetHandler = mw(ctx, targetHandler)
	}

	if targetHandler == nil {
		err := ctx.WriteJson(http.StatusNotFound, NewErrorResponse("", fmt.Sprintf(`Обработчик роута "%s" не найден`, rq.RequestURI)))
		if err != nil {
			r.log.Debugf("Не удалось определить handler после применения middleware и записать ответ с ошибкой: %v", err)
		}

		return
	}

	if err := targetHandler(ctx); err != nil {
		rsProxy.WriteHeader(http.StatusBadRequest)

		return
	}
}
