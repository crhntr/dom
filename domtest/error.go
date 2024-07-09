package domtest

import (
	"net/http"

	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/spec"
)

func DocumentResponseErrorMessage(t T, res *http.Response, httpStatusCode int, errMessageElementSelector string) (spec.Element, bool) {
	t.Helper()
	if res.StatusCode != httpStatusCode {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	document := Response(t, res)
	if document == nil {
		return nil, false
	}

	errEl := document.QuerySelector(errMessageElementSelector)
	if errEl == nil {
		t.Errorf("error message element not found for query: %s", errMessageElementSelector)
		return nil, false
	}

	return errEl, true
}

func DocumentFragmentResponseErrorMessage(t T, res *http.Response, httpStatusCode int, errMessageElementSelector string, parent atom.Atom) (spec.Element, bool) {
	t.Helper()
	if res.StatusCode != httpStatusCode {
		t.Errorf("unexpected status code: %d", res.StatusCode)
	}

	document := DocumentFragmentResponse(t, res, parent)
	if document == nil {
		return nil, false
	}

	errEl := document.QuerySelector(errMessageElementSelector)
	if errEl == nil {
		t.Errorf("error message element not found for query: %s", errMessageElementSelector)
		return nil, false
	}

	return errEl, true
}
