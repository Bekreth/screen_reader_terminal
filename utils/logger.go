package utils

import (
	"fmt"
	"testing"
)

// Basic logging interface
type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

// NoOpLogger matches the Logger interface and does nothing if no logging is desired
type NoOpLogger struct {
}

func (NoOpLogger) Infof(format string, args ...interface{})  {}
func (NoOpLogger) Debugf(format string, args ...interface{}) {}

// TestLogger is to be used with the *testing.T.Log function for logging in test
type TestLogger struct {
	TestPrefix string
	Tester     *testing.T
}

func (logger TestLogger) Infof(format string, args ...interface{}) {
	logger.Tester.Logf(fmt.Sprintf(
		"IN: %v - %v",
		logger.TestPrefix,
		fmt.Sprintf(format, args...),
	))
}

func (logger TestLogger) Debugf(format string, args ...interface{}) {
	logger.Tester.Logf(fmt.Sprintf(
		"DE: %v - %v",
		logger.TestPrefix,
		fmt.Sprintf(format, args...),
	))
}
