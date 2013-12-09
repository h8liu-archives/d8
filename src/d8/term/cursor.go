package term

import (
	"d8/domain"
	"net"
)

type Cursor interface {
	Do(t Task)
	Query(d *domain.Domain, t uint16, server *net.IP)
}
