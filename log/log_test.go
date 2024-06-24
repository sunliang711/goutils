package log

import (
	"os"
	"testing"
)

func TestConsoleLog(t *testing.T) {

	writer := NewConsoleWriter(true, true)
	logger := New(WithLevel("debug"), WithWriter(writer), WithCaller(true))

	logger.Info("info message")
	SetLoglevel("fatal")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestJsonLog(t *testing.T) {

	logger := New(WithLevel("debug"), WithWriter(os.Stdout), WithCaller(true))
	logger.logger.With()

	logger.Info("info message")
	// SetLoglevel("fatal")
	logger.Warn("warn message")
	newLogger := logger.With("function", "TestJsonLog")
	logger.Error("error message")
	newLogger.Error("new error message")

	newLogger2 := logger.With("function", "TestJsonLog2")
	newLogger2.Error("new error message2")
}

func TestDefaultLog(t *testing.T) {
	err := SetLoglevel("info")
	if err != nil {
		t.Fatal(err)
	}

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
}
