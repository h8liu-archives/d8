package tasks

import (
	"fmt"
	"net"

	. "d8/domain"
)

type NameServer struct {
	Zone   *Domain
	Domain *Domain
	IP     net.IP
}

func (self *NameServer) String() string {
	if self.IP == nil {
		return fmt.Sprintf("%v ns %v",
			self.Zone, self.Domain,
		)
	}

	return fmt.Sprintf("%v ns %v(%v)",
		self.Zone, self.Domain, self.IP,
	)
}

func (self *NameServer) Key() string {
	if self.IP == nil {
		panic("unresolved")
	}

	return fmt.Sprintf("%v@%v", self.IP, self.Zone)
}
