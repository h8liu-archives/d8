package client

import (
	"fmt"
	"net"

	"d8/domain"
	. "d8/packet/consts"
	"printer"
)

const (
	PrintAll = iota
	PrintReply
)

type Query struct {
	Domain    *domain.Domain
	Type      uint16
	Server    *net.UDPAddr
	Printer   *printer.Printer
	PrintFlag int
}

func (self *Query) String() string {
	return fmt.Sprintf("%v %s @%s",
		self.Domain,
		TypeString(self.Type),
		addrString(self.Server),
	)
}
