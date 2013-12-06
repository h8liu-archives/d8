package packet

import (
	"bytes"
	"d8/domain"
)

type Question struct {
	Domain *domain.Domain
	Type   uint16
	Class  uint16
}

func (self *Question) packFlags(out *bytes.Buffer) {
	var buf [4]byte
	putU16(buf[0:2], self.Type)
	putU16(buf[2:4], self.Class)

	out.Write(buf[:])
}

func (self *Question) pack(out *bytes.Buffer) {
	self.Domain.Pack(out)
	self.packFlags(out)
}
