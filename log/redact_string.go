package log

import "strings"

func redactString(s string) string {
	const minChars = 10
	const showChars = 10
	const numStars = 10

	l := len(s)
	if l == 0 {
		return ""
	}

	var b strings.Builder
	b.Grow(l)

	aQuarter := l / 4
	if aQuarter > showChars {
		aQuarter = showChars
	}

	// write first "real" chars if we have more chars than "minChars"
	if l > minChars {
		b.WriteString(s[:aQuarter])
	}

	// write the middles "*"
	for i := 0; i < numStars; i++ {
		b.WriteString("*")
	}

	// write the end "real" chars if we have more chars then "minChars"
	if l > minChars {
		i := l - aQuarter
		b.WriteString(s[i:])
	}

	return b.String()
}
