package domtest

import (
	"bytes"
	"io"
	"net/http"

	"golang.org/x/net/html"

	"github.com/crhntr/dom"
	"github.com/crhntr/dom/spec"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate --fake-name T -o ../internal/fakes/t.go . T

type T interface {
	Helper()
	Error(...any)
	Log(...any)
}

func ResponseDocument(t T, res *http.Response) spec.Document {
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

func Reader(t T, r io.Reader) spec.Document {
	t.Helper()
	node, err := html.Parse(r)
	if err != nil {
		t.Error(err)
		return nil
	}
	return dom.NewNode(node).(spec.Document)
}
