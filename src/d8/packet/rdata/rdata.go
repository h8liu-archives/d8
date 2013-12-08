package rdata

import (
	"bytes"
	"encoding/binary"
)

var enc = binary.BigEndian

type Rdata interface {
	PrintTo(out *bytes.Buffer)
	Pack(out *bytes.Buffer)
}

func Unpack(t, c uint16, in *bytes.Reader, p []byte) (Rdata, error) {
	// TODO:
	return UnpackBytes(in)
}
