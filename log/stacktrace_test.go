package log

import (
	"fmt"
	"testing"

	ge "github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func TestStackTrace(t *testing.T) {
	standard_error := fmt.Errorf("standard err")
	assert.NotNil(t, standard_error)

	trace := stackTracer(standard_error)
	assert.Contains(t, trace, "runtime/asm")

	stacktraced_error := ge.Errorf("stack traced err")
	assert.NotNil(t, stacktraced_error)

	trace = stackTracer(stacktraced_error)
	assert.Contains(t, trace, "ca-go/log/stacktrace_test.go")
}
