package domx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

func TestDocument_NodeType(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	document := Document{
		node: parsedDocument,
	}

	assert.Equal(t, dom.NodeTypeDocument, document.NodeType())
}

func TestDocument_CloneNode(t *testing.T) {
	t.Run("deep", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)

		document := Document{node: parsedDocument}
		document.CloneNode(false)
	})
	t.Run("not deep", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)

		document := Document{node: parsedDocument}
		document.CloneNode(false)
	})
}
func TestDocument_IsSameNode(t *testing.T) {
	t.Run("same", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)

		document := &Document{node: parsedDocument}
		require.True(t, document.IsSameNode(document))

		parsedDocument2, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		require.False(t, document.IsSameNode(&Document{node: parsedDocument2}))
	})

	t.Run("same html different address", func(t *testing.T) {
		in := `<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`
		parsedDocument1, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)
		parsedDocument2, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: parsedDocument1}).IsSameNode(&Document{node: parsedDocument2}))
	})

	t.Run("receiver node is nil", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`
		parsedDocument, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: nil}).IsSameNode(&Document{node: parsedDocument}))
	})

	t.Run("param is nil", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`
		parsedDocument, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: parsedDocument}).IsSameNode(nil))
	})

	t.Run("param node is nil", func(t *testing.T) {
		// language=html
		in := `<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`
		parsedDocument, err := html.Parse(strings.NewReader(in))
		require.NoError(t, err)

		require.False(t, (&Document{node: parsedDocument}).IsSameNode(&Document{node: nil}))
	})
}

//func TestDocument_GetElementsByTagName(t *testing.T) {
//	document := Document{}
//	document.GetElementsByTagName()
//}
//func TestDocument_GetElementsByClassName(t *testing.T) {
//	document := Document{}
//	document.GetElementsByClassName()
//}
//func TestDocument_QuerySelector(t *testing.T) {
//	document := Document{}
//	document.QuerySelector()
//}
//func TestDocument_QuerySelectorAll(t *testing.T) {
//	document := Document{}
//	document.QuerySelectorAll()
//}
//func TestDocument_Contains(t *testing.T) {
//	document := Document{}
//	document.Contains()
//}
//func TestDocument_TextContent(t *testing.T) {
//	document := Document{}
//	document.TextContent()
//}

func TestDocument_CreateElement(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.LastChild.FirstChild
	exp.Parent = nil

	var document *Document
	got := document.CreateElement("SPAN").(*Element).node

	assert.Equal(t, exp, got)
}

func TestDocument_CreateElementIs(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><div is="fruit"></div></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.LastChild.FirstChild
	exp.Parent = nil

	var document *Document
	got := document.CreateElementIs("div", "fruit").(*Element).node

	assert.Equal(t, exp, got)
}

func TestDocument_CreateTextNode(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title>peach</title></head><body></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild
	exp.Parent = nil

	var document *Document
	got := document.CreateTextNode("peach").(*Text).node

	assert.Equal(t, exp, got)
}
