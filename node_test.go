package dom

import (
	"errors"
	"strings"
	"testing"

	"github.com/andybalholm/cascadia"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/crhntr/dom/spec"
)

func Test_nodeType(t *testing.T) {
	for _, tt := range []struct {
		HTMLNodeType html.NodeType
		DOMNodeType  spec.NodeType
	}{
		{html.ElementNode, spec.NodeTypeElement},
		{html.CommentNode, spec.NodeTypeComment},
		{html.DocumentNode, spec.NodeTypeDocument},
		{html.DoctypeNode, spec.NodeTypeDocumentType},

		{html.NodeType(100000), spec.NodeTypeUnknown},
		{html.ErrorNode, spec.NodeTypeUnknown},
		{html.RawNode, spec.NodeTypeUnknown},
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

func Test_recursiveTextContent_write_failure(t *testing.T) {
	w := &writeError{}
	text := &html.Node{
		Type: html.TextNode,
		Data: "failure",
	}
	require.Panics(t, func() {
		recursiveTextContent(w, text)
	})
}

type writeError struct{}

func (w writeError) WriteString(string) (n int, err error) {
	return 0, errors.New("lemon")
}
