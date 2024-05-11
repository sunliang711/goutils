package log

import (
	"testing"
)

func TestConsoleLog(t *testing.T) {

	writer := NewConsoleWriter(true, true)
	logger := New(WithLevel("info"), WithWriter(writer), WithCaller(true))

	logger.Debug("debug message")
	logger.Info("info message")
	SetLoglevel("fatal")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestLog(t *testing.T) {
	err := SetLoglevel("error")
	if err != nil {
		t.Fatal(err)
	}

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
}
