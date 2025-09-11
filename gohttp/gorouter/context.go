package gorouter

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gohttp"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/reflect"
	"github.com/sunshineOfficial/golib/golog"
	"go.opentelemetry.io/otel/trace"
)

type Context struct {
	response ResponseWriter
	request  *http.Request

	ctx goctx.Context
	log golog.Logger

	tracer trace.Tracer
}

func NewContext(log golog.Logger, rs ResponseWriter, rq *http.Request, options ...ContextOption) Context {
	var holder contextOptionHolder
	for _, option := range options {
		holder = option.apply(holder)
	}

	ctx, _ := goctx.GetContext(rq)

	if holder.userInfo {
		log = log.WithUserInfo(ctx.Authorize.UserId)
	}

	if holder.traces {
		log = log.WithTraceId(ctx.TraceId)
	}

	return Context{
		response: rs,
		request:  rq,
		ctx:      ctx,
		log:      log,
	}
}

func (c Context) Log() golog.Logger {
	return c.log
}

func (c Context) Ctx() goctx.Context {
	return c.ctx
}

// CheckAuthorization
// Проверяет авторизацию по указанным правилам:
// - если авторизация не указана или невалиден токен - записывает ответ 401 и возвращает false
// в остальных случаях - не записывает ничего и возвращает true
func (c Context) CheckAuthorization() bool {
	if !c.ctx.IsAuthorized() {
		c.Write(http.StatusUnauthorized)
		return false
	}

	return true
}

func (c Context) Vars(a any) error {
	header := c.request.Header
	query := c.request.URL.Query()

	if err := reflect.SetValuesToItem(getPathVars(c.request), "path", a); err != nil {
		return err
	}
	if err := reflect.SetValuesToItem(header, "header", a); err != nil {
		return err
	}
	if err := reflect.SetValuesToItem(query, "query", a); err != nil {
		return err
	}

	return nil
}

func (c Context) ReadJson(a any) error {
	return json.NewDecoder(c.request.Body).Decode(a)
}

func (c Context) ReadText() (string, error) {
	b, err := io.ReadAll(c.request.Body)
	if err == nil {
		return "", err
	}

	return string(b), nil
}

func (c Context) ReadBytes() ([]byte, error) {
	return io.ReadAll(c.request.Body)
}

func (c Context) Reader() io.ReadCloser {
	return c.request.Body
}

func (c Context) Write(code int) {
	c.response.WriteHeader(code)
}

func (c Context) WriteJson(code int, a any) error {
	return gohttp.WriteResponseJson(c.response, code, a)
}

func (c Context) WriteXml(code int, a any) error {
	return gohttp.WriteResponseXml(c.response, code, a)
}

func (c Context) WriteText(code int, s string) error {
	c.response.WriteHeader(code)
	_, err := io.WriteString(c.response, s)
	return err
}

func (c Context) WriteBinary(code int, b []byte) error {
	c.response.WriteHeader(code)
	_, err := c.response.Write(b)
	return err
}

func (c Context) WriteHeader(name, value string) {
	c.response.Header().Set(name, value)
}

func (c Context) Writer() io.WriteCloser {
	return c.response
}

func (c Context) Request() *http.Request {
	return c.request
}

func (c Context) Response() ResponseWriter {
	return c.response
}

// Tracer возвращает trace.Tracer, ассоциированный с gorouter
// ВАЖНО! Не стоит использовать этот трейсер для трассировки бизнес-логики, т. к. он ассоциирован именно с роутингом.
// Для бизнес-логики используйте трейсеры, созданные в app.go
func (c Context) Tracer() trace.Tracer {
	return c.tracer
}

func (c Context) Close() error {
	if c.request.Body != nil {
		_ = c.request.Body.Close()
	}

	return c.response.Close()
}

func getPathVars(r *http.Request) map[string][]string {
	vars := mux.Vars(r)
	if len(vars) == 0 {
		return nil
	}

	result := make(map[string][]string, len(vars))
	for key, value := range vars {
		result[key] = []string{value}
	}

	return result
}

func (c Context) PathTemplate() string {
	route := mux.CurrentRoute(c.request)
	if route == nil {
		return ""
	}

	template, err := route.GetPathTemplate()
	if err != nil {
		return ""
	}

	return template
}
