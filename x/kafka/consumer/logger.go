package consumer

// DebugLogger interface for consumer debug logging.
type DebugLogger interface {
	Print(string, ...any)
}

// LoggerFunc is a bridge between DebugLogger and any third party logger.
// Usage:
//   l := NewLogger() // some logger
//   c consumer.GroupConfig{
//       DebugLogger: consumer.LoggerFunc(l.Debug),
//   }
type LoggerFunc func(msg string, args ...any)

func (f LoggerFunc) Print(msg string, keyvals ...any) { f(msg, keyvals...) }

type noopDebugLogger struct{}

func (f noopDebugLogger) Print(_ string, _ ...any) {}
