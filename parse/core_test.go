package parse

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/markhayden/bblwrap/util"
)

func writeFile(file, payload string) {
	// write the file
	err := ioutil.WriteFile(file, []byte(payload), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCss(t *testing.T) {
	var fileName, file, payload, confirm string

	// // *************************************************************
	// // TEST ELEMENTS : elements.html
	// // *************************************************************
	// fileName = "elements.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// // *************************************************************
	// // TEST ELEMENTS : elements.html
	// // *************************************************************
	// fileName = "elementdotclass.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// // *************************************************************
	// // TEST SINGLE CLASS : singleclass.html
	// // *************************************************************
	// fileName = "singleclass.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// // *************************************************************
	// // TEST SINGLE ID : singleid.html
	// // *************************************************************
	// fileName = "singleid.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// // *************************************************************
	// // TEST NESTED CLASSES : nestedclass.html
	// // *************************************************************
	// fileName = "nestedclass.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// // *************************************************************
	// // TEST NESTED IDS : nestedid.html
	// // *************************************************************
	// fileName = "nestedid.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// // *************************************************************
	// // TEST ADVANCED ATTR SELECTORS : attributeselector.html
	// // *************************************************************
	// fileName = "attributeselector.html"
	// file, _ = loadLocalFile("tests/case/" + fileName)
	// payload = MakeInline(file)
	// confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	// assert.Equal(t, confirm, payload)
	// util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))

	// *************************************************************
	// TEST ADVANCED ATTR SELECTORS : attributeselector.html
	// *************************************************************
	fileName = "lytics.html"
	file, _ = loadLocalFile("tests/case/" + fileName)
	payload = MakeInline(file)
	fmt.Println(payload)
	confirm, _ = loadLocalFile("tests/confirm/" + fileName)
	writeFile("tests/confirm/"+fileName, payload)
	assert.Equal(t, confirm, payload)
	util.PassLog(fmt.Sprintf("inlined styles for %s successfully", fileName))
}

// func TestCase(t *testing.T) {
// 	//var fileName, file, payload string

// 	// *************************************************************
// 	// TEST ELEMENTS : elements.html
// 	// *************************************************************
// 	// fileName = "attributeselector.html"
// 	// file, _ = loadLocalFile("tests/case/" + fileName)
// 	// payload = MakeInline(file)
// 	// writeFile("tests/confirm/"+fileName, payload)
// }
