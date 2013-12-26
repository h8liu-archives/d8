package tasks

import (
	. "d8/domain"
	"time"
)

const (
	cacheLifeSpan = time.Hour
)

type CacheEntry struct {
	zone       *Domain
	ips        map[uint32]*NameServer
	resolved   map[string]*Domain
	unresolved map[string]*Domain
	expires    time.Time
}

func (self *CacheEntry) Expired() bool {
	return time.Now().After(self.expires)
}

func (self *CacheEntry) addResolved(d *Domain) {
	s := d.String()
	self.resolved[s] = d
	if self.unresolved[s] != nil {
		delete(self.unresolved, s)
	}
}

func NewEmptyCacheEntry(zone *Domain) *CacheEntry {
	return &CacheEntry{
		zone,
		make(map[uint32]*NameServer),
		make(map[string]*Domain),
		make(map[string]*Domain),
		time.Now().Add(cacheLifeSpan),
	}
}

func NewCacheEntry(zs *ZoneServers) *CacheEntry {
	ret := NewEmptyCacheEntry(zs.zone)
	ret.Add(zs)
	return ret
}

func (self *CacheEntry) Add(zs *ZoneServers) {
	if !zs.zone.Equal(self.zone) {
		panic("zone mismatch")
	}

	for key, ns := range zs.ips {
		self.ips[key] = ns
		self.addResolved(ns.Domain)
	}

	for key, d := range zs.unresolved {
		s := d.String()
		if key != s {
			panic("bug")
		}

		self.unresolved[s] = d
	}
}

func (self *CacheEntry) ZoneServers() *ZoneServers {
	panic("todo")
}
