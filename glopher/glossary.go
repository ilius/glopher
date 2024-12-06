package glopher

import (
	"container/heap"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
)

type LimitedGlossary interface {
	Info() *StrOrderedMap
	Iter() <-chan *Entry
}

type Glossary interface {
	LimitedGlossary
	Filename() string
	SetFilename(string)
	// DefaultDefiFormat() DefiFormat
	SetProgressBar(ProgressBar)
	Len() int
	Read(filename string, format string) error
	Write(filename string, format string) error
}

type EntryReader struct {
	Plug     PluginType1
	Next     func() *Entry
	Filename string
}

func NewGlossary() Glossary {
	pluginByExt := map[string]PluginType1{}
	for _, name := range PluginNames() {
		plug := PluginByName(name)
		for _, ext := range plug.Extensions() {
			pluginByExt[ext] = plug
		}
	}
	return &glossaryImp{
		filename:    "",
		info:        NewStrOrderedMap(),
		readers:     []EntryReader{},
		pluginByExt: pluginByExt,
	}
}

type glossaryImp struct {
	pbar              ProgressBar
	info              *StrOrderedMap
	pluginByExt       map[string]PluginType1
	filename          string
	firstEntries      EntryHeap
	readers           []EntryReader
	entryCount        int
	iterBufferSize    int
	iterating         bool
	defaultDefiFormat DefiFormat
}

func (g *glossaryImp) Filename() string {
	return g.filename
}

func (g *glossaryImp) SetFilename(filename string) {
	g.filename = filename
}

func (g *glossaryImp) Info() *StrOrderedMap {
	return g.info
}

func (g *glossaryImp) SetProgressBar(pbar ProgressBar) {
	g.pbar = pbar
}

func (g *glossaryImp) DefaultDefiFormat() DefiFormat {
	return g.defaultDefiFormat
}

func (g *glossaryImp) findPlugin(filename string, format string) PluginType1 {
	var plug PluginType1
	format = strings.TrimSpace(strings.ToLower(format))
	if format != "" {
		plug = PluginByName(format)
		if plug != nil {
			return plug
		}
		plug = g.pluginByExt["."+strings.TrimLeft(format, ".")]
		if plug != nil {
			return plug
		}
	}
	ext := strings.ToLower(filepath.Ext(filename)) // already prefixed with "."
	plug = g.pluginByExt[ext]
	if plug != nil {
		return plug
	}
	return nil
}

func (g *glossaryImp) Read(filename string, format string) error {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err // FIXME
	}
	filename = filepath.Clean(filename)
	plug := g.findPlugin(filename, format)
	if plug == nil {
		return fmt.Errorf("could not read file %#v, unknown format/extention", filename)
	}
	log.Printf("Reading from %v: %#v\n", plug.Description(), filename)
	count, err := plug.Count(filename)
	if err != nil {
		return err
	}
	g.entryCount += count
	nextFunc, err := plug.Read(filename) // TODO: options
	if err != nil {
		return err
	}
	if nextFunc == nil {
		panic("plug.Read returned nil func with no error")
	}
	reader := EntryReader{
		Plug:     plug,
		Filename: filename,
		Next:     nextFunc,
	}
	g.readers = append(g.readers, reader)
	maxNonInfo := 10 // TODO: get from options
	info, nonInfo, err := ReadInfo(reader, maxNonInfo)
	if err != nil {
		return err
	}
	for _, row := range info.Items() {
		g.info.Set(row[0], row[1])
	}
	for _, entry := range nonInfo {
		heap.Push(&g.firstEntries, entry)
	}
	if g.filename == "" {
		g.filename = filename
	}
	return nil
}

func (g *glossaryImp) Write(filename string, format string) error {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return err // FIXME
	}
	filename = filepath.Clean(filename)
	plug := g.findPlugin(filename, format)
	if plug == nil {
		return fmt.Errorf("could not write to file %#v, unknown format/extention", filename)
	}
	log.Printf("Writing to %v: %#v\n", plug.Description(), filename)
	err = plug.Write(g, filename)
	if err != nil {
		return err
	}
	return nil
}

func (g *glossaryImp) Len() int {
	return g.entryCount
}

func (g *glossaryImp) Iter() <-chan *Entry {
	if g.iterating {
		panic("Glossary.Iter: already iterating somewhere")
	}
	g.iterating = true
	out := make(chan *Entry, g.iterBufferSize)
	// sendError := func(err error) {
	// 	out <- &Entry{
	// 		Error: err,
	// 	}
	// }
	go func() {
		defer func() {
			close(out)
			g.iterating = false
			g.entryCount = 0
		}()
		for len(g.firstEntries) > 0 {
			out <- heap.Pop(&g.firstEntries).(*Entry)
		}
		if g.pbar != nil {
			total := 0
			for _, reader := range g.readers {
				c, err := reader.Plug.Count(reader.Filename)
				if err != nil {
					log.Println(err)
				}
				total += c
			}
			g.pbar.SetTotal(total)
			g.pbar.Start("Converting")
		}
		index := 0
		for _, reader := range g.readers {
			for {
				entry := reader.Next()
				if g.pbar != nil {
					g.pbar.Update(index)
				}
				index++
				if entry == nil {
					continue
				}
				if entry.Error == io.EOF {
					break
				}
				// TODO: sort: Push to EntryHeap, and Pop from it (re-assign entry)
				out <- entry
			}
		}
	}()
	return out
}
