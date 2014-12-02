package parse

import (
// "fmt"
// "github.com/markhayden/bblwrap/util"
)

var (
	sampleCss  = `div[abc="hup"]{color:#FFF;}`
	sampleHtml = `<div stuff="hup">hi</div>`
)

type Style struct {
	Origin          string
	RawSelectors    string
	RawDeclarations string
	Selectors       []Selector
	Declarations    []Declaration
	Specificity     int
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
	styles := parseStyles(sampleCss)
	processHtml(sampleHtml, styles)
	//fmt.Println(util.PrettyJson(styles))
}
