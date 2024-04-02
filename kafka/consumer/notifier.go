package consumer

type Notifier interface {
	Info(event string, msg string)
	Error(err error, event string, msg string)
}

type noopNotifier struct{}

func (noop *noopNotifier) Info(event string, msg string)             {}
func (noop *noopNotifier) Error(err error, event string, msg string) {}
