package stardict

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
)

// dictionaryImp stardict dictionary
type StarDictReader struct {
	*Info

	dict     *Dict
	ifoPath  string
	idxPath  string
	dictPath string
	synPath  string
	resDir   string
	resURL   string

	decodeData func(data []byte) []*SearchResultItem
}

func (d *StarDictReader) Disabled() bool {
	return d.disabled
}

func (d *StarDictReader) Loaded() bool {
	return d.dict != nil
}

func (d *StarDictReader) SetDisabled(disabled bool) {
	d.disabled = disabled
}

func (d *StarDictReader) ResourceDir() string {
	return d.resDir
}

func (d *StarDictReader) ResourceURL() string {
	return d.resURL
}

func (d *StarDictReader) IndexPath() string {
	return d.idxPath
}

func (d *StarDictReader) InfoPath() string {
	return d.ifoPath
}

func (d *StarDictReader) Close() {
	d.dict.Close()
}

func (d *StarDictReader) decodeWithSametypesequence(data []byte) (items []*SearchResultItem) {
	seq := d.Options["sametypesequence"]

	seqLen := len(seq)

	var dataPos int
	dataSize := len(data)

	for i, t := range seq {
		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			// if last seq item
			if i == seqLen-1 {
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				end := bytes.IndexRune(data[dataPos:], '\000')
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			if i == seqLen-1 {
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
				items = append(items, &SearchResultItem{Type: t, Data: data[dataPos+4 : dataPos+int(size)+5]})
				dataPos += int(size) + 5
			}
		}
	}

	return
}

func (d *StarDictReader) decodeWithoutSametypesequence(data []byte) (items []*SearchResultItem) {
	var dataPos int
	dataSize := len(data)

	for {
		t := data[dataPos]

		dataPos++

		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			end := bytes.IndexRune(data[dataPos:], '\000')

			if end < 0 { // last item
				items = append(items, &SearchResultItem{Type: rune(t), Data: data[dataPos:dataSize]})
				dataPos = dataSize
			} else {
				items = append(items, &SearchResultItem{Type: rune(t), Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
			items = append(items, &SearchResultItem{Type: rune(t), Data: data[dataPos+4 : dataPos+int(size)+5]})
			dataPos += int(size) + 5
		}

		if dataPos >= dataSize {
			break
		}
	}

	return
}

// DictName returns book name
func (d *StarDictReader) DictName() string {
	return d.Options["bookname"]
}

// NewReader returns a new Dictionary
// path - path to dictionary files
// name - name of dictionary to parse
func NewReader(path string, name string) (*StarDictReader, error) {
	d := &StarDictReader{}

	path = filepath.Clean(path)

	ifoPath := filepath.Join(path, name+".ifo")
	idxPath := filepath.Join(path, name+".idx")
	synPath := filepath.Join(path, name+".syn")

	dictDzPath := filepath.Join(path, name+".dict.dz")
	dictPath := filepath.Join(path, name+".dict")

	if _, err := os.Stat(ifoPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(idxPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(synPath); err != nil {
		synPath = ""
	}

	// we should have either .dict or .dict.dz file
	if _, err := os.Stat(dictPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		if _, errDz := os.Stat(dictDzPath); errDz != nil {
			return nil, err
		}
		dictPath = dictDzPath
	}

	info, err := ReadInfo(ifoPath)
	if err != nil {
		return nil, err
	}
	d.Info = info

	d.ifoPath = ifoPath
	d.idxPath = idxPath
	d.synPath = synPath
	d.dictPath = dictPath

	if _, ok := info.Options["sametypesequence"]; ok {
		d.decodeData = d.decodeWithSametypesequence
	} else {
		d.decodeData = d.decodeWithoutSametypesequence
	}

	return d, nil
}

type t_state int

const (
	termState t_state = iota
	offsetState
	sizeState
)

func (d *StarDictReader) Read() (func() ([]string, []*SearchResultItem), error) {
	info := d.Info
	dict, err := ReadDict(d.dictPath, info)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(d.idxPath)
	// unable to read index
	if err != nil {
		return nil, err
	}

	altsMap := map[int][]string{}
	if d.synPath != "" {
		var err error
		altsMap, err = readSyn(d.synPath)
		if err != nil {
			return nil, err
		}
	}
	// {
	// 	jsonB, _ := json.MarshalIndent(altsMap, "", "\t")
	// 	os.WriteFile("syn.json", jsonB, 0644)
	// }

	var buf [255]byte // temporary buffer
	var bufPos int
	state := termState

	var term string
	var dataOffset uint64

	maxIntBytes := info.MaxIdxBytes()

	pos := 0
	entryIndex := 0

	return func() ([]string, []*SearchResultItem) {
		alts := altsMap[int(entryIndex)]
		entryIndex++
		for {
			if pos >= len(data) {
				return nil, nil
			}
			b := data[pos]
			pos++
			buf[bufPos] = b
			if state == termState {
				if b > 0 {
					bufPos++
					continue
				}
				term = string(buf[:bufPos])
				bufPos = 0
				state = offsetState
				continue
			}
			if bufPos < maxIntBytes-1 {
				bufPos++
				continue
			}
			var num uint64
			if info.Is64 {
				num = binary.BigEndian.Uint64(buf[:maxIntBytes])
			} else {
				num = uint64(binary.BigEndian.Uint32(buf[:maxIntBytes]))
			}
			if state == offsetState {
				dataOffset = num
				bufPos = 0
				state = sizeState
				continue
			}
			// finished with one record
			bufPos = 0
			state = termState
			// log.Println("entryIndex:", entryIndex, ", alt count:", len(alts))
			terms := append([]string{term}, alts...)
			return terms, d.decodeData(dict.GetSequence(dataOffset, num))
		}
	}, nil
}
