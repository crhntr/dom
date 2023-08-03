package dom

import (
	"golang.org/x/net/html"
)

type TextHTMLNode struct {
	node *html.Node
}

func (t *TextHTMLNode) Data() string     { return t.node.Data }
func (t *TextHTMLNode) SetData(d string) { t.node.Data = d }

func (t *TextHTMLNode) NodeType() NodeType         { return nodeType(t.node.Type) }
func (t *TextHTMLNode) IsConnected() bool          { return isConnected(t.node) }
func (t *TextHTMLNode) OwnerDocument() Document    { return ownerDocument(t.node) }
func (t *TextHTMLNode) Length() int                { return len(t.node.Data) }
func (t *TextHTMLNode) ParentNode() Node           { return parentNode(t.node) }
func (t *TextHTMLNode) ParentElement() Element     { return parentElement(t.node) }
func (t *TextHTMLNode) PreviousSibling() ChildNode { return previousSibling(t.node) }
func (t *TextHTMLNode) NextSibling() ChildNode     { return nextSibling(t.node) }
func (t *TextHTMLNode) TextContent() string        { return t.node.Data }
func (t *TextHTMLNode) CloneNode(_ bool) Node {
	return &TextHTMLNode{
		node: &html.Node{
			Type: html.TextNode,
			Data: t.node.Data,
		},
	}
}
func (t *TextHTMLNode) IsSameNode(other Node) bool { return isSameNode(t.node, other) }
