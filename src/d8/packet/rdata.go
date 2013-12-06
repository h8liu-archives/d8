package packet

import (
	"bytes"
)

type Rdata struct {
	data   []byte
	packet []byte
}

func (self *Rdata) pack(out *bytes.Buffer) {
	panic("todo")
}
