package log

import (
	"testing"
)

func TestLog(t *testing.T) {

	writer := NewConsoleWriter(true, true)
	logger := New(WithLevel("info"), WithWriter(writer), WithCaller(true))

	logger.Debug("debug message")
	logger.Info("info message")
	SetLoglevel("error")
	logger.Warn("warn message")
	logger.Error("error message")
}
