package cacheadapter

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	// DefaultExpireDuration default duration
	DefaultExpireDuration = 2 * time.Minute
)

// CacheAdapter structure
type CacheAdapter struct {
	cache *cache.Cache
}

// NewCacheAdapter creates new instance
func NewCacheAdapter() *CacheAdapter {
	cache := cache.New(DefaultExpireDuration, DefaultExpireDuration)

	return &CacheAdapter{
		cache: cache,
	}
}

// GetTwitterPosts Gets twitter posts
func (s *CacheAdapter) GetTwitterPosts(count int) map[string]interface{} {
	obj, _ := s.cache.Get(fmt.Sprintf("twitter.posts.%d", count))
	if obj != nil {
		return obj.(map[string]interface{})
	}
	return nil
}

// SetTwitterPosts Sets twitter posts
func (s *CacheAdapter) SetTwitterPosts(count int, posts map[string]interface{}) map[string]interface{} {
	key := fmt.Sprintf("twitter.posts.%d", count)
	if posts == nil {
		s.cache.Delete(key)
	} else {
		s.cache.Set(key, posts, cache.DefaultExpiration)
	}
	return posts
}
