package blog_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html/atom"

	"github.com/crhntr/dom/domtest"
	"github.com/crhntr/dom/spec"

	"github.com/crhntr/dom/domtest/example/blog"
	"github.com/crhntr/dom/domtest/example/blog/internal/fake"
)

func TestCase(t *testing.T) {
	for _, tt := range []domtest.Case[*testing.T, *fake.App]{
		{
			Name: "viewing the home page",
			Given: func(t *testing.T) *fake.App {
				app := new(fake.App)
				app.ArticleReturns(blog.Article{
					Title:   "Greetings!",
					Content: "Hello, friends!",
					Error:   nil,
				})
				return app
			},
			When: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/1", nil)
			},
			Then: domtest.Document(func(t *testing.T, document spec.Document, app *fake.App) {
				require.Equal(t, 1, app.ArticleArgsForCall(0))
				if heading := document.QuerySelector("h1"); assert.NotNil(t, heading) {
					require.Equal(t, "Greetings!", heading.TextContent())
				}
				if content := document.QuerySelector("p"); assert.NotNil(t, content) {
					require.Equal(t, "Hello, friends!", content.TextContent())
				}
			}),
		},
		{
			Name: "the page has an error",
			Given: func(t *testing.T) *fake.App {
				app := new(fake.App)
				app.ArticleReturns(blog.Article{
					Error: fmt.Errorf("lemon"),
				})
				return app
			},
			When: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/1", nil)
			},
			Then: domtest.QuerySelector("#error-message", func(t *testing.T, msg spec.Element, app *fake.App) {
				require.Equal(t, "lemon", msg.TextContent())
			}),
		},
		{
			Name: "the page has an error and is requested by HTMX",
			Given: func(t *testing.T) *fake.App {
				app := new(fake.App)
				app.ArticleReturns(blog.Article{
					Error: fmt.Errorf("lemon"),
				})
				return app
			},
			When: func(t *testing.T) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/article/1", nil)
				req.Header.Set("HX-Request", "true")
				return req
			},
			Then: domtest.Fragment(atom.Body, func(t *testing.T, fragment spec.DocumentFragment, app *fake.App) {
				el := fragment.FirstElementChild()
				require.Equal(t, "lemon", el.TextContent())
				require.Equal(t, "*errors.errorString", el.GetAttribute("data-type"))
			}),
		},
		{
			Name: "when the id is not an integer",
			When: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodGet, "/article/banana", nil)
			},
			Then: func(t *testing.T, res *http.Response, f *fake.App) {
				require.Equal(t, http.StatusBadRequest, res.StatusCode)
			},
		},
	} {
		t.Run(tt.Name, tt.Run(func(fakes *fake.App) http.Handler {
			mux := http.NewServeMux()
			blog.Routes(mux, fakes)
			return mux
		}))
	}
}
