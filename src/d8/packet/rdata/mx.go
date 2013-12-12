package rdata

import (
	"d8/domain"

	"bytes"
	"errors"
	"fmt"
	"strings"
)

type MailEx struct {
	Priority uint16
	Domain   []string
}

func (self *MailEx) PrintTo(out *bytes.Buffer) {
	fmt.Fprintf(out, "%s/%d",
		strings.Join(self.Domain, "."), self.Priority)
}

func unpackLabels(in *bytes.Reader, n uint16, p []byte) ([]string, error) {
	if n == 0 {
		return nil, errors.New("zero labels len")
	}

	was := in.Len()
	d, e := domain.UnpackLabels(in, p)
	now := in.Len()
	if was-now != int(n) {
		return nil, fmt.Errorf("domain length expect %d, actual %d",
			n, was-now)
	}

	return d, e
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
	labels, e := unpackLabels(in, n-2, p)
	if e != nil {
		return nil, e
	}
	ret.Domain = labels

	return ret, nil
}

func (self *MailEx) Pack() []byte {
	buf := new(bytes.Buffer)
	b := make([]byte, 2)
	enc.PutUint16(b, self.Priority)
	buf.Write(b)
	domain.PackLabels(buf, self.Domain)
	return buf.Bytes()
}
