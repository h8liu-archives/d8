package tasks

import (
	. "d8/domain"
	pa "d8/packet"
	. "d8/packet/consts"
	. "d8/term"
)

type Info struct {
	Domain    *Domain
	StartWith *ZoneServers
	HeadLess  bool
	Shallow   bool

	EndWith *ZoneServers
	Records []*pa.RR

	Followed map[string]*ZoneServers
}

func NewInfo(d *Domain) *Info {
	return &Info{Domain: d}
}

func (self *Info) Run(c Cursor) {
	if !self.HeadLess {
		c.Printf("records %v {", self.Domain)
		c.ShiftIn()
		defer ShiftOutWith(c, "}")
	}

	self.run(c)
	if c.Error() != nil {
		return
	}

	self.dedup()

	c.Print()
	for _, rr := range self.Records {
		c.Printf("// %s", rr.Digest())
	}
}

func (self *Info) dedup() {
	m := make(map[string]bool)
	records := make([]*pa.RR, 0, len(self.Records))

	for _, rr := range self.Records {
		s := rr.Digest()
		if m[s] {
			continue
		}
		m[s] = true
		records = append(records, rr)
	}

	self.Records = records
}

func appendAll(list []*pa.RR, rrs []*pa.RR) []*pa.RR {
	for _, rr := range rrs {
		list = append(list, rr)
	}
	return list
}

func (self *Info) appendAll(rrs []*pa.RR) {
	self.Records = appendAll(self.Records, rrs)
}

func (self *Info) run(c Cursor) {
	ips := NewIPs(self.Domain)
	ips.StartWith = self.StartWith
	ips.HideResult = true

	_, e := c.T(ips)
	if e != nil {
		return
	}

	self.EndWith = ips.EndWith
	self.Records = make([]*pa.RR, 0, 100)

	self.appendAll(ips.CnameRecords)
	self.appendAll(ips.Records)
	self.appendAll(ips.ServerRecords)

	self.Followed = make(map[string]*ZoneServers)
	self.followUp(ips, c)

	ips.PrintResult(c)
}

var otherTypes = []uint16{NS, MX, SOA, TXT}

func (self *Info) check(z *ZoneServers) bool {
	zoneStr := z.Zone().String()
	if self.Followed[zoneStr] != nil {
		return true
	}
	self.Followed[zoneStr] = z
	return false
}

func (self *Info) _followUp(ips *IPs, c Cursor) error {
	if ips.Return == Okay || ips.Return == NotExists {
		z := ips.EndWith
		if self.check(z) {
			return nil
		}

		for _, t := range otherTypes {
			recur := NewRecurType(z.Zone(), t)
			recur.StartWith = z
			_, e := c.T(recur)
			if e != nil {
				return e
			}

			self.appendAll(recur.Answers)
			self.appendAll(recur.Records)
		}
	}

	return nil
}

func (self *Info) followUp(ips *IPs, c Cursor) error {
	e := self._followUp(ips, c)
	if e != nil {
		return e
	}

	if !self.Shallow {
		for _, cnameIPs := range ips.CnameEndIPs {
			e = self.followUp(cnameIPs, c)
			if e != nil {
				return e
			}
		}
	}

	return nil
}
