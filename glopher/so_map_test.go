package glopher

import (
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
)

func Test_StrOrderedMap(t *testing.T) {
	var ok bool
	var value string

	m := NewStrOrderedMap()
	assert.Equal(t, 0, m.Len())

	value, ok = m.Get("foo")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", value)
	assert.Equal(t, 0, m.Len())

	value = m.GetDefault("foo", "defaultBar")
	assert.Equal(t, "defaultBar", value)

	m.Set("foo", "bar")
	value, ok = m.Get("foo")
	assert.Equal(t, true, ok)
	assert.Equal(t, "bar", value)
	assert.Equal(t, 1, m.Len())

	m.Set("name", "unknown")
	value, ok = m.Get("name")
	assert.Equal(t, true, ok)
	assert.Equal(t, "unknown", value)
	assert.Equal(t, 2, m.Len())
	value = m.GetDefault("name", "defaultName")
	assert.Equal(t, "unknown", value)

	itemStrList := []string{}
	for _, row := range m.Items() {
		itemStrList = append(itemStrList, row[0]+":"+row[1])
	}
	assert.Equal(t, "foo:bar | name:unknown", strings.Join(itemStrList, " | "))

	// {
	// 	keys := []string{}
	// 	values := []string{}
	// 	for row := range m.IterItems() {
	// 		keys = append(keys, row[0])
	// 		values = append(values, row[1])
	// 	}
	// 	assert.Equal(t, "foo | name", strings.Join(keys, " | "))
	// 	assert.Equal(t, "bar | unknown", strings.Join(values, " | "))
	// }
	// {
	// 	keys := []string{}
	// 	for key := range m.IterKeys() {
	// 		keys = append(keys, key)
	// 	}
	// 	assert.Equal(t, "foo | name", strings.Join(keys, " | "))
	// }

	value, ok = m.Pop("foo")
	assert.Equal(t, "bar", value)
	assert.Equal(t, true, ok)

	value, ok = m.Pop("foo")
	assert.Equal(t, "", value)
	assert.Equal(t, false, ok)

	value, ok = m.Pop("abcd")
	assert.Equal(t, "", value)
	assert.Equal(t, false, ok)

}
