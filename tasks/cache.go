package tasks

import (
	. "github.com/h8liu/d8/domain"
	"time"
	// "log"
)

type Cache struct {
	RegistrarOnly bool
	entries       map[string]*cacheEntry

	puts chan *cachePut
	gets chan *cacheGet
}

type cachePut struct {
	zs    *ZoneServers
	reply chan bool
}

type cacheGet struct {
	zone  *Domain
	reply chan *ZoneServers
}

func NewCache() *Cache {
	ret := new(Cache)
	ret.entries = make(map[string]*cacheEntry)
	ret.RegistrarOnly = true

	ret.puts = make(chan *cachePut)
	ret.gets = make(chan *cacheGet)

	go ret.serve()

	return ret
}

func (self *Cache) serve() {
	ticker := time.Tick(time.Minute * 5)

	for {
		select {
		case put := <-self.puts:
			put.reply <- self.put(put.zs)
		case get := <-self.gets:
			get.reply <- self.get(get.zone)
		case <-ticker:
			self.clean()
		}
	}
}

func (self *Cache) put(z *ZoneServers) bool {
	zone := z.Zone()
	if zone.IsRoot() {
		return false // we never cache root
	}
	if self.RegistrarOnly && !zone.IsRegistrar() {
		return false
	}

	key := zone.String()
	entry := self.entries[key]
	if entry == nil {
		self.entries[key] = newCacheEntry(z)
	} else {
		entry.Add(z)
	}

	return true
}

func (self *Cache) cleanZone(z *Domain) {
	zstr := z.String()
	entry := self.entries[z.String()]
	if entry == nil {
		return
	}

	if entry.Expired() {
		delete(self.entries, zstr)
	}
}

func (self *Cache) clean() {
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

func (self *Cache) get(z *Domain) *ZoneServers {
	self.cleanZone(z)

	entry := self.entries[z.String()]
	if entry == nil {
		return nil
	}

	return entry.ZoneServers()
}

func (self *Cache) Get(z *Domain) *ZoneServers {
	c := make(chan *ZoneServers)
	self.gets <- &cacheGet{z, c}
	return <-c
}

func (self *Cache) Put(zs *ZoneServers) bool {
	c := make(chan bool)
	self.puts <- &cachePut{zs, c}
	return <-c
}
