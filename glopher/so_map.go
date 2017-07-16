package glopher

import (
	"sync"
)

func NewStrOrderedMap() *StrOrderedMap {
	return &StrOrderedMap{
		data: map[string]*string{},
	}
}

type StrOrderedMap struct {
	data      map[string]*string
	dataMutex sync.RWMutex
	keys      []string
	keysMutex sync.RWMutex
}

// hasKey: does not take the aliases into account
func (m *StrOrderedMap) hasKey(key string) bool {
	m.dataMutex.RLock()
	_, ok := m.data[key]
	m.dataMutex.RUnlock()
	return ok
}

func (m *StrOrderedMap) Get(key string) (string, bool) {
	m.dataMutex.RLock()
	valuePtr, ok := m.data[key]
	m.dataMutex.RUnlock()
	if !ok {
		return "", false
	}
	if valuePtr == nil {
		return "", false
	}
	return *valuePtr, true
}

func (m *StrOrderedMap) GetDefault(key string, defaultVal string) string {
	value, ok := m.Get(key)
	if ok {
		return value
	}
	return defaultVal
}

// Set does not change the order of keys if the key is already present
func (m *StrOrderedMap) Set(key string, value string) {
	if !m.hasKey(key) {
		m.keysMutex.Lock()
		m.keys = append(m.keys, key)
		m.keysMutex.Unlock()
	}
	value2 := value
	m.dataMutex.Lock()
	m.data[key] = &value2
	m.dataMutex.Unlock()
}

func (m *StrOrderedMap) Pop(key string) (string, bool) {
	value, ok := m.Get(key)
	if !ok {
		return "", false
	}

	m.dataMutex.Lock()
	delete(m.data, key)
	m.dataMutex.Unlock()

	// index := -1
	m.keysMutex.Lock()
	for i, tmpKey := range m.keys {
		if tmpKey == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			// index = i
			break
		}
	}
	m.keysMutex.Unlock()

	// if index == -1 {}

	return value, true
}

func (m *StrOrderedMap) Len() int {
	m.keysMutex.RLock()
	defer m.keysMutex.RUnlock()
	return len(m.keys)
}

func (m *StrOrderedMap) Items() [][2]string {
	m.keysMutex.RLock()
	defer m.keysMutex.RUnlock()
	m.dataMutex.RLock()
	defer m.dataMutex.RUnlock()
	items := make([][2]string, len(m.keys))
	for i, key := range m.keys {
		valuePtr, ok := m.data[key]
		if !ok || valuePtr == nil {
			continue
		}
		items[i] = [2]string{key, *valuePtr}
	}
	return items
}

// func (m *StrOrderedMap) IterKeys() chan string {
// 	out := make(chan string, 10)
// 	go func() {
// 		defer close(out)
// 		m.keysMutex.RLock()
// 		defer m.keysMutex.RUnlock()
// 		for _, key := range m.keys {
// 			out <- key
// 		}
// 	}()
// 	return out
// }

// func (m *StrOrderedMap) IterItems() <-chan [2]string {
// 	out := make(chan [2]string, 10)
// 	go func() {
// 		defer close(out)
// 		m.keysMutex.RLock()
// 		defer m.keysMutex.RUnlock()
// 		m.dataMutex.RLock()
// 		defer m.dataMutex.RUnlock()
// 		for _, key := range m.keys {
// 			valuePtr := m.data[key]
// 			if valuePtr == nil {
// 				continue
// 			}
// 			out <- [2]string{key, *valuePtr}
// 		}
// 	}()
// 	return out
// }
