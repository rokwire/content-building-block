package cacheadapter

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"strconv"
	"strings"
	"time"
)

// CacheAdapter structure
type CacheAdapter struct {
	cache *cache.Cache
}

// NewCacheAdapter creates new instance
func NewCacheAdapter(defaultCacheExpirationSeconds string) *CacheAdapter {

	val, err := strconv.ParseInt(defaultCacheExpirationSeconds, 0, 64)
	var duration time.Duration
	if val == 0 || err != nil {
		duration = 120 * time.Second
	} else {
		duration = time.Duration(val) * time.Second
	}

	cache := cache.New(duration, duration)

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
