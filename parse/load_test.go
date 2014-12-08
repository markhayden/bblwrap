package parse

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/markhayden/bblwrap/util"
)

func TestLoadLocalFile(t *testing.T) {
	var fileName, file string

	fileName = "singleclass.html"
	file, _ = loadLocalFile("tests/case/" + fileName)
	assert.T(t, file != "")
	util.PassLog("loaded file tests/case/singleclass.html successfully")
}
