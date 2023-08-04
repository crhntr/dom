package domx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

func TestText_NodeType(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title>peach</title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	document := Text{
		node: parsedDocument.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild,
	}

	assert.Equal(t, dom.NodeTypeText, document.NodeType())
}
