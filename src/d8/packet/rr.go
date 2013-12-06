package packet

import (
	"bytes"

	"d8/domain"
	"d8/packet/rdata"
)

type RR struct {
	Domain *domain.Domain
	Type   uint16
	Class  uint16
	TTL    uint32
	Rdata  rdata.Rdata
}

func (self *RR) packFlags(out *bytes.Buffer) {
	var buf [8]byte
	enc.PutUint16(buf[0:2], self.Type)
	enc.PutUint16(buf[2:4], self.Class)
	enc.PutUint32(buf[4:8], self.TTL)
	out.Write(buf[:])
}

func (self *RR) pack(out *bytes.Buffer) {
	self.Domain.Pack(out)
	self.packFlags(out)
	self.Rdata.Pack(out)
}

func (self *RR) unpackFlags(in *bytes.Reader) error {
	var buf [8]byte
	if _, e := in.Read(buf[:]); e != nil {
		return e
	}
	self.Type = enc.Uint16(buf[0:2])
	self.Class = enc.Uint16(buf[2:4])
	self.TTL = enc.Uint32(buf[4:8])

	return nil
}

func (self *RR) unpackRdata(in *bytes.Reader, p []byte) error {
	var e error
	self.Rdata, e = rdata.Unpack(self.Type, self.Class, in, p)
	return e
}

func (self *RR) unpack(in *bytes.Reader, p []byte) error {
	var e error

	self.Domain, e = domain.Unpack(in, p)
	if e != nil {
		return e
	}

	if e = self.unpackFlags(in); e != nil {
		return e
	}

	return self.unpackRdata(in, p)
}

func unpackRR(in *bytes.Reader, p []byte) (*RR, error) {
	ret := new(RR)
	return ret, ret.unpack(in, p)
}
