package rdata

import (
	"d8/domain"

	"bytes"
	"errors"
	"fmt"
)

type SrcOfAuth struct {
	Mname                                   *domain.Domain
	Rname                                   *domain.Domain
	Serial, Refresh, Retry, Expire, Minimum uint32
}

func (self *SrcOfAuth) PrintTo(out *bytes.Buffer) {
	fmt.Fprintf(out, "%v/%v serial=%d refresh=%d retry=%d exp=%d min=%d",
		self.Mname, self.Rname,
		self.Serial, self.Refresh, self.Retry, self.Expire, self.Minimum)
}

func UnpackSrcOfAuth(in *bytes.Reader, n uint16, p []byte) (*SrcOfAuth, error) {
	if n <= 22 {
		return nil, fmt.Errorf("soa with %d bytes", n)
	}

	ret := new(SrcOfAuth)
	var e error
	was := in.Len()
	ret.Mname, e = domain.Unpack(in, p)
	if e != nil {
		return nil, e
	}

	ret.Rname, e = domain.Unpack(in, p)
	if e != nil {
		return nil, e
	}
	now := in.Len()
	if was-now+20 != int(n) {
		return nil, errors.New("invalid soa field length")
	}

	buf := make([]byte, 20)
	_, e = in.Read(buf)
	if e != nil {
		return nil, e
	}
	ret.Serial = enc.Uint32(buf[0:4])
	ret.Refresh = enc.Uint32(buf[4:8])
	ret.Retry = enc.Uint32(buf[8:12])
	ret.Expire = enc.Uint32(buf[12:16])
	ret.Minimum = enc.Uint32(buf[16:20])

	return ret, nil
}

func (self *SrcOfAuth) Pack() []byte {
	buf := new(bytes.Buffer)
	self.Mname.Pack(buf)
	self.Rname.Pack(buf)

	b := make([]byte, 20)
	enc.PutUint32(b[0:4], self.Serial)
	enc.PutUint32(b[4:8], self.Refresh)
	enc.PutUint32(b[8:12], self.Retry)
	enc.PutUint32(b[12:16], self.Expire)
	enc.PutUint32(b[16:20], self.Minimum)
	buf.Write(b)

	return buf.Bytes()

}
