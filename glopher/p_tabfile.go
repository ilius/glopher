package glopher

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func init() {
	RegisterPlugin(&tabfilePlug{})
}

type tabfilePlug struct{}

func (p *tabfilePlug) Name() string {
	return "tabfile"
}

func (p *tabfilePlug) Description() string {
	return "Tabfile (.txt)"
}

func (p *tabfilePlug) Extentions() []string {
	return []string{
		".txt",
		".tab",
		".dic",
	}
}

func (p *tabfilePlug) ReadOptionTypes() []*OptionType {
	return []*OptionType{}
}

func (p *tabfilePlug) WriteOptionsTypes() []*OptionType {
	return []*OptionType{}
}

func (p *tabfilePlug) Read(filename string, options ...Option) (<-chan *Entry, error) {
	bufferSize := 10 // TODO: get from options
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	out := make(chan *Entry, bufferSize)
	sendError := func(err error) {
		out <- &Entry{
			Error: err,
		}
	}
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			parts := strings.Split(line, "\t")
			if len(parts) < 1 {
				sendError(fmt.Errorf("Tabfile: bad line: %v", line))
				continue
			}
			word := parts[0]
			defi := ""
			if len(parts) == 1 {
			} else if len(parts) == 2 {
				defi = parts[1]
			} else {
				defi = strings.Join(parts[1:], "\t")
			}
			out <- &Entry{
				Word: word,
				Defi: defi,
			}
		}
		if err := scanner.Err(); err != nil {
			sendError(err)
		}
		close(out)
	}()
	return out, nil
}

func (p *tabfilePlug) Write(filename string, reader <-chan *Entry, info *StrOrderedMap, nonInfo []*Entry, options ...Option) error {
	file, err := os.Create(filename)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return err
	}
	for entry := range reader {
		if entry.Error != nil {
			return entry.Error
		}
		line := entry.Word + "\t" + entry.Defi
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
