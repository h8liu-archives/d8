package tasks

import (
	. "d8/domain"
	pa "d8/packet"
	. "d8/packet/consts"
	. "d8/term"
)

type Info struct {
	Domain     *Domain
	StartWith  *ZoneServers
	HeadLess   bool
	Shallow    bool
	HideResult bool

	EndWith *ZoneServers

	Cnames  []*pa.RR
	Results []*pa.RR

	Records    []*pa.RR
	RecordsMap map[string]*pa.RR

	NameServers    []*NameServer
	NameServersMap map[string]*NameServer

	Zones map[string]*ZoneServers
}

func NewInfo(d *Domain) *Info {
	return &Info{Domain: d}
}

func (self *Info) Run(c Cursor) {
	if !self.HeadLess {
		c.Printf("info %v {", self.Domain)
		c.ShiftIn()
		defer ShiftOutWith(c, "}")
	}

	ips := self.run(c)
	if c.Error() != nil {
		return
	}

	if !self.HideResult {
		ips.PrintResult(c)

		if len(self.NameServers) > 0 {
			c.Print()
			for _, ns := range self.NameServers {
				c.Printf("// %v", ns)
			}
		}

		if len(self.Records) > 0 {
			c.Print()
			for _, rr := range self.Records {
				c.Printf("// %s", rr.Digest())
			}
		}
	}
}

func (self *Info) appendAll(rrs []*pa.RR) {
	for _, rr := range rrs {
		k := rr.Digest()
		if self.RecordsMap[k] != nil {
			continue
		}
		self.RecordsMap[k] = rr
		self.Records = append(self.Records, rr)
	}
}

func (self *Info) run(c Cursor) *IPs {
	ips := NewIPs(self.Domain)
	ips.StartWith = self.StartWith
	ips.HideResult = true

	_, e := c.T(ips)
	if e != nil {
		return nil
	}

	self.EndWith = ips.EndWith

	self.Cnames, self.Results = ips.Results()

	self.RecordsMap = make(map[string]*pa.RR)
	self.Records = make([]*pa.RR, 0, 100)
	self.Zones = make(map[string]*ZoneServers)
	self.NameServers = make([]*NameServer, 0, 100)
	self.NameServersMap = make(map[string]*NameServer)

	self.appendAll(self.Cnames)
	self.appendAll(self.Results)

	self.collectInfo(ips)

	for _, z := range self.Zones {
		self.queryZone(z, c)
	}

	return ips
}

var infoTypes = []uint16{NS, MX, SOA, TXT}

func (self *Info) collectInfo(ips *IPs) {
	self._collectInfo(ips)

	if self.Shallow {
		return
	}

	for _, ips := range ips.CnameIPs {
		self._collectInfo(ips)
	}
}

func (self *Info) _collectInfo(ips *IPs) {
	for _, z := range ips.Zones {
		if z.Zone().IsRegistrar() {
			continue
		}

		for _, s := range z.List() {
			if s.IP == nil {
				continue
			}
			k := s.Key()
			if self.NameServersMap[k] != nil {
				continue
			}
			self.NameServersMap[k] = s
			self.NameServers = append(self.NameServers, s)
		}

		self.appendAll(z.Records())

		zoneStr := z.Zone().String()
		if self.Zones[zoneStr] == nil {
			self.Zones[zoneStr] = z
		}
	}
}

func (self *Info) queryZone(z *ZoneServers, c Cursor) error {
	for _, t := range infoTypes {
		recur := NewRecurType(z.Zone(), t)
		recur.StartWith = z
		_, e := c.T(recur)
		if e != nil {
			return e
		}

		self.appendAll(recur.Answers)
	}
	return nil
}
