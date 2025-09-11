package middleware

import (
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/status"
)

// IsAnyAuthorized
// Проверяет указана ли в запросе любая валидная авторизация. В случае провала проверки - вызывает failHandler, если он указан
func IsAnyAuthorized(failHandler gorouter.Handler) gorouter.Middleware {
	if failHandler == nil {
		failHandler = status.ForbiddenHandler
	}

	return func(c gorouter.Context, h gorouter.Handler) gorouter.Handler {
		if !c.CheckAuthorization() {
			return failHandler
		}

		return h
	}
}
