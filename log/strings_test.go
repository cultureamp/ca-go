package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakeCase(t *testing.T) {
	testCases := []struct {
		desc           string
		input          string
		expectedOutput string
	}{
		{
			desc:           "Test Case 1",
			input:          "CAP WITH SPACES",
			expectedOutput: "cap_with_spaces",
		},
		{
			desc:           "Test Case 2",
			input:          "CAP_WITH_UNDERSCORES",
			expectedOutput: "cap_with_underscores",
		},
		{
			desc:           "Test Case 3",
			input:          "Pascal Case With Spaces",
			expectedOutput: "pascal_case_with_spaces",
		},
		{
			desc:           "Test Case 4",
			input:          "lower with spaces",
			expectedOutput: "lower_with_spaces",
		},
		{
			desc:           "Test Case 5",
			input:          "miXed",
			expectedOutput: "mi_xed",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			sc := toSnakeCase(tC.input)
			assert.Equal(t, tC.expectedOutput, sc, tC.desc)
		})
	}
}
