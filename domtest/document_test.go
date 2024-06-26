package domtest_test

import (
	_ "embed"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/domtest"
	"github.com/crhntr/dom/internal/fakes"
	"github.com/crhntr/dom/spec"
)

var (
	//go:embed testdata/index.html
	indexHTML string

	//go:embed testdata/fragment.html
	fragmentHTML string
)

func TestResponse(t *testing.T) {
	t.Run("when a valid html document is passed", func(t *testing.T) {
		testingT := new(fakes.T)
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(indexHTML)),
		}
		document := domtest.Response(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())

		require.NotNil(t, document)
		p := document.QuerySelector(`p`)
		assert.Equal(t, p.TextContent(), "Hello, world!")
	})

	t.Run("when a fragment is returned", func(t *testing.T) {
		testingT := new(fakes.T)
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(fragmentHTML)),
		}
		document := domtest.Response(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())

		require.NotNil(t, document)
		list := document.QuerySelectorAll(`p`)
		require.Equal(t, 2, list.Length())
		require.Equal(t, list.Item(0).TextContent(), "Hello, world!")
		require.Equal(t, list.Item(1).TextContent(), "Greeting!")
	})

	t.Run("when read fails and close is ok", func(t *testing.T) {
		testingT := new(fakes.T)
		fakeBody := &errClose{
			Reader:   iotest.ErrReader(errors.New("banana")),
			closeErr: nil,
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.Response(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 1, "it should report an error")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})

	t.Run("when read is ok but close fails", func(t *testing.T) {
		testingT := new(fakes.T)
		fakeBody := &errClose{
			Reader:   strings.NewReader(indexHTML),
			closeErr: errors.New("banana"),
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.Response(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 1, "it should report two errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})

	t.Run("when both read and close fail", func(t *testing.T) {
		testingT := new(fakes.T)
		fakeBody := &errClose{
			Reader:   iotest.ErrReader(errors.New("banana")),
			closeErr: errors.New("lemon"),
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.Response(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 2, "it should report two errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.Equal(t, testingT.HelperCallCount(), 1)
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})

	t.Run("when both read and close fail", func(t *testing.T) {
		testingT := new(fakes.T)
		fakeBody := &errClose{
			Reader:   iotest.ErrReader(errors.New("banana")),
			closeErr: errors.New("lemon"),
		}
		res := &http.Response{
			Body: fakeBody,
		}
		document := domtest.Response(testingT, res)

		assert.Equal(t, testingT.ErrorCallCount(), 2, "it should report two errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())
		assert.Equal(t, fakeBody.closeCallCount, 1)

		assert.Nil(t, document)
	})
}

func TestDocumentFragment(t *testing.T) {
	t.Run("when a valid html document is passed", func(t *testing.T) {
		testingT := new(fakes.T)
		res := &http.Response{
			Body: io.NopCloser(strings.NewReader(fragmentHTML)),
		}
		fragment := domtest.DocumentFragmentResponse(testingT, res, atom.Body)

		assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
		assert.Equal(t, testingT.LogCallCount(), 0)
		assert.NotZero(t, testingT.HelperCallCount())

		require.NotNil(t, fragment)
		require.Equal(t, spec.NodeTypeDocumentFragment, fragment.NodeType())
	})
}

func TestReader(t *testing.T) {
	testingT := new(fakes.T)
	r := iotest.ErrReader(errors.New("banana"))

	document := domtest.Reader(testingT, r)

	assert.Equal(t, testingT.ErrorCallCount(), 1, "it should report two errors")
	assert.NotZero(t, testingT.HelperCallCount())
	assert.Nil(t, document)
}

func TestString(t *testing.T) {
	testingT := new(fakes.T)

	document := domtest.String(testingT, "<p>Hello, world!</p>")

	assert.Equal(t, testingT.ErrorCallCount(), 0, "it should not report errors")
	assert.Equal(t, testingT.LogCallCount(), 0)
	assert.NotZero(t, testingT.HelperCallCount())

	assert.NotNil(t, document)
	p := document.QuerySelector(`p`)
	assert.Equal(t, p.TextContent(), "Hello, world!")
}

type errClose struct {
	io.Reader
	closeCallCount int
	closeErr       error
}

func (e *errClose) Close() error {
	e.closeCallCount++
	return e.closeErr
}
