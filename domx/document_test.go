package domx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestDocumentHTMLNode_CreateTextNode(t *testing.T) {
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html><head><title>peach</title></head><body></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild
	exp.Parent = nil

	var document *Document
	got := document.CreateTextNode("peach").(*Text).node

	assert.Equal(t, exp, got)
}

func TestDocumentHTMLNode_CreateElementIs(t *testing.T) {
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html><head><title></title></head><body><div is="fruit"></div></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.LastChild.FirstChild
	exp.Parent = nil

	var document *Document
	got := document.CreateElementIs("div", "fruit").(*Element).node

	assert.Equal(t, exp, got)
}

func TestDocumentHTMLNode_CreateElement(t *testing.T) {
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	exp := parsedDocument.FirstChild.NextSibling.LastChild.FirstChild
	exp.Parent = nil

	var document *Document
	got := document.CreateElement("SPAN").(*Element).node

	assert.Equal(t, exp, got)
}
