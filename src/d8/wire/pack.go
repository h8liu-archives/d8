package wire

import (
	"bytes"
)

func packTypeClass(out *bytes.Buffer, t uint16, c uint16) {
	var buf [4]byte
	putU16(buf[0:2], t)
	putU16(buf[2:4], c)

	out.Write(buf[:])
}
