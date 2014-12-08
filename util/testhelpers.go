package util

import "fmt"

// PassLog outputs a formatted log for more detailed testing feedback
func PassLog(s string) {
	fmt.Println(fmt.Sprintf("    √ %s", s))
}
