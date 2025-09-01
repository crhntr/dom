package dom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/typelate/dom/spec"
)

func TestText_Data(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><span>Hello!</span></body>
</html>`
	_, bodyNode := parseDocument(t, textHTML, "body span")
	textNode := &Text{
		node: bodyNode.node.FirstChild,
	}

	require.Equal(t, "Hello!", textNode.Data())
	textNode.SetData("Greetings!")
	require.Equal(t, "Greetings!", textNode.Data())
}

func TestText_NodeType(t *testing.T) {
	t.Run("cromulent", func(t *testing.T) {
		textNode := &Text{
			node: &html.Node{Type: html.TextNode},
		}

		assert.Equal(t, spec.NodeTypeText, textNode.NodeType())
	})

	t.Run("mismatched node type", func(t *testing.T) {
		textNode := &Text{
			node: &html.Node{Type: html.DocumentNode},
		}

		assert.Equal(t, spec.NodeTypeDocument, textNode.NodeType())
	})
}

func TestText_IsConnected(t *testing.T) {
	t.Run("connected", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body>Hello!</body>
</html>`
		_, bodyNode := parseDocument(t, textHTML, "body")
		textNode := &Text{
			node: bodyNode.node.FirstChild,
		}

		require.True(t, textNode.IsConnected())
	})

	t.Run("not connected", func(t *testing.T) {
		textNode := &Text{
			node: &html.Node{},
		}
		require.False(t, textNode.IsConnected())
	})
}

func TestText_OwnerDocument(t *testing.T) {
	t.Run("connected", func(t *testing.T) {
		// language=html
		textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body>Hello!</body>
</html>`
		document, bodyNode := parseDocument(t, textHTML, "body")
		textNode := &Text{
			node: bodyNode.node.FirstChild,
		}

		owner := textNode.OwnerDocument()
		require.NotNil(t, owner)

		require.True(t, owner.(*Document).node == document.node)
	})
	t.Run("not connected", func(t *testing.T) {
		textNode := &Text{
			node: &html.Node{},
		}
		require.Nil(t, textNode.OwnerDocument())
	})
}

func TestText_Length(t *testing.T) {
	textNode := &Text{
		node: &html.Node{
			Data: "Hello!",
		},
	}

	require.Equal(t, textNode.Length(), len("Hello!"))
}

func TestText_ParentNode(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body id="app">Hello!</body>
</html>`
	_, bodyNode := parseDocument(t, textHTML, "body")
	textNode := &Text{
		node: bodyNode.node.FirstChild,
	}
	parent, ok := textNode.ParentNode().(*Element)
	require.True(t, ok)
	require.True(t, parent.node == textNode.node.Parent)
}

func TestText_ParentElement(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body id="app">Hello!</body>
</html>`
	_, bodyNode := parseDocument(t, textHTML, "body")
	textNode := &Text{
		node: bodyNode.node.FirstChild,
	}
	parent, ok := textNode.ParentElement().(*Element)
	require.True(t, ok)
	require.True(t, parent.node == textNode.node.Parent)
}

func TestText_PreviousSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><span id="oldest"></span>Hello!</body>
</html>`
	_, bodyNode := parseDocument(t, textHTML, "body")
	textNode := &Text{
		node: bodyNode.node.LastChild,
	}
	sibling := textNode.PreviousSibling()
	siblingNode, ok := sibling.(*Element)
	require.Truef(t, ok, "wrong type %T", sibling)

	require.Equal(t, "oldest", getAttribute(siblingNode.node, "id"))
}

func TestText_NextSibling(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body>Hello!<span id="youngest"></span></body>
</html>`
	_, bodyNode := parseDocument(t, textHTML, "body")
	textNode := &Text{
		node: bodyNode.node.FirstChild,
	}
	sibling := textNode.NextSibling()
	siblingNode, ok := sibling.(*Element)
	require.Truef(t, ok, "wrong type %T", sibling)

	require.Equal(t, "youngest", getAttribute(siblingNode.node, "id"))
}

func TestText_TextContent(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><span>Hello!</span></body>
</html>`
	_, span := parseDocument(t, textHTML, "body span")
	textNode := &Text{
		node: span.node.FirstChild,
	}
	assert.Equal(t, "Hello!", textNode.TextContent())
}

func TestText_CloneNode(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><p>Hello!</p></body>
</html>`
	_, p := parseDocument(t, textHTML, "body p")
	textNode := &Text{
		node: p.node.FirstChild,
	}

	cloned := textNode.CloneNode(true)

	require.Equal(t, "Hello!", cloned.(*Text).node.Data)

	require.Nil(t, cloned.(*Text).node.Parent)
}

func TestText_IsSameNode(t *testing.T) {
	// language=html
	textHTML := `<!DOCTYPE html>
<html lang="us-en">
<head><title></title></head>
<body><p>Hello!</p></body>
</html>`
	_, p := parseDocument(t, textHTML, "body p")
	textNode := &Text{
		node: p.node.FirstChild,
	}

	require.True(t, textNode.IsSameNode(&Text{
		node: p.node.FirstChild,
	}))
}
