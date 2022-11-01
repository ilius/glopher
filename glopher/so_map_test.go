package glopher

import (
	"testing"

	"strings"

	"github.com/ilius/is/v2"
)

func Test_StrOrderedMap(t *testing.T) {
	is := is.New(t)
	var ok bool
	var value string

	m := NewStrOrderedMap()
	is.Equal(0, m.Len())

	value, ok = m.Get("foo")
	is.Equal(false, ok)
	is.Equal("", value)
	is.Equal(0, m.Len())

	value = m.GetDefault("foo", "defaultBar")
	is.Equal("defaultBar", value)

	m.Set("foo", "bar")
	value, ok = m.Get("foo")
	is.Equal(true, ok)
	is.Equal("bar", value)
	is.Equal(1, m.Len())

	m.Set("name", "unknown")
	value, ok = m.Get("name")
	is.Equal(true, ok)
	is.Equal("unknown", value)
	is.Equal(2, m.Len())
	value = m.GetDefault("name", "defaultName")
	is.Equal("unknown", value)

	itemStrList := []string{}
	for _, row := range m.Items() {
		itemStrList = append(itemStrList, row[0]+":"+row[1])
	}
	is.Equal("foo:bar | name:unknown", strings.Join(itemStrList, " | "))

	// {
	// 	keys := []string{}
	// 	values := []string{}
	// 	for row := range m.IterItems() {
	// 		keys = append(keys, row[0])
	// 		values = append(values, row[1])
	// 	}
	// 	is.Equal("foo | name", strings.Join(keys, " | "))
	// 	is.Equal("bar | unknown", strings.Join(values, " | "))
	// }
	// {
	// 	keys := []string{}
	// 	for key := range m.IterKeys() {
	// 		keys = append(keys, key)
	// 	}
	// 	is.Equal("foo | name", strings.Join(keys, " | "))
	// }

	value, ok = m.Pop("foo")
	is.Equal("bar", value)
	is.Equal(true, ok)

	value, ok = m.Pop("foo")
	is.Equal("", value)
	is.Equal(false, ok)

	value, ok = m.Pop("abcd")
	is.Equal("", value)
	is.Equal(false, ok)

}
