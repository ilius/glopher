package glopher

import (
	"sync"
)

type GPlugin interface {
	Name() string
	Description() string
	Extentions() []string
	ReadOptionTypes() []*OptionType
	WriteOptionsTypes() []*OptionType
	Read(filename string, options ...Option) (chan *Entry, error)
	Write(filename string, reader chan *Entry, options ...Option) error
}

var pluginMap = map[string]GPlugin{}
var pluginMapMutex sync.RWMutex

func RegisterPlugin(p GPlugin) {
	pluginMapMutex.Lock()
	defer pluginMapMutex.Unlock()
	pluginMap[p.Name()] = p
}

func PluginNames() []string {
	pluginMapMutex.RLock()
	defer pluginMapMutex.RUnlock()
	names := make([]string, 0, len(pluginMap))
	for name, _ := range pluginMap {
		names = append(names, name)
	}
	return names
}

func PluginByName(name string) GPlugin {
	pluginMapMutex.RLock()
	defer pluginMapMutex.RUnlock()
	p, ok := pluginMap[name]
	if !ok {
		return nil
	}
	return p
}
