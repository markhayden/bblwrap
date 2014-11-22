package util

// StripFirst removes the first character of a string
func StripFirst(s string) string {
	return string([]rune(s)[1:])
}

// StripLast removes the last character of a string
func StripLast(s string) string {
	return string([]rune(s)[:len(s)-1])
}
