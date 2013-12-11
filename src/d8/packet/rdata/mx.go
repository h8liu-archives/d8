package rdata

import (
	"d8/domain"

	"bytes"
	"fmt"
)

type MailEx struct {
	Priority uint16
	Domain   *domain.Domain
}

func (self *MailEx) PrintTo(out *bytes.Buffer) {
	fmt.Fprintf(out, "%v/%d", self.Domain, self.Priority)
}

func UnpackMailEx(in *bytes.Reader, n uint16, p []byte) (*MailEx, error) {
	if n <= 2 {
		return nil, fmt.Errorf("mx with %d bytes", n)
	}

	buf := make([]byte, 2)
	_, e := in.Read(buf)
	if e != nil {
		return nil, e
	}

	ret := new(MailEx)
	ret.Priority = enc.Uint16(buf)
	ret.Domain, e = unpackDomain(in, n-2, p)
	if e != nil {
		return nil, e
	}

	return ret, nil
}

func (self *MailEx) Pack() []byte {
	buf := new(bytes.Buffer)
	b := make([]byte, 2)
	enc.PutUint16(b, self.Priority)
	buf.Write(b)
	self.Domain.Pack(buf)
	return buf.Bytes()
}
