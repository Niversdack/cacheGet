package cache

import (
	"sort"
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	items    map[string]Item
	maxValue uint64
}

type Item struct {
	Value    []byte
	CreateAt time.Time
}

func New(maxValue uint64) *Cache {
	items := make(map[string]Item)
	cache := Cache{
		items:    items,
		maxValue: maxValue,
	}
	return &cache
}

func (c *Cache) Set(key string, value []byte) {

	if int(c.maxValue) == len(c.items) {
		c.del()
	}
	c.Lock()

	defer c.Unlock()

	c.items[key] = Item{
		Value:    value,
		CreateAt: time.Now(),
	}

}

type reviewsData struct {
	key  string
	time time.Time
}

type timeSlice []reviewsData

func (p timeSlice) len() int {
	return len(p)
}

func (p timeSlice) less(i, j int) bool {
	return p[i].time.Before(p[j].time)
}

func (p timeSlice) swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (c *Cache) sortByTime() []reviewsData {
	c.RLock()

	defer c.RUnlock()
	result := make(timeSlice, 0, len(c.items))
	for key, d := range c.items {
		result = append(result, reviewsData{
			key:  key,
			time: d.CreateAt,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].time.Before(result[j].time)
	})
	return result
}
func (c *Cache) del() {
	c.RLock()

	defer c.RUnlock()
	sortItems := c.sortByTime()
	delete(c.items, sortItems[0].key)
}
func (c *Cache) Get(key string) (*[]byte, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	// cache not found
	if !found {
		return nil, false
	}

	return &item.Value, true
}
