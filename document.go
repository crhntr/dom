package dom

import (
	"golang.org/x/net/html"
)

type DocumentHTMLNode struct {
	node *html.Node
}

// Node

func (d *DocumentHTMLNode) NodeType() NodeType         { return nodeType(d.node.Type) }
func (d *DocumentHTMLNode) CloneNode(deep bool) Node   { return cloneNode(d.node, deep) }
func (d *DocumentHTMLNode) IsSameNode(other Node) bool { return isSameNode(d.node, other) }
func (d *DocumentHTMLNode) GetElementsByTagName(name string) ElementCollection {
	return getElementsByTagName(d.node, name)
}
func (d *DocumentHTMLNode) GetElementsByClassName(name string) ElementCollection {
	return getElementsByClassName(d.node, name)
}
func (d *DocumentHTMLNode) QuerySelector(query string) Element { return querySelector(d.node, query) }
func (d *DocumentHTMLNode) QuerySelectorAll(query string) NodeList {
	return querySelectorAll(d.node, query)
}
func (d *DocumentHTMLNode) Contains(other Node) bool { return contains(d.node, other) }

func (d *DocumentHTMLNode) TextContent() string { return textContent(d.node) }

// Document

func (d *DocumentHTMLNode) CreateElement(localName string) Element {
	//TODO implement me
	panic("implement me")
}

func (d *DocumentHTMLNode) CreateElementIs(localName, is string) Element {
	//TODO implement me
	panic("implement me")
}

func (d *DocumentHTMLNode) CreateTextNode(text string) Text {
	//TODO implement me
	panic("implement me")
}
