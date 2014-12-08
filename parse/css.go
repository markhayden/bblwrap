package parse

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	u "github.com/araddon/gou"
	"github.com/markhayden/bblwrap/util"
)

type ParseDefinition struct {
	name  string
	descr string
	rx    *regexp.Regexp
}

var (
	// definitions holds the master definitions for each of the supported selector types
	definitions = map[string]ParseDefinition{
		"attr": ParseDefinition{
			name:  "Attribute Selector",
			descr: "test",
			rx:    regexp.MustCompile(`\[`),
		},

		"id": ParseDefinition{
			name:  "Id Selector",
			descr: "test",
			rx:    regexp.MustCompile(`^#`),
		},

		"class": ParseDefinition{
			name:  "Class Selector",
			descr: "test",
			rx:    regexp.MustCompile(`^\.`),
		},

		"element": ParseDefinition{
			name:  "Element Selector",
			descr: "test",
			rx:    regexp.MustCompile(`^\S+\.\S+$`),
		},

		"wildcardElement": ParseDefinition{
			name:  "Wildcard Element Selector",
			descr: "test",
			rx:    regexp.MustCompile(`^\*$`),
		},
	}
)

// praseStyles transforms a full string of styles into various stucts and definiions for matching
// and ultimately reconstruction into an inline style for a specific html tag
func parseStyles(subject string) []Style {
	// clean up the style string be removing goofy characters, comments, line breaks, etc.
	subject = prettyStyles(subject)

	// break the master style string up into a slice of intidivual styles for processing
	allSplit := strings.Split(subject, "}")

	var parsed []Style
	for pos, val := range allSplit {
		// if the split results in any empty strings skip them
		if val == "" {
			u.Debugf("style split resulted in empty string, skipping: %d", pos)
			continue
		}

		// separate the selectors from declarations
		selDecSplit := strings.Split(val, "{")
		if len(selDecSplit) != 2 {
			u.Errorf("invalid CSS, could not separate selectors from declarations: %s", subject)
		}

		// if there are multiple selectors associated with a single set of declarations break them apart
		// so that we can analyze them as the individuals that they are
		multiSelSplit := strings.Split(selDecSplit[0], ",")
		if len(multiSelSplit) == 0 {
			u.Errorf("invalid CSS, could not separate declarations from values: %s", subject)
		}

		// iterate over selectors and create new style structs
		for _, sel := range multiSelSplit {
			s := Style{
				Origin:          fmt.Sprintf("%s}", val),
				RawSelectors:    strings.TrimSpace(sel),
				RawDeclarations: selDecSplit[1],
				Position:        pos,
			}

			// break the individual selectors out so that we can determine if it should be applied to an html tag
			err := s.parseSelectors()
			if err != nil {
				u.Errorf("invalid selector found: %v", sel)
			}

			// break the individual declarations out so that we can match, update and output them
			s.parseDeclarations()

			// add the successfully parsed style to the master slice
			parsed = append(parsed, s)
		}
	}

	// sort styles ascendin by specificity score so that definition replacement works
	sort.Sort(styleBySpecificity(parsed))

	return parsed
}

// parseSelectors handles primary parsing of selector for processing calculates specificity
// score, selector type, key, values, etc and rolls is up into a master definition
func (s *Style) parseSelectors() error {
	// split individual selectors
	split := strings.Split(s.RawSelectors, " ")

	// reset the specificity score upon start
	specificity := 0

	var selectors []Selector
	for _, o := range split {
		// if the origin is empty we have nothing to define so continue to next selector
		if o == "" {
			continue
		}

		// set empty type
		var sType, sKey, sElement, sVal string

		// determine what type of selector we are using and properly set defintion
		switch {
		case definitions["attr"].rx.MatchString(o):
			sType, sElement, sKey, sVal = parseAdvancedAttrSelector(o)
			specificity = specificity + 1000
		case definitions["id"].rx.MatchString(o):
			sType = "id"
			sKey = "id"
			sVal = util.StripFirst(o)
			specificity = specificity + 100
		case definitions["element"].rx.MatchString(o):
			sType = "class"
			sKey = "class"
			sElement = strings.Split(o, ".")[0]
			sVal = strings.Split(o, ".")[1]
			specificity = specificity + 11
		case definitions["class"].rx.MatchString(o):
			sType = "class"
			sKey = "class"
			sVal = util.StripFirst(o)
			specificity = specificity + 10
		case definitions["wildcardElement"].rx.MatchString(o):
			sType = "wildcard"
			sKey = "*"
			sVal = "*"
			specificity = specificity + 0
		default:
			sType = "element"
			sKey = o
			sVal = o
			specificity = specificity + 1
		}

		u.Debugf("found %s: element:%s, key:%s, value:%s, specificity:%d | %v", sType, sElement, sKey, sVal, specificity, o)

		s := Selector{
			Origin:  o,
			Type:    strings.TrimSpace(sType),
			Key:     strings.TrimSpace(sKey),
			Element: strings.TrimSpace(sElement),
			Value:   strings.TrimSpace(sVal),
		}

		selectors = append(selectors, s)
	}

	// set the selectors
	s.Selectors = selectors

	// set the depth (number)
	s.Depth = len(s.Selectors)

	// set the specificity score
	s.Specificity = specificity

	return nil
}

// parseAdvancedAttrSelector handles the primary parsing of advanced attribute selectors such as div[name="taco"] & div[name]
func parseAdvancedAttrSelector(s string) (string, string, string, string) {
	var sType, sElement, sKey, sVal string

	// need to handle both div[name="taco"] & div[name] so we use two regex definitions for accuracy
	attrElementValue, _ := regexp.Compile(`(?P<element>[a-zA-Z0-9\-\_]+)\[(?P<attr>[a-zA-Z0-9\-\_]+)=\"(?P<value>[a-zA-Z0-9\-\_]+)\"\]`)
	attrElement, _ := regexp.Compile(`(?P<element>[a-zA-Z0-9\-\_]+)\[(?P<attr>[a-zA-Z0-9\-\_]+)\]`)

	// check if the string has both an attribute and a value to determine how to parse
	if attrElementValue.MatchString(s) {
		// found both an attribute and a value
		parsed := attrElementValue.FindStringSubmatch(s)
		if len(parsed) < 4 {
			u.Errorf("failed to parse advanced selector: %v", s)
			return sType, sElement, sKey, sVal
		}

		sElement = parsed[1]
		sKey = parsed[2]
		sVal = parsed[3]
	} else {
		// found only an attribute
		parsed := attrElement.FindStringSubmatch(s)
		if len(parsed) < 3 {
			u.Errorf("failed to parse advanced selector: %v", s)
			return sType, sElement, sKey, sVal
		}

		sElement = parsed[1]
		sKey = parsed[2]
	}

	// it is likely that this will be an attribute but in the off chance a user has defined
	// something like div[class="something"] handle setting the proper type
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

// parseDeclarations handles the primary parsing of style declarations
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
			msg := fmt.Sprintf("Invalid declaration found: %v | %s", o, s.Origin)
			u.Errorf(msg)
			return errors.New(msg)
		}

		// make sure the value doesn't contains important
		if important.MatchString(b[1]) {
			b[1] = strings.Replace(b[1], " !important", "", -1)
			b[1] = strings.Replace(b[1], "!important", "", -1)
		}

		d := Declaration{
			Origin:    o,
			Property:  strings.TrimSpace(b[0]),                                 // make sure all properties have white space removed
			Value:     strings.Replace(strings.TrimSpace(b[1]), "\"", "'", -1), // make sure all values use single quotes and have no white space
			Important: important.MatchString(o),
		}

		declarations = append(declarations, d)
	}

	// set selectors
	s.Declarations = declarations

	return nil
}

// prettyStyles does some cleanup on a raw string of CSS. it will clean up hanging semi-colons, whitespace
// and ensure that all the tags are formatted the same way to improve accuracy of parsing.
func prettyStyles(subject string) string {
	replace := map[string]string{
		"};": "}",
		"< ": "<",
		" >": ">",
		"\n": "",
		"	": "",
		"  ": "",
	}

	// handle cleaning out the comments
	comments := regexp.MustCompile(`\/\*[^*]+\*\/`).FindAllString(subject, -1)
	for _, c := range comments {
		replace[c] = ""
	}

	for fin, rep := range replace {
		subject = strings.Replace(subject, fin, rep, -1)
	}

	return subject
}

// styleBySpecificity handles sorting the slice of styles by specificity ascending
type styleBySpecificity []Style

func (s styleBySpecificity) Len() int           { return len(s) }
func (s styleBySpecificity) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s styleBySpecificity) Less(i, j int) bool { return s[i].Specificity < s[j].Specificity }

// styleByPosition handles sorting the slice of styles by load order (position) ascending
type styleByPosition []Style

func (s styleByPosition) Len() int           { return len(s) }
func (s styleByPosition) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s styleByPosition) Less(i, j int) bool { return s[i].Position < s[j].Position }
