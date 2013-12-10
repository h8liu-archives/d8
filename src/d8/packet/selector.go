package packet

import (
	"d8/domain"
	. "d8/packet/consts"
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

type RedirectSelector struct {
	Zone *domain.Domain
}

func (self *RedirectSelector) Select(rr *RR, _ int) bool {
	if rr.Type != NS {
		return false
	}
	return rr.Domain.IsChildOf(self.Zone)
}

type IPSelector struct {
	Domain *domain.Domain
}

func (self *IPSelector) Select(rr *RR, _ int) bool {
	return rr.Type == A && rr.Domain.Equal(self.Domain)
}
