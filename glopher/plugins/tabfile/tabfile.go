package tabfile

import (
	"bufio"
	"fmt"
	"github.com/ilius/glopher/glopher"
	"os"
	"strings"
)

func init() {
	glopher.RegisterPlugin(&pluginImp{})
}

var extentions = []string{
	".txt",
	".tab",
	".dic",
}

var readOptionTypes = []*glopher.OptionType{}
var writeOptionsTypes = []*glopher.OptionType{}

type pluginImp struct{}

func (p *pluginImp) Name() string {
	return "tabfile"
}

func (p *pluginImp) Description() string {
	return "Tabfile (.txt)"
}

func (p *pluginImp) Extentions() []string {
	return extentions
}

func (p *pluginImp) ReadOptionTypes() []*glopher.OptionType {
	return readOptionTypes
}

func (p *pluginImp) WriteOptionsTypes() []*glopher.OptionType {
	return writeOptionsTypes
}

func (p *pluginImp) Read(filename string, options ...glopher.Option) (chan *glopher.Entry, error) {
	bufferSize := 10 // TODO: get from options
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	out := make(chan *glopher.Entry, bufferSize)
	sendError := func(err error) {
		out <- &glopher.Entry{
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
			out <- &glopher.Entry{
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

func (p *pluginImp) Write(filename string, reader chan *glopher.Entry, options ...glopher.Option) error {
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
