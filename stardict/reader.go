package stardict

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/ilius/glopher/common"
)

// ReadInfo reads ifo file and collects dictionary options
func ReadInfo(filename string) (info *Info, err error) {
	reader, err := os.Open(filename)
	if err != nil {
		return
	}

	defer common.Close(reader)

	r := bufio.NewReader(reader)

	_, err = r.ReadString('\n')

	if err != nil {
		return
	}

	version, err := r.ReadString('\n')
	if err != nil {
		return
	}

	key, value, err := decodeOption(version[:len(version)-1])
	if err != nil {
		return
	}

	if key != "version" {
		err = errors.New("version missing (should be on second line)")
		return
	}

	if value != "2.4.2" && value != "3.0.0" {
		err = errors.New("stardict version should be either 2.4.2 or 3.0.0")
		return
	}

	info = &Info{}

	info.Version = value

	info.Options = make(map[string]string)

	for {
		option, err := r.ReadString('\n')

		if err != nil && err != io.EOF {
			return info, err
		}

		if err == io.EOF && len(option) == 0 {
			break
		}

		key, value, err = decodeOption(option[:len(option)-1])

		if err != nil {
			return info, err
		}

		info.Options[key] = value

		if err == io.EOF {
			break
		}
	}

	if bits, ok := info.Options[I_idxoffsetbits]; ok {
		if bits == "64" {
			info.Is64 = true
		}
	} else {
		info.Is64 = false
	}

	return
}

// dictionaryImp stardict dictionary
type StarDictReader struct {
	*Info

	dict     *Dict
	ifoPath  string
	idxPath  string
	dictPath string
	synPath  string
	// resDir   string

	decodeData func(data []byte) []*ArticleItem
}

func (d *StarDictReader) Loaded() bool {
	return d.dict != nil
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

func (d *StarDictReader) decodeWithSametypesequence(data []byte) (items []*ArticleItem) {
	seq := d.Options[I_sametypesequence]

	seqLen := len(seq)

	var dataPos int
	dataSize := len(data)

	for i, t := range seq {
		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			// if last seq item
			if i == seqLen-1 {
				items = append(items, &ArticleItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				end := bytes.IndexRune(data[dataPos:], '\000')
				items = append(items, &ArticleItem{Type: t, Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			if i == seqLen-1 {
				items = append(items, &ArticleItem{Type: t, Data: data[dataPos:dataSize]})
			} else {
				size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
				items = append(items, &ArticleItem{Type: t, Data: data[dataPos+4 : dataPos+int(size)+5]})
				dataPos += int(size) + 5
			}
		}
	}

	return
}

func (d *StarDictReader) decodeWithoutSametypesequence(data []byte) (items []*ArticleItem) {
	var dataPos int
	dataSize := len(data)

	for {
		t := data[dataPos]

		dataPos++

		switch t {
		case 'm', 'l', 'g', 't', 'x', 'y', 'k', 'w', 'h', 'r':
			end := bytes.IndexRune(data[dataPos:], '\000')

			if end < 0 { // last item
				items = append(items, &ArticleItem{Type: rune(t), Data: data[dataPos:dataSize]})
				dataPos = dataSize
			} else {
				items = append(items, &ArticleItem{Type: rune(t), Data: data[dataPos : dataPos+end+1]})
				dataPos += end + 1
			}
		case 'W', 'P':
			size := binary.BigEndian.Uint32(data[dataPos : dataPos+4])
			items = append(items, &ArticleItem{Type: rune(t), Data: data[dataPos+4 : dataPos+int(size)+5]})
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
	return d.Options[I_bookname]
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

	if _, ok := info.Options[I_sametypesequence]; ok {
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

func (r *StarDictReader) readSyn() (map[int][]string, error) {
	if r.synPath == "" {
		return nil, nil
	}
	return readSyn(r.synPath)
}

func (r *StarDictReader) Read() (func() ([]string, []*ArticleItem), error) {
	info := r.Info
	dict, err := ReadDict(r.dictPath, info)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(r.idxPath)
	if err != nil {
		return nil, err
	}

	synMap, err := r.readSyn()
	if err != nil {
		return nil, err
	}

	var buf [255]byte // temporary buffer
	var bufPos int
	state := termState

	var term string
	var dataOffset uint64

	maxIntBytes := info.MaxIdxBytes()

	pos := 0
	entryIndex := 0

	return func() ([]string, []*ArticleItem) {
		synTerms := synMap[entryIndex]
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
			terms := append([]string{term}, synTerms...)
			return terms, r.decodeData(dict.GetSequence(dataOffset, num))
		}
	}, nil
}
