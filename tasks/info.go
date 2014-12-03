package tasks

import (
	"bytes"
	"encoding/json"
	"fmt"

	. "github.com/h8liu/d8/domain"
	pa "github.com/h8liu/d8/packet"
	. "github.com/h8liu/d8/packet/consts"
	"github.com/h8liu/d8/packet/rdata"
	"github.com/h8liu/d8/printer"
	. "github.com/h8liu/d8/term"
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
		defer c.ShiftOut("}")
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

func (self *Info) PrintTo(p printer.Interface) {
	if len(self.Cnames) > 0 {
		p.Print("cnames {")
		p.ShiftIn()
		for _, r := range self.Cnames {
			p.Printf("%v -> %v", r.Domain, rdata.ToDomain(r.Rdata))
		}
		p.ShiftOut("}")
	}

	if len(self.Results) == 0 {
		p.Print("(unresolvable)")
	} else {
		p.Print("ips {")
		p.ShiftIn()

		for _, r := range self.Results {
			d := r.Domain
			ip := rdata.ToIPv4(r.Rdata)
			if d.Equal(self.Domain) {
				p.Printf("%v", ip)
			} else {
				p.Printf("%v(%v)", ip, d)
			}
		}

		p.ShiftOut("}")
	}

	if len(self.NameServers) > 0 {
		p.Print("servers {")
		p.ShiftIn()

		for _, ns := range self.NameServers {
			p.Printf("%v", ns)
		}

		p.ShiftOut("}")
	}

	if len(self.Records) > 0 {
		p.Print("records {")
		p.ShiftIn()

		for _, rr := range self.Records {
			p.Printf("%v", rr.Digest())
		}

		p.ShiftOut("}")
	}
}

func (self *Info) Out() string {
	ret := new(bytes.Buffer)
	p := printer.New(ret)
	self.PrintTo(p)
	return ret.String()
}

type infoResult struct {
	domain  string
	ips     []string
	cnames  []string
	servers []string
	records []string
}

func newInfoResult(domain string) *infoResult {
	ret := new(infoResult)
	ret.domain = domain
	ret.ips = make([]string, 0, 10)
	ret.cnames = make([]string, 0, 10)
	ret.servers = make([]string, 0, 10)
	ret.records = make([]string, 0, 10)
	return ret
}

func jmarsh(v interface{}) []byte {
	ret, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}
	return ret
}

func (self *Info) Result() string {
	ret := newInfoResult(self.Domain.String())

	for _, r := range self.Cnames {
		s := fmt.Sprintf("%v -> %v", r.Domain, rdata.ToDomain(r.Rdata))
		ret.cnames = append(ret.cnames, s)
	}

	for _, r := range self.Results {
		d := r.Domain
		ip := rdata.ToIPv4(r.Rdata)
		var s string
		if d.Equal(self.Domain) {
			s = fmt.Sprintf("%v", ip)
		} else {
			s = fmt.Sprintf("%v(%v)", ip, d)
		}

		ret.ips = append(ret.ips, s)
	}

	for _, ns := range self.NameServers {
		s := fmt.Sprintf("%v", ns)
		ret.servers = append(ret.servers, s)
	}
	for _, rr := range self.Records {
		s := fmt.Sprintf("%v", rr.Digest())
		ret.records = append(ret.records, s)
	}

	out := new(bytes.Buffer)

	wr := func(obj interface{}, post string) {
		out.Write(jmarsh(obj))
		out.WriteString(post)
	}

	wr(ret.domain, "\t")
	wr(ret.ips, "\t")
	wr(ret.cnames, "\t")
	wr(ret.servers, "\t")
	wr(ret.records, "\n")

	return out.String()
}
