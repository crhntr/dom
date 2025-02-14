package domtest

import (
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/spec"
)

type ThenFunc[T TestingT, F any] func(t T, res *http.Response, app *F)

type Case[T TestingT, F any] struct {
	Name  string
	Given func(t T, app *F)
	When  func(t T) *http.Request
	Then  ThenFunc[T, F]
}

func (tc Case[TestingT, F]) Run(handler func(f *F) http.Handler) func(t TestingT) {
	return func(t TestingT) {
		t.Helper()
		if tc.When == nil {
			t.Errorf("case field 'When' is not set")
			return
		}
		f := new(F)
		if tc.Given != nil {
			tc.Given(t, f)
		}
		req := tc.When(t)
		rec := httptest.NewRecorder()

		handler(f).ServeHTTP(rec, req)

		res := rec.Result()

		if tc.Then != nil {
			tc.Then(t, res, f)
		}
	}
}

type DocumentTestFunc[T TestingT, F any] func(t T, document spec.Document, f *F)

func Document[T TestingT, F any](then DocumentTestFunc[T, F]) ThenFunc[T, F] {
	return func(t T, res *http.Response, f *F) {
		t.Helper()
		doc := ParseResponseDocument(t, res)
		then(t, doc, f)
	}
}

type FragmentTestFunc[T TestingT, F any] func(t T, fragment spec.DocumentFragment, f *F)

func Fragment[T TestingT, F any](parent atom.Atom, then FragmentTestFunc[T, F]) ThenFunc[T, F] {
	return func(t T, res *http.Response, f *F) {
		t.Helper()
		fragment := ParseResponseDocumentFragment(t, res, parent)
		then(t, fragment, f)
	}
}

type QuerySelectorFunc[T TestingT, F any] func(t T, element spec.Element, f *F)

func QuerySelector[T TestingT, F any](query string, then QuerySelectorFunc[T, F]) ThenFunc[T, F] {
	return func(t T, res *http.Response, f *F) {
		t.Helper()
		document := ParseResponseDocument(t, res)
		el := document.QuerySelector(query)
		if !assert.NotNilf(t, el, "querySelector(%q) did not select any elements", query) {
			t.Log("document", document)
		}
		then(t, el, f)
	}
}
