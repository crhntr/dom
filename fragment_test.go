package dom_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom"
	"github.com/crhntr/dom/spec"
)

func parseDocumentFragment(t *testing.T, input string) *dom.DocumentFragment {
	t.Helper()
	children, err := html.ParseFragment(strings.NewReader(input), &html.Node{
		Type:     html.ElementNode,
		Data:     atom.Div.String(),
		DataAtom: atom.Div,
	})
	require.NoError(t, err)
	return dom.NewDocumentFragment(children)
}

func TestDocumentFragment_NodeType(t *testing.T) {
	fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
	require.Equal(t, spec.NodeTypeDocumentFragment, fragment.NodeType())
}

func TestDocumentFragment_CloneNode(t *testing.T) {
	t.Run("deep", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)

		clone := fragment.CloneNode(true)
		clonedFragment := clone.(*dom.DocumentFragment)

		require.False(t, fragment.IsSameNode(clone))
		children := fragment.Children()
		clonedChildren := clonedFragment.Children()
		for i := 0; i < children.Length(); i++ {
			require.False(t, clonedChildren.Item(i).IsSameNode(children.Item(i)))
		}
	})
	t.Run("shallow", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)

		clone := fragment.CloneNode(false)
		clonedFragment := clone.(*dom.DocumentFragment)

		require.False(t, fragment.IsSameNode(clone))
		children := fragment.Children()
		clonedChildren := clonedFragment.Children()
		for i := 0; i < children.Length(); i++ {
			require.True(t, clonedChildren.Item(i).IsSameNode(children.Item(i)))
		}
	})
}

func TestDocumentFragment_IsSameNode(t *testing.T) {
	t.Run("another fragment", func(t *testing.T) {
		a := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
		require.True(t, a.IsSameNode(a))

		b := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)

		require.False(t, a.IsSameNode(b))
		require.False(t, b.IsSameNode(a))

		c := b
		require.True(t, b.IsSameNode(c))
		require.True(t, c.IsSameNode(b))
		require.True(t, c.IsSameNode(c))
	})
	t.Run("another node type", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
		require.False(t, fragment.IsSameNode(nil))
	})
}

func TestDocumentFragment_TextContent(t *testing.T) {
	fragment := parseDocumentFragment(t, `Hello, <em>world</em>!<br>`)
	require.Equal(t, "Hello, world!", fragment.TextContent())
}

func TestDocumentFragment_Children(t *testing.T) {
	fragment := parseDocumentFragment(t, `<div>1</div><div>2</div><div>3</div>`)
	children := fragment.Children()
	require.Equal(t, 3, children.Length())
	require.Equal(t, "1", children.Item(0).TextContent())
	require.Equal(t, "2", children.Item(1).TextContent())
	require.Equal(t, "3", children.Item(2).TextContent())
}

func TestDocumentFragment_FirstElementChild(t *testing.T) {
	t.Run("has elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `a<span>x</span>b<span>y</span>c`)
		require.Equal(t, "x", fragment.FirstElementChild().TextContent())
	})
	t.Run("no elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello!`)
		require.Nil(t, fragment.FirstElementChild())
	})
	t.Run("one element", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `hello<div>x</div>`)
		require.Equal(t, "x", fragment.FirstElementChild().TextContent())
	})
}

func TestDocumentFragment_LastElementChild(t *testing.T) {
	t.Run("has elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `a<span>x</span>b<span>y</span>c`)
		require.Equal(t, "y", fragment.LastElementChild().TextContent())
	})
	t.Run("no elements", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `Hello!`)
		require.Nil(t, fragment.LastElementChild())
	})
	t.Run("one element", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div>x</div>world`)
		require.Equal(t, "x", fragment.LastElementChild().TextContent())
	})
}

func TestDocumentFragment_ChildElementCount(t *testing.T) {
	fragment := parseDocumentFragment(t, `a<div>x</div>b<div>y</div>`)
	require.Equal(t, 2, fragment.ChildElementCount())
}

func TestDocumentFragment_Append(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fragment := parseDocumentFragment(t, ``)
		fragment.Append()
		require.Equal(t, 0, fragment.ChildElementCount())
	})
	t.Run("add one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 1)
		node := dom.NewNode(other[0]).(spec.Element)

		fragment.Append(node)

		require.Equal(t, 2, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "a", children.Item(0).ID())
		assert.Equal(t, "b", children.Item(1).ID())
	})
	t.Run("add two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div><div id="c"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 2)

		a, ok := dom.NewNode(other[0]).(spec.Element)
		require.True(t, ok)

		b, ok := dom.NewNode(other[1]).(spec.Element)
		require.True(t, ok)

		fragment.Append(a, b)

		require.Equal(t, 3, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "a", children.Item(0).ID())
		assert.Equal(t, "b", children.Item(1).ID())
		assert.Equal(t, "c", children.Item(2).ID())
	})
}

func TestDocumentFragment_Prepend(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fragment := parseDocumentFragment(t, ``)
		fragment.Prepend()
		require.Equal(t, 0, fragment.ChildElementCount())
	})
	t.Run("add one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 1)
		node := dom.NewNode(other[0]).(spec.Element)

		fragment.Prepend(node)

		require.Equal(t, 2, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
		assert.Equal(t, "a", children.Item(1).ID())
	})
	t.Run("add two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div><div id="c"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 2)

		a, ok := dom.NewNode(other[0]).(spec.Element)
		require.True(t, ok)

		b, ok := dom.NewNode(other[1]).(spec.Element)
		require.True(t, ok)

		fragment.Prepend(a, b)

		require.Equal(t, 3, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
		assert.Equal(t, "c", children.Item(1).ID())
		assert.Equal(t, "a", children.Item(2).ID())
	})
}

func TestDocumentFragment_ReplaceChildren(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fragment := parseDocumentFragment(t, ``)
		fragment.ReplaceChildren()
		require.Equal(t, 0, fragment.ChildElementCount())
	})
	t.Run("replace one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 1)
		node := dom.NewNode(other[0]).(spec.Element)

		fragment.ReplaceChildren(node)

		require.Equal(t, 1, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
	})
	t.Run("replace two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<div id="a"></div>`)
		require.Equal(t, 1, fragment.ChildElementCount())

		other, err := html.ParseFragment(strings.NewReader(`<div id="b"></div><div id="c"></div>`), &html.Node{
			Type:     html.ElementNode,
			Data:     atom.Div.String(),
			DataAtom: atom.Div,
		})
		require.NoError(t, err)
		require.Len(t, other, 2)

		a, ok := dom.NewNode(other[0]).(spec.Element)
		require.True(t, ok)

		b, ok := dom.NewNode(other[1]).(spec.Element)
		require.True(t, ok)

		fragment.ReplaceChildren(a, b)

		require.Equal(t, 2, fragment.ChildElementCount())

		children := fragment.Children()
		assert.Equal(t, "b", children.Item(0).ID())
		assert.Equal(t, "c", children.Item(1).ID())
	})
}

func TestDocumentFragment_QuerySelector(t *testing.T) {
	t.Run("found one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.NotNil(t, fragment.QuerySelector("#a"))
	})
	t.Run("found two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `
<section><div class="a" id="i1"></div></section>
<section></section>
<section><div class="a" id="i2"></div><div class="a" id="i3"></div></section>`)
		result := fragment.QuerySelector(".a")
		require.NotNil(t, result)
		assert.Equal(t, "i1", result.ID())
	})
	t.Run("not found", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.Nil(t, fragment.QuerySelector("#not-found"))
	})
}

func TestDocumentFragment_QuerySelectorAll(t *testing.T) {
	t.Run("found one", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.NotNil(t, fragment.QuerySelectorAll("#a"))
	})
	t.Run("found two", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `
<section><div class="a" id="i1"></div></section>
<section></section>
<section><div class="a" id="i2"></div><div class="a" id="i3"></div></section>`)
		results := fragment.QuerySelectorAll(".a")
		assert.Equal(t, "i1", results.Item(0).ID())
		assert.Equal(t, "i2", results.Item(1).ID())
		assert.Equal(t, "i3", results.Item(2).ID())
	})
	t.Run("not found", func(t *testing.T) {
		fragment := parseDocumentFragment(t, `<section><div id="a"></div></section><section></section>`)
		require.Nil(t, fragment.QuerySelectorAll("#not-found"))
	})
}
