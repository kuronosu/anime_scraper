package scrape

import "sync"

type DetailsCollector interface {
	// Collect(url string, data T) error
	CollectField(url string, field string, data interface{}) error
}

type MemoryDetailsCollector struct {
	Items map[string]map[string]interface{}
	mutex *sync.RWMutex
}

func NewMemoryDetailsCollector() *MemoryDetailsCollector {
	return &MemoryDetailsCollector{
		Items: make(map[string]map[string]interface{}),
		mutex: &sync.RWMutex{},
	}
}

func (c *MemoryDetailsCollector) CollectField(url string, field string, data interface{}) error {
	c.mutex.Lock()
	if c.Items[url] == nil {
		c.Items[url] = make(map[string]interface{})
	}
	c.Items[url][field] = data
	c.mutex.Unlock()
	return nil
}
