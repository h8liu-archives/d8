package client

import (
	"net"

	"d8/domain"
)

type Query struct {
	Domain *domain.Domain
	Type   uint16
	Server *net.UDPAddr
}
