package dom

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
)

type ElementHTMLNode struct {
	node *html.Node
}

// Node

func (e *ElementHTMLNode) NodeType() NodeType         { return nodeType(e.node.Type) }
func (e *ElementHTMLNode) IsConnected() bool          { return isConnected(e.node) }
func (e *ElementHTMLNode) OwnerDocument() Document    { return ownerDocument(e.node) }
func (e *ElementHTMLNode) ParentNode() Node           { return parentNode(e.node) }
func (e *ElementHTMLNode) ParentElement() Element     { return parentElement(e.node) }
func (e *ElementHTMLNode) PreviousSibling() ChildNode { return previousSibling(e.node) }
func (e *ElementHTMLNode) NextSibling() ChildNode     { return nextSibling(e.node) }
func (e *ElementHTMLNode) TextContent() string        { return textContent(e.node) }
func (e *ElementHTMLNode) CloneNode(deep bool) Node   { return cloneNode(e.node, deep) }
func (e *ElementHTMLNode) IsSameNode(other Node) bool { return isSameNode(e.node, other) }
func (e *ElementHTMLNode) Length() int {
	c := e.node
	result := 0
	for c != nil {
		result++
		c = c.NextSibling
	}
	return result
}

// ParentNode

func (e *ElementHTMLNode) Children() ElementCollection        { return children(e.node) }
func (e *ElementHTMLNode) FirstElementChild() Element         { return firstElementChild(e.node) }
func (e *ElementHTMLNode) LastElementChild() Element          { return lastElementChild(e.node) }
func (e *ElementHTMLNode) ChildElementCount() int             { return childElementCount(e.node) }
func (e *ElementHTMLNode) Prepend(nodes ...ChildNode)         { prependNodes(e.node, nodes) }
func (e *ElementHTMLNode) Append(nodes ...ChildNode)          { appendNodes(e.node, nodes) }
func (e *ElementHTMLNode) ReplaceChildren(nodes ...ChildNode) { replaceChildren(e.node, nodes) }
func (e *ElementHTMLNode) GetElementsByTagName(name string) ElementCollection {
	return getElementsByTagName(e.node, name)
}
func (e *ElementHTMLNode) GetElementsByClassName(name string) ElementCollection {
	return getElementsByClassName(e.node, name)
}

func (e *ElementHTMLNode) QuerySelector(query string) Element { return querySelector(e.node, query) }
func (e *ElementHTMLNode) QuerySelectorAll(query string) NodeList {
	return querySelectorAll(e.node, query)
}
func (e *ElementHTMLNode) Closest(selector string) Element { return closest(e.node, selector) }
func (e *ElementHTMLNode) Matches(selector string) bool    { return matches(e.node, selector) }

func (e *ElementHTMLNode) HasChildNodes() bool      { return hasChildNodes(e.node) }
func (e *ElementHTMLNode) ChildNodes() NodeList     { return childNodes(e.node) }
func (e *ElementHTMLNode) FirstChild() ChildNode    { return firstChild(e.node) }
func (e *ElementHTMLNode) LastChild() ChildNode     { return lastChild(e.node) }
func (e *ElementHTMLNode) Contains(other Node) bool { return contains(e.node, other) }
func (e *ElementHTMLNode) InsertBefore(node, child ChildNode) ChildNode {
	return insertBefore(e.node, node, child)
}
func (e *ElementHTMLNode) AppendChild(node ChildNode) ChildNode { return appendChild(e.node, node) }
func (e *ElementHTMLNode) ReplaceChild(node, child ChildNode) ChildNode {
	return replaceChild(e.node, node, child)
}
func (e *ElementHTMLNode) RemoveChild(node ChildNode) ChildNode { return removeChild(e.node, node) }

// Element

func (e *ElementHTMLNode) TagName() string                 { return strings.ToUpper(e.node.Data) }
func (e *ElementHTMLNode) ID() string                      { return getAttribute(e.node, "id") }
func (e *ElementHTMLNode) ClassName() string               { return getAttribute(e.node, "class") }
func (e *ElementHTMLNode) GetAttribute(name string) string { return getAttribute(e.node, name) }

func getAttribute(node *html.Node, name string) string {
	name = strings.ToLower(name)
	for _, att := range node.Attr {
		if att.Key == name {
			return att.Val
		}
	}
	return ""
}

func (e *ElementHTMLNode) SetAttribute(name, value string) {
	name = strings.ToLower(name)
	for index, att := range e.node.Attr {
		if att.Key == name {
			e.node.Attr[index].Val = value
		}
	}
	e.node.Attr = append(e.node.Attr, html.Attribute{
		Key: name, Val: value,
	})
}

func (e *ElementHTMLNode) RemoveAttribute(name string) {
	name = strings.ToLower(name)
	filtered := e.node.Attr[:0]
	for _, att := range e.node.Attr {
		if att.Key == name {
			continue
		}
		filtered = append(filtered, att)
	}
	e.node.Attr = filtered
}

func (e *ElementHTMLNode) ToggleAttribute(name string) bool {
	name = strings.ToLower(name)
	if e.HasAttribute(name) {
		e.RemoveAttribute(name)
		return false
	}
	e.SetAttribute(name, "")
	return true
}

func (e *ElementHTMLNode) HasAttribute(name string) bool {
	name = strings.ToLower(name)
	for _, att := range e.node.Attr {
		if att.Key == name {
			return true
		}
	}
	return false
}

func (e *ElementHTMLNode) isNamed(name string) bool {
	return isNamed(e.node, name)
}

func isNamed(node *html.Node, name string) bool {
	if node == nil || node.Type != html.ElementNode {
		return false
	}
	id := getAttribute(node, "id")
	nm := getAttribute(node, "name")
	return (id != "" && id == name) || (nm != "" && nm == name)
}

func (e *ElementHTMLNode) SetInnerHTML(s string) {
	nodes, err := html.ParseFragment(strings.NewReader(s), &html.Node{Type: html.ElementNode})
	if err != nil {
		panic(err)
	}
	clearChildren(e.node)
	for _, n := range nodes {
		e.node.AppendChild(n)
	}
}

func (e *ElementHTMLNode) InnerHTML() string {
	var buf bytes.Buffer
	c := e.node.FirstChild
	for c != nil {
		err := html.Render(&buf, c)
		if err != nil {
			panic(err)
		}
		c = c.NextSibling
	}
	return buf.String()
}

func (e *ElementHTMLNode) SetOuterHTML(s string) {
	nodes, err := html.ParseFragment(strings.NewReader(s), &html.Node{Type: html.ElementNode})
	if err != nil {
		panic(err)
	}
	if len(nodes) == 0 {
		return
	}
	if e.node.Parent == nil {
		panic("browser: SetOuterHTML called on an unattached node")
	}
	for _, node := range nodes {
		e.node.Parent.InsertBefore(node, e.node)
	}
	e.node.Parent.RemoveChild(e.node)
}

func (e *ElementHTMLNode) OuterHTML() string {
	var buf bytes.Buffer
	err := html.Render(&buf, e.node)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
