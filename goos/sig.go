package goos

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WaitTerminate ждет сигнал на завершение и вызывает quitFn
// В mainCtx должен быть задан общий таймаут для shutdown!!!
func WaitTerminate(mainCtx context.Context, quitFn func(ctx context.Context)) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	if quitFn == nil {
		return
	}

	quitFn(mainCtx)
}
