package rdata

import (
	"bytes"
	"fmt"
)

type Text string

func UnpackString(in *bytes.Reader, n uint16) (Text, error) {
	buf := make([]byte, n)
	_, e := in.Read(buf)
	if e != nil {
		return "", e
	}
	return Text(string(buf)), nil
}

func (self Text) PrintTo(out *bytes.Buffer) {
	fmt.Fprintf(out, "%#v", string(self))
}

func (self Text) Pack() []byte {
	return []byte(self)
}
