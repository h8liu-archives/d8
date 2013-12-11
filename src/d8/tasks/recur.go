package tasks

import (
	// "log"
	"fmt"
	"net"

	"d8/client"
	. "d8/domain"
	pa "d8/packet"
	"d8/packet/consts"
	. "d8/term"
)

type Recur struct {
	Domain    *Domain
	Type      uint16
	StartWith *ZoneServers
	HeadLess  bool

	Return  int          // valid when Error is not null
	Packet  *pa.Packet   // valid when Return is Okay
	EndWith *ZoneServers // valid when Return is Okay
	Answers []*pa.RR     // the records in Packet that ends the query
	Records []*pa.RR     // all the non-trivial records that helps
}

func NewRecur(d *Domain) *Recur {
	return NewRecurType(d, consts.A)
}

func NewRecurType(d *Domain, t uint16) *Recur {
	return &Recur{
		Domain: d,
		Type:   t,
	}
}

var _ Task = new(Recur)

var roots = MakeRoots()

const (
	Okay      = iota
	NotExists // domain not exists
	Lost      // no valid server reachable
)

func MakeRoots() *ZoneServers {
	ret := NewZoneServers(D("."))

	ns := func(n, ip string) {
		ret.Add(
			D(fmt.Sprintf("%s.root-servers.net", n)),
			net.ParseIP(ip),
		)
	}

	// see en.wikipedia.org/wiki/Root_name_server for reference
	// (last update: year 2012)
	ns("a", "198.41.0.4")     // Verisign
	ns("b", "192.228.79.201") // USC-ISI
	ns("c", "192.33.4.12")    // Cogent
	ns("d", "128.8.10.90")    // U Maryland
	ns("e", "192.203.230.10") // NASA
	ns("f", "192.5.5.241")    // Internet Systems Consortium
	ns("g", "192.112.36.4")   // DISA
	ns("h", "128.63.2.53")    // U.S. Army Research Lab
	ns("i", "192.36.148.17")  // Netnod
	ns("j", "198.41.0.10")    // Verisign
	ns("k", "193.0.14.129")   // RIPE NCC
	ns("l", "199.7.83.42")    // ICANN
	ns("m", "202.12.27.33")   // WIDE Project

	return ret
}

func (self *Recur) begin() *ZoneServers {
	if self.StartWith != nil {
		return self.StartWith
	}

	// cached = self.cache.Find(self.Domain.Registrar())
	return roots
}

func (self *Recur) Run(c Cursor) {
	if !self.HeadLess {
		c.Printf("recur %v %s {", self.Domain, consts.TypeString(self.Type))
		c.ShiftIn()
		defer ShiftOutWith(c, "}")
	}

	zone := self.begin()
	var e error
	self.Records = make([]*pa.RR, 0, 100)

	for zone != nil {
		zone, e = self.query(c, zone)
		if e != nil {
			return
		}
	}
}

func (self *Recur) query(c Cursor, z *ZoneServers) (*ZoneServers, error) {
	servers := z.prepareOrder()
	tried := make(map[uint32]bool)

	c.Printf("// zone: %v", z.Zone())

	for _, server := range servers {
		ips := server.IPs
		if len(ips) == 0 {
			// server is a useless alias
			continue
		}

		if server.Resolved() {
			// TODO: no glue IPs, do some query here
			// only add untried ones
			if ips == nil {
				continue
			}
		}

		for _, ip := range ips {
			ipIndex := IP2Uint(ip)
			if tried[ipIndex] {
				panic("bug")
			}
			tried[ipIndex] = true

			q := &client.Query{
				Domain:     self.Domain,
				Type:       self.Type,
				Server:     client.Server(ip),
				Zone:       z.Zone(),
				ServerName: server.Domain,
			}

			reply, e := c.Q(q)
			if e != nil {
				return nil, e // some resource limit reached
			}

			attempt := reply.Last()

			if attempt.Error != nil {
				c.Printf("// unreachable: %v, last error %v",
					server.Domain, attempt.Error)
				continue
			}

			p := attempt.Recv.Packet

			rcode := p.Rcode()
			if !(rcode == pa.RcodeOkay || rcode == pa.RcodeNameError) {
				c.Printf("// server error %s, rcode=%d", server.Domain, rcode)
			}

			ans := p.SelectAnswers(self.Domain, self.Type)
			if len(ans) > 0 {
				self.Return = Okay
				self.Packet = p
				self.Answers = ans
				self.EndWith = z

				return nil, nil
			}

			next, rrs := ExtractServers(p, z.Zone(), self.Domain, c)
			self.Records = appendAll(self.Records, rrs)
			if next == nil {
				self.Return = NotExists
				c.Print("// record does not exist")
			}
			return next, nil
		}
	}

	c.Print("// no reachable server")
	self.Return = Lost
	return nil, nil
}
