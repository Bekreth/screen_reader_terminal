package screen_reader_terminal

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
