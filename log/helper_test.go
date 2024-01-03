package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name     string
		snakeStr string
	}{
		{
			name:     "should remove spaces",
			snakeStr: "should_remove_spaces",
		},
		{
			name:     "should Remove CAP cases",
			snakeStr: "should_remove_cap_cases",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedStr := ToSnakeCase(tt.name)
			assert.Equal(t, tt.snakeStr, parsedStr)
		})
	}
}
