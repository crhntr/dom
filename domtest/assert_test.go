package domtest_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crhntr/dom/domtest"
	"github.com/crhntr/dom/internal/fakes"
)

func TestQuerySelect(t *testing.T) {

	const documentText =
	/* language=html */
	`<!DOCTYPE html>
<html lang="en-us">
<head>
    <meta charset="utf-8">
    <title>Test Data</title>
</head>
<body>
    <p id="greeting">Hello, world!</p>
	<ul>
		<li class="happy-item" data-selected="false" data-index="0">Item 1</li>
		<li class="happy-item" data-selected="false" data-index="1">Item 2</li>
		<li class="happy-item" data-selected="false" data-index="2">Item 3</li>
	</ul>
	
	<div id="texts">
		<div><span>text</span></div>
		<div><span>tex</span>t</div>
		<div><span>te</span>xt</div>
		<div><span>t</span>ext</div>
	</div>
</body>
</html>`

	t.Run("query select text equals", func(t *testing.T) {
		tst := new(fakes.T)
		document := domtest.String(t, documentText)

		domtest.QuerySelect(tst, document, `head title`).
			TextContentEquals(tst, "Test Data")

		assert.Zero(t, tst.FailNowCallCount())
	})

	t.Run("query select chained assertions", func(t *testing.T) {
		tst := new(fakes.T)
		document := domtest.String(t, documentText)

		domtest.QuerySelect(tst, document, `body p`).
			TextContentEquals(tst, "Hello, world!").
			HasAttributeValue(tst, "id", "greeting")

		assert.Zero(t, tst.FailNowCallCount())
	})

	t.Run("query select all chained assertions", func(t *testing.T) {
		tst := new(fakes.T)
		document := domtest.String(t, documentText)

		domtest.QuerySelectAll(tst, document, `body ul li`).
			LengthEquals(tst, 3).
			HasAttribute(tst, "data-index").
			HasAttributeValue(tst, "data-selected", "false").
			TextContentContains(tst, "Item")

		assert.Zero(t, tst.FailNowCallCount())
	})

	t.Run("query select all text content", func(t *testing.T) {
		tst := new(fakes.T)
		document := domtest.String(t, documentText)

		domtest.QuerySelectAll(tst, document, `#texts>div`).
			LengthEquals(tst, 4).
			TextContentEquals(tst, "text")

		assert.Zero(t, tst.FailNowCallCount())
	})

	t.Run("query select all unexpected count", func(t *testing.T) {
		tst := new(fakes.T)
		document := domtest.String(t, documentText)

		domtest.QuerySelectAll(tst, document, `#texts>div`).
			LengthEquals(tst, 10000)

		assert.Equal(t, 1, tst.FailNowCallCount())
	})

	t.Run("query select attribute values checks", func(t *testing.T) {
		tst := new(fakes.T)
		document := domtest.String(t, documentText)

		domtest.QuerySelectAll(tst, document, `[data-index]`).
			HasAttributeValues(tst, map[string]string{
				"data-selected": "false",
				"class":         "happy-item",
			})

		assert.Equal(t, 0, tst.FailNowCallCount())
	})
}
