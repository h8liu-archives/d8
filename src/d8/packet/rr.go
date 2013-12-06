package packet

import (
	"bytes"

	"d8/domain"
)

type RR struct {
	Domain *domain.Domain
	Type   uint16
	Class  uint16
	TTL    uint32
	Rdata  *Rdata
}

func (self *RR) packFlags(out *bytes.Buffer) {
	var buf [8]byte
	putU16(buf[0:2], self.Type)
	putU16(buf[2:4], self.Class)
	putU32(buf[4:8], self.TTL)
	out.Write(buf[:])
}

func (self *RR) pack(out *bytes.Buffer) {
	self.Domain.Pack(out)
	self.packFlags(out)
	self.Rdata.pack(out)
}
