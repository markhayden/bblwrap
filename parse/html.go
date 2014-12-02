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
}

type Tag struct {
	Origin    string
	Token     string
	Parent    string
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

func processHtml(body string, styles []Style) {
	// clean up body
	body = prettyHtml(body)

	// parse tags
	elementRegex, _ := regexp.Compile(`<(\/)?\b[^>]*>`)
	closingElementRegex, _ := regexp.Compile(`<\/`)
	inlineStyleRegex, _ := regexp.Compile(`style\=\"\b[^>]*\"`)
	skipTokenRegex, _ := regexp.Compile(`<(\/|img|area|base|br|col|command|embed|hr|input|keygen|link|meta|param|source|track|wbr)`)
	tags := elementRegex.FindAllString(body, -1)
	m := Html{
		Body: body,
	}

	var parents []string
	for tagKey, tagValue := range tags {
		var token, parent string
		var selectors []Selector
		var inlineRaw Style

		// handle logic for self closing and closing tags
		if !skipTokenRegex.MatchString(tagValue) {
			token = fmt.Sprintf("bblwrap-%d", tagKey)

			if len(parents) > 0 {
				parent = parents[len(parents)-1]
			}

			parents = append(parents, token)
			selectors, inlineRaw = parseTag(tagValue)

			m.Body = strings.Replace(m.Body, tagValue, fmt.Sprintf("<%s>", token), 1)
		}

		// handle parent logic on closing tag
		if closingElementRegex.MatchString(tagValue) && len(parents) > 0 {
			parents = parents[:len(parents)-1]
		}

		t := Tag{
			Origin:    tagValue,
			Token:     token,
			Parent:    parent,
			Selectors: selectors,
		}

		if token != "" && parents[len(parents)-1] != token {
			parents = append(parents, token)
		}

		// handle logic for self closing and closing tags
		var styleString string
		if !skipTokenRegex.MatchString(tagValue) {
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
			t.Outbound = tagReplace
		}

		m.Tags = append(m.Tags, t)
	}

	// replace all the tags and prepare output
	for _, tag := range m.Tags {
		m.Body = strings.Replace(m.Body, "<"+tag.Token+">", tag.Outbound, -1)
	}

	// temporary output
	// fmt.Println("")
	// fmt.Println(util.PrettyJson(styles))

	// fmt.Println("")
	// fmt.Println(util.PrettyJson(m.Tags))

	fmt.Println("")
	fmt.Println("Output")
	fmt.Println(m.Body)

	// fmt.Println(m.Body)
	// fmt.Println(util.PrettyJson(m))
}

func (m *Html) fullfilled(climber Tag, mountain []Style, inline Style) (bool, []Declaration) {
	// iterate over all styles
	var fStyles []Style
	for _, mtn := range mountain {
		matchCount := 0
		rolloverCount := false

		// fmt.Println(util.PrettyJson(mtn))

		// iterate the selectors of each style in reverse order all selectors must be fulfilled
		// to validate a style, on each new style the counts should be reset
		for i := len(mtn.Selectors) - 1; i >= 0; i-- {
			v := mtn.Selectors[i]

			if !rolloverCount {
				matchCount = 0
			}

			for _, c := range climber.Selectors {
				if c.Value == v.Value && c.Type == v.Type && c.Key == v.Key {
					if v.Element != "" {
						if v.Element == c.Element {
							matchCount++
							break
						}
					} else {
						matchCount++
						break
					}
				}
			}

			if matchCount == len(mtn.Selectors) {
				// fuck yah, it was fulfilled
				fStyles = append(fStyles, mtn)
				break
			} else if matchCount > 0 && matchCount < len(mtn.Selectors) {
				// if there is no parent continue, not a match
				if climber.Parent == "" {
					continue
				}

				// begin cascading through parents
				parent, err := m.findTagByToken(climber.Parent)
				if err != nil {
					continue
				}

				climber = parent
				rolloverCount = true
			} else {
				rolloverCount = false
			}
		}
	}

	if len(fStyles) > 0 {
		return true, declarationBattle(fStyles, inline)
	}

	return false, []Declaration{}
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
