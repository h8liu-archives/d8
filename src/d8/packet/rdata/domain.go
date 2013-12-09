package rdata

import (
	"d8/domain"

	"bytes"
	"fmt"
)

type Domain domain.Domain

func (self *Domain) PrintTo(out *bytes.Buffer) {
	fmt.Fprint(out, (*domain.Domain)(self))
}

func (self *Domain) Pack() []byte {
	buf := new(bytes.Buffer)
	(*domain.Domain)(self).Pack(buf)
	return buf.Bytes()
}

func UnpackDomain(in *bytes.Reader, n uint16, p []byte) (*Domain, error) {
	d, e := domain.Unpack(in, p)
	return (*Domain)(d), e

}
