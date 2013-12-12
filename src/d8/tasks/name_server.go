package tasks

import (
	"net"

	. "d8/domain"
)

type NameServer struct {
	Zone   *Domain
	Domain *Domain
	IP     net.IP
	Glued  bool
}
