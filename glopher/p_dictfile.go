package glopher

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func init() {
	RegisterPluginType1(&dictfilePlug{})
}

type dictfilePlug struct{}

func (p *dictfilePlug) Name() string {
	return "dictfile"
}

func (p *dictfilePlug) Description() string {
	return "dictfile (.df)"
}

func (p *dictfilePlug) Extentions() []string {
	return []string{
		".df",
	}
}

func (p *dictfilePlug) ReadOptionTypes() []*OptionType {
	return []*OptionType{}
}

func (p *dictfilePlug) WriteOptionsTypes() []*OptionType {
	return []*OptionType{}
}

func (p *dictfilePlug) Count(filename string) (int, error) {
	// file, err := os.Open(filename)
	// if err != nil {
	// 	return 0, err
	// }
	// return CountBlocks(file, Newline)
	return 0, nil
}

func (p *dictfilePlug) fixDefi(
	defi string,
) string {
	defi = strings.Replace(defi, "\n @", "\n@", -1)
	defi = strings.Replace(defi, "\n :", "\n:", -1)
	defi = strings.Replace(defi, "\n &", "\n&", -1)
	// defi = strings.Trim(defi)
	return defi
}

func (p *dictfilePlug) Read(filename string, options ...Option) (func() *Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	bufferLine := ""
	return func() *Entry {
		var words []string = nil
		defiLines := []string{}
	Loop:
		for {
			var line string
			if bufferLine != "" {
				line = bufferLine
				bufferLine = ""
			} else {
				if !scanner.Scan() {
					break
				}
				line = scanner.Text()
				if err := scanner.Err(); err != nil {
					return &Entry{
						Error: err,
					}
				}
				line = strings.TrimRight(line, "\n\r")
				if line == "" {
					continue
				}
			}
			switch line[0] {
			case '@':
				if words != nil {
					bufferLine = line
					return &Entry{
						Word:    words[0],
						Defi:    p.fixDefi(strings.Join(defiLines, "\n")),
						AltWord: words[1:],
					}
				}
				words = []string{
					strings.TrimLeft(line[1:], " "),
				}
				continue Loop
			case '&':
				words = append(words, strings.TrimLeft(line[1:], " "))
				continue Loop
			case ':':
				if len(line) < 2 {
					// Warning: bad line
					continue Loop
				}
				switch line[1] {
				case ' ':
					defiLines = append(defiLines, line[2:])
					continue Loop
				case ':':
					continue Loop
				}
			}

			if strings.HasPrefix(line, "<html>") {
				line = line[6:]
			}

			defiLines = append(defiLines, line)
		}
		if len(words) == 0 {
			return &Entry{
				Error: io.EOF,
			}
		}
		return &Entry{
			Word:    words[0],
			Defi:    p.fixDefi(strings.Join(defiLines, "\n")),
			AltWord: words[1:],
		}
	}, nil
}

// for writer
func (p *dictfilePlug) escapeDefi(defi string) string {
	defi = strings.Replace(defi, "\n@", "\n @", -1)
	defi = strings.Replace(defi, "\n:", "\n :", -1)
	defi = strings.Replace(defi, "\n&", "\n &", -1)
	return defi
}

func (p *dictfilePlug) Write(glos LimitedGlossary, filename string, options ...Option) error {
	// file, err := os.Create(filename)
	// if file != nil {
	// 	defer file.Close()
	// }
	// if err != nil {
	// 	return err
	// }
	// for entry := range glos.Iter() {
	// 	if entry.Error != nil {
	// 		return entry.Error
	// 	}
	// 	line := entry.Word + "\t" + entry.Defi
	// 	_, err := file.WriteString(line + "\n")
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
