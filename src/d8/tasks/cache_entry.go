package tasks

import (
	. "d8/domain"
	"time"
)

const (
	cacheLifeSpan = time.Hour
)

type cacheEntry struct {
	zone       *Domain
	ips        map[uint32]*NameServer
	resolved   map[string]*Domain
	unresolved map[string]*Domain
	expires    time.Time
}

func (self *cacheEntry) Expired() bool {
	return time.Now().After(self.expires)
}

func (self *cacheEntry) addResolved(d *Domain) {
	s := d.String()
	self.resolved[s] = d
	if self.unresolved[s] != nil {
		delete(self.unresolved, s)
	}
}

func emptyCacheEntry(zone *Domain) *cacheEntry {
	return &cacheEntry{
		zone,
		make(map[uint32]*NameServer),
		make(map[string]*Domain),
		make(map[string]*Domain),
		time.Now().Add(cacheLifeSpan),
	}
}

func newCacheEntry(zs *ZoneServers) *cacheEntry {
	ret := emptyCacheEntry(zs.zone)
	ret.Add(zs)
	return ret
}

func (self *cacheEntry) Add(zs *ZoneServers) {
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

func (self *cacheEntry) ZoneServers() *ZoneServers {
	ret := NewZoneServers(self.zone)

	for _, ns := range self.ips {
		ret.Add(ns.Domain, ns.IP)
	}

	for _, d := range self.unresolved {
		ret.Add(d)
	}

	return ret
}
