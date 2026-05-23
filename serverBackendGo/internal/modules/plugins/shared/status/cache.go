package status

import "sync"

// Cache tracks per-customer disabled plugin IDs (Java PluginStatusCache subset).
type Cache struct {
	mu   sync.RWMutex
	data map[int64]map[int64]struct{}
}

func NewCache() *Cache {
	return &Cache{data: make(map[int64]map[int64]struct{})}
}

func (c *Cache) SetDisabled(customerID int64, pluginIDs []int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	m := make(map[int64]struct{}, len(pluginIDs))
	for _, id := range pluginIDs {
		m[id] = struct{}{}
	}
	c.data[customerID] = m
}

func (c *Cache) IsDisabled(customerID, pluginID int64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	m, ok := c.data[customerID]
	if !ok {
		return false
	}
	_, ok = m[pluginID]
	return ok
}
