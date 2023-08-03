package domx

import (
	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

type SiblingElements struct {
	firstChild *html.Node
}

func (list SiblingElements) Length() int {
	result := 0
	for c := list.firstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		result++
	}
	return result
}

func (list SiblingElements) Item(index int) dom.Element {
	childIndex := 0
	for c := list.firstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if childIndex == index {
			return &Element{node: c}
		}
		childIndex++
	}
	return nil
}

func (list SiblingElements) NamedItem(name string) dom.Element {
	for c := list.firstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if isNamed(c, name) {
			return &Element{node: c}
		}
	}
	return nil
}

type ElementList []*html.Node

func (list ElementList) Length() int { return len(list) }

func (list ElementList) Item(index int) dom.Element {
	if index < 0 || index >= len(list) {
		return nil
	}
	return &Element{node: list[index]}
}

func (list ElementList) NamedItem(name string) dom.Element {
	for _, el := range list {
		if isNamed(el, name) {
			return &Element{node: el}
		}
	}
	return nil
}
