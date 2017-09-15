package cache

import (
	"fmt"
	"sync"
)

// Record contains a DNS record's info
type Record struct {
	Qname    string `json:"qname"`
	Qtype    string `json:"qtype"`
	TTL      int    `json:"ttl"`
	Content  string `json:"content"`
	DomainID int    `json:"domain_id"`
}

// Cache contains the methods to save and get records to/from cache
type Cache struct {
	content map[string][]Record
	mutex   *sync.Mutex
}

// New sends back a new Cache
func New() Cache {
	returnValue := Cache{}
	returnValue.mutex = &sync.Mutex{}
	returnValue.content = map[string][]Record{}
	return returnValue
}

// NewRecord sends back a new Record
func NewRecord() Record {
	return Record{DomainID: -1}
}

// Add adds an array of records to a specific name
func (c *Cache) Add(name string, records []Record) {
	fmt.Println(name)
	fmt.Println(records)
	if c.mutex == nil {
		c.mutex = &sync.Mutex{}
	}
	c.mutex.Lock()
	c.content[name] = records
	c.mutex.Unlock()
}

// Get fetches content from the cache
func (c *Cache) Get(name string) []Record {
	c.mutex.Lock()
	returnVal := c.content[name]
	c.mutex.Unlock()
	return returnVal
}
