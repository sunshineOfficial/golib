package gorouter

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"sync/atomic"
)

var _ ResponseWriter = NewResponseWriter(nil)

// ResponseWriter
// Является оберткой над стандартным http.ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	io.Closer

	IsCommitted() bool
	Status() int
}

type ResponseWriterProxy struct {
	status    *atomic.Int32
	committed *atomic.Bool

	responseWriter http.ResponseWriter
}

func NewResponseWriter(rs http.ResponseWriter) *ResponseWriterProxy {
	return &ResponseWriterProxy{
		status:         &atomic.Int32{},
		committed:      &atomic.Bool{},
		responseWriter: rs,
	}
}

func (r *ResponseWriterProxy) Header() http.Header {
	return r.responseWriter.Header()
}

func (r *ResponseWriterProxy) Write(bytes []byte) (int, error) {
	r.writeStatus()

	return r.responseWriter.Write(bytes)
}

func (r *ResponseWriterProxy) WriteHeader(statusCode int) {
	r.status.Store(int32(statusCode))
}

func (r *ResponseWriterProxy) IsCommitted() bool {
	return r.committed.Load()
}

func (r *ResponseWriterProxy) Status() int {
	code := int(r.status.Load())
	if code == 0 {
		return http.StatusOK
	}

	return code
}

func (r *ResponseWriterProxy) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.responseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, ErrNoHijacker
	}

	return hijacker.Hijack()
}

func (r *ResponseWriterProxy) Flush() {
	flusher, ok := r.responseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

func (r *ResponseWriterProxy) Close() error {
	closer, ok := r.responseWriter.(io.Closer)
	if ok {
		defer closer.Close()
	}

	r.writeStatus()
	return nil
}

func (r *ResponseWriterProxy) writeStatus() {
	if r.committed.Load() {
		return
	}

	if status := int(r.status.Load()); status > 0 {
		r.committed.Store(true)
		r.responseWriter.WriteHeader(status)
	}
}
