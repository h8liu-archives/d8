package rdata

import (
	"bytes"
)

type Bytes []byte

var _ Rdata = Bytes(nil)

func (self Bytes) Pack(out *bytes.Buffer) {
	n := len(self)
	if n > 255 {
		panic("rdata too long")
	}

	buf := make([]byte, 2)
	enc.PutUint16(buf, uint16(n))
	out.Write(buf)
	out.Write(self)
}

func UnpackBytes(in *bytes.Reader) (Bytes, error) {
	buf := make([]byte, 2)
	if _, e := in.Read(buf); e != nil {
		return nil, e
	}
	n := enc.Uint16(buf)

	ret := make([]byte, n)
	if _, e := in.Read([]byte(ret)); e != nil {
		return nil, e
	}

	return Bytes(ret), nil
}
