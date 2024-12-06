package glopher

import (
	"fmt"
	"strings"
	"sync"
)

type PluginBase interface {
	Name() string
	Description() string
	Extensions() []string
	ReadOptionTypes() []*OptionType
	WriteOptionsTypes() []*OptionType
	Count(filename string) (int, error)
	Write(glos LimitedGlossary, filename string, options ...Option) error
}

type PluginType1 interface {
	PluginBase
	Read(filename string, options ...Option) (func() *Entry, error)
}

type PluginType2 interface {
	PluginBase
	Read2(filename string, options ...Option) (<-chan *Entry, error)
}

var (
	pluginMap      = map[string]PluginType1{}
	pluginMapMutex sync.RWMutex
)

func RegisterPluginType1(p PluginType1) {
	name := p.Name()
	if strings.ToLower(name) != name {
		panic(fmt.Sprintf("RegisterPluginType1(%#v): plugin name must be lowercase", name))
	}
	pluginMapMutex.Lock()
	defer pluginMapMutex.Unlock()
	pluginMap[name] = p
}

func RegisterPluginType2(p PluginType1) {
	panic("Not Implemented")
}

func PluginNames() []string {
	pluginMapMutex.RLock()
	defer pluginMapMutex.RUnlock()
	names := make([]string, 0, len(pluginMap))
	for name := range pluginMap {
		names = append(names, name)
	}
	return names
}

func PluginByName(name string) PluginType1 {
	pluginMapMutex.RLock()
	defer pluginMapMutex.RUnlock()
	p, ok := pluginMap[name]
	if !ok {
		return nil
	}
	return p
}
