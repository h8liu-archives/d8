package tasks

import (
	"encoding/binary"
	"math/rand"
	"net"
	"time"

	. "d8/domain"
	"d8/packet"
	"d8/packet/rdata"
	. "d8/term"
)

// ZoneServers keep records name servers and their IPs if any
type ZoneServers struct {
	d       *Domain
	servers map[string]*server
	ips     map[uint32]*Domain
}

func (self *ZoneServers) Zone() *Domain { return self.d }

func NewZoneServers(zone *Domain) *ZoneServers {
	return &ZoneServers{
		zone,
		make(map[string]*server),
		make(map[uint32]*Domain),
	}
}

func IP2Uint(ip net.IP) uint32 {
	bytes := []byte(ip.To4())
	if bytes == nil {
		panic("not ipv4")
	}
	return binary.BigEndian.Uint32(bytes)
}

func (self *ZoneServers) add(s *server, ip net.IP) bool {
	ipIndex := IP2Uint(ip)
	if self.ips[ipIndex] != nil {
		return false
	}

	s.Add(ip)
	self.ips[ipIndex] = s.Domain
	return true
}

func (self *ZoneServers) addServer(server *Domain) *server {
	serverStr := server.String()
	s := self.servers[serverStr]
	if s == nil {
		s = newServer(server)
		self.servers[serverStr] = s
	}

	return s
}

func (self *ZoneServers) Add(server *Domain, ips ...net.IP) bool {
	s := self.addServer(server)

	anyAdded := false
	for _, ip := range ips {
		if self.add(s, ip) {
			anyAdded = true
		}
	}

	// if no ips glued, then leave it open
	if len(ips) > 0 {
		s.setResolved()
	}

	return anyAdded
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func shuffleAppend(ret, list []*server) []*server {
	n := len(list)
	order := random.Perm(n)
	for i := 0; i < n; i++ {
		ret = append(ret, list[order[i]])
	}
	return ret
}

func (self *ZoneServers) prepareOrder() []*server {
	resolved := make([]*server, 0, len(self.servers))
	unresolved := make([]*server, 0, len(self.servers))
	ret := make([]*server, 0, len(self.servers))

	for _, s := range self.servers {
		if s.Resolved() {
			resolved = append(resolved, s)
		} else {
			unresolved = append(unresolved, s)
		}
	}

	ret = shuffleAppend(ret, resolved)
	ret = shuffleAppend(ret, unresolved)
	return ret
}

func (self *ZoneServers) Serves(d *Domain) bool {
	return self.d.IsZoneOf(d)
}

func ExtractServers(p *packet.Packet, z *Domain, d *Domain, c Cursor) *ZoneServers {
	redirects := p.SelectRedirects(z, d)
	if len(redirects) == 0 {
		return nil
	}

	next := redirects[0].Domain
	ret := NewZoneServers(next)

	for _, rr := range redirects {
		if !rr.Domain.Equal(next) {
			c.Printf("// ignore different subzone: %v", rr.Domain)
			continue
		}

		ns := rdata.ToDomain(rr.Rdata)
		rrs := p.SelectIPs(ns) // glued IPs
		for _, iprr := range rrs {
			ret.Add(ns, rdata.ToIPv4(iprr.Rdata))
		}
	}

	return ret
}
