package rdata

import (
	"bytes"
	. "github.com/h8liu/d8/packet/consts"
	"encoding/binary"
)

var enc = binary.BigEndian

type Rdata interface {
	PrintTo(out *bytes.Buffer)
	Pack() []byte
}

func unpack(t, c uint16, in *bytes.Reader, p []byte) (Rdata, error) {
	n := uint16(in.Len())
	if c == IN {
		switch t {
		case A:
			return UnpackIPv4(in, n)
		case NS, CNAME:
			return UnpackDomain(in, n, p)
		case AAAA:
			return UnpackIPv6(in, n)
		case TXT:
			return UnpackString(in, n)
		case MX:
			return UnpackMailEx(in, n, p)
		case SOA:
			return UnpackSrcOfAuth(in, n, p)
		}
	}
	return UnpackBytes(in, n)
}

func Unpack(t, c uint16, in *bytes.Reader, p []byte) (Rdata, error) {
	buf := make([]byte, 2)
	if _, e := in.Read(buf); e != nil {
		return nil, e
	}
	n := enc.Uint16(buf)

	buf = make([]byte, n)
	if _, e := in.Read(buf); e != nil {
		return nil, e
	}

	in = bytes.NewReader(buf)
	ret, e := unpack(t, c, bytes.NewReader(buf), p)
	return ret, e
}

func Pack(out *bytes.Buffer, rdata Rdata) {
	pack := rdata.Pack()
	n := len(pack)
	if n > 255 {
		panic("rdata too long")
	}

	buf := make([]byte, 2)
	enc.PutUint16(buf, uint16(n))
	out.Write(buf)
	out.Write(pack)
}
