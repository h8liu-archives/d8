package tasks

import (
	. "d8/domain"
	pa "d8/packet"
	"d8/packet/consts"
	"d8/packet/rdata"
	. "d8/term"
)

type IPs struct {
	Domain     *Domain
	StartWith  *ZoneServers
	HeadLess   bool
	HideResult bool

	// inherit from the initializing Recur Task
	Return  int
	Packet  *pa.Packet
	EndWith *ZoneServers

	// All the cnames that have ever tracked
	CnameEndpoints []*Domain
	CnameTraceBack map[string]*Domain // in and out, inherit from father IPs
	CnameEndIPs    []*IPs

	CnameRecords  []*pa.RR
	Records       []*pa.RR
	ServerRecords []*pa.RR
}

func NewIPs(d *Domain) *IPs {
	return &IPs{Domain: d}
}

// Look for Query error or A records in Answer
func (self *IPs) findResults(recur *Recur) bool {
	if recur.Return != Okay {
		return true
	}

	for _, rr := range recur.Answers {
		if rr.Type == consts.A {
			self.Records = append(self.Records, rr)
		}
	}

	if len(self.Records) > 0 {
		return true
	}

	return false
}

func (self *IPs) findCnameResults(recur *Recur) (unresolved []*Domain) {
	unresolved = make([]*Domain, 0, len(self.CnameEndpoints))

	for _, cname := range self.CnameEndpoints {
		rrs := recur.Packet.SelectRecords(cname, consts.A)
		if len(rrs) == 0 {
			unresolved = append(unresolved, cname)
			continue
		}
		for _, rr := range rrs {
			self.Records = append(self.Records, rr)
		}
	}

	return
}

// Returns true when if finds any endpoints
func (self *IPs) extractCnames(recur *Recur, d *Domain, c Cursor) bool {
	if _, found := self.CnameTraceBack[d.String()]; !found {
		panic("bug")
	}

	if !self.EndWith.Serves(d) {
		// domain not in the zone
		// so even there were cname records about this domain
		// they cannot be trusted
		return false
	}

	rrs := recur.Packet.SelectRecords(d, consts.CNAME)
	ret := false

	for _, rr := range rrs {
		cname := rdata.ToDomain(rr.Rdata)
		cnameStr := cname.String()
		if self.CnameTraceBack[cnameStr] != nil {
			// some error cnames, pointing to self or forming circles
			continue
		}

		c.Printf("// cname: %v -> %v", d, cname)
		self.CnameRecords = append(self.CnameRecords, rr)
		self.CnameTraceBack[cname.String()] = d

		// see if it follows another CNAME
		if self.extractCnames(recur, cname, c) {
			// see so, then we only tracks the end point
			ret = true // we added an endpoint in the recursion
			continue
		}

		c.Printf("// cname endpoint: %v", cname)
		// these are end points that needs to be crawled
		self.CnameEndpoints = append(self.CnameEndpoints, cname)
		ret = true
	}

	return ret
}

func (self *IPs) PrintResult(c Cursor) {
	for _, r := range self.CnameRecords {
		c.Printf("// %v -> %v", r.Domain, rdata.ToDomain(r.Rdata))
	}
	for _, r := range self.Records {
		c.Printf("// %v(%v)", r.Domain, rdata.ToIPv4(r.Rdata))
	}
}

func (self *IPs) Run(c Cursor) {
	if !self.HeadLess {
		c.Printf("ips %v {", self.Domain)
		c.ShiftIn()
		defer ShiftOutWith(c, "}")
	}

	self.run(c)
	if c.Error() != nil {
		return
	}

	if !self.HideResult {
		self.PrintResult(c)
	}
}

func (self *IPs) run(c Cursor) {
	recur := NewRecur(self.Domain)
	recur.HeadLess = true
	recur.StartWith = self.StartWith

	_, e := c.T(recur)
	if e != nil {
		return
	}

	// inherit from recur
	self.Return = recur.Return
	self.EndWith = recur.EndWith
	self.Packet = recur.Packet
	self.ServerRecords = recur.Records

	self.Records = make([]*pa.RR, 0, 10)
	self.findResults(recur)

	// even if we find results, we still track cnames if any
	self.CnameEndpoints = make([]*Domain, 0, 10)
	if self.CnameTraceBack == nil {
		self.CnameTraceBack = make(map[string]*Domain)
		self.CnameTraceBack[self.Domain.String()] = nil
	} else {
		_, found := self.CnameTraceBack[self.Domain.String()]
		if !found {
			panic("bug")
		}
	}

	self.CnameRecords = make([]*pa.RR, 0, 10)
	if !self.extractCnames(recur, self.Domain, c) {
		return
	}

	if len(self.CnameEndpoints) == 0 {
		panic("bug")
	}

	unresolved := self.findCnameResults(recur)
	if len(unresolved) == 0 {
		return
	}

	// trace down the cnames
	p := self.Packet
	z := self.EndWith
	self.CnameEndIPs = make([]*IPs, 0, len(unresolved))

	for _, cname := range unresolved {
		// search for redirects
		servers, rrs := ExtractServers(p, z.Zone(), cname, c)
		self.ServerRecords = appendAll(self.ServerRecords, rrs)

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
		cnameIPs.CnameTraceBack = self.CnameTraceBack
		self.CnameEndIPs = append(self.CnameEndIPs, cnameIPs)

		_, e := c.T(cnameIPs)
		if e != nil {
			return
		}

		self.Records = appendAll(self.Records, cnameIPs.Records)
		self.CnameRecords = appendAll(self.CnameRecords, cnameIPs.CnameRecords)
		self.ServerRecords = appendAll(self.ServerRecords, cnameIPs.ServerRecords)
	}
}
