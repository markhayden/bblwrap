package parse

import (
	"fmt"
	// "io/ioutil"
	"github.com/markhayden/bblwrap/util"
	"regexp"
	"sort"
	"strings"
)

func StartCss() {
	sample := `
		body #taco .wee someting[name="adf"]{
			background:#F00;
		}
		a b, d c{
			background:#000;
		};
		.parent a{
			background:#666;
			color:#F00 !important;
			text-align:center;
		}`

	sample2 := `body {
			font: 80% arial, helvetica, sans-serif;
		}

		h1 {
			font-size: 1.5em;
		}

		h2 {
			font-size: 1em;
		}

		code {
			font-family: courier;
		}

		#example1, #example2 {
			background: #ccc;
			border: 2px solid black;
		}

		span {
			background: white;
			display: block;
			border: 0.5em solid red;
			padding: 1em;
			margin: 0.5em;
		}

		span.altern8 {
			background: #5b5;
		}

		#example2 span {
			display: inline;
		}`

	one := prettyStyles(sample)
	two := parseStyles(one)

	fmt.Sprintf("%v", sample2)
	fmt.Println(util.PrettyJson(two))
}

// loadLocalFile loads content from local path and returns as cleaned up string
// func loadLocalFile(tmp string) (string, error) {
// 	content, err := ioutil.ReadFile(tmp)
// 	if err != nil {
// 		return "", err
// 	}
// 	bodyPretty := util.PrepHtmlForJson(string(content), false)
// 	return bodyPretty, nil
// }

// func matchRegex(match, body string) string {
// 	var replaced []Styles
// 	regex := `<link rel=\"stylesheet\" type=\"text/css\" href=\"(?P<path>[A-Za-z0-9<>\&\/.;:\-_,= ]+)\">`

// 	r := regexp.MustCompile(regex)

// 	// find all external styles
// 	var needles [][]string
// 	if len(r.FindAllStringSubmatch(body, -1)) > 1 {
// 		needles = r.FindAllStringSubmatch(body, -1)
// 	}

// 	return body, replaced, nil
// }

type Style struct {
	Origin          string
	RawSelectors    string
	RawDeclarations string
	Selectors       []Selector
	Declarations    []Declaration
	Specificity     int
}

type Selector struct {
	Origin string
	Type   string
	Key    string
	Value  string
}

type Declaration struct {
	Origin    string
	Property  string
	Value     string
	Important bool
}

// praseStyles transforms a full string of styles into individual definitions
func parseStyles(subject string) []Style {
	// split everything up into a massive slice for processing
	// allSplit := regexp.MustCompile("{|}").Split(subject, -1)
	allSplit := strings.Split(subject, "}")

	var parsed []Style
	for _, val := range allSplit {
		// if the split results in any empty strings skip them
		if val == "" {
			continue
		}

		// separate the selectors from declarations
		selDecSplit := strings.Split(val, "{")
		if len(selDecSplit) != 2 {
			fmt.Println("Invalid CSS")
		}

		// handle multiple sets of selectors
		multiSelSplit := strings.Split(selDecSplit[0], ",")
		if len(multiSelSplit) == 0 {
			fmt.Println("More Invalid CSS")
		}

		// iterate over selectors and create new style structs
		for _, sel := range multiSelSplit {
			s := Style{
				Origin:          fmt.Sprintf("%s}", val),
				RawSelectors:    strings.TrimSpace(sel),
				RawDeclarations: selDecSplit[1],
			}

			err := s.parseSelectors()
			if err != nil {
				fmt.Println("Invalid Selectors")
			}

			err = s.parseDeclarations()
			if err != nil {
				fmt.Println("Invalid Declarations")
			}

			parsed = append(parsed, s)
		}
	}

	// sort styles by specificity
	sort.Sort(styleBySpecificity(parsed))

	return parsed
}

// parseSelectors handles primary parsing of selector for processing
// calculates specificity score and sorts descending order
func (s *Style) parseSelectors() error {
	// split individual selectors
	split := strings.Split(s.RawSelectors, " ")

	// base score
	specificity := 0

	// selector type definitions
	attr := regexp.MustCompile(`\[`)
	id := regexp.MustCompile(`^#`)
	class := regexp.MustCompile(`^\.`)
	// https://docs.google.com/a/markhayden.me/spreadsheets/d/19eMZ9bPB7rDsWnT0UZQFO5q7UxUXPjmhvepR7Edf--Y/edit#gid=0

	var selectors []Selector
	for _, o := range split {
		// if the origin is empty string we dont care about it
		if o == "" {
			continue
		}

		// set empty type
		var sType, sKey, sVal string

		switch {
		case attr.MatchString(o):
			sType = "attr"
			sKey = "need to parse key"
			sVal = util.StripFirst(o)
			specificity = specificity + 1000
		case id.MatchString(o):
			sType = "id"
			sKey = "id"
			sVal = util.StripFirst(o)
			specificity = specificity + 100
		case class.MatchString(o):
			sType = "class"
			sKey = "class"
			sVal = util.StripFirst(o)
			specificity = specificity + 10
		default:
			sType = "element"
			sVal = o
			specificity = specificity + 1
		}

		s := Selector{
			Origin: o,
			Type:   strings.TrimSpace(sType),
			Key:    strings.TrimSpace(sKey),
			Value:  strings.TrimSpace(sVal),
		}

		selectors = append(selectors, s)
	}

	// set selectors
	s.Selectors = selectors

	// set specificity score
	s.Specificity = specificity

	return nil
}

// parseDeclarations handles parsing the declaration string to a struct
func (s *Style) parseDeclarations() error {
	// split individual selectors
	split := strings.Split(s.RawDeclarations, ";")

	// selector type definitions
	important := regexp.MustCompile(`\!important`)

	var declarations []Declaration
	for _, o := range split {
		// if the origin is empty string we dont care about it
		if o == "" {
			continue
		}

		// split properties from values
		b := strings.Split(o, ":")
		if len(b) != 2 {
			fmt.Println("Invalid Declaration")
		}

		// make sure the value doesn't contains important
		if important.MatchString(b[1]) {
			b[1] = strings.Replace(b[1], " !important", "", -1)
			b[1] = strings.Replace(b[1], "!important", "", -1)
		}

		d := Declaration{
			Origin:    o,
			Property:  strings.TrimSpace(b[0]),
			Value:     strings.TrimSpace(b[1]),
			Important: important.MatchString(o),
		}

		declarations = append(declarations, d)
	}

	// set selectors
	s.Declarations = declarations

	return nil
}

// prettyStyles does some cleanup on the main style strign to force consistency
func prettyStyles(subject string) string {
	replace := map[string]string{
		"};": "}",
		"< ": "<",
		" >": ">",
		"\n": "",
		"	": "",
		"  ": "",
	}

	for fin, rep := range replace {
		subject = strings.Replace(subject, fin, rep, -1)
	}

	return subject
}

// styleBySpecificity handles sorting the slice of styles by specificity descending
type styleBySpecificity []Style

func (s styleBySpecificity) Len() int           { return len(s) }
func (s styleBySpecificity) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s styleBySpecificity) Less(i, j int) bool { return s[i].Specificity < s[j].Specificity }
