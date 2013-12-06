package wire

import (
	"d8/domain"
)

type RR struct {
	Domain *domain.Domain
	Type   uint16
	Class  uint16
	TTL    uint32
	Rdata  *Rdata
}
