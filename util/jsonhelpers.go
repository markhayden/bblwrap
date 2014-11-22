package util

import (
	"encoding/json"
)

// PrettyJson
func PrettyJson(i interface{}) string {
	fltB, _ := json.MarshalIndent(i, "", "	")
	return string(fltB)
}
