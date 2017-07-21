package glopher

import (
	"fmt"
	"io"
)

func ReadInfo(reader func() *Entry, maxNonInfo int) (*StrOrderedMap, []*Entry, error) {
	if maxNonInfo < 1 {
		return nil, nil, fmt.Errorf("bad maxNonInfo = %v, must be at least 1", maxNonInfo)
	}
	info := NewStrOrderedMap()
	nonInfo := make([]*Entry, 0, maxNonInfo)
	for {
		entry := reader()
		if entry.Error == io.EOF {
			break
		}
		if entry.Error != nil {
			return info, nonInfo, entry.Error
		}
		if entry.IsInfo {
			info.Set(entry.Word, entry.Defi)
			continue
		}
		nonInfo = append(nonInfo, entry)
		if len(nonInfo) >= maxNonInfo {
			break
		}
	}
	return info, nonInfo, nil
}
