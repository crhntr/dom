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
