package storage

import (
	"log2metric/types"
	"sync"
)

type CacheStorage struct {
	sync.RWMutex
	Data map[types.ServiceId]types.MetricLabel
}

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{Data: make(map[types.ServiceId]types.MetricLabel)}
}

func (c *CacheStorage) Set(id types.ServiceId, labels types.MetricLabel) {
	c.Lock()
	defer c.Unlock()

	labels.Value = c.GetValue(id) + 1
	c.Data[id] = labels
}

func (c *CacheStorage) List() map[types.ServiceId]types.MetricLabel {
	c.RLock()
	defer c.RUnlock()

	return c.Data
}

func (c *CacheStorage) GetValue(id types.ServiceId) float64 {
	return c.Data[id].Value
}
