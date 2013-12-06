package wire

import (
	"bytes"
	"d8/domain"
)

type Question struct {
	Domain *domain.Domain
	Type   uint16
	Class  uint16
}

func (self *Question) pack(out *bytes.Buffer) {
	self.Domain.Pack(out)
	packTypeClass(out, self.Type, self.Class)
}
