package log

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var (
	defaultLogger *Logger
)

func init() {
	defaultLogger = New()
}

type Option func(*options)

type options struct {
	level     string
	timestamp bool
	caller    bool
	writer    io.Writer
	pairs     []Pair
}

func correctLevel(level string) string {
	switch level {
	case "trace":
	case "debug":
	case "info":
	case "warn":
	case "error":
	case "fatal":
	default:
		return "error"
	}
	return level
}

func WithLevel(level string) Option {
	return func(o *options) {
		lvl := correctLevel(level)
		o.level = lvl
	}
}

func WithTimestamp(t bool) Option {
	return func(o *options) {
		o.timestamp = t
	}
}

func WithCaller(c bool) Option {
	return func(o *options) {
		o.caller = c
	}
}

func WithPairs(pairs ...Pair) Option {
	return func(o *options) {
		o.pairs = append(o.pairs, pairs...)
	}
}

func WithWriter(w io.Writer) Option {
	return func(o *options) {
		o.writer = w
	}
}

type Logger struct {
	opts   options
	logger *zerolog.Logger
	ctx    context.Context
}

func New(opts ...Option) *Logger {
	loggerOptions := options{
		level:     "error",
		timestamp: true,
		caller:    false,
		writer:    os.Stdout,
		pairs:     []Pair{},
	}

	for _, opt := range opts {
		opt(&loggerOptions)
	}

	lg := NewLogger(loggerOptions.writer, loggerOptions.level, loggerOptions.timestamp, loggerOptions.caller, loggerOptions.pairs...)

	logger := &Logger{
		opts:   loggerOptions,
		logger: lg,
	}

	return logger
}

func (l *Logger) Ctx(ctx context.Context) *Logger {
	l.ctx = ctx

	return l
}

func (l *Logger) Trace(format string, v ...any) {
	evt := l.logger.Trace()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func (l *Logger) Debug(format string, v ...any) {
	evt := l.logger.Debug()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func (l *Logger) Info(format string, v ...any) {
	evt := l.logger.Info()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	evt := l.logger.Warn()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	evt := l.logger.Error()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func (l *Logger) Fatal(format string, v ...any) {
	evt := l.logger.Fatal()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func (l *Logger) Panic(format string, v ...any) {
	evt := l.logger.Panic()

	if l.ctx != nil {
		evt = evt.Ctx(l.ctx)
	}

	evt.Msgf(format, v...)
}

func Trace(format string, v ...any) {
	defaultLogger.Trace(format, v...)
}

func Debug(format string, v ...any) {
	defaultLogger.Debug(format, v...)
}

func Info(format string, v ...any) {
	defaultLogger.Info(format, v...)
}

func Warn(format string, v ...any) {
	defaultLogger.Warn(format, v...)
}

func Error(format string, v ...any) {
	defaultLogger.Error(format, v...)
}

func Fatal(format string, v ...any) {
	defaultLogger.Fatal(format, v...)
}

func Panic(format string, v ...any) {
	defaultLogger.Panic(format, v...)
}
