package domx

import (
	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

func nodeType(nodeType html.NodeType) dom.NodeType {
	switch nodeType {
	case html.TextNode:
		return dom.NodeTypeText
	case html.DocumentNode:
		return dom.NodeTypeDocument
	case html.ElementNode:
		return dom.NodeTypeElement
	case html.CommentNode:
		return dom.NodeTypeComment
	case html.DoctypeNode:
		return dom.NodeTypeDocumentType
	default:
		fallthrough
	case html.ErrorNode, html.RawNode:
		return dom.NodeTypeUnknown
	}
}

func htmlNodeToDomNode(node *html.Node) dom.Node {
	if node == nil {
		return nil
	}
	switch node.Type {
	case html.ElementNode:
		return &Element{node: node}
	case html.TextNode:
		return &Text{node: node}
	case html.DocumentNode:
		return &Document{node: node}
	default:
		panic("not supported")
	}
}

func htmlNodeToDomChildNode(node *html.Node) dom.ChildNode {
	if node == nil {
		return nil
	}
	switch node.Type {
	case html.ElementNode:
		return &Element{node: node}
	case html.TextNode:
		return &Text{node: node}
	default:
		panic("not supported")
	}
}

func htmlNodeToDomElement(node *html.Node) dom.Element {
	if node == nil {
		return nil
	}
	return &Element{node: node}
}

func domNodeToHTMLNode(node dom.Node) *html.Node {
	switch ot := node.(type) {
	case *Element:
		return ot.node
	case *Text:
		return ot.node
	case *Document:
		return ot.node
	default:
		panic("not implemented")
	}
}

func walkNodes(start *html.Node, fn func(node *html.Node) (done bool)) bool {
	if fn(start) {
		return true
	}

	c := start.FirstChild
	for c != nil {
		if walkNodes(c, fn) {
			return true
		}
		c = c.NextSibling
	}

	return false
}

type SiblingNodeList html.Node

func (node *SiblingNodeList) Length() int {
	c := (*html.Node)(node)
	result := 0
	for c != nil {
		result++
		c = c.NextSibling
	}
	return result
}

func (node *SiblingNodeList) Item(index int) dom.Node {
	c := (*html.Node)(node)
	offset := 0
	for c != nil {
		if offset == index {
			return htmlNodeToDomNode(c)
		}
		offset++
		c = c.NextSibling
	}
	return nil
}
