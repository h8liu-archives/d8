package term

import (
	"net"
	"d8/domain"
)

type Cursor interface {
	Do(t Task)
	Query(d *domain.Domain, t uint16, server *net.IP)
}
