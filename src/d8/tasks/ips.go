package tasks

import (
	"net"

	. "d8/domain"
	"d8/packet/consts"
	"d8/packet/rdata"
	. "d8/term"
)

type IPs struct {
	Domain     *Domain
	StartWith  *ZoneServers
	HeadLess   bool
	HideResult bool

	Error   error
	Return  int
	EndWith *ZoneServers
	Cnames  map[string]*Domain
	IPs     []net.IP
}

func NewIPs(d *Domain) *IPs {
	return &IPs{Domain: d}
}

func (self *IPs) findResults(recur *Recur) bool {
	if recur.Error != nil {
		self.Error = recur.Error
		return true
	}

	self.Return = recur.Return
	if recur.Return != Okay {
		return true
	}

	for _, rr := range recur.Answers {
		if rr.Type == consts.A {
			self.IPs = append(self.IPs, rdata.ToIPv4(rr.Rdata))
		}
	}

	if len(self.IPs) > 0 {
		return true
	}

	return false
}

func (self *IPs) findCnameResults(recur *Recur) bool {
	ret := false

	for _, cname := range self.Cnames {
		rrs := recur.Packet.SelectRecords(cname, consts.A)
		if len(rrs) > 0 {
			ret = true
		}
		for _, rr := range rrs {
			self.IPs = append(self.IPs, rdata.ToIPv4(rr.Rdata))
		}
	}

	return ret
}

func (self *IPs) Run(c Cursor) {
	if !self.HeadLess {
		c.Printf("ips %v {", self.Domain)
		c.ShiftIn()
	}

	self.run(c)

	if !self.HideResult {
		for _, ip := range self.IPs {
			c.Printf("// ip: %v", ip)
		}
	}

	if !self.HeadLess {
		c.ShiftOut()
		c.Print("}")
	}
}

func (self *IPs) extractCnames(recur *Recur, d *Domain) {
	rrs := recur.Packet.SelectRecords(d, consts.CNAME)

	for _, rr := range rrs {
		cname := rdata.ToDomain(rr.Rdata)
		if self.Domain.Equal(cname) || self.Cnames[cname.String()] != nil {
			continue
		}
		self.Cnames[cname.String()] = cname
		self.extractCnames(recur, cname)
	}
}

func (self *IPs) run(c Cursor) {
	recur := NewRecur(self.Domain)
	recur.HeadLess = true
	recur.StartWith = self.StartWith

	_, e := c.T(recur)
	if e != nil {
		c.Printf("error %v", e)
		self.Error = e
		return
	}
	self.EndWith = recur.EndWith

	self.IPs = make([]net.IP, 0, 10)
	if self.findResults(recur) {
		return
	}

	self.Cnames = make(map[string]*Domain)
	self.extractCnames(recur, self.Domain)

	for _, cname := range self.Cnames {
		c.Printf("// cname: %v", cname)
	}

	if len(self.Cnames) == 0 {
		panic("bug")
	}

	if self.findCnameResults(recur) {
		return
	}

	p := recur.Packet
	z := recur.EndWith

	for _, cname := range self.Cnames {
		// search for redirects
		servers := ExtractServers(p, z.Zone(), cname, c)

		// check for last result
		if servers == nil {
			if z.Serves(cname) {
				servers = z
			}
		}

		if servers == nil {
			if self.StartWith != nil && self.StartWith.Serves(cname) {
				servers = self.StartWith
			}
		}

		cnameIPs := NewIPs(cname)
		cnameIPs.HideResult = true
		cnameIPs.StartWith = servers

		_, e := c.T(cnameIPs)
		if e != nil {
			self.Error = e
			return
		}

		if cnameIPs.Error != nil {
			self.Error = cnameIPs.Error
			return
		}

		if cnameIPs.IPs != nil {
			for _, ip := range cnameIPs.IPs {
				self.IPs = append(self.IPs, ip)
			}
		}
	}
}
