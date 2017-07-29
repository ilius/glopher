package glopher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func repr(v interface{}) string {
	return fmt.Sprintf("%#v", v)
}

func Test_UnescapeNTB_EscapeNTB(t *testing.T) {
	expectedMap := map[string]string{
		"":             "",
		"\\t":          "\t",
		"\\n":          "\n",
		"\\\\":         "\\",
		"\\\\\\\\":     "\\\\",
		"\\\\\\\\\\\\": "\\\\\\",
		"\\\\n":        "\\n",
		"\\\\\\n":      "\\\n",
		"\\\\t":        "\\t",
		"\\\\\\t":      "\\\t",
		"a\\\\nb":      "a\\nb",
		"a\\\\tb\\\\n": "a\\tb\\n",
		"a\\\\\\\\tb":  "a\\\\tb",
		"a\\\\\\\\nb":  "a\\\\nb",

		"a\"b": "a\"b",
		"a'b":  "a'b",

		"خط اول\\nخط دوم":         "خط اول\nخط دوم",
		"ستون ۱\\tستون۲\\tستون ۳": "ستون ۱\tستون۲\tستون ۳",
	}
	for escaped, expectedRaw := range expectedMap {
		actualRaw := UnescapeNTB(escaped)
		assert.Equal(t, expectedRaw, actualRaw)
		newEscaped := EscapeNTB(expectedRaw, false)
		assert.Equal(t, escaped, newEscaped)
	}
}

func Test_SplitByBarUnescapeNTB(t *testing.T) {
	expectedMap := map[string][]string{
		"a|b":               {"a", "b"},
		"a|":                {"a", ""},
		"|a":                {"", "a"},
		"a||b":              {"a", "", "b"},
		"a\"b|a'b":          {"a\"b", "a'b"},
		"a\\nb|c\\td|f\\|g": {"a\nb", "c\td", "f|g"},
		"a\\\\nb":           {"a\\nb"},
		"a\\\\nb|a\\\\tb\\\\n|a\\\\\\\\tb|a\\\\\\\\nb": {"a\\nb", "a\\tb\\n", "a\\\\tb", "a\\\\nb"},
	}
	for orig, expectedResult := range expectedMap {
		result := SplitByBarUnescapeNTB(orig)
		expectedResultRepr := repr(expectedResult)
		resultRepr := repr(result)
		assert.Equal(t, expectedResultRepr, resultRepr)
	}
}
