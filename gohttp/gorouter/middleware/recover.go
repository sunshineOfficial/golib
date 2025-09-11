package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/sunshineOfficial/golib/apierr"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func Recover(_ gorouter.Context, h gorouter.Handler) gorouter.Handler {
	return func(c gorouter.Context) error {
		var (
			rq  = c.Request()
			log = c.Log()
			id  = fmt.Sprintf("%s %s", rq.Method, rq.URL.Path)
		)

		defer func() {
			if panicErr := recover(); panicErr != nil {
				log.Errorf("Паника при выполнении %s: %+v\n%s", id, panicErr, debug.Stack())
				_ = c.WriteJson(http.StatusInternalServerError,
					gorouter.NewErrorResponse(apierr.CodePanic, "Something went wrong, we have a panic!"))
			}
		}()

		return h(c)
	}
}
