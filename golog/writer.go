package golog

import (
	"io"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type writeSyncer struct {
	w  io.Writer
	mx sync.RWMutex
}

func newWriteSyncer(w io.Writer) zapcore.WriteSyncer {
	return &writeSyncer{
		w: w,
	}
}

func (s *writeSyncer) Write(p []byte) (int, error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	return s.w.Write(p)
}

func (s *writeSyncer) Sync() error {
	return nil
}

var (
	_levelEnabler = zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == zapcore.DebugLevel || l == zapcore.ErrorLevel
	})
	_stacktraceEnabler = zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l == zapcore.ErrorLevel
	})
)
