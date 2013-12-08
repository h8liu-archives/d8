package packet

import (
	"bytes"
	"fmt"

	"d8/domain"
	. "d8/packet/consts"
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
	rdata.Pack(out, self.Rdata)
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

func (self *RR) String() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s %s ", self.Domain.String(), TypeString(self.Type))
	if self.Class != IN {
		fmt.Fprintf(buf, "%s ", ClassString(self.Class))
	}
	self.Rdata.PrintTo(buf)
	fmt.Fprintf(buf, " %s", ttlString(self.TTL))

	return buf.String()
}

func ttlString(t uint32) string {
	if t == 0 {
		return "0"
	}

	buf := new(bytes.Buffer)
	second := t % 60
	minute := t / 60 % 60
	hour := t / 3600 % 24
	day := t / 3600 / 24
	if day > 0 {
		fmt.Fprintf(buf, "%dd", day)
	}
	if hour > 0 {
		fmt.Fprintf(buf, "%dh", hour)
	}
	if minute > 0 {
		fmt.Fprintf(buf, "%dm", minute)
	}
	if second > 0 {
		fmt.Fprintf(buf, "%ds", second)
	}

	return buf.String()
}
