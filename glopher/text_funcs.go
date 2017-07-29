package glopher

import (
	"bytes"
	"io"
)

var Newline = []byte{'\n'}

const (
	backslash = byte('\\')
	byte_n    = byte('n')
	byte_r    = byte('r')
	byte_t    = byte('t')
	byte_NL   = byte('\n') // Newline
	byte_CR   = byte('\r') // Carriage return
	byte_TAB  = byte('\t')
	byte_BAR  = byte('|')
)

func CountBlocks(r io.Reader, sep []byte) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], sep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func escapeNTB_Bytes(orig []byte, bar bool) []byte {
	result := make([]byte, 0, 2*len(orig))
	for _, c := range orig {
		switch c {
		case backslash:
			result = append(result, backslash, backslash)
		case byte_NL:
			result = append(result, backslash, byte_n)
		case byte_CR:
			result = append(result, backslash, byte_r)
		case byte_TAB:
			result = append(result, backslash, byte_t)
		case byte_BAR:
			if bar {
				result = append(result, backslash, byte_BAR)
			} else {
				result = append(result, byte_BAR)
			}
		default:
			result = append(result, c)
		}
	}
	return result
}

func unescapeNTB_Bytes(orig []byte) []byte {
	result := make([]byte, 0, len(orig))
	openBS := false
	for _, c := range orig {
		switch c {
		case backslash:
			if openBS {
				result = append(result, backslash)
				openBS = false
			} else {
				openBS = true
			}
		case byte_n:
			if openBS {
				result = append(result, byte_NL)
				openBS = false
			} else {
				result = append(result, c)
			}
		case byte_r:
			if openBS {
				result = append(result, byte_CR)
				openBS = false
			} else {
				result = append(result, c)
			}
		case byte_t:
			if openBS {
				result = append(result, byte_TAB)
				openBS = false
			} else {
				result = append(result, c)
			}
		case byte_BAR:
			if openBS {
				openBS = false
			}
			result = append(result, c)
		default:
			if openBS {
				result = append(result, backslash)
				openBS = false
			}
			result = append(result, c)
		}
	}
	return result
}

func splitByBar_Bytes(orig []byte) [][]byte {
	result := [][]byte{
		[]byte{},
	}
	openBS := false
	for _, c := range orig {
		n := len(result)
		switch c {
		case backslash:
			if openBS {
				// we must not escape backslash here
				result[n-1] = append(result[n-1], backslash, backslash)
				openBS = false
			} else {
				openBS = true
			}
		case byte_BAR:
			if openBS {
				result[n-1] = append(result[n-1], byte_BAR)
				openBS = false
			} else {
				result = append(result, []byte{})
			}
		default:
			if openBS {
				result[n-1] = append(result[n-1], backslash)
				openBS = false
			}
			result[n-1] = append(result[n-1], c)
		}
	}
	return result
}

// UnescapeNTB: escapes Newline, Tab, Baskslash, and vertical Bar (if bar=True)
func EscapeNTB(orig string, bar bool) string {
	return string(escapeNTB_Bytes([]byte(orig), bar))
}

// UnescapeNTB: unescapes Newline, Tab, Baskslash, and vertical Bar (if bar=True)
func UnescapeNTB(orig string) string {
	return string(unescapeNTB_Bytes([]byte(orig)))
}

func SplitByBarUnescapeNTB(orig string) []string {
	result := []string{}
	rawParts := splitByBar_Bytes([]byte(orig))
	for _, rawPart := range rawParts {
		result = append(result, string(unescapeNTB_Bytes(rawPart)))
	}
	return result
}
