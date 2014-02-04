package packet

import (
	"github.com/h8liu/d8/domain"
	. "github.com/h8liu/d8/packet/consts"
)

const (
	SecAnsw = 1 << iota
	SecAuth
	SecAddi
)

type Selector interface {
	Select(rr *RR, section int) bool
}

type AnswerSelector struct {
	Domain *domain.Domain
	Type   uint16
}

func (self *AnswerSelector) Select(rr *RR, _ int) bool {
	if !rr.Domain.Equal(self.Domain) {
		return false
	}
	return self.Type == rr.Type || (self.Type == A && rr.Type == CNAME)
}

type RecordSelector struct {
	Domain *domain.Domain
	Type   uint16
}

func (self *RecordSelector) Select(rr *RR, _ int) bool {
	return rr.Domain.Equal(self.Domain) && self.Type == rr.Type
}

type RedirectSelector struct{ Zone, Domain *domain.Domain }

func (self *RedirectSelector) Select(rr *RR, _ int) bool {
	return rr.Type == NS && rr.Domain.IsChildOf(self.Zone) &&
		rr.Domain.IsZoneOf(self.Domain)
}

type IPSelector struct{ Domain *domain.Domain }

func (self *IPSelector) Select(rr *RR, _ int) bool {
	return rr.Type == A && rr.Domain.Equal(self.Domain)
}
