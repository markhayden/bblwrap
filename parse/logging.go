package parse

import (
	"fmt"
)

var (
	MasterLog = ""
)

func addToLog(entry string) error {
	MasterLog += fmt.Sprintf("%s\n", entry)
	fmt.Println(MasterLog)

	return nil
}
