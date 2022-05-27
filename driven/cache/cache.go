// Copyright 2022 Board of Trustees of the University of Illinois.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cacheadapter

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
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
