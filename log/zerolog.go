package zrlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	Logger        *zerolog.Logger
	ConsoleLogger *zerolog.Logger
)

func init() {
	Logger, _ = NewLogger(os.Stdout, "trace", true, true)

	ConsoleLogger, _ = NewConsoleLogger("trace", true, true)
}

// SetLoglevel sets zerolog global logger level
// available level:
// "trace"
// "debug"
// "info"
// "warn"
// "error"
// "fatal"
// "panic"
// "disabled"
func SetLoglevel(level string) error {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(lvl)
	return nil
}

// SetTimeFormat set zerolog time format, not working with zerolog.ConsoleWriter
// available values:  zerolog.TimeFormatUnix ... and time.RFC3339 ...
func SetTimeFormat(format string) {
	zerolog.TimeFieldFormat = format
}

// zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

func NewConsoleWriter(withTimestamp, withCaller bool) zerolog.ConsoleWriter {

	writer := zerolog.ConsoleWriter{
		Out: os.Stdout,
		// FormatLevel: custom formatLevel makes console color not work
		// FormatLevel: func(i interface{}) string {
		// 	return strings.ToUpper(fmt.Sprintf("[%s]", i))
		// },
		// FormatFieldName: ,
		// FormatFieldValue: ,
		FormatMessage: func(i interface{}) string {
			if s, ok := i.(string); ok {
				return fmt.Sprintf("| %s |", s)
			} else {
				return "| |"
			}
		},
		PartsExclude: []string{},
	}

	if withTimestamp {
		writer.TimeFormat = time.RFC3339
	}

	if withCaller {
		writer.FormatCaller = func(i interface{}) string {
			paths := strings.Split(i.(string), "/")
			l := len(paths)
			if l > 2 {
				return strings.Join([]string{paths[l-2], paths[l-1]}, "/")
			}
			return filepath.Base(fmt.Sprintf("%s", i))
		}

	}

	return writer
}

type Pair struct {
	Key   string
	Value string
}

func NewLogger(writer io.Writer, level string, withTimestamp, withCaller bool, pairs ...Pair) (*zerolog.Logger, error) {
	loglevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	loggerContext := zerolog.New(writer).Level(loglevel).With()
	if withTimestamp {
		loggerContext = loggerContext.Timestamp()
	}
	if withCaller {
		loggerContext = loggerContext.Caller()
	}

	for _, pair := range pairs {
		loggerContext = loggerContext.Str(pair.Key, pair.Value)
	}

	logger := loggerContext.Logger()
	return &logger, nil
}

func NewConsoleLogger(level string, withTimestamp, withCaller bool, pairs ...Pair) (*zerolog.Logger, error) {
	consoleWriter := NewConsoleWriter(withTimestamp, withCaller)
	return NewLogger(consoleWriter, level, withTimestamp, withCaller, pairs...)
}
