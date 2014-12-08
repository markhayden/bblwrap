package parse

import (
	"errors"
	"fmt"
	//"github.com/markhayden/bblwrap/util"
	"regexp"
	"strings"

	u "github.com/araddon/gou"
)

// Html contains a core html object. all tags, the origin body and a set
// of temporary parent tokens are stored at any given point while parsing
type Html struct {
	Tags       []Tag
	Body       string
	TempParent []string
}

// Tag contains all the necessary data to perform the inlinine of styles
type Tag struct {
	Origin    string
	Token     string
	Parent    []string
	Outbound  string
	Selectors []Selector
}

// stringifyInlineStyle takes the final style declarations and converts them to an inline
// ready format for final token replacement and output
func stringifyInlineStyle(d []Declaration) string {
	var inlineOut string

	for _, v := range d {
		var important string
		if v.Important {
			important = " !important"
		}
		inlineOut = fmt.Sprintf(`%s %s:%s%s;`, inlineOut, v.Property, v.Value, important)
	}

	return fmt.Sprintf(`style="%s"`, strings.TrimSpace(inlineOut))
}

// processHtml handles the main parsing & inlining of all styles, a variety of tasks are performed
// such as general string cleanup, tag parsing, tag tokenization, style matching and token replacement
func processHtml(body string, styles []Style) string {
	// clean up body by removing breaks, unnecessary space, etc
	body = prettyHtml(body)

	// primary regex definitions necessary for parsing top level html and tokenizing
	elementRegex, _ := regexp.Compile(`<(\/)?\b[^>]*>`)
	closingElementRegex, _ := regexp.Compile(`<\/`)
	inlineStyleRegex, _ := regexp.Compile(`style\=\"\b[^>]*\"`)
	skipTokenRegex, _ := regexp.Compile(`<(\/|img|area|base|br|col|command|embed|hr|input|keygen|link|meta|param|source|track|wbr)`)

	// find all html tags so that we can process each individually
	tags := elementRegex.FindAllString(body, -1)
	m := Html{
		Body: body,
	}

	// process each html tag found and tokenize
	var parents []string
	for tagKey, tagValue := range tags {
		var token string
		var selectors []Selector
		var inlineRaw Style

		// handle initial logic for creating a token, adding to parent slice, etc if tag is not a closing tag
		if !closingElementRegex.MatchString(tagValue) {
			u.Debugf("opening tag: %s", tagValue)
			token = fmt.Sprintf("bblwrap-%d", tagKey)
			selectors, inlineRaw = parseTag(tagValue)
			m.Body = strings.Replace(m.Body, tagValue, fmt.Sprintf("<%s>", token), 1)
			parents = append(parents, token)
		}

		// if it is a closing tag make sure we remove the tag that it is in face closing from parent tracking
		if closingElementRegex.MatchString(tagValue) && len(parents) > 0 {
			u.Debugf("closing tag: %s", tagValue)
			parents = parents[:len(parents)-1]
		}

		// create a new tag definition for the tag in question
		t := Tag{
			Origin:    tagValue,
			Token:     token,
			Parent:    parents,
			Selectors: selectors,
		}

		// if it is a self closing tag we dont need to continue looking for a close, so go ahead and close
		if !closingElementRegex.MatchString(tagValue) && skipTokenRegex.MatchString(tagValue) && len(parents) > 0 {
			u.Debugf("removing self closing tag from parents: %s", tagValue)
			parents = parents[:len(parents)-1]
		}

		// add the tag to the master tag slice for future processing and replacement
		m.Tags = append(m.Tags, t)

		// for everything but a closing element we need to determine what styles should be applied
		var styleString string
		if !closingElementRegex.MatchString(tagValue) {
			// process fulfillment and stringify the results so we can drop in the inline styles if we have some
			match, inline := m.fullfilled(t, styles, inlineRaw)
			if match {
				styleString = stringifyInlineStyle(inline)
			}

			// if we already have inline styles we need to make sure we do not drop in a duplicate so find them
			currentInline := inlineStyleRegex.FindString(tagValue)
			tagReplace := tagValue
			if currentInline != "" {
				// an inline style exists so we need to replace it
				tagReplace = strings.Replace(tagReplace, currentInline, styleString, 1)
			} else {
				// no inline style exists so we need to add one
				tagReplace = strings.Replace(tagReplace, ">", fmt.Sprintf(" %s>", styleString), 1)
			}

			// save the final outbound tag to the master struct and move on to the next tag, we apply
			// them all in one big batch at the end...magic
			m.Tags[len(m.Tags)-1].Outbound = tagReplace
		}
	}

	// replace all the tokens with shiny new inlined styles
	for _, tag := range m.Tags {
		m.Body = strings.Replace(m.Body, "<"+tag.Token+">", tag.Outbound, -1)
	}

	// have a beer gopher, you done did good
	return m.Body
}

// fullfilled determines what styles need to be applied to a particular html tag, the climber
func (m *Html) fullfilled(climber Tag, mountain []Style, inline Style) (bool, []Declaration) {
	// prepare defaults
	var fStyles []Style

	// loop through all css to see which styles to add
	for _, mtn := range mountain {
		// loop through all selectors on style in reverse order, only loop if we match as all selectors need to be fulfilled
		isFulfilled := m.checkSelectorFulfillment(mtn, climber)

		// if the style is fulfilled add it to the master fstyle array for declaration processing later
		if isFulfilled {
			u.Debugf("★ FULFILLED: %v | %v \n", mtn.Origin, climber.Origin)
			fStyles = append(fStyles, mtn)
		}
	}

	if len(fStyles) > 0 {
		return true, declarationBattle(fStyles, inline)
	}

	return false, []Declaration{}
}

// checkSelectorFulfillment determines all necessary selectors on a can be fulfilled by a tag and its parents
func (m *Html) checkSelectorFulfillment(style Style, tag Tag) bool {
	m.TempParent = tag.Parent
	// iterate the selectors of each style in reverse order all selectors must be fulfilled
	// to validate a style, on each new style the counts should be reset
	for i := len(style.Selectors) - 1; i >= 0; i-- {
		// which selector are we testing on this loop
		ms := style.Selectors[i]

		// if any of the individual selectors are not fulfilled then return false, do not apply to tag
		isOkToContinue := m.doSelectorsMatch(ms, tag, len(style.Selectors))
		if !isOkToContinue {
			return false
		}
	}

	return true
}

// doSelectorsMatch determines if a single selector has been fulfilled by a tag or its parents
// each time a fulfillment is located the master temp parent slice is updated to ensure that we dont
// get false fulfillments, if a top level parent does not return it fails forever
func (m *Html) doSelectorsMatch(styleSel Selector, tag Tag, l int) bool {
	// loop through all parents on the style looking for a match
	// as we loop, remove them from the temp slice so we dont iterate over them
	// more than one time
	for i := len(m.TempParent) - 1; i >= 0; i-- {
		subject, err := m.findTagByToken(m.TempParent[i])
		if err != nil {
			u.Errorf("could not find parent: %s", m.TempParent[i])
			return false
		}

		// pop off parent we are looping
		m.TempParent = m.TempParent[:len(m.TempParent)-1]

		for _, tagSel := range subject.Selectors {
			// in the case of wildcard selector, apply to every element
			if styleSel.Type == "wildcard" {
				return true
			}

			if tagSel.Value == styleSel.Value && tagSel.Type == styleSel.Type && tagSel.Key == styleSel.Key {
				// in the case of element.class, make sure we fulfill all demands
				if styleSel.Element != "" {
					if styleSel.Element == tagSel.Element {
						return true
					}
				} else {
					return true
				}
			}
		}
		return false
	}
	return false
}

// findTagByToken finds a specific tokenized html tag based on its unique token id
func (m *Html) findTagByToken(token string) (Tag, error) {
	for _, t := range m.Tags {
		if t.Token == token {
			return t, nil
		}
	}

	fmt.Printf("No tag found with token: %s\n", token)
	return Tag{}, errors.New(fmt.Sprintf("No tag found with token: %s", token))
}

// declarationBattle loops through all fulfilled styles to determine which ones should be surfaced
// in the final output. this ensures we dont have root level styles overriding those with a higher
// specificity level
func declarationBattle(challengers []Style, inline Style) []Declaration {
	var champions []Declaration

	// if the inline style has declarations append it to challengers for final override
	if len(inline.Declarations) > 0 {
		challengers = append(challengers, inline)
	}

	for _, challenger := range challengers {
		// iterate over the challengers declarations
		var found bool
		for _, declaration := range challenger.Declarations {
			// for each one, see if it exists already. if it does update if not important
			for k, champion := range champions {
				if champion.Property == declaration.Property {
					found = true
					if declaration.Important {
						// this declaration wins, update it
						champions[k].Value = declaration.Value
						champions[k].Important = true
					} else if !champion.Important {
						// this declaration wins, update it
						champions[k].Value = declaration.Value
					}
				}
			}

			if !found {
				// this declaration does not exist in the master slice yet, add it
				champions = append(champions, declaration)
			}

			found = false
		}
	}

	return champions
}

// parseTag handles the dirty work of parsing an html tag out into its various selectors
// inline styles, etc
func parseTag(tag string) ([]Selector, Style) {
	elementRegex, _ := regexp.Compile(`<(?P<element>[a-zA-Z0-9\-\_]+)[ |>|/>]`)
	defRegex, _ := regexp.Compile(`(?P<declaration>[a-zA-Z0-9\-\_]+)="(?P<value>[a-zA-Z0-9\-\_\/;:# ]+)"`)

	var master, element, class, id, attr []Selector
	var inline Style

	// get the element name
	el := elementRegex.FindStringSubmatch(tag)
	if len(el) != 2 {
		u.Errorf("invalid tag found: %v", tag)
		return master, Style{}
	}

	e := Selector{
		Origin: el[0],
		Type:   "element",
		Key:    el[1],
		Value:  el[1],
	}

	element = append(element, e)

	// get the other definitions
	needles := defRegex.FindAllStringSubmatch(tag, -1)
	for _, v := range needles {
		if len(v) != 3 {
			u.Errorf("invalid CSS, got wrong count(%d) for declaration length: %v", len(v), v)
			continue
		}

		declaration := v[1]

		if declaration == "style" {
			v[2] = strings.Replace(v[2], "; ", ";", -1)
		}

		values := strings.Split(v[2], " ")

		for _, value := range values {
			s := Selector{
				Origin:  v[0],
				Element: e.Value,
				Key:     strings.Split(v[0], "=")[0],
				Value:   strings.TrimSpace(value),
			}

			switch declaration {
			case "class":
				s.Type = "class"
				class = append(class, s)
			case "id":
				s.Type = "id"
				id = append(id, s)
			case "style":
				s.Type = "style"
				inline = parseStyles("inline{" + s.Value + "}")[0]
			default:
				s.Type = "attr"
				attr = append(attr, s)
			}
		}

	}

	// merge it all back togehter into one master declaration slice
	master = append(master, element...)
	master = append(master, class...)
	master = append(master, id...)
	master = append(master, attr...)

	return master, inline
}

// prettyHtml does some cleanup on the main style strign to force consistency
func prettyHtml(subject string) string {
	replace := map[string]string{
		"\n": "",
		"	": "",
		"  ": "",
	}

	for fin, rep := range replace {
		subject = strings.Replace(subject, fin, rep, -1)
	}

	return subject
}
