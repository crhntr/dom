package domx

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
