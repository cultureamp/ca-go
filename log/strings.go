package log

import (
	"strings"
	"unicode"
)

// toSnakeCase returns a new string in the format word_word
func toSnakeCase(s string) string {
	var sb strings.Builder

	in := []rune(strings.TrimSpace(s))
	for i, r := range in {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(in[i-1]) {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			if unicode.IsSpace(r) {
				if !unicode.IsSpace(in[i-1]) {
					sb.WriteRune('_')
				}
			} else {
				sb.WriteRune(r)
			}
		}
	}

	return sb.String()
}
