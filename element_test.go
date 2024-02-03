package dom

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/crhntr/dom/spec"
)

func TestElement_NodeType(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	document := Element{
		node: parsedDocument.FirstChild.NextSibling,
	}

	assert.Equal(t, spec.NodeTypeElement, document.NodeType())
}

func TestElement_IsConnected(t *testing.T) {
	t.Run("connected", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body>Hello!</body>
</html>`
		_, body := parseDocument(t, textHTML, "body")

		require.True(t, body.IsConnected())
	})

	t.Run("not connected", func(t *testing.T) {
		el := &Element{
			node: &html.Node{},
		}
		require.False(t, el.IsConnected())
	})
}

func TestElement_OwnerDocument(t *testing.T) {
	t.Run("connected", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body>Hello!</body>
</html>`
		document, body := parseDocument(t, textHTML, "body")

		owner := body.OwnerDocument()
		require.NotNil(t, owner)

		require.True(t, owner.(*Document).node == document.node)
	})
	t.Run("not connected", func(t *testing.T) {
		textNode := &Element{
			node: &html.Node{},
		}
		require.Nil(t, textNode.OwnerDocument())
	})
}

func TestElement_Length(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body>
	<div id="target"><br id="1">2<div id="3"><span id="does-not-count"></span></div></div>
</body>
</html>`
	_, body := parseDocument(t, textHTML, "#target")
	assert.Equal(t, 3, body.Length())
}

func TestElement_ParentNode(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><div id="app"></div></body>
</html>`
	_, app := parseDocument(t, textHTML, "html")

	parent := app.ParentNode()
	require.NotNil(t, parent)
	require.Equal(t, spec.NodeTypeDocument, parent.NodeType())
	_, ok := parent.(spec.Document)
	require.True(t, ok)
}

func TestElement_ParentElement(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><div id="app"></div></body>
</html>`
	_, app := parseDocument(t, textHTML, "#app")

	parent, ok := app.ParentElement().(*Element)
	require.True(t, ok)
	assert.Equal(t, "body", parent.node.Data)
}

func TestElement_PreviousSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><span id="oldest"></span><div id="target"></div></body>
</html>`
	_, target := parseDocument(t, textHTML, "#target")
	sibling := target.PreviousSibling()
	siblingNode, ok := sibling.(*Element)
	require.Truef(t, ok, "wrong type %T", sibling)

	require.Equal(t, "oldest", getAttribute(siblingNode.node, "id"))
}

func TestElement_NextSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><div id="target"></div><span id="youngest"></span></body>
</html>`
	_, target := parseDocument(t, textHTML, "#target")
	sibling := target.NextSibling()
	siblingNode, ok := sibling.(*Element)
	require.Truef(t, ok, "wrong type %T", sibling)

	require.Equal(t, "youngest", getAttribute(siblingNode.node, "id"))
}

func TestElement_TextContent(t *testing.T) {
	t.Run("body no text", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title>Greetings!</title></head>
<body><div id="target"><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "body")
		text := target.TextContent()
		require.Equal(t, "\n", text)
	})

	t.Run("element no text", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title>Greetings!</title></head>
<body><div id="target"><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		text := target.TextContent()
		require.Equal(t, "", text)
	})

	t.Run("multiple lines with whitespace", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title>Greetings!</title></head>
<body>
<div id="target">a
b<span>c </span>d
	e
</div>
</body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		text := target.TextContent()
		require.Equal(t, `a
bc d
	e
`, text)
	})
}

func TestElement_CloneNode(t *testing.T) {
	t.Run("not deep", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><div id="target"><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		clonedElement, ok := target.CloneNode(false).(*Element)
		require.True(t, ok)
		require.Nil(t, clonedElement.node.FirstChild)
		require.Nil(t, clonedElement.node.Parent)
	})

	t.Run("deep", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><div id="target"><span></span></div></body>
</html>`
		_, target := parseDocument(t, textHTML, "#target")
		clonedElement, ok := target.CloneNode(true).(*Element)
		require.True(t, ok)
		require.NotNil(t, clonedElement.node.FirstChild)
		require.Nil(t, clonedElement.node.Parent)
		require.EqualValues(t, clonedElement.node, clonedElement.node.FirstChild.Parent)
	})
}

func TestElement_IsSameNode(t *testing.T) {
	t.Run("different nodes", func(t *testing.T) {
		e1 := &Element{node: &html.Node{Type: html.ElementNode}}
		e2 := &Element{node: &html.Node{Type: html.ElementNode}}
		assert.False(t, e1.IsSameNode(e2))
	})
	t.Run("receiver node is nil", func(t *testing.T) {
		e1 := &Element{node: nil}
		e2 := &Element{node: &html.Node{Type: html.ElementNode}}
		assert.False(t, e1.IsSameNode(e2))
	})
	t.Run("param node is nil", func(t *testing.T) {
		e1 := &Element{node: &html.Node{Type: html.ElementNode}}
		e2 := &Element{node: nil}
		assert.False(t, e1.IsSameNode(e2))
	})
	t.Run("param is nil", func(t *testing.T) {
		e1 := &Element{node: &html.Node{Type: html.ElementNode}}
		assert.False(t, e1.IsSameNode(nil))
	})
	t.Run("true", func(t *testing.T) {
		node := &html.Node{Type: html.ElementNode}
		e1 := &Element{node: node}
		e2 := &Element{node: node}
		assert.True(t, e1.IsSameNode(e2))
	})
}
