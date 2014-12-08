package parse

import (
	"io/ioutil"
)

// loadLocalFile opens a local file and returns the body of said file
func loadLocalFile(tmp string) (string, error) {
	content, err := ioutil.ReadFile(tmp)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
