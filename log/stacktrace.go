package log

import (
	"bytes"
	"runtime"
	"strings"

	"github.com/go-errors/errors"
)

// logStackTracer implements the zerolog.ErrorStackMarshaler func signature.
func logStackTracer(err error) interface{} {
	return stackTracer(err)
}

func stackTracer(err error) string {
	// is it the "github.com/go-errors/errors" error type? If so, it has a stack trace we can use
	var e *errors.Error
	if errors.As(err, &e) {
		s := string(e.Stack())
		return cleanStackTrace(s)
	}

	// Otherwise we just get the current stack trace (minus the ca-go calls)
	return cleanStackTrace(currentStack(5))
}

func currentStack(skip int) string {
	stack := make([]uintptr, errors.MaxStackDepth)
	length := runtime.Callers(skip, stack)
	stack = stack[:length]

	buf := bytes.Buffer{}
	for _, pc := range stack {
		frame := errors.NewStackFrame(pc)
		buf.WriteString(frame.String())
	}

	return buf.String()
}

func cleanStackTrace(stack string) string {
	// since we log in JSON make sure that the stack trace does NOT have any "{", "}" or "\""
	nr := strings.NewReplacer(
		"{", "",
		"}", "",
		"\"", "",
	)

	return nr.Replace(stack)
}
