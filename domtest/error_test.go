package domtest_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/domtest"
	"github.com/crhntr/dom/internal/fakes"
	"github.com/crhntr/dom/spec"
)

func TestDocumentResponseErrorMessage(t *testing.T) {
	for _, tt := range []struct {
		Name    string
		ResCode int
		ExpCode int
		Body    string
		Query   string
		Assert  func(t *testing.T, fakeT *fakes.T, el spec.Element, ok bool)
	}{
		{
			Name:    "wrong status code",
			ExpCode: http.StatusBadRequest,
			ResCode: http.StatusOK,
			Query:   "#error",
			Assert: func(t *testing.T, fakeT *fakes.T, el spec.Element, ok bool) {
				assert.NotZero(t, errorMethodCallCount(fakeT))
				assert.Nil(t, el)
				assert.False(t, ok)
			},
		},
		{
			Name:    "query succeeds ",
			ExpCode: http.StatusBadRequest,
			ResCode: http.StatusBadRequest,
			Query:   "#error",
			Body:    `<!DOCTYPE html><html lang="us-en"><head></head><body><p id="error">bad input value</p></body></html>`,
			Assert: func(t *testing.T, fakeT *fakes.T, el spec.Element, ok bool) {
				assert.Zero(t, errorMethodCallCount(fakeT))
				assert.NotNil(t, el)
				assert.True(t, ok)
				assert.Equal(t, "bad input value", el.TextContent())
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(tt.ResCode)
			_, _ = rec.WriteString(tt.Body)
			res := rec.Result()
			fakeT := new(fakes.T)
			el, ok := domtest.DocumentResponseErrorMessage(fakeT, res, tt.ExpCode, tt.Query)
			tt.Assert(t, fakeT, el, ok)
		})
	}

	t.Run("fail to read document", func(t *testing.T) {
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(iotest.ErrReader(errors.New("banana"))),
		}
		fakeT := new(fakes.T)
		el, ok := domtest.DocumentResponseErrorMessage(fakeT, res, http.StatusOK, "#error")

		assert.NotZero(t, errorMethodCallCount(fakeT))
		assert.Nil(t, el)
		assert.False(t, ok)
	})
}

func errorMethodCallCount(t *fakes.T) int {
	return t.ErrorfCallCount() + t.ErrorCallCount() + t.FailNowCallCount()
}

func TestDocumentFragmentResponseErrorMessage(t *testing.T) {
	for _, tt := range []struct {
		Name    string
		ResCode int
		ExpCode int
		Body    string
		Query   string
		Assert  func(t *testing.T, fakeT *fakes.T, el spec.Element, ok bool)
	}{
		{
			Name:    "wrong status code",
			ExpCode: http.StatusBadRequest,
			ResCode: http.StatusOK,
			Query:   "#error",
			Assert: func(t *testing.T, fakeT *fakes.T, el spec.Element, ok bool) {
				assert.NotZero(t, errorMethodCallCount(fakeT))
				assert.Nil(t, el)
				assert.False(t, ok)
			},
		},
		{
			Name:    "query succeeds ",
			ExpCode: http.StatusBadRequest,
			ResCode: http.StatusBadRequest,
			Query:   "#error",
			Body:    `<p id="error">bad input value</p>`,
			Assert: func(t *testing.T, fakeT *fakes.T, el spec.Element, ok bool) {
				assert.Zero(t, errorMethodCallCount(fakeT))
				assert.NotNil(t, el)
				assert.True(t, ok)
				assert.Equal(t, "bad input value", el.TextContent())
			},
		},
	} {
		t.Run(tt.Name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(tt.ResCode)
			_, _ = rec.WriteString(tt.Body)
			res := rec.Result()
			fakeT := new(fakes.T)
			el, ok := domtest.DocumentFragmentResponseErrorMessage(fakeT, res, tt.ExpCode, tt.Query, atom.Body)
			tt.Assert(t, fakeT, el, ok)
		})
	}

	t.Run("fail to read document", func(t *testing.T) {
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(iotest.ErrReader(errors.New("banana"))),
		}
		fakeT := new(fakes.T)
		el, ok := domtest.DocumentFragmentResponseErrorMessage(fakeT, res, http.StatusOK, "#error", atom.Body)

		assert.NotZero(t, errorMethodCallCount(fakeT))
		assert.Nil(t, el)
		assert.False(t, ok)
	})
}
