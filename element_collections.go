package dom

import (
	"golang.org/x/net/html"
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

func (list SiblingElements) Item(index int) Element {
	childIndex := 0
	for c := list.firstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if childIndex == index {
			return &ElementHTMLNode{node: c}
		}
		childIndex++
	}
	return nil
}

func (list SiblingElements) NamedItem(name string) Element {
	for c := list.firstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if isNamed(c, name) {
			return &ElementHTMLNode{node: c}
		}
	}
	return nil
}

type ElementList []*html.Node

func (list ElementList) Length() int { return len(list) }

func (list ElementList) Item(index int) Element {
	if index < 0 || index >= len(list) {
		return nil
	}
	return &ElementHTMLNode{node: list[index]}
}

func (list ElementList) NamedItem(name string) Element {
	for _, el := range list {
		if isNamed(el, name) {
			return &ElementHTMLNode{node: el}
		}
	}
	return nil
}
