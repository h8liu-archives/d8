package tasks

import (
	"net"

	. "d8/domain"
	pa "d8/packet"
	"d8/packet/consts"
	"d8/packet/rdata"
	. "d8/term"
	"printer"
)

type ipsResult struct {
	cnames  []*pa.RR
	results []*pa.RR
}

type IPs struct {
	Domain     *Domain
	StartWith  *ZoneServers
	HeadLess   bool
	HideResult bool

	// inherit from the initializing Recur Task
	Return  int
	Packet  *pa.Packet
	EndWith *ZoneServers
	Zones   []*ZoneServers

	CnameTraceBack map[string]*Domain // in and out, inherit from father IPs

	CnameEndpoints []*Domain       // new endpoint cnames discovered
	CnameIPs       map[string]*IPs // sub IPs for each unresolved end point

	CnameRecords []*pa.RR // new cname records
	Records      []*pa.RR // new end point ip records

	resultSave *ipsResult
}

func NewIPs(d *Domain) *IPs {
	return &IPs{Domain: d}
}

// Look for Query error or A records in Answer
func (self *IPs) collectResults(recur *Recur) {
	if recur.Return != Okay {
		panic("bug")
	}

	for _, rr := range recur.Answers {
		switch rr.Type {
		case consts.A:
			self.Records = append(self.Records, rr)
		case consts.CNAME:
			// okay
		default:
			panic("bug")
		}
	}
}

func (self *IPs) findCnameResults(recur *Recur) (unresolved []*Domain) {
	unresolved = make([]*Domain, 0, len(self.CnameEndpoints))

	for _, cname := range self.CnameEndpoints {
		rrs := recur.Packet.SelectRecords(cname, consts.A)
		if len(rrs) == 0 {
			unresolved = append(unresolved, cname)
			continue
		}
		self.Records = append(self.Records, rrs...)
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
	cnames, results := self.Results()

	for _, r := range cnames {
		c.Printf("// %v -> %v", r.Domain, rdata.ToDomain(r.Rdata))
	}

	if len(results) == 0 {
		c.Printf("// (%v is unresolvable)", self.Domain)
	}

	for _, r := range results {
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

func (self *IPs) Results() (cnames, results []*pa.RR) {
	if self.resultSave != nil {
		return self.resultSave.cnames, self.resultSave.results
	}

	cnames = make([]*pa.RR, 0, 20)
	results = make([]*pa.RR, 0, 20)
	cnames, results = self.results(cnames, results)
	self.resultSave = &ipsResult{cnames, results}

	return
}

func (self *IPs) ResultAndIPs() (cnames, results []*pa.RR, ips []net.IP) {
	cnames, results = self.Results()
	if len(results) == 0 {
		return
	}

	hits := make(map[uint32]bool)
	ips = make([]net.IP, 0, len(results))

	for _, rr := range results {
		ip := rdata.ToIPv4(rr.Rdata)
		index := IP2Uint(ip)
		if hits[index] {
			continue
		}
		hits[index] = true
		ips = append(ips, ip)
	}

	return
}

func (self *IPs) IPs() []net.IP {
	_, _, ret := self.ResultAndIPs()
	return ret
}

func (self *IPs) results(cnames, results []*pa.RR) (c, r []*pa.RR) {
	cnames = append(cnames, self.CnameRecords...)
	results = append(results, self.Records...)

	for _, ips := range self.CnameIPs {
		cnames, results = ips.results(cnames, results)
	}

	return cnames, results
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
	self.Zones = recur.Zones

	if self.Return != Okay {
		return
	}

	self.Records = make([]*pa.RR, 0, 10)
	self.collectResults(recur)

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
	self.CnameIPs = make(map[string]*IPs)

	for _, cname := range unresolved {
		// search for redirects
		servers := Servers(p, z.Zone(), cname, c)

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

		self.CnameIPs[cname.String()] = cnameIPs

		_, e := c.T(cnameIPs)
		if e != nil {
			return
		}
	}
}

func (self *IPs) PrintTo(p printer.Interface) {
	cnames, results := self.Results()

	if len(cnames) > 0 {
		for _, r := range cnames {
			p.Printf("%v -> %v", r.Domain, rdata.ToDomain(r.Rdata))
		}
		p.Println()
	}

	if len(results) == 0 {
		p.Print("(unresolvable)")
	} else {

		for _, r := range results {
			d := r.Domain
			ip := rdata.ToIPv4(r.Rdata)
			if d.Equal(self.Domain) {
				p.Printf("%v", ip)
			} else {
				p.Printf("%v(%v)", ip, d)
			}
		}
	}
}
