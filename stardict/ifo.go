package stardict

import (
	"errors"
	"strconv"
	"strings"
)

const (
	I_bookname    = "bookname"
	I_wordcount   = "wordcount"
	I_description = "description"
	I_idxfilesize = "idxfilesize"
)

// Info contains dictionary options
type Info struct {
	Options map[string]string
	Version string
	Is64    bool
}

func (info Info) DictName() string {
	return info.Options[I_bookname]
}

// EntryCount returns number of words in the dictionary
func (info Info) EntryCount() (int, error) {
	num, err := strconv.ParseUint(info.Options[I_wordcount], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

func (info Info) Description() string {
	return info.Options[I_description]
}

func (info Info) IndexFileSize() uint64 {
	num, _ := strconv.ParseUint(info.Options[I_idxfilesize], 10, 64)
	return num
}

func (info Info) MaxIdxBytes() int {
	if info.Is64 {
		return 8
	}
	return 4
}

func decodeOption(str string) (key string, value string, err error) {
	a := strings.Split(str, "=")

	if len(a) < 2 {
		return "", "", errors.New("Invalid file format: " + str)
	}

	return a[0], a[1], nil
}
