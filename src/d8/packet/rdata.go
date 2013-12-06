package packet

import (
	"bytes"
)

type Rdata []byte

func (self Rdata) pack(out *bytes.Buffer) {
	n := len(self)
	if n > 255 {
		panic("too long data")
	}

	var buf [2]byte
	enc.PutUint16(buf[:], uint16(n))

	out.Write(buf[:])
	out.Write(self)
}
