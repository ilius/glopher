package stardict

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func readSyn(synPath string) (map[int][]string, error) {
	synBytes, err := os.ReadFile(synPath)
	if err != nil {
		return nil, err
	}
	synByteN := len(synBytes)
	pos := 0
	altsMap := map[int][]string{}
	for pos < synByteN {
		beg := pos
		// Python: pos = data.find("\x00", beg)
		offset := bytes.Index(synBytes[beg:], []byte{0})
		if offset < 0 {
			return nil, fmt.Errorf("synonym file is corrupted")
		}
		pos = offset + beg
		b_alt := synBytes[beg:pos]
		pos += 1
		if pos+4 > len(synBytes) {
			return nil, fmt.Errorf("synonym file is corrupted")
		}
		entryIndex := int(binary.BigEndian.Uint32(synBytes[pos : pos+4]))
		pos += 4
		alt := string(b_alt)
		altsMap[entryIndex] = append(altsMap[entryIndex], alt)
	}
	return altsMap, nil
}
