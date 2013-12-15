package tasks

import (
	"time"
)

const (
	cacheLifeSpan = time.Hour
)

type CacheEntry struct {
	zone       *ZoneServers
	expireTime time.Time
}

func NewCacheEntry(z *ZoneServers) *CacheEntry {
	return &CacheEntry{z, time.Now().Add(cacheLifeSpan)}
}

func (self *CacheEntry) Expired() bool {
	return time.Now().After(self.expireTime)
}
