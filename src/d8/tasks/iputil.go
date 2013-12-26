package tasks

import (
	"encoding/binary"
	"net"
)

func IP2Uint(ip net.IP) uint32 {
	bytes := []byte(ip.To4())
	if bytes == nil {
		panic("not ipv4")
	}
	return binary.BigEndian.Uint32(bytes)
}
