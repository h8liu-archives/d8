package rdata

import (
	"bytes"
	"fmt"
	"net"
)

type IPv4 net.IP

func UnpackIPv4(in *bytes.Reader, n uint16) (IPv4, error) {
	if n != 4 {
		return nil, fmt.Errorf("IPv4 with %d bytes", n)
	}

	buf := make([]byte, 4)
	_, e := in.Read(buf)
	if e != nil {
		return nil, e
	}

	return IPv4(net.IPv4(buf[0], buf[1], buf[2], buf[3])), nil
}

func (self IPv4) PrintTo(out *bytes.Buffer) {
	fmt.Fprint(out, net.IP(self))
}

func (self IPv4) Pack() []byte {
	return net.IP(self).To4()
}
