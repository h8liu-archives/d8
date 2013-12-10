package tasks

import (
	"net"

	. "d8/domain"
)

type server struct {
	Domain *Domain
	IPs    []net.IP
}

func newServer(d *Domain) *server {
	ret := new(server)
	ret.Domain = d
	return ret
}

func (self *server) Add(ip net.IP) {
	self.setResolved()
	self.IPs = append(self.IPs, ip)
}

func (self *server) setResolved() {
	if self.IPs == nil {
		self.IPs = make([]net.IP, 0, 10)
	}
}

func (self *server) Resolved() bool {
	return self.IPs == nil
}
