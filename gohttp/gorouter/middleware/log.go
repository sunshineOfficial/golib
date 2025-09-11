package middleware

import (
	"time"

	"github.com/sunshineOfficial/golib/goerr"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func LogError(_ gorouter.Context, h gorouter.Handler) gorouter.Handler {
	return func(c gorouter.Context) error {
		var (
			log = c.Log()
			err = h(c)
		)

		switch {
		case err == nil:
			return nil

		case goerr.IsCriticalError(err):
			log.Error(err.Error())

		default:
			log.Debug(err.Error())
		}

		return err
	}
}

func VerboseLog(_ gorouter.Context, h gorouter.Handler) gorouter.Handler {
	return func(c gorouter.Context) error {
		start := time.Now()

		defer func() {
			var (
				rq    = c.Request()
				rs    = c.Response()
				log   = c.Log()
				total = time.Since(start)
			)

			if rs.IsCommitted() {
				log.Debugf("Запрос %s %s выполнен за %s", rq.Method, rq.URL.Path, total)
			} else {
				log.Debugf("Запрос %s %s (статус %d) выполнен за %s", rq.Method, rq.URL.Path, rs.Status(), total)
			}
		}()

		return h(c)
	}
}
