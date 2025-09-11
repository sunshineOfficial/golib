package status

import (
	"net/http"

	"github.com/sunshineOfficial/golib/apierr"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func BadRequestHandler(c gorouter.Context) error {
	return c.WriteJson(http.StatusBadRequest,
		gorouter.NewErrorResponse(apierr.CodeUnknownError, "Произошла неизвестная ошибка"))
}

func UnauthorizedHandler(c gorouter.Context) error {
	return c.WriteJson(http.StatusUnauthorized,
		gorouter.NewErrorResponse(apierr.CodeUnauthorized, "Авторизация не выполнена"))
}

func ForbiddenHandler(c gorouter.Context) error {
	return c.WriteJson(http.StatusForbidden,
		gorouter.NewErrorResponse(apierr.CodeForbidden, "Доступ запрещен"))
}
