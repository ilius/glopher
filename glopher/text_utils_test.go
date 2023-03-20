package glopher

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestUnescapeNTB(t *testing.T) {
	is := is.New(t)

	is.Equal("a", UnescapeNTB("a", false))
	is.Equal("a\t", UnescapeNTB("a\\t", false))
	is.Equal("a\n", UnescapeNTB("a\\n", false))
	is.Equal("\ta", UnescapeNTB("\\ta", false))
	is.Equal("\na", UnescapeNTB("\\na", false))
	is.Equal("a\tb\n", UnescapeNTB("a\\tb\\n", false))
	is.Equal("a\\b", UnescapeNTB("a\\\\b", false))
	is.Equal("a\\\tb", UnescapeNTB("a\\\\\\tb", false))
	is.Equal("a|b\tc", UnescapeNTB("a|b\\tc", false))
	is.Equal("a\\|b\tc", UnescapeNTB("a\\|b\\tc", false))
	is.Equal("a\\|b\tc", UnescapeNTB("a\\\\|b\\tc", false))
	is.Equal("|", UnescapeNTB("\\|", true))
	is.Equal("a|b", UnescapeNTB("a\\|b", true))
	is.Equal("a|b\tc", UnescapeNTB("a\\|b\\tc", true))

	is.Equal(`\a`, UnescapeNTB(`\a`, false))
}

func TestSplitByBarUnescapeNTB(t *testing.T) {
	is := is.New(t)
	f := SplitByBarUnescapeNTB
	is.Equal(f(""), []string{""})
	is.Equal(f("|"), []string{"", ""})
	is.Equal(f("a"), []string{"a"})
	is.Equal(f("a|"), []string{"a", ""})
	is.Equal(f("|a"), []string{"", "a"})
	is.Equal(f("a|b"), []string{"a", "b"})
	is.Equal(f("a\\|b|c"), []string{"a|b", "c"})
	is.Equal(f("a\\\\1|b|c"), []string{"a\\1", "b", "c"})
	is.Equal(f("a\\\\|b|c"), []string{"a\\", "b", "c"})
	is.Equal(f("a\\\\1|b\\n|c\\t"), []string{"a\\1", "b\n", "c\t"})
}
