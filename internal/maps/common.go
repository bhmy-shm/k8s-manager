package maps

import "sync"

type MapItems []*mapItem

type mapItem struct {
	key   string
	value interface{}
}

func convertToMapItems(m *sync.Map) (int64, MapItems) {
	var count int64
	items := make(MapItems, 0)
	m.Range(func(key, value interface{}) bool {
		items = append(items, &mapItem{key: key.(string), value: value})
		count++
		return true
	})
	return count, items
}

func (m MapItems) Len() int {
	return len(m)
}

func (m MapItems) Less(i, j int) bool {
	return m[i].key < m[j].key
}

func (m MapItems) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
