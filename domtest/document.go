package domtest

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom"
	"github.com/crhntr/dom/spec"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate --fake-name T -o ../internal/fakes/t.go . T

type T interface {
	Helper()
	Error(...any)
	Log(...any)
	Errorf(format string, args ...interface{})
	FailNow()
	Failed() bool
	SkipNow()
}

func Response(t T, res *http.Response) spec.Document {
	t.Helper()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
		if err := res.Body.Close(); err != nil {
			t.Error(err)
		}
		return nil
	}
	if err := res.Body.Close(); err != nil {
		t.Error(err)
		return nil
	}
	return Reader(t, bytes.NewReader(buf))
}

func String(t T, s string) spec.Document {
	t.Helper()
	return Reader(t, strings.NewReader(s))
}

func Reader(t T, r io.Reader) spec.Document {
	t.Helper()
	node, err := html.Parse(r)
	if err != nil {
		t.Error(err)
		return nil
	}
	return dom.NewNode(node).(spec.Document)
}

func DocumentFragmentReader(t T, r io.Reader, parent atom.Atom) spec.DocumentFragment {
	t.Helper()

	body, err := io.ReadAll(r)
	if err != nil {
		t.Error(err)
		return nil
	}
	nodes, err := html.ParseFragment(bytes.NewReader(body), &html.Node{
		Type:     html.ElementNode,
		Data:     parent.String(),
		DataAtom: parent,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	return dom.NewDocumentFragment(nodes)
}

func DocumentFragmentResponse(t T, res *http.Response, parent atom.Atom) spec.DocumentFragment {
	t.Helper()
	defer closeAndCheckError(t, res.Body)
	return DocumentFragmentReader(t, res.Body, parent)
}

func closeAndCheckError(t T, c io.Closer) {
	if err := c.Close(); err != nil {
		t.Error(err)
	}
}
