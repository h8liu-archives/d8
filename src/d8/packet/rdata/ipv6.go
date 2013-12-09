package rdata

import (
	"bytes"
	"fmt"
	"net"
)

type IPv6 net.IP

func UnpackIPv6(in *bytes.Reader, n uint16) (IPv6, error) {
	if n != 16 {
		return nil, fmt.Errorf("IPv6 with %d bytes", n)
	}
	buf := make([]byte, 16)
	_, e := in.Read(buf)
	if e != nil {
		return nil, e
	}

	return IPv6(buf), nil
}

func (self IPv6) PrintTo(out *bytes.Buffer) {
	fmt.Fprint(out, net.IP(self))
}

func (self IPv6) Pack() []byte {
	return net.IP(self).To16()
}
