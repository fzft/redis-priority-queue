package redisPriorityQueue

type LogLevel int

const (
	Info LogLevel = iota
	Error
	Debug
	Trace
)

// Logger is an interface for logging
// in the redis-priority-queue package
// if you want log this package, you need implement this interface
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Trace(format string, args ...interface{})
	GetLevel() LogLevel
	SetLevel(level LogLevel)
}
