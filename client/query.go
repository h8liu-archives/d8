package client

import (
	"fmt"
	"net"

	"github.com/h8liu/d8/domain"
	. "github.com/h8liu/d8/packet/consts"
	"github.com/h8liu/d8/printer"
)

const (
	PrintAll = iota
	PrintReply
)

type QueryPrinter struct {
	*Query

	Printer   *printer.Printer
	PrintFlag int
}

type Query struct {
	Domain *domain.Domain
	Type   uint16
	Server *net.UDPAddr

	Zone       *domain.Domain
	ServerName *domain.Domain
}

func Server(ip net.IP) *net.UDPAddr {
	return &net.UDPAddr{
		IP:   ip,
		Port: DNSPort,
	}
}

func Q(d *domain.Domain, t uint16, at net.IP) *Query {
	return &Query{
		Domain: d,
		Type:   t,
		Server: Server(at),
	}
}

func Qs(d string, t uint16, at string) *Query {
	return Q(domain.D(d), t, net.ParseIP(at))
}

func (self *Query) addrString() string {
	if self.ServerName == nil {
		return addrString(self.Server)
	}

	p := self.Server.Port
	if p == 0 || p == DNSPort {
		return fmt.Sprintf("%v(%v)", self.ServerName, self.Server.IP)
	}
	return fmt.Sprintf("%v(%v):%d", self.ServerName, self.Server.IP, p)
}

func (self *Query) String() string {
	return fmt.Sprintf("%v %s @%s",
		self.Domain,
		TypeString(self.Type),
		self.addrString(),
	)
}
