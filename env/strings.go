package env

import (
	"strings"
)

// redact returns a string with some of the string replaced with *.
func redact(s string) string {
	const stars = 6
	const literals = 4

	l := len(s)
	var b strings.Builder
	b.Grow(l + stars)

	// no matter how long the string, show at least 6 "*"
	for i := 0; i < stars; i++ {
		b.WriteString("*")
	}

	if l <= stars {
		// For small strings, always return "******" don't suffix with any literal chars
		return b.String()
	}

	// For larger strings, redact the first n-chars, and keep the last 4 as is
	r := l - stars // how many of the last chars to print
	if r > literals {
		r = literals // we never print more than the last 4 chars
	}
	r = l - r // index of last chars

	b.WriteString(s[r:])
	return b.String()
}
