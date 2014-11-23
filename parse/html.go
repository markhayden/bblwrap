package parse

import (
	"fmt"
	"github.com/markhayden/bblwrap/util"
	"regexp"
	"strings"
)

type Html struct {
	Master []string
}

func parseHtml() {
	dummy := `<div class="bonk"><span><img src="w3schools.jpg" alt="W3Schools.com" width="104" height="142"></span><div id="taco"><span><div class="idaho">here</div></span></div></div>`
	r, _ := regexp.Compile(`<(\/)?\b[^>]*>`)
	tags := r.FindAllString(dummy, -1)
	var m Html

	fmt.Println(tags)
	for _, x := range tags {
		m.whatClass(x)
		fmt.Println(m.Master)
		// this is where i would check a style match
	}
}

func (m *Html) whatClass(tag string) {
	// last key
	last := len(m.Master) - 1

	// handle previous img tag
	if last > -1 {
		imgMatch, _ := regexp.Compile(`<img`)
		if imgMatch.Match([]byte(m.Master[last])) {
			fmt.Println("  ** closing image")
			m.Master = m.Master[:last]
		}
	}

	// handle closing tags
	closeMatch, _ := regexp.Compile(`\/`)
	if closeMatch.Match([]byte(tag)) {
		fmt.Println(fmt.Sprintf("  ** closing tag %s", tag))
		m.Master = m.Master[:last]
		return
	}

	// new tag found
	m.Master = append(m.Master, tag) // save to master for navigation
	style := renderStyle(tag)        // determine the style for tag

	fmt.Println(len(style))
	fmt.Println(util.PrettyJson(style))

	return
}

func renderStyle(tag string) []Selector {
	// for each element determine what can fulfill it, bottom to top
	// loop through that slice looking for matches in styles and store temporary style output
	// on complete set the styles and move to next element
	// <div class="taco bell" id="stuff" hup="test">
	// <div
	// create 4 slices. loop through element attributes and store class to one, id to one, attributes to one and style to one
	// rejoin all 4 slices in order class, id, attr, style so that we can handle priority order
	// loop through each looking for a match
	// save style

	r, _ := regexp.Compile(`(?P<declaration>[a-zA-Z0-9\-\_]+)="(?P<value>[a-zA-Z0-9\-\_\/ ]+)"`)

	var master []Selector
	var class []Selector
	var id []Selector
	var style []Selector
	var attr []Selector

	needles := r.FindAllStringSubmatch(tag, -1)
	for _, v := range needles {
		if len(v) != 3 {
			fmt.Println("invalid class, declaration")
		}

		declaration := v[1]
		values := strings.Split(v[2], " ")

		for _, value := range values {
			s := Selector{
				Origin: v[0],
				Key:    strings.Split(v[0], "=")[0],
				Value:  strings.TrimSpace(value),
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
				style = append(style, s)
			default:
				s.Type = "attr"
				attr = append(attr, s)
			}
		}

	}

	// merge it all back togehter into one master declaration slice
	master = append(master, class...)
	master = append(master, id...)
	master = append(master, attr...)
	master = append(master, style...)

	return master
}
