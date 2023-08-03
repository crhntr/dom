package dom

import (
	"golang.org/x/net/html"
)

func nodeType(nodeType html.NodeType) NodeType {
	switch nodeType {
	case html.TextNode:
		return NodeTypeText
	case html.DocumentNode:
		return NodeTypeDocument
	case html.ElementNode:
		return NodeTypeElement
	case html.CommentNode:
		return NodeTypeComment
	case html.DoctypeNode:
		return NodeTypeDocumentType
	default:
		fallthrough
	case html.ErrorNode, html.RawNode:
		return NodeTypeUnknown
	}
}

func htmlNodeToDomNode(node *html.Node) Node {
	if node == nil {
		return nil
	}
	switch node.Type {
	case html.ElementNode:
		return &ElementHTMLNode{node: node}
	case html.TextNode:
		return &TextHTMLNode{node: node}
	case html.DocumentNode:
		return &DocumentHTMLNode{node: node}
	default:
		panic("not supported")
	}
}

func htmlNodeToDomChildNode(node *html.Node) ChildNode {
	if node == nil {
		return nil
	}
	switch node.Type {
	case html.ElementNode:
		return &ElementHTMLNode{node: node}
	case html.TextNode:
		return &TextHTMLNode{node: node}
	default:
		panic("not supported")
	}
}

func htmlNodeToDomElement(node *html.Node) Element {
	if node == nil {
		return nil
	}
	return &ElementHTMLNode{node: node}
}

func domNodeToHTMLNode(node Node) *html.Node {
	switch ot := node.(type) {
	case *ElementHTMLNode:
		return ot.node
	case *TextHTMLNode:
		return ot.node
	case *DocumentHTMLNode:
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

func (node *SiblingNodeList) Item(index int) Node {
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
