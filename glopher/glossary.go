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
	Read(filename string, format string) error
	Write(filename string, format string) error
}

func NewGlossary() Glossary {
	pluginByExt := map[string]PluginType1{}
	for _, name := range PluginNames() {
		plug := PluginByName(name)
		for _, ext := range plug.Extentions() {
			pluginByExt[ext] = plug
		}
	}
	return &glossaryImp{
		filename:    "",
		info:        NewStrOrderedMap(),
		readers:     []func() *Entry{},
		pluginByExt: pluginByExt,
	}
}

type glossaryImp struct {
	filename     string
	info         *StrOrderedMap
	firstEntries EntryHeap // minimal, ideally empty
	readers      []func() *Entry

	defaultDefiFormat DefiFormat
	// entryFilters = []*EntryFilter
	// sortKey *func....
	// sortCacheSize int
	iterBufferSize int

	pluginByExt map[string]PluginType1
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
		return fmt.Errorf("Could not read file %#v, unknown format/extention", filename)
	}
	log.Printf("Reading from %v: %#v\n", plug.Description(), filename)
	reader, err := plug.Read(filename) // TODO: options
	if err != nil {
		return err
	}
	if reader == nil {
		panic("plug.Read returned nil func with no error")
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
		return fmt.Errorf("Could not write to file %#v, unknown format/extention", filename)
	}
	log.Printf("Writing to %v: %#v\n", plug.Description(), filename)
	err = plug.Write(g, filename)
	if err != nil {
		return err
	}
	return nil
}

func (g *glossaryImp) Iter() <-chan *Entry {
	out := make(chan *Entry, g.iterBufferSize)
	// sendError := func(err error) {
	// 	out <- &Entry{
	// 		Error: err,
	// 	}
	// }
	go func() {
		defer close(out)
		for len(g.firstEntries) > 0 {
			out <- heap.Pop(&g.firstEntries).(*Entry)
		}
		for _, reader := range g.readers {
			for {
				entry := reader()
				// TODO: update progressbar
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
