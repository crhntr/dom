package domx

import (
	"strings"
	"testing"

	"github.com/andybalholm/cascadia"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

func Test_nodeType(t *testing.T) {
	for _, tt := range []struct {
		HTMLNodeType html.NodeType
		DOMNodeType  dom.NodeType
	}{
		{html.ElementNode, dom.NodeTypeElement},
		{html.CommentNode, dom.NodeTypeComment},
		{html.DocumentNode, dom.NodeTypeDocument},
		{html.DoctypeNode, dom.NodeTypeDocumentType},

		{html.NodeType(100000), dom.NodeTypeUnknown},
		{html.ErrorNode, dom.NodeTypeUnknown},
		{html.RawNode, dom.NodeTypeUnknown},
	} {
		assert.Equal(t, tt.DOMNodeType, nodeType(tt.HTMLNodeType))
	}
}

func parseDocument(t *testing.T, document, selector string) (*Document, *Element) {
	t.Helper()
	parsedDocument, err := html.Parse(strings.NewReader(document))
	require.NoError(t, err)
	var result *Element
	if selector != "" {
		result = &Element{
			node: cascadia.Query(parsedDocument, cascadia.MustCompile(selector)),
		}
	}
	return &Document{
		node: parsedDocument,
	}, result
}
