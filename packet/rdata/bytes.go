package rdata

import (
	"bytes"
	"fmt"
)

type Bytes []byte

var _ Rdata = Bytes(nil)

func (self Bytes) Pack() []byte {
	return self
}

func UnpackBytes(in *bytes.Reader, n uint16) (Bytes, error) {
	ret := make([]byte, n)
	if _, e := in.Read([]byte(ret)); e != nil {
		return nil, e
	}

	return Bytes(ret), nil
}

func (self Bytes) PrintTo(out *bytes.Buffer) {
	fmt.Fprintf(out, "[")
	for i, b := range self {
		if i > 0 && i%4 == 0 {
			fmt.Fprintf(out, " ")
		}
		fmt.Fprintf(out, "%02x", b)
	}
	fmt.Fprintf(out, "]")
}
