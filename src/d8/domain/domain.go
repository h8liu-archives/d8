package domain

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Domain struct {
	name   string
	labels []string
}

func (self *Domain) Equals(other *Domain) bool {
	if self == nil {
		return other == nil
	}
	if other == nil {
		return false
	}
	return self.name == other.name
}

func (self *Domain) String() string {
	return self.name
}

func err(s, r string) error {
	return fmt.Errorf("'%s': %s", s, r)
}

func checkLabel(s string, label string) error {
	nl := len(label)
	if nl == 0 {
		return err(s, "empty label")
	}
	if nl >= 64 {
		return err(s, "label too long")
	}
	if label[0] == '-' {
		return err(s, "label starts with dash")
	}
	if label[nl-1] == '-' {
		return err(s, "label ends with dash")
	}

	for _, c := range label {
		if 'a' <= c && c <= 'z' {
			continue
		}
		if '0' <= c && c <= '9' {
			continue
		}
		if c == '_' || c == '-' {
			continue
		}

		return err(s, "invalid char")
	}

	return nil
}

var Root *Domain

func init() {
	Root = &Domain{"", []string{}}
}

func Parse(s string) (*Domain, error) {
	orig := s

	ip := net.ParseIP(s)
	if ip != nil {
		return nil, err(orig, "IP addr")
	}

	n := len(s)

	if n > 255 {
		return nil, err(orig, "name too long")
	}

	if n > 0 && s[n-1] == '.' {
		s = s[:n-1]
	}

	if s == "" {
		return Root, nil
	}

	s = strings.ToLower(s)
	labels := strings.Split(s, ".")

	for _, label := range labels {
		e := checkLabel(orig, label)
		if e != nil {
			return nil, e
		}
	}

	return &Domain{s, labels}, nil
}

func D(s string) *Domain {
	ret, e := Parse(s)
	if e != nil {
		panic(e)
	}
	return ret
}

func (self *Domain) IsRoot() bool {
	return self.name == ""
}

func (self *Domain) IsParentOf(other *Domain) bool {
	return other.IsChildOf(self)
}

func (self *Domain) IsChildOf(other *Domain) bool {
	if self.Equals(other) {
		return false
	}
	return strings.HasSuffix(self.name, other.name)
}

func (self *Domain) Parent() *Domain {
	if self.IsRoot() {
		return nil
	}

	if len(self.labels) == 1 {
		return Root
	}

	labels := self.labels[1:]
	name := self.name[len(self.labels[0])+1:]
	return &Domain{name, labels}
}

func (self *Domain) RegParts() (registered *Domain, registrar *Domain) {
	var last *Domain
	cur := self
	parent := self.Parent()
	for {
		if parent == nil {
			return last, cur
		}
		if superRegs[parent.name] && nonRegs[cur.name] {
			return last, cur
		}
		if regs[cur.name] {
			return last, cur
		}

		last = cur
		cur = parent
		parent = parent.Parent()
	}
}

func (self *Domain) IsRegistrar() bool {
	reged, _ := self.RegParts()
	return reged == nil
}

func (self *Domain) Pack(buf *bytes.Buffer) {
	for _, lab := range self.labels {
		_lab := []byte(lab)
		buf.WriteByte(byte(len(_lab)))
		buf.Write(_lab)
	}
	buf.WriteByte(0)
}
