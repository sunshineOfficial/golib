package middleware

import (
	"net/http"

	"github.com/sunshineOfficial/golib/goerr"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func WithServiceErrors(status int) gorouter.Middleware {
	if status == 0 {
		status = http.StatusBadRequest
	}

	return func(_ gorouter.Context, h gorouter.Handler) gorouter.Handler {
		return func(c gorouter.Context) error {
			err := h(c)

			serviceErr := goerr.AsServiceError(err)
			if serviceErr == nil {
				return err
			}

			code, message := serviceErr.Info()

			writeErr := c.WriteJson(status, gorouter.ErrorResponse{
				Error: gorouter.ErrorInfo{
					Code:    code,
					Message: message,
				},
			})
			if writeErr != nil {
				c.Log().Debugf("не удалось записать ошибку в ответ: %v", writeErr)
			}

			return serviceErr.Internal()
		}
	}
}
