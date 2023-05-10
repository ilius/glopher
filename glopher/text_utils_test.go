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

func TestEscapeNTB(t *testing.T) {
	is := is.New(t)
	f := EscapeNTB
	is.Equal(f("a", false), "a")
	is.Equal(f("a\t", false), "a\\t")
	is.Equal(f("a\n", false), "a\\n")
	is.Equal(f("\ta", false), "\\ta")
	is.Equal(f("\na", false), "\\na")
	is.Equal(f("a\tb\n", false), "a\\tb\\n")
	is.Equal(f("a\\b", false), "a\\\\b")
	is.Equal(f("a\\\tb", false), "a\\\\\\tb")
	is.Equal(f("a|b\tc", false), "a|b\\tc")
	is.Equal(f("a\\|b\tc", false), "a\\\\|b\\tc")
	is.Equal(f("|", true), "\\|")
	is.Equal(f("a|b", true), "a\\|b")
	is.Equal(f("a|b\tc", true), "a\\|b\\tc")
}
