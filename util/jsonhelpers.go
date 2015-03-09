package util

import (
	"encoding/json"
)

// PrettyJson
func PrettyJson(i interface{}) string {
	fltB, _ := json.MarshalIndent(i, "", "	")
	return string(fltB)
}

// ToJson
func ToJson(i interface{}) string {
	fltB, _ := json.Marshal(i)
	return string(fltB)
}
