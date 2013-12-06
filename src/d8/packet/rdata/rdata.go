package rdata

import (
	"bytes"
	"encoding/binary"
)

var enc = binary.BigEndian

type Rdata interface {
	Pack(in *bytes.Buffer)
}

func Unpack(t, c uint16, in *bytes.Reader, p []byte) (Rdata, error) {
	return UnpackBytes(in)
}
