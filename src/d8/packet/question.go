package packet

import (
	"bytes"
	"fmt"

	"d8/domain"
	. "d8/packet/consts"
)

type Question struct {
	Domain *domain.Domain
	Type   uint16
	Class  uint16
}

func (self *Question) packFlags(out *bytes.Buffer) {
	buf := make([]byte, 4)
	enc.PutUint16(buf[0:2], self.Type)
	enc.PutUint16(buf[2:4], self.Class)
	out.Write(buf)
}

func (self *Question) pack(out *bytes.Buffer) {
	self.Domain.Pack(out)
	self.packFlags(out)
}

func (self *Question) unpack(in *bytes.Reader, p []byte) error {
	d, e := domain.Unpack(in, p)
	if e != nil {
		return e
	}
	self.Domain = d

	return self.unpackFlags(in)
}

func (self *Question) unpackFlags(in *bytes.Reader) error {
	buf := make([]byte, 4)
	if _, e := in.Read(buf); e != nil {
		return e
	}

	self.Type = enc.Uint16(buf[0:2])
	self.Class = enc.Uint16(buf[2:4])

	return nil
}

func (self *Question) String() string {
	ret := fmt.Sprintf("%s %s", self.Domain.String(), TypeString(self.Type))
	if self.Class != IN {
		ret += fmt.Sprintf(" %s", ClassString(self.Class))
	}
	return ret
}
