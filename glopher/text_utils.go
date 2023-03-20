package glopher

// UnescapeNTB: unscapes Newline, Tab, Baskslash, and vertical Bar (if bar=True)
func UnescapeNTB(st string, bar bool) string {
	res := []rune{}
	backslash := false
	for _, c := range st {
		switch c {
		case '\\':
			if backslash {
				res = append(res, '\\')
				backslash = false
			} else {
				backslash = true
			}
		case 'n':
			if backslash {
				res = append(res, '\n')
				backslash = false
			} else {
				res = append(res, 'n')
			}
		case 't':
			if backslash {
				res = append(res, '\t')
				backslash = false
			} else {
				res = append(res, 't')
			}
		case '|':
			if backslash {
				if !bar {
					res = append(res, '\\')
				}
				backslash = false
			}
			res = append(res, '|')
		default:
			if backslash {
				res = append(res, '\\')
				backslash = false
			}
			res = append(res, c)
		}
	}
	if backslash {
		res = append(res, '\\')
	}
	return string(res)
}

// SplitByBarUnescapeNTB: splits by "|" (and not "\\|")
// then unescapes Newline (\\n), Tab (\\t), Baskslash (\\)
// and Bar (\\|) in each part
func SplitByBarUnescapeNTB(st string) []string {
	if st == "" {
		return []string{""}
	}
	parts := []string{}
	buf := []rune{}
	backslash := false
	for _, c := range st {
		switch c {
		case '\\':
			if backslash {
				buf = append(buf, '\\')
				backslash = false
			} else {
				backslash = true
			}
		case '|':
			if backslash {
				buf = append(buf, '|')
				backslash = false
			} else {
				parts = append(parts, UnescapeNTB(string(buf), false))
				buf = nil
			}
		default:
			if backslash {
				buf = append(buf, '\\')
				backslash = false
			}
			buf = append(buf, c)
		}
	}
	if backslash {
		buf = append(buf, '\\')
	}
	parts = append(parts, UnescapeNTB(string(buf), false))
	return parts
}
