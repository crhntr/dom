package domx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

func TestElement_NodeType(t *testing.T) {
	// language=html
	parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title></title></head><body><span></span></body</html>`))
	require.NoError(t, err)
	document := Element{
		node: parsedDocument.FirstChild.NextSibling,
	}

	assert.Equal(t, dom.NodeTypeElement, document.NodeType())
}
