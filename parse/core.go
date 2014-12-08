package parse

import (
	"flag"
	"regexp"

	u "github.com/araddon/gou"
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

var logLevel *string = flag.String("logging", "info", "Which log level: [debug,info,warn,error,fatal]")

func MakeInline(source string) string {
	// set up logging
	flag.Parse()
	u.SetupLogging(*logLevel)

	// prepare main regex to parse incoming file
	bodyRegex, _ := regexp.Compile(`<body([^>]+)?>([\s\S]*)<\/body>`)
	stylesRegex, _ := regexp.Compile(`<style([^>]+)?>(?P<styles>[\s\S]+)<\/style>`)

	// clean up spacing, comments, breaks and other goofy things with style string
	source = prettyStyles(source)

	// parse the styles
	htmlStyleSlice := stylesRegex.FindAllStringSubmatch(source, -1)
	s := parseStyles(htmlStyleSlice[0][len(htmlStyleSlice[0])-1])

	// parse and inline the html
	htmlBodySlice := bodyRegex.FindAllString(source, -1)

	// add doctype and closing tag to payload
	output := `<!DOCTYPE html><html>` + processHtml(htmlBodySlice[0], s) + `</html>`

	// final inlined styles, browser ready
	return output
}
