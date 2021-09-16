package cacheadapter

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"strings"
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
func (s *CacheAdapter) GetTwitterPosts(userID string, twitterQueryParams string) map[string]interface{} {
	var key = fmt.Sprintf("twitter.%s.params.%s", userID, twitterQueryParams)
	obj, _ := s.cache.Get(key)
	if obj != nil {
		return obj.(map[string]interface{})
	}
	return nil
}

// SetTwitterPosts Sets twitter posts
func (s *CacheAdapter) SetTwitterPosts(userID string, twitterQueryParams string, posts map[string]interface{}) map[string]interface{} {
	var key = fmt.Sprintf("twitter.%s.params.%s", userID, twitterQueryParams)

	if posts == nil {
		s.cache.Delete(key)
	} else {
		s.cache.Set(key, posts, cache.DefaultExpiration)
	}
	return posts
}

// ClearTwitterCacheForUser clears cache for specified user
func (s *CacheAdapter) ClearTwitterCacheForUser(userID string) {
	var prefix = fmt.Sprintf("twitter.%s", userID)
	for key := range s.cache.Items() {
		if strings.HasPrefix(key, prefix) {
			s.cache.Delete(key)
		}
	}
}
