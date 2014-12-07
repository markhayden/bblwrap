package parse

import (
	"errors"
	"fmt"
	//"github.com/markhayden/bblwrap/util"
	"regexp"
	"strings"
)

type Html struct {
	Tags []Tag
	Body string
	TempParent   []string
}

type Tag struct {
	Origin    string
	Token     string
	Parent    []string
	Outbound  string
	Selectors []Selector
}

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

func processHtml(body string, styles []Style) string {
	// clean up body
	body = prettyHtml(body)

	// parse tags
	elementRegex, _ := regexp.Compile(`<(\/)?\b[^>]*>`)
	closingElementRegex, _ := regexp.Compile(`<\/`)
	inlineStyleRegex, _ := regexp.Compile(`style\=\"\b[^>]*\"`)
	skipTokenRegex, _ := regexp.Compile(`<(\/|img|area|base|br|col|command|embed|hr|input|keygen|link|meta|param|source|track|wbr)`)
	// isBodyRegex, _ := regexp.Compile(`<(\/body|body)`)
	tags := elementRegex.FindAllString(body, -1)
	m := Html{
		Body: body,
	}

	var parents []string
	for tagKey, tagValue := range tags {
		var token string
		var selectors []Selector
		var inlineRaw Style

		// handle logic for self closing and closing tags
		if !closingElementRegex.MatchString(tagValue) {
			fmt.Println("adding", tagValue)
			token = fmt.Sprintf("bblwrap-%d", tagKey)
			selectors, inlineRaw = parseTag(tagValue)
			m.Body = strings.Replace(m.Body, tagValue, fmt.Sprintf("<%s>", token), 1)
			parents = append(parents, token)
		}

		// handle parent logic on closing tag
		if closingElementRegex.MatchString(tagValue) && len(parents) > 0 {
			fmt.Println("closing", tagValue)
			parents = parents[:len(parents)-1]
		}

		// if token != "" && parents[len(parents)-1] != token {
		// 	fmt.Println("does this ever run")
		// 	parents = append(parents, token)
		// }

		t := Tag{
			Origin:    tagValue,
			Token:     token,
			Parent:    parents,
			Selectors: selectors,
		}

		// handle removing the self closing stuff from the main parent slice
		if !closingElementRegex.MatchString(tagValue) && skipTokenRegex.MatchString(tagValue) && len(parents) > 0 {
			fmt.Println("killing", tagValue)
			parents = parents[:len(parents)-1]
		}

		// add to master tags slice
		m.Tags = append(m.Tags, t)

		// handle logic for self closing and closing tags
		var styleString string
		if !closingElementRegex.MatchString(tagValue) {
			match, inline := m.fullfilled(t, styles, inlineRaw)
			if match {
				styleString = stringifyInlineStyle(inline)
			}

			currentInline := inlineStyleRegex.FindString(tagValue)
			tagReplace := tagValue

			if currentInline != "" {
				// an inline style exists so we need to replace it
				tagReplace = strings.Replace(tagReplace, currentInline, styleString, 1)
			} else {
				// no inline style exists so we need to add one
				tagReplace = strings.Replace(tagReplace, ">", fmt.Sprintf(" %s>", styleString), 1)
			}

			// finally replace the string in the body
			m.Tags[len(m.Tags)-1].Outbound = tagReplace
		}
	}

	// replace all the tags and prepare output
	for _, tag := range m.Tags {
		m.Body = strings.Replace(m.Body, "<"+tag.Token+">", tag.Outbound, -1)
	}

	// fmt.Println(util.PrettyJson(styles))
	// fmt.Println("")
	// fmt.Println(util.PrettyJson(m.Tags))
	// fmt.Println(m.Body)
	// fmt.Println(util.PrettyJson(m))

	return m.Body
}

func (m *Html) fullfilled(climber Tag, mountain []Style, inline Style) (bool, []Declaration) {
	// prepare defaults
	var fStyles []Style

	// loop through all css to see which styles to add
	for _, mtn := range mountain {
		// loop through all selectors on style in reverse order, only loop if we match as all selectors need to be fulfilled
		isFulfilled := m.checkSelectorFulfillment(mtn, climber)

		// if the style is fulfilled add it to the master fstyle array for declaration processing later
		if isFulfilled {
			//fmt.Printf("++++++++++ FULFILLED: %v | %v \n", mtn.Origin, climber.Origin)
			fStyles = append(fStyles, mtn)
		}
	}

	if len(fStyles) > 0 {
		return true, declarationBattle(fStyles, inline)
	}

	//fmt.Printf("---------- NOT FULFILLED: %v\n", climber.Origin)
	return false, []Declaration{}
}


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


func (m *Html) doSelectorsMatch(styleSel Selector, tag Tag, l int) bool {

	//fmt.Println(util.PrettyJson(tag))

	// loop through all parents on the style looking for a match
	// as we loop, remove them from the temp slice so we dont iterate over them
	// more than one time
	for i := len(m.TempParent) - 1; i >= 0; i-- {
		subject, err := m.findTagByToken(m.TempParent[i])
		if err != nil {
			fmt.Println("Could not locate that parent. Failed.")
			return false
		}

		// pop off parent we are looping
		m.TempParent = m.TempParent[:len(m.TempParent)-1]

		for _, tagSel := range subject.Selectors {
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






func (m *Html) findTagByToken(token string) (Tag, error) {
	for _, t := range m.Tags {
		if t.Token == token {
			return t, nil
		}
	}

	fmt.Printf("No tag found with token: %s\n", token)
	return Tag{}, errors.New(fmt.Sprintf("No tag found with token: %s", token))
}

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
						// update record
						champions[k].Value = declaration.Value
						champions[k].Important = true
					} else if !champion.Important {
						// update record
						champions[k].Value = declaration.Value
					}
				}
			}

			if !found {
				// add record
				champions = append(champions, declaration)
			}

			found = false
		}
	}

	return champions
}

func parseTag(tag string) ([]Selector, Style) {
	elementRegex, _ := regexp.Compile(`<(?P<element>[a-zA-Z0-9\-\_]+)[ |>|/>]`)
	defRegex, _ := regexp.Compile(`(?P<declaration>[a-zA-Z0-9\-\_]+)="(?P<value>[a-zA-Z0-9\-\_\/;:# ]+)"`)

	var master, element, class, id, attr []Selector
	var inline Style

	// get the element name
	el := elementRegex.FindStringSubmatch(tag)
	if len(el) != 2 {
		fmt.Println("Invalid Tag")
		fmt.Println(tag)
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
			fmt.Println("invalid class, declaration")
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

// prettyStyles does some cleanup on the main style strign to force consistency
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
