package domtest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crhntr/dom/spec"
)

type ElementAssertion struct {
	element spec.Element
}

func QuerySelect(t T, parent spec.ElementQueries, query string) ElementAssertion {
	t.Helper()
	element := parent.QuerySelector(query)
	require.NotNil(t, element)
	return ElementAssertion{element: element}
}

func (a ElementAssertion) TextContentEquals(t T, expected string) ElementAssertion {
	t.Helper()
	require.Equal(t, expected, a.element.TextContent())
	return a
}

func (a ElementAssertion) TextContentContains(t T, expected string) ElementAssertion {
	t.Helper()
	require.Contains(t, a.element.TextContent(), expected)
	return a
}

func (a ElementAssertion) HasAttributeValue(t T, name, value string) ElementAssertion {
	t.Helper()
	require.True(t, a.element.HasAttribute(name))
	assert.Equal(t, value, a.element.GetAttribute(name))
	return a
}

func (a ElementAssertion) HasAttribute(t T, name string) ElementAssertion {
	t.Helper()
	assert.True(t, a.element.HasAttribute(name))
	return a
}

func (a ElementAssertion) HasAttributeValues(t T, attributes map[string]string) ElementAssertion {
	for key, value := range attributes {
		a.HasAttributeValue(t, key, value)
	}
	return a
}

type ElementNodeListAssertions struct {
	list spec.NodeList[spec.Element]
}

func QuerySelectAll(t T, parent spec.ElementQueries, query string) ElementNodeListAssertions {
	t.Helper()
	list := parent.QuerySelectorAll(query)
	require.NotZero(t, list.Length())
	return ElementNodeListAssertions{list: list}
}

func (a ElementNodeListAssertions) LengthEquals(t T, expected int) ElementNodeListAssertions {
	t.Helper()
	require.Equal(t, expected, a.list.Length())
	return a
}

func (a ElementNodeListAssertions) Each(t T, assertion func(T, spec.Element)) ElementNodeListAssertions {
	t.Helper()
	for i := 0; i < a.list.Length(); i++ {
		element := a.list.Item(i)
		assertion(t, element)
	}
	return a
}

func (a ElementNodeListAssertions) TextContentEquals(t T, expected string) ElementNodeListAssertions {
	t.Helper()
	a.Each(t, func(t T, element spec.Element) {
		ElementAssertion{element: element}.TextContentEquals(t, expected)
	})
	return a
}

func (a ElementNodeListAssertions) TextContentContains(t T, substring string) ElementNodeListAssertions {
	t.Helper()
	a.Each(t, func(t T, element spec.Element) {
		ElementAssertion{element: element}.TextContentContains(t, substring)
	})
	return a
}

func (a ElementNodeListAssertions) HasAttribute(t T, name string) ElementNodeListAssertions {
	t.Helper()
	a.Each(t, func(t T, element spec.Element) {
		ElementAssertion{element: element}.HasAttribute(t, name)
	})
	return a
}

func (a ElementNodeListAssertions) HasAttributeValue(t T, name, value string) ElementNodeListAssertions {
	t.Helper()
	a.Each(t, func(t T, element spec.Element) {
		ElementAssertion{element: element}.HasAttributeValue(t, name, value)
	})
	return a
}

func (a ElementNodeListAssertions) HasAttributeValues(t T, attributes map[string]string) ElementNodeListAssertions {
	t.Helper()
	a.Each(t, func(t T, element spec.Element) {
		ElementAssertion{element: element}.HasAttributeValues(t, attributes)
	})
	return a
}
