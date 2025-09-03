package cache

import (
	"sync"
	"iranAccess/internal/models"
)

// Cache holds domain check results with thread-safe access
type Cache struct {
	mu      sync.RWMutex
	results []models.DomainCheckResult
}

// DomainCache is the global cache instance
var DomainCache = &Cache{
	results: []models.DomainCheckResult{},
}

// GetResults returns a copy of the cached results
func (c *Cache) GetResults() []models.DomainCheckResult {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	resultCopy := make([]models.DomainCheckResult, len(c.results))
	copy(resultCopy, c.results)
	return resultCopy
}

// SetResults updates the cached results
func (c *Cache) SetResults(results []models.DomainCheckResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.results = results
}
