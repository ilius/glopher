package glopher

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func init() {
	RegisterPluginType1(&tabfilePlug{})
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

func (p *tabfilePlug) Read(filename string, options ...Option) (func() *Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	return func() *Entry {
		if !scanner.Scan() {
			return &Entry{
				Error: io.EOF,
			}
		}
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return &Entry{
				Error: err,
			}
		}
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "\t")
		if len(parts) < 1 {
			return &Entry{
				Error: fmt.Errorf("Tabfile: bad line: %v", line),
			}
		}
		word := parts[0]
		defi := ""
		if len(parts) == 1 {
		} else if len(parts) == 2 {
			defi = parts[1]
		} else {
			defi = strings.Join(parts[1:], "\t")
		}
		return &Entry{
			Word: word,
			Defi: defi,
		}
	}, nil
}

func (p *tabfilePlug) Write(glos LimitedGlossary, filename string, options ...Option) error {
	file, err := os.Create(filename)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return err
	}
	for entry := range glos.Iter() {
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