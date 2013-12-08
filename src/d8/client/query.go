package client

import (
	"fmt"
	"net"

	"d8/domain"
	. "d8/packet/consts"
	"printer"
)

type Query struct {
	Domain *domain.Domain
	Type   uint16
	Server *net.UDPAddr
	Printer *printer.Printer
}

func (self *Query) String() string {
	return fmt.Sprintf("%v %s @%s",
		self.Domain,
		TypeString(self.Type),
		addrString(self.Server),
	)
}
