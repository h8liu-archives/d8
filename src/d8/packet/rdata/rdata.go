package rdata

import (
	"bytes"
	. "d8/packet/consts"
	"encoding/binary"
)

var enc = binary.BigEndian

type Rdata interface {
	PrintTo(out *bytes.Buffer)
	Pack() []byte
}

func Unpack(t, c uint16, in *bytes.Reader, p []byte) (Rdata, error) {
	buf := make([]byte, 2)
	if _, e := in.Read(buf); e != nil {
		return nil, e
	}
	n := enc.Uint16(buf)

	if c == IN {
		switch t {
		case A:
			return UnpackIPv4(in, n)
		}
	}
	return UnpackBytes(in, n)
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
