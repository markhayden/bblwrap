package parse

import (
	"fmt"
	"io/ioutil"
	"regexp"
	//"strings"
	//"github.com/markhayden/bblwrap/util"
)

type Style struct {
	Origin          string
	RawSelectors    string
	RawDeclarations string
	Selectors       []Selector
	Declarations    []Declaration
	Specificity     int
	Position        int
	Depth           int
}

type Selector struct {
	Origin   string
	Position int `json:",omitempty"`
	Type     string
	Key      string
	Element  string
	Value    string
}

type Declaration struct {
	Origin    string
	Property  string
	Value     string
	Important bool
}

func Kickoff() {
	bodyRegex, _ := regexp.Compile(`<body([^>]+)?>([\s\S]*)<\/body>`)
	stylesRegex, _ := regexp.Compile(`<style([^>]+)?>(?P<styles>[\s\S]+)<\/style>`)

	file, _ := loadLocalFile("sample.html")

	file = prettyStyles(file)

	htmlStyleSlice := stylesRegex.FindAllStringSubmatch(file, -1)
	s := parseStyles(htmlStyleSlice[0][len(htmlStyleSlice[0])-1])

	htmlBodySlice := bodyRegex.FindAllString(file, -1)
	data := []byte(`<!DOCTYPE html><html>`+processHtml(htmlBodySlice[0], s)+`</html>`)
	err := ioutil.WriteFile("output.html", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}
