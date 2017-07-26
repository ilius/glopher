package glopher

import (
	"bytes"
	"io"
)

var Newline = []byte{'\n'}

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
