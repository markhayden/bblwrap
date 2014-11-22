package util

// LeadingWhiteSpace removes the first character of a string
func LeadingWhiteSpace(s string) string {
	if string([]rune(s)[:1]) == " " {
		return string([]rune(s)[1:])
	}
	return s
}

// StripFirst removes the first character of a string
func StripFirst(s string) string {
	return string([]rune(s)[1:])
}

// StripLast removes the last character of a string
func StripLast(s string) string {
	return string([]rune(s)[:len(s)-1])
}
