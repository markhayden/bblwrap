package parse

import (
	"io/ioutil"
)

// loadLocalFile loads content from local path and returns as cleaned up string
func loadLocalFile(tmp string) (string, error) {
	content, err := ioutil.ReadFile(tmp)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
