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
	zone       *Domain
	resolved   map[uint32]*NameServer
	unresolved map[string]*Domain
}

func (self *ZoneServers) Zone() *Domain { return self.zone }

func NewZoneServers(zone *Domain) *ZoneServers {
	return &ZoneServers{
		zone,
		make(map[uint32]*NameServer),
		make(map[string]*Domain),
	}
}

func IP2Uint(ip net.IP) uint32 {
	bytes := []byte(ip.To4())
	if bytes == nil {
		panic("not ipv4")
	}
	return binary.BigEndian.Uint32(bytes)
}

func (self *ZoneServers) addUnresolved(server *Domain) bool {
	s := server.String()
	if _, found := self.unresolved[s]; found {
		return false
	}

	self.unresolved[s] = server
	return true
}

func (self *ZoneServers) add(server *Domain, ip net.IP) bool {
	index := IP2Uint(ip)
	if _, found := self.resolved[index]; found {
		return false
	}

	self.resolved[index] = &NameServer{
		Zone:   self.zone,
		Domain: server,
		IP:     ip,
		Glued:  true,
	}

	return true
}

func (self *ZoneServers) Add(server *Domain, ips ...net.IP) bool {
	if len(ips) == 0 {
		return self.addUnresolved(server)
	}

	anyAdded := false
	for _, ip := range ips {
		if self.add(server, ip) {
			anyAdded = true
		}
	}

	return anyAdded
}

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

func shuffleAppend(ret, list []*NameServer) []*NameServer {
	n := len(list)
	order := random.Perm(n)
	for i := 0; i < n; i++ {
		ret = append(ret, list[order[i]])
	}
	return ret
}

func (self *ZoneServers) prepareOrder() []*NameServer {
	ret := make([]*NameServer, 0, len(self.resolved)+len(self.unresolved))

	resolved := make([]*NameServer, 0, len(self.resolved))
	for _, s := range self.resolved {
		resolved = append(resolved, s)
	}

	unresolved := make([]*NameServer, 0, len(self.unresolved))
	for _, d := range self.unresolved {
		unresolved = append(unresolved, &NameServer{
			Zone:   self.zone,
			Domain: d,
			IP:     nil,
			Glued:  false,
		})
	}

	ret = shuffleAppend(ret, resolved)
	ret = shuffleAppend(ret, unresolved)

	return ret
}

func (self *ZoneServers) Serves(d *Domain) bool {
	return self.zone.IsZoneOf(d)
}

func ExtractServers(p *packet.Packet, z *Domain, d *Domain,
	c Cursor) (*ZoneServers, []*packet.RR) {

	redirects := p.SelectRedirects(z, d)
	if len(redirects) == 0 {
		return nil, nil
	}

	next := redirects[0].Domain

	var records []*packet.RR
	if !next.IsRegistrar() {
		records = make([]*packet.RR, 0, 100)
		records = appendAll(records, redirects)
	}

	ret := NewZoneServers(next)

	for _, rr := range redirects {
		if !rr.Domain.Equal(next) {
			c.Printf("// ignore different subzone: %v", rr.Domain)
			continue
		}

		ns := rdata.ToDomain(rr.Rdata)

		rrs := p.SelectIPs(ns) // glued IPs
		ips := make([]net.IP, 0, len(rrs))
		for _, rr := range rrs {
			ips = append(ips, rdata.ToIPv4(rr.Rdata))
		}
		ret.Add(ns, ips...)

		if !next.IsRegistrar() {
			records = appendAll(records, rrs)
		}
	}

	return ret, records
}
