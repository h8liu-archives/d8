// Package wire parses DNS messages into records and pack DNS queries
package wire

type Msg struct {
	b []byte
	
}

func Unpack(buf []byte) *Msg {
	panic("todo")
}

