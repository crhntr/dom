package domx

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom"
)

type Document struct {
	node *html.Node
}

func (d *Document) NodeType() dom.NodeType         { return nodeType(d.node.Type) }
func (d *Document) CloneNode(deep bool) dom.Node   { return NewNode(cloneNode(d.node, deep)) }
func (d *Document) IsSameNode(other dom.Node) bool { return isSameNode(d.node, other) }
func (d *Document) GetElementsByTagName(name string) dom.ElementCollection {
	return getElementsByTagName(d.node, name)
}
func (d *Document) GetElementsByClassName(name string) dom.ElementCollection {
	return getElementsByClassName(d.node, name)
}
func (d *Document) QuerySelector(query string) dom.Element {
	return querySelector(d.node, query)
}
func (d *Document) QuerySelectorAll(query string) dom.NodeList {
	return querySelectorAll(d.node, query)
}
func (d *Document) Contains(other dom.Node) bool { return contains(d.node, other) }

// TextContent returns an empty string.
// The spec says it should return null
// https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent
func (d *Document) TextContent() string { return "" }

// Document

func (*Document) CreateElement(localName string) dom.Element {
	localName = strings.ToLower(localName)
	return &Element{
		node: &html.Node{
			DataAtom: atom.Lookup([]byte(localName)),
			Type:     html.ElementNode,
			Data:     localName,
		},
	}
}

func (*Document) CreateElementIs(localName, is string) dom.Element {
	localName = strings.ToLower(localName)
	return &Element{
		node: &html.Node{
			DataAtom: atom.Lookup([]byte(localName)),
			Type:     html.ElementNode,
			Data:     localName,
			Attr:     []html.Attribute{{Key: "is", Val: is}},
		},
	}
}

func (*Document) CreateTextNode(text string) dom.Text {
	return &Text{
		node: &html.Node{
			Type: html.TextNode,
			Data: text,
		},
	}
}
