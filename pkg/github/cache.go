package github

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cache represents a simple file-based cache for GitHub API responses
type Cache struct {
	dir      string
	mutex    sync.RWMutex
	maxAge   time.Duration
	enabled  bool
}

type cacheEntry struct {
	Data      interface{}
	Timestamp time.Time
}

// NewCache creates a new cache instance
func NewCache(dir string, maxAge time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{
		dir:     dir,
		maxAge:  maxAge,
		enabled: true,
	}, nil
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string, result interface{}) bool {
	if !c.enabled {
		return false
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	path := filepath.Join(c.dir, fmt.Sprintf("%x.json", key))
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return false
	}

	// Check if cache entry has expired
	if time.Since(entry.Timestamp) > c.maxAge {
		os.Remove(path) // Clean up expired entry
		return false
	}

	// Unmarshal cached data into result
	b, err := json.Marshal(entry.Data)
	if err != nil {
		return false
	}

	return json.Unmarshal(b, result) == nil
}

// Set stores an item in the cache
func (c *Cache) Set(key string, data interface{}) error {
	if !c.enabled {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	entry := cacheEntry{
		Data:      data,
		Timestamp: time.Now(),
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	path := filepath.Join(c.dir, fmt.Sprintf("%x.json", key))
	return os.WriteFile(path, b, 0644)
}

// Clear removes all cached items
func (c *Cache) Clear() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	dir, err := os.ReadDir(c.dir)
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if filepath.Ext(entry.Name()) == ".json" {
			os.Remove(filepath.Join(c.dir, entry.Name()))
		}
	}

	return nil
}

// Enable enables the cache
func (c *Cache) Enable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.enabled = true
}

// Disable disables the cache
func (c *Cache) Disable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.enabled = false
} 