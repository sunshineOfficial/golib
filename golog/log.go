package golog

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	WithSkip(s int) Logger

	// WithUserInfo Возвращает Logger с сохраненным полем userId.
	// Эти поля будут попадать в каждую запись, сделанную полученным логгером
	WithUserInfo(userId int) Logger

	// WithTraceId Возвращает Logger с сохраненным полем traceId.
	WithTraceId(traceId trace.TraceID) Logger

	// WithBookId Возвращает Logger с сохраненным полем bookId.
	WithBookId(bookId string) Logger

	// WithTags Возвращает Logger c новыми тегами, имеющиеся теги из оригинального Logger будут скопированы в новый
	WithTags(tags ...Tag) Logger

	// Message возвращает LogEntry ассоциированный с указанным уровнем
	Message(level MessageLevel, message string) LogEntry

	// DebugEntry возвращает LogEntry ассоциированный с указанным сообщением и уровнем лога LevelDebug
	DebugEntry(message string) LogEntry
	// DebugEntryf форматирует текст и возвращает LogEntry ассоциированный с указанным сообщением и уровнем лога LevelDebug
	DebugEntryf(format string, args ...interface{}) LogEntry

	// Debug записывает сообщение в лог в момент вызова
	Debug(message string)
	// Debugf форматирует сообщение и записывает его в момент вызова
	Debugf(format string, args ...interface{})

	// ErrorEntry возвращает LogEntry ассоциированный с указанным сообщением и уровнем лога LevelError
	ErrorEntry(message string) LogEntry
	// ErrorEntryf форматирует текст и возвращает LogEntry ассоциированный с указанным сообщением и уровнем лога LevelError
	ErrorEntryf(format string, args ...interface{}) LogEntry

	// Error записывает сообщение в лог в момент вызова
	Error(message string)
	// Errorf форматирует сообщение и записывает его в момент вызова
	Errorf(format string, args ...interface{})
}

var (
	_global Logger
)

func Global() Logger {
	if _global == nil {
		_global = NewLogger("default")
	}

	return _global
}

// SwLogger реализация Logger основанная на zap.Logger
type SwLogger struct {
	log    *zap.Logger
	nowLog *zap.Logger
	tags   []Tag
}

// NewLogger возвращает новый экземпляр логгера
func NewLogger(appName string, options ...Option) SwLogger {
	jsonConfig := zapcore.EncoderConfig{
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.RFC3339TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "date",
		CallerKey:     "caller",
		StacktraceKey: "stackTrace",
	}

	optionsHolder := optionHolder{
		out: os.Stdout,
		err: os.Stderr,
	}
	for _, option := range options {
		optionsHolder = option.apply(optionsHolder)
	}

	skip := 1 // log.Debugf -> log.writeEntry
	if optionsHolder.skip > 0 {
		skip += optionsHolder.skip
	}

	zapOptions := []zap.Option{
		zap.ErrorOutput(newWriteSyncer(optionsHolder.err)),
		zap.AddCaller(),
		zap.AddCallerSkip(skip),
	}

	if optionsHolder.stacktrace {
		zapOptions = append(zapOptions, zap.AddStacktrace(_stacktraceEnabler))
	}
	containerId, _ := os.Hostname()

	log := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(jsonConfig),
			newWriteSyncer(optionsHolder.out),
			_levelEnabler),
		zapOptions...,
	).With(
		zap.String("appName", appName),
		zap.String("containerId", containerId),
	)

	logger := SwLogger{
		log:    log.WithOptions(zap.AddCallerSkip(1)),
		nowLog: log,
		tags:   optionsHolder.tags,
	}

	if optionsHolder.global {
		_global = logger
	}

	return logger
}

func (s SwLogger) WithSkip(n int) Logger {
	s.nowLog = s.nowLog.WithOptions(zap.AddCallerSkip(n))
	s.log = s.nowLog.WithOptions(zap.AddCallerSkip(1))
	return s
}

func (s SwLogger) WithUserInfo(userId int) Logger {
	if userId < 0 {
		return s
	}

	fields := []zap.Field{
		zap.Int("userId", userId),
	}
	s.log = s.log.With(fields...)
	s.nowLog = s.nowLog.With(fields...)

	return s
}

func (s SwLogger) WithTraceId(traceId trace.TraceID) Logger {
	if !traceId.IsValid() {
		return s
	}

	field := zap.String("traceId", traceId.String())
	s.log = s.log.With(field)
	s.nowLog = s.nowLog.With(field)

	return s
}

func (s SwLogger) WithBookId(bookId string) Logger {
	field := zap.String("bookId", bookId)
	s.log = s.log.With(field)
	s.nowLog = s.nowLog.With(field)

	return s
}

func (s SwLogger) WithTags(tags ...Tag) Logger {
	s.tags = append(s.tags, tags...)
	return s
}

func (s SwLogger) Message(level MessageLevel, message string) LogEntry {
	return Entry{
		level:     level,
		message:   message,
		writeFunc: s.writeEntry,
	}
}

func (s SwLogger) DebugEntry(message string) LogEntry {
	return Entry{
		level:     LevelDebug,
		message:   message,
		tags:      s.tags,
		writeFunc: s.writeEntry,
	}
}

func (s SwLogger) DebugEntryf(format string, args ...interface{}) LogEntry {
	return s.DebugEntry(fmt.Sprintf(format, args...))
}

func (s SwLogger) Debug(message string) {
	s.nowLog.Debug(message, marshalTags(s.tags))
}

func (s SwLogger) Debugf(format string, args ...interface{}) {
	s.nowLog.Debug(fmt.Sprintf(format, args...), marshalTags(s.tags))
}

func (s SwLogger) ErrorEntry(message string) LogEntry {
	return Entry{
		level:     LevelError,
		message:   message,
		tags:      s.tags,
		writeFunc: s.writeEntry,
	}
}

func (s SwLogger) ErrorEntryf(format string, args ...interface{}) LogEntry {
	return s.ErrorEntry(fmt.Sprintf(format, args...))
}

func (s SwLogger) Error(message string) {
	s.nowLog.Error(message, marshalTags(s.tags))
}

func (s SwLogger) Errorf(format string, args ...interface{}) {
	s.nowLog.Error(fmt.Sprintf(format, args...), marshalTags(s.tags))
}

func (s SwLogger) writeEntry(entry Entry) {
	switch entry.level {
	case LevelDebug:
		s.log.Debug(entry.message, marshalTags(entry.tags))
	case LevelError:
		s.log.Error(entry.message, marshalTags(entry.tags))
	}
}

func marshalTags(tags []Tag) zap.Field {
	return zap.Array("tags", zapcore.ArrayMarshalerFunc(func(e zapcore.ArrayEncoder) error {
		for _, tag := range tags {
			if len(tag) == 0 {
				continue
			}

			e.AppendString(string(tag))
		}

		return nil
	}))
}
