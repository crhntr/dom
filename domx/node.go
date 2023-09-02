package domx

import (
	"bytes"
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
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

func NewNode(node *html.Node) dom.Node {
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
			return NewNode(c)
		}
		offset++
		c = c.NextSibling
	}
	return nil
}

func isConnected(node *html.Node) bool {
	p := node.Parent
	for p != nil {
		if p.Type == html.DocumentNode {
			return true
		}
		p = p.Parent
	}
	return false
}

func ownerDocument(node *html.Node) dom.Document {
	p := node.Parent
	for p != nil {
		if p.Type == html.DocumentNode {
			return &Document{node: p}
		}
		p = p.Parent
	}
	return nil
}
func parentNode(node *html.Node) dom.Node           { return NewNode(node.Parent) }
func parentElement(node *html.Node) dom.Element     { return htmlNodeToDomElement(node.Parent) }
func hasChildNodes(node *html.Node) bool            { return node.FirstChild != nil }
func childNodes(node *html.Node) dom.NodeList       { return (*SiblingNodeList)(node.FirstChild) }
func firstChild(node *html.Node) dom.ChildNode      { return htmlNodeToDomChildNode(node.FirstChild) }
func lastChild(node *html.Node) dom.ChildNode       { return htmlNodeToDomChildNode(node.LastChild) }
func previousSibling(node *html.Node) dom.ChildNode { return htmlNodeToDomChildNode(node.PrevSibling) }
func nextSibling(node *html.Node) dom.ChildNode     { return htmlNodeToDomChildNode(node.NextSibling) }

func textContent(node *html.Node) string {
	var buf bytes.Buffer
	recursiveTextContent(&buf, node)
	return buf.String()
}

func recursiveTextContent(sw io.StringWriter, n *html.Node) {
	if n.Type == html.TextNode {
		_, err := sw.WriteString(n.Data)
		if err != nil {
			panic(err)
		}
	}
	c := n.FirstChild
	for c != nil {
		recursiveTextContent(sw, c)
		c = c.NextSibling
	}
}

func cloneNode(node *html.Node, deep bool) *html.Node {
	result := &html.Node{
		Type:      node.Type,
		Namespace: node.Namespace,
		Data:      node.Data,
		DataAtom:  node.DataAtom,
	}
	if deep {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			result.AppendChild(cloneNode(c, true))
		}
	}
	if node.Attr != nil {
		result.Attr = make([]html.Attribute, len(node.Attr))
		for i, at := range node.Attr {
			result.Attr[i].Key = at.Key
			result.Attr[i].Val = at.Val
			result.Attr[i].Namespace = at.Namespace
		}
	}
	return result
}

func isSameNode(node *html.Node, other dom.Node) bool {
	if node == nil || other == nil {
		return false
	}
	n := domNodeToHTMLNode(other)
	return n != nil && node == n
}

func contains(node *html.Node, other dom.Node) bool {
	o := domNodeToHTMLNode(other)
	if o == nil {
		return false
	}

	found := false
	walkNodes(node, func(n *html.Node) bool {
		found = n == o
		return found
	})

	return found
}

func insertBefore(parent *html.Node, node, child dom.ChildNode) dom.ChildNode {
	n := domNodeToHTMLNode(node)
	c := domNodeToHTMLNode(child)
	if n.Parent != nil {
		n.Parent.RemoveChild(n)
	}
	parent.InsertBefore(n, c)
	return htmlNodeToDomChildNode(n)
}

func appendChild(parent *html.Node, node dom.ChildNode) dom.ChildNode {
	n := domNodeToHTMLNode(node)
	if n.Parent != nil {
		n.Parent.RemoveChild(n)
	}
	parent.AppendChild(n)
	return htmlNodeToDomChildNode(n)
}

func replaceChild(parent *html.Node, node, child dom.ChildNode) dom.ChildNode {
	n := domNodeToHTMLNode(node)
	c := domNodeToHTMLNode(child)
	if c.Parent != parent {
		panic("browser: ReplaceChild called for an attached child node")
	}
	if c.PrevSibling != nil {
		c.PrevSibling.NextSibling = n
	}
	if c.NextSibling != nil {
		c.NextSibling.PrevSibling = n
	}
	if parent.FirstChild == c {
		parent.FirstChild = n
	}
	if parent.LastChild == c {
		parent.LastChild = n
	}
	n.PrevSibling = c.PrevSibling
	n.NextSibling = c.NextSibling
	n.Parent = c.Parent

	c.PrevSibling = nil
	c.NextSibling = nil
	c.Parent = nil

	return htmlNodeToDomChildNode(c)
}

func removeChild(parent *html.Node, node dom.ChildNode) dom.ChildNode {
	n := domNodeToHTMLNode(node)
	parent.RemoveChild(n)
	return htmlNodeToDomChildNode(n)
}

func children(parent *html.Node) dom.ElementCollection {
	return SiblingElements{firstChild: parent.FirstChild}
}

func firstElementChild(node *html.Node) dom.Element {
	child := node.FirstChild
	for child != nil {
		if child.Type == html.ElementNode {
			return &Element{node: child}
		}
		child = child.NextSibling
	}
	return nil
}

func lastElementChild(node *html.Node) dom.Element {
	child := node.LastChild
	for child != nil {
		if child.Type == html.ElementNode {
			return &Element{node: child}
		}
		child = child.PrevSibling
	}
	return nil
}

func childElementCount(node *html.Node) int {
	var (
		result = 0
		child  = node.FirstChild
	)
	for child != nil {
		if child.Type == html.ElementNode {
			result++
		}
		child = child.NextSibling
	}
	return result
}

func prependNodes(node *html.Node, nodes []dom.ChildNode) {
	for i := range nodes {
		dn := nodes[len(nodes)-1-i]
		n := domNodeToHTMLNode(dn)

		fc := node.FirstChild
		if fc != nil {
			fc.PrevSibling = n
			n.NextSibling = fc
		}
		n.Parent = node
		node.FirstChild = n

		if node.LastChild == nil {
			node.LastChild = n
		}
	}
}

func appendNodes(parent *html.Node, nodes []dom.ChildNode) {
	for _, node := range nodes {
		n := domNodeToHTMLNode(node)
		parent.AppendChild(n)
	}
}

func replaceChildren(parent *html.Node, nodes []dom.ChildNode) {
	clearChildren(parent)
	for _, node := range nodes {
		n := domNodeToHTMLNode(node)
		parent.AppendChild(n)
	}
}

func clearChildren(node *html.Node) {
	if fc := node.FirstChild; fc != nil {
		fc.Parent = nil
	}
	if lc := node.LastChild; lc != nil {
		lc.Parent = nil
	}
	node.FirstChild = nil
	node.LastChild = nil
}

func getElementsByTagName(node *html.Node, name string) dom.ElementCollection {
	name = strings.ToUpper(name)
	var list ElementList
	walkNodes(node, func(n *html.Node) bool {
		if strings.ToUpper(n.Data) == name {
			list = append(list, n)
		}
		return false
	})
	return list
}

func getElementsByClassName(node *html.Node, name string) dom.ElementCollection {
	var list ElementList
	walkNodes(node, func(n *html.Node) bool {
		if hasClasses(getAttribute(n, "class"), name) {
			list = append(list, n)
		}
		return false
	})
	return list
}

func hasClasses(elementClassesStr, classesStr string) bool {
	elementClasses := strings.Fields(elementClassesStr)
	classes := strings.Fields(classesStr)

	set := make(map[string]struct{}, len(classesStr))
	for _, c := range classes {
		set[c] = struct{}{}
	}

	for _, c := range elementClasses {
		delete(set, c)
	}

	return len(set) == 0
}

func querySelector(node *html.Node, query string) dom.Element {
	result := cascadia.Query(node, cascadia.MustCompile(query))
	if result == nil {
		return nil
	}
	return &Element{node: result}
}

func querySelectorAll(node *html.Node, query string) dom.NodeList {
	return NodeListHTMLElements(cascadia.QueryAll(node, cascadia.MustCompile(query)))
}

var _ dom.NodeList = NodeListHTMLElements(nil)

type NodeListHTMLElements []*html.Node

func (n NodeListHTMLElements) Length() int { return len(n) }

func (n NodeListHTMLElements) Item(i int) dom.Node { return NewNode(n[i]) }

func closest(node *html.Node, selector string) dom.Element {
	s := cascadia.MustCompile(selector)
	for p := node; p != nil; p = p.Parent {
		if s.Match(p) {
			return htmlNodeToDomElement(p)
		}
	}
	return nil
}

func matches(node *html.Node, selector string) bool {
	s := cascadia.MustCompile(selector)
	return s.Match(node)
}

func isNamed(node *html.Node, name string) bool {
	if node == nil || node.Type != html.ElementNode {
		return false
	}
	id := getAttribute(node, "id")
	nm := getAttribute(node, "name")
	return (id != "" && id == name) || (nm != "" && nm == name)
}

func getAttribute(node *html.Node, name string) string {
	name = strings.ToLower(name)
	for _, att := range node.Attr {
		if att.Key == name {
			return att.Val
		}
	}
	return ""
}
