package domtest

import (
	"net/http"

	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/spec"
)

func DocumentResponseQuery(t T, res *http.Response, statusCode int, selector string) spec.Element {
	t.Helper()
	if res.StatusCode != statusCode {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	document := Response(t, res)
	if document == nil {
		return nil
	}

	errEl := document.QuerySelector(selector)
	if errEl == nil {
		t.Errorf("error message element not found for query: %s", selector)
		return nil
	}
	return errEl
}

func DocumentFragmentQuery(t T, res *http.Response, statusCode int, selector string, parent atom.Atom) spec.Element {
	t.Helper()
	if res.StatusCode != statusCode {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	document := DocumentFragmentResponse(t, res, parent)
	if document == nil {
		return nil
	}

	errEl := document.QuerySelector(selector)
	if errEl == nil {
		t.Errorf("error message element not found for query: %s", selector)
		return nil
	}
	return errEl
}

func DocumentResponseQueryAll(t T, res *http.Response, statusCode int, selector string) spec.NodeList[spec.Element] {
	t.Helper()
	if res.StatusCode != statusCode {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	document := Response(t, res)
	if document == nil {
		return nil
	}

	return document.QuerySelectorAll(selector)
}

func DocumentFragmentQueryAll(t T, res *http.Response, statusCode int, selector string, parent atom.Atom) spec.NodeList[spec.Element] {
	t.Helper()
	if res.StatusCode != statusCode {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	document := DocumentFragmentResponse(t, res, parent)
	if document == nil {
		return nil
	}

	return document.QuerySelectorAll(selector)
}
