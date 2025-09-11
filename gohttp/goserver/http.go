package goserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/sunshineOfficial/golib/golog"
)

const (
	listenTimeout          = 3 * time.Second
	defaultShutdownTimeout = 4 * time.Minute
)

type Server interface {
	Start()
	Stop()
	StopContext(ctx context.Context) error
	UseHandler(http.Handler)
}

type HTTPServer struct {
	log     golog.Logger
	ctx     context.Context
	server  *http.Server
	running *atomic.Bool
}

func NewHTTPServer(mainCtx context.Context, log golog.Logger, addr string) *HTTPServer {
	return &HTTPServer{
		log: log.WithTags("server"),
		ctx: mainCtx,
		server: &http.Server{
			Addr: addr,
			BaseContext: func(_ net.Listener) context.Context {
				return mainCtx
			},
		},
		running: &atomic.Bool{},
	}
}

func (h *HTTPServer) Start() {
	if h.running.Load() {
		return
	}

	h.running.Store(true)
	go h.listen()
}

func (h *HTTPServer) listen() {
	h.log.Debugf("Сервер запускается")

	for h.running.Load() {
		err := h.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Debugf("Не удалось запустить сервер: %v. Повторная попытка через %s", err, listenTimeout)
			continue
		}
	}

	h.log.Debugf("Сервер остановлен")
}

func (h *HTTPServer) Stop() {
	shutdownCtx, cancel := context.WithTimeout(h.ctx, defaultShutdownTimeout)
	defer cancel()

	if err := h.StopContext(shutdownCtx); err != nil {
		h.log.Debugf("Не удалось остановить сервер за %s: %v", defaultShutdownTimeout, err)
	}
}

func (h *HTTPServer) StopContext(ctx context.Context) error {
	if !h.running.Load() {
		return nil
	}

	h.running.Store(false)

	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("остановка сервера: %w", err)
	}

	return nil
}

func (h *HTTPServer) UseHandler(handler http.Handler) {
	h.server.Handler = handler
}
