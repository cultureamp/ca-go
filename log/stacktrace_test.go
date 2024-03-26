package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-errors/errors"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc        string
		err         error
		errContains string
	}{
		{
			desc:        "standard error",
			err:         fmt.Errorf("standard err"),
			errContains: "runtime/asm",
		},
		{
			desc:        "library error",
			err:         errors.Errorf("stack traced err"),
			errContains: "ca-go/log/stacktrace_test.go",
		},
		{
			desc:        "json escape",
			err:         errors.Errorf("\"} {\" \"}\"}"),
			errContains: "ca-go/log/stacktrace_test.go",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			trace := stackTracer(tC.err)
			assert.Contains(t, trace, tC.errContains)

			encoded := bytes.NewBufferString("")
			json.NewEncoder(encoded).Encode(trace)

			fmt.Printf("trace: %v\n", trace)
			fmt.Printf("encoded: %v\n", encoded)
		})
	}
}
