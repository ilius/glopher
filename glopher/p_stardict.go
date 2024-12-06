package glopher

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/ilius/glopher/stardict"
)

func init() {
	RegisterPluginType1(&stardictPlug{})
}

type stardictPlug struct {
	reader *stardict.StarDictReader
}

func (p *stardictPlug) Name() string {
	return "stardict"
}

func (p *stardictPlug) Description() string {
	return "StarDict (.ifo)"
}

func (p *stardictPlug) Extensions() []string {
	return []string{
		".ifo",
	}
}

func (p *stardictPlug) ReadOptionTypes() []*OptionType {
	return []*OptionType{}
}

func (p *stardictPlug) WriteOptionsTypes() []*OptionType {
	return []*OptionType{}
}

func (p *stardictPlug) open(filename string) error {
	dictDir := filepath.Dir(filename)
	name := filepath.Base(filename)
	r, err := stardict.NewReader(
		dictDir,
		name[:len(name)-len(".ifo")],
	)
	if err != nil {
		return err
	}
	p.reader = r
	return nil
}

func (p *stardictPlug) Count(filename string) (int, error) {
	if p.reader == nil {
		err := p.open(filename)
		if err != nil {
			return 0, err
		}
	}
	return p.reader.EntryCount()
}

func (p *stardictPlug) Read(filename string, options ...Option) (func() *Entry, error) {
	err := p.open(filename)
	if err != nil {
		return nil, err
	}
	infoList := [][2]string{}
	for key, value := range p.reader.Info.Options {
		if value == "" {
			continue
		}
		switch key {
		case stardict.I_bookname:
			key = "name"
		}
		infoList = append(infoList, [2]string{key, value})
	}
	next, err := p.reader.Read()
	if err != nil {
		return nil, err
	}
	return func() *Entry {
		if len(infoList) > 0 {
			pair := infoList[0]
			infoList = infoList[1:]
			return &Entry{
				Word:   pair[0],
				Defi:   pair[1],
				IsInfo: true,
			}
		}
		terms, defiItems := next()
		if terms == nil {
			return &Entry{
				Error: io.EOF,
			}
		}
		if len(defiItems) == 0 {
			return nil
		}
		// defiFormats := map[DefiFormat]bool{}
		defiFormat := DefiFormatHTML
		for _, item := range defiItems {
			defiFormat = DefiFormat(item.Type)
			break
		}
		defiParts := make([]string, len(defiItems))
		for i, item := range defiItems {
			defiParts[i] = string(item.Data)
		}
		return &Entry{
			Word:       terms[0],
			Defi:       strings.Join(defiParts, "\n"),
			AltWord:    terms[1:],
			DefiFormat: defiFormat,
		}
	}, nil
}

func (p *stardictPlug) Write(glos LimitedGlossary, filename string, options ...Option) error {
	return fmt.Errorf("writing stardict is not implemented")
}
