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
	enc.PutUint16(buf[0:2], self.Type)
	enc.PutUint16(buf[2:4], self.Class)

	out.Write(buf[:])
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
	var buf [4]byte
	if _, e := in.Read(buf[:]); e != nil {
		return e
	}

	self.Type = enc.Uint16(buf[0:2])
	self.Class = enc.Uint16(buf[2:4])

	return nil
}
