//go:build js

package browser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crhntr/dom/browser"
	"github.com/crhntr/dom/spec"
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
