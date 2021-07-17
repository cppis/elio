package elio

import (
	"encoding/json"

	cmap "github.com/orcaman/concurrent-map"
)

// // upsertCallback update or insert
// func upsertCallback(exist bool, old interface{}, new interface{}) interface{} {
// 	if !exist {
// 		ss := NewSafeSlice()
// 		ss.Append(new)
// 		return ss
// 	}
// 	old.(*SafeSlice).Append(new)
// 	return old
// }

// UnsafeMap safe map
type UnsafeMap struct {
	Map map[uint64]interface{}
}

// UnsafeMapItem contains a key/value pair item of a concurrent map
type UnsafeMapItem struct {
	Key   uint64
	Value interface{}
}

// NewUnsafeMap creates a new unsafe map
func NewUnsafeMap() *UnsafeMap {
	m := new(UnsafeMap)
	if nil != m {
		m.Init()
	}

	return m
}

// Init init
func (m *UnsafeMap) Init() {
	m.Map = make(map[uint64]interface{})
}

// Count count
func (m *UnsafeMap) Count() int {
	return len(m.Map)
}

// Get retrieves the value for a concurrent map item
func (m *UnsafeMap) Get(key uint64) (interface{}, bool) {
	value, ok := m.Map[key]

	return value, ok
}

// Set adds an item to a concurrent map
func (m *UnsafeMap) Set(key uint64, value interface{}) {
	m.Map[key] = value
}

// Del deletes an item to a concurrent map
func (m *UnsafeMap) Del(key uint64) {
	delete(m.Map, key)
}

// // Upsert update or insert
// func (m *UnsafeMap) Upsert(key interface{}, new interface{}) {
// 	m.UpsertByCallback(key, new, upsertCallback)
// }

// Upsert update or insert
func (m *UnsafeMap) Upsert(key uint64, new interface{}, callback cmap.UpsertCb) (value interface{}) {
	var old interface{}
	var ok bool

	old, ok = m.Map[key]

	value = callback(ok, old, new)
	m.Map[key] = value
	return value
}

// Fetch fetch map and renew
func (m *UnsafeMap) Fetch() map[uint64]interface{} {
	old := m.Map
	m.Map = make(map[uint64]interface{})
	return old
}

// Iterate iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (m *UnsafeMap) Iterate() <-chan UnsafeMapItem {
	c := make(chan UnsafeMapItem)

	f := func() {
		for k, v := range m.Map {
			c <- UnsafeMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}

// ToSlice to slice
func (m *UnsafeMap) ToSlice() []UnsafeMapItem {
	s := []UnsafeMapItem{}

	for k, v := range m.Map {
		s = append(s, UnsafeMapItem{k, v})
	}

	return s
}

// MarshalJSON marshal json
func (m *UnsafeMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Map)
}
