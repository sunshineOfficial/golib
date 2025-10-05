package gohttp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/golog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/atomic"
)

type Client interface {
	// SetVerbose включает логирование всех запросов и ответов (без тела)
	SetVerbose(v bool)
	Do(request *http.Request) (*http.Response, error)
	DoJson(ctx goctx.Context, method, url string, in, out any) (int, error)
	SendFormData(ctx goctx.Context, url string, fields []FormDataField, files []FormDataFile, out any) (int, error)
}

type HTTPClient struct {
	verbose    *atomic.Bool
	log        golog.Logger
	before     func(r *http.Request) error
	after      func(r *http.Response) error
	httpClient http.Client
}

func NewClient(options ...ClientOption) HTTPClient {
	var holder clientOptionHolder
	for _, opt := range options {
		holder = opt.apply(holder)
	}

	client := HTTPClient{}

	if holder.client == nil {
		client.httpClient = http.Client{}
	} else {
		client.httpClient = *holder.client
	}

	if holder.timeout > 0 {
		client.httpClient.Timeout = holder.timeout
	}

	if holder.logger != nil {
		client.log = holder.logger
	}

	if holder.transport != nil {
		client.httpClient.Transport = holder.transport
	}

	if holder.before != nil {
		client.before = holder.before
	}

	if holder.after != nil {
		client.after = holder.after
	}

	if holder.traces {
		originTransport := client.httpClient.Transport
		client.httpClient.Transport = otelhttp.NewTransport(originTransport)
	}

	client.verbose = atomic.NewBool(holder.verbose)

	return client
}

func (c HTTPClient) SetVerbose(verbose bool) {
	c.verbose.Store(verbose)
}

func (c HTTPClient) Do(httpRequest *http.Request) (*http.Response, error) {
	if c.before != nil {
		if err := c.before(httpRequest); err != nil {
			return nil, fmt.Errorf("не удалось выполнить дополнительную подготовку запроса: %v", err)
		}
	}

	oldLog := c.log
	defer func() {
		c.log = oldLog
	}()

	spanCtx := trace.SpanContextFromContext(httpRequest.Context())
	if c.log != nil {
		if spanCtx.HasTraceID() {
			c.log = c.log.WithTraceId(spanCtx.TraceID())
		}

		if spanCtx.HasSpanID() {
			c.log = c.log.WithSpanId(spanCtx.SpanID())
		}
	}

	id := fmt.Sprintf("%s %s", httpRequest.Method, httpRequest.URL)

	c.logVerbose("request", "Запрос %s запущен", id)

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		c.logVerbose("", "Выполнить запрос %s не удалось: %v", id, err)
		return nil, err
	}

	count := response.ContentLength
	if count < 0 && response.Header != nil {
		count, _ = strconv.ParseInt(response.Header.Get(ContentLengthHeader), 10, 64)
	}

	c.logVerbose("response", "Запрос %s выполнен: %s, получено %d байт", id, response.Status, count)

	if c.after != nil {
		if err = c.after(response); err != nil {
			return nil, fmt.Errorf("не удалось выполнить дополнительную обработку ответа: %v", err)
		}
	}

	return response, nil
}

func (c HTTPClient) DoJson(ctx goctx.Context, method, url string, in, out any) (int, error) {
	rq, err := NewRequest(ctx, method, url, nil)
	if err != nil {
		return 0, err
	}

	if err = WriteRequestJson(rq, in); err != nil {
		return 0, err
	}

	rs, err := c.Do(rq)
	if err != nil {
		if rs != nil && rs.Body != nil {
			closeErr := rs.Body.Close()
			err = errors.Join(err, closeErr)
		}

		return 0, err
	}

	if err = ReadResponseJson(rs, out); err != nil {
		return rs.StatusCode, err
	}

	return rs.StatusCode, nil
}

func (c HTTPClient) logVerbose(tag golog.Tag, format string, params ...interface{}) {
	if !c.verbose.Load() {
		return
	}
	if c.log == nil {
		log.Printf(format, params...)
		return
	}

	tags := make([]golog.Tag, 0, 2)
	tags = append(tags, "http")
	if len(tag) > 0 {
		tags = append(tags, tag)
	}

	c.log.DebugEntryf(format, params...).WithTags(tags...).Write()
}

func (c HTTPClient) SendFormData(ctx goctx.Context, url string, fields []FormDataField, files []FormDataFile, out any) (int, error) {
	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
	)

	for _, field := range fields {
		err := writer.WriteField(field.Name, field.Value)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("can't write files: %w", err)
		}
	}

	for _, file := range files {
		part, err := writer.CreatePart(file.MIMEHeader)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("can't create part: %w", err)
		}

		_, err = io.Copy(part, file.Payload)
		if err != nil {
			return http.StatusBadRequest, fmt.Errorf("can't copy file: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return http.StatusBadRequest, fmt.Errorf("can't close multipart-writer: %w", err)
	}

	req, err := NewRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("can't create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	rs, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("can't create request: %w; rs status code: %d", err, rs.StatusCode)
	}

	if err = ReadResponseJson(rs, out); err != nil {
		return rs.StatusCode, fmt.Errorf("can't read response: %w", err)
	}

	return rs.StatusCode, nil
}
