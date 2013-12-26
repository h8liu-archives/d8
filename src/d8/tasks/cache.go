package tasks

import (
	"d8/domain"
)

type Cache struct {
	RegistrarOnly bool
	entries       map[string]*CacheEntry
}

func NewCache() *Cache {
	ret := new(Cache)
	ret.entries = make(map[string]*CacheEntry)
	ret.RegistrarOnly = true

	return ret
}

func (self *Cache) put(z *ZoneServers) {
	// self.entries[z.Zone().String()] = NewCacheEntry(z)
}

func (self *Cache) Put(z *ZoneServers) bool {
	zone := z.Zone()
	if zone.IsRoot() {
		return false // we never cache root
	}
	if self.RegistrarOnly && !zone.IsRegistrar() {
		return false
	}

	if self.Get(zone) != nil {
		return false // zone already in cache
	}
	self.put(z)
	return true
}

func (self *Cache) clean(z *domain.Domain) {
	zstr := z.String()
	entry := self.entries[z.String()]
	if entry == nil {
		return
	}

	if entry.Expired() {
		delete(self.entries, zstr)
	}
}

func (self *Cache) Clean() {
	toClean := make([]string, 0, 100)
	for k, v := range self.entries {
		if v.Expired() {
			toClean = append(toClean, k)
		}
	}

	for _, k := range toClean {
		delete(self.entries, k)
	}
}

func (self *Cache) Get(z *domain.Domain) *ZoneServers {
	self.clean(z)
	panic("todo")
	//return self.entries[z.String()].zone
}
