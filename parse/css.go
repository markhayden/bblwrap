package parse

import (
	"fmt"
	// "io/ioutil"
	"github.com/markhayden/bblwrap/util"
	"regexp"
	"sort"
	"strings"
	"errors"
	//"reflect"
)

// praseStyles transforms a full string of styles into individual definitions
func parseStyles(subject string) []Style {
	// clean up the string
	subject = prettyStyles(subject)

	// split everything up into a massive slice for processing
	// allSplit := regexp.MustCompile("{|}").Split(subject, -1)
	allSplit := strings.Split(subject, "}")

	var parsed []Style
	for pos, val := range allSplit {
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
		//fmt.Println(selDecSplit[0])

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
				Position: pos,
			}

			err := s.parseSelectors()
			if err != nil {
				fmt.Println("Invalid Selectors")
			}

			err = s.parseDeclarations()
			if err != nil {
				fmt.Printf("Invalid Declaration: %v\n", err)
			}

			parsed = append(parsed, s)
		}
	}

	// sort styles by specificity
	sort.Sort(styleBySpecificity(parsed))

	// dedupe the final parsed styles and alert of potential code goofs
	dedupeStyles(parsed)

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
	elementClass := regexp.MustCompile(`^\S+\.\S+$`)
	// https://docs.google.com/a/markhayden.me/spreadsheets/d/19eMZ9bPB7rDsWnT0UZQFO5q7UxUXPjmhvepR7Edf--Y/edit#gid=0

	var selectors []Selector
	for _, o := range split {
		// if the origin is empty string we dont care about it
		if o == "" {
			continue
		}

		// set empty type
		var sType, sKey, sElement, sVal string

		switch {
		case attr.MatchString(o):
			sType, sElement, sKey, sVal = parseAdvancedAttrSelector(o)
			specificity = specificity + 1000
		case id.MatchString(o):
			sType = "id"
			sKey = "id"
			sVal = util.StripFirst(o)
			specificity = specificity + 100
		case elementClass.MatchString(o):
			sType = "class"
			sKey = "class"
			sElement = strings.Split(o, ".")[0]
			sVal = strings.Split(o, ".")[1]
			specificity = specificity + 11
		case class.MatchString(o):
			sType = "class"
			sKey = "class"
			sVal = util.StripFirst(o)
			specificity = specificity + 10
		default:
			sType = "element"
			sKey = o
			sVal = o
			specificity = specificity + 1
		}

		s := Selector{
			Origin:  o,
			Type:    strings.TrimSpace(sType),
			Key:     strings.TrimSpace(sKey),
			Element: strings.TrimSpace(sElement),
			Value:   strings.TrimSpace(sVal),
		}

		//fmt.Println(util.PrettyJson(s))

		selectors = append(selectors, s)
	}

	// set selectors
	s.Selectors = selectors

	// set depth
	s.Depth = len(s.Selectors)

	// set specificity score
	s.Specificity = specificity

	return nil
}

func parseAdvancedAttrSelector(s string) (string, string, string, string) {
	var sType, sElement, sKey, sVal string

	attrElementValue, _ := regexp.Compile(`(?P<element>[a-zA-Z0-9\-\_]+)\[(?P<attr>[a-zA-Z0-9\-\_]+)=\"(?P<value>[a-zA-Z0-9\-\_]+)\"\]`)
	attrElement, _ := regexp.Compile(`(?P<element>[a-zA-Z0-9\-\_]+)\[(?P<attr>[a-zA-Z0-9\-\_]+)\]`)

	// check if we have both a attribute and a value
	if attrElementValue.MatchString(s) {
		parsed := attrElementValue.FindStringSubmatch(s)
		if len(parsed) < 4 {
			fmt.Println("Failed to parse advanced selector.")
			return sType, sElement, sKey, sVal
		}

		sElement = parsed[1]
		sKey = parsed[2]
		sVal = parsed[3]
	} else {
		// only have an attribute so run through defaults
		parsed := attrElement.FindStringSubmatch(s)
		if len(parsed) < 3 {
			fmt.Println("Failed to parse advanced selector (key only).")
			return sType, sElement, sKey, sVal
		}

		sElement = parsed[1]
		sKey = parsed[2]
	}

	// most likely it will be an attr tag but if its id or class, smack them with a stick for yucky code and handle it
	switch {
	case sKey == "class":
		sType = "class"
	case sKey == "id":
		sType = "id"
	default:
		sType = "attr"
	}

	return sType, sElement, sKey, sVal
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
		b := strings.SplitN(o, ":", 2)
		if len(b) != 2 {
			msg := fmt.Sprintf("Invalid declaration: %v", o)
			if err := addToLog(msg); err != nil {
				fmt.Println(err)
			}

			return errors.New(msg)
		}

		// make sure the value doesn't contains important
		if important.MatchString(b[1]) {
			b[1] = strings.Replace(b[1], " !important", "", -1)
			b[1] = strings.Replace(b[1], "!important", "", -1)
		}

		d := Declaration{
			Origin:    o,
			Property:  strings.TrimSpace(b[0]), // make sure all properties have white space removed
			Value:     strings.Replace(strings.TrimSpace(b[1]), "\"", "'", -1), // make sure all values use single quotes and have no white space
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

// dedupeStyles does some cleanup on the main style strign to force consistency
func dedupeStyles(dedupe []Style){
	// break the styles out into slices by specificity
	sortMe := map[int][]int{}
	for _, style := range dedupe{
		sortMe[style.Specificity] = append(sortMe[style.Specificity], style.Position)
	}
	//fmt.Println(sortMe)

	// dedupe the individual slices

	// reconstruct everything after final sort for output
	// var out []Style
	// for _, appendMe := range sortMe {
	// 	sort.Sort(styleByPosition(appendMe))
	// 	out = append(out, appendMe...)
	// }

	//return out
}

// styleBySpecificity handles sorting the slice of styles by specificity descending
type styleBySpecificity []Style

func (s styleBySpecificity) Len() int           { return len(s) }
func (s styleBySpecificity) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s styleBySpecificity) Less(i, j int) bool { return s[i].Specificity < s[j].Specificity }

// styleByPosition handles sorting the slice of styles by specificity descending
type styleByPosition []Style

func (s styleByPosition) Len() int           { return len(s) }
func (s styleByPosition) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s styleByPosition) Less(i, j int) bool { return s[i].Position < s[j].Position }