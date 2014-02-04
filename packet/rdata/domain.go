package rdata

import (
	"github.com/h8liu/d8/domain"

	"bytes"
	"errors"
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

func unpackDomain(in *bytes.Reader, n uint16,
	p []byte) (*domain.Domain, error) {
	if n == 0 {
		return nil, errors.New("zero domain len")
	}

	was := in.Len()
	d, e := domain.Unpack(in, p)
	now := in.Len()
	if was-now != int(n) {
		return nil, fmt.Errorf("domain length expect %d, actual %d",
			n, was-now)
	}

	return d, e
}

func UnpackDomain(in *bytes.Reader, n uint16, p []byte) (*Domain, error) {
	d, e := unpackDomain(in, n, p)
	return (*Domain)(d), e
}

func ToDomain(r Rdata) *domain.Domain {
	return (*domain.Domain)(r.(*Domain))
}
