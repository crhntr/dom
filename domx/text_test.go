package domx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"

	"github.com/crhntr/dom"
)

func TestText_Data(t *testing.T)    {}
func TestText_SetData(t *testing.T) {}

func TestText_NodeType(t *testing.T) {
	t.Run("cromulent", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title>peach</title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		document := Text{
			node: parsedDocument.FirstChild.NextSibling.FirstChild.FirstChild.FirstChild,
		}

		assert.Equal(t, dom.NodeTypeText, document.NodeType())
	})

	t.Run("mismatched node type", func(t *testing.T) {
		// language=html
		parsedDocument, err := html.Parse(strings.NewReader(`<!DOCTYPE html><html lang="us-en"><head><title>peach</title></head><body><span></span></body</html>`))
		require.NoError(t, err)
		document := Text{
			node: parsedDocument,
		}

		assert.Equal(t, dom.NodeTypeDocument, document.NodeType())
	})
}

func TestText_IsConnected(t *testing.T)     {}
func TestText_OwnerDocument(t *testing.T)   {}
func TestText_Length(t *testing.T)          {}
func TestText_ParentNode(t *testing.T)      {}
func TestText_ParentElement(t *testing.T)   {}
func TestText_PreviousSibling(t *testing.T) {}
func TestText_NextSibling(t *testing.T)     {}
func TestText_TextContent(t *testing.T)     {}
func TestText_CloneNode(t *testing.T)       {}
func TestText_IsSameNode(t *testing.T)      {}
func TestText_(t *testing.T)                {}
