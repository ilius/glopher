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

func (p *tabfilePlug) Count(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	return CountBlocks(file, Newline)
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
				Error: fmt.Errorf("tabfile: bad line: %v", line),
			}
		}
		wordsStr := parts[0]
		defi := ""
		if len(parts) == 1 {
		} else if len(parts) == 2 {
			defi = parts[1]
		} else {
			defi = strings.Join(parts[1:], "\t")
		}
		words := SplitByBarUnescapeNTB(wordsStr)
		defi = UnescapeNTB(defi, false)
		return &Entry{
			Word:    words[0],
			Defi:    defi,
			AltWord: words[1:],
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
	for _, pair := range glos.Info().Items() {
		line := EscapeNTB("##"+pair[0], false) + "\t" + EscapeNTB(pair[1], false)
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	for entry := range glos.Iter() {
		if entry.Error != nil {
			return entry.Error
		}
		words := append([]string{entry.Word}, entry.AltWord...)
		line := JoinByBarEscapeNTB(words) + "\t" + EscapeNTB(entry.Defi, false)
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
