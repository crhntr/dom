//go:build js

package browser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/typelate/dom/browser"
	"github.com/typelate/dom/spec"
)

func TestDocument(t *testing.T) {
	document := browser.OpenDocument()
	div := document.CreateElement("div")

	head := document.Head()
	require.NotNil(t, head)

	body := document.Body()
	require.NotNil(t, body)

	const childClass = "child"

	childOne := document.CreateElement("div")
	childOne.SetAttribute("id", "one")
	childOne.SetAttribute("class", childClass)
	assert.NotNil(t, body.AppendChild(childOne))

	childTwo := document.CreateElement("div")
	childTwo.SetAttribute("id", "two")
	childTwo.SetAttribute("class", childClass)
	assert.NotNil(t, body.AppendChild(childTwo))

	t.Run("NodeType", func(t *testing.T) {
		assert.Equal(t, spec.NodeTypeDocument, document.NodeType())
	})
	t.Run("CloneNode", func(t *testing.T) {
		assert.NotPanics(t, func() { document.CloneNode(false) })
		assert.NotPanics(t, func() { document.CloneNode(true) })
	})
	t.Run("IsSameNode", func(t *testing.T) {
		assert.True(t, document.IsSameNode(document))
		assert.False(t, document.IsSameNode(div))
	})
	t.Run("TextContent", func(t *testing.T) {
		assert.NotPanics(t, func() { document.TextContent() })
	})
	t.Run("Contains", func(t *testing.T) {
		assert.True(t, document.Contains(head))
		assert.False(t, document.Contains(div))
	})
	t.Run("GetElementsByTagName", func(t *testing.T) {
		assert.NotZero(t, document.GetElementsByTagName("div").Length())
	})
	t.Run("GetElementsByClassName", func(t *testing.T) {
		assert.Equal(t, 2, document.GetElementsByClassName(childClass).Length())
	})
}

func TestElement_CompareDocumentPosition(t *testing.T) {
	document := browser.OpenDocument()

	a := document.CreateElement("div")
	a.SetAttribute("id", "a")

	b := document.CreateElement("span")
	b.SetAttribute("id", "b")
	a.AppendChild(b)

	c := document.CreateElement("div")
	c.SetAttribute("id", "c")

	document.Body().AppendChild(a)
	document.Body().AppendChild(c)

	t.Run("same node", func(t *testing.T) {
		assert.Equal(t, spec.DocumentPosition(0), a.CompareDocumentPosition(a))
	})

	t.Run("contains", func(t *testing.T) {
		pos := a.CompareDocumentPosition(b)
		assert.Equal(t, spec.DocumentPositionContainedBy|spec.DocumentPositionFollowing, pos)
	})

	t.Run("contained by", func(t *testing.T) {
		pos := b.CompareDocumentPosition(a)
		assert.Equal(t, spec.DocumentPositionContains|spec.DocumentPositionPreceding, pos)
	})

	t.Run("preceding", func(t *testing.T) {
		pos := a.CompareDocumentPosition(c)
		assert.Equal(t, spec.DocumentPositionFollowing, pos)
	})

	t.Run("following", func(t *testing.T) {
		pos := c.CompareDocumentPosition(a)
		assert.Equal(t, spec.DocumentPositionPreceding, pos)
	})

	t.Run("disconnected", func(t *testing.T) {
		d := document.CreateElement("div")
		// not appended to the DOM
		pos := a.CompareDocumentPosition(d)
		assert.True(t, spec.DocumentPositionDisconnected&pos != 0)
		assert.True(t, spec.DocumentPositionImplementationSpecific&pos != 0)
	})
}
