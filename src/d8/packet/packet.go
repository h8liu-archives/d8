// Package wire parses DNS messages into records and pack DNS queries
package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"

	"d8/domain"
	. "d8/packet/consts"
	"printer"
)

type Packet struct {
	Bytes []byte

	Id        uint16
	Flag      uint16
	Question  *Question
	Answer    Section
	Authority Section
	Addition  Section
}

func (self *Packet) Rcode() uint16 {
	return Rcode(self.Flag)
}

func randomId() uint16 { return uint16(rand.Uint32()) }

var enc = binary.BigEndian

func Unpack(p []byte) (*Packet, error) {
	m := new(Packet)
	m.Bytes = p
	m.Question = new(Question)

	e := m.unpack()

	return m, e
}

func (self *Packet) unpack() error {
	if self.Bytes == nil {
		return errors.New("nil packet")
	}

	in := bytes.NewReader(self.Bytes)

	if e := self.unpackHeader(in); e != nil {
		return e
	}

	if e := self.Question.unpack(in, self.Bytes); e != nil {
		return e
	}

	if e := self.Answer.unpack(in, self.Bytes); e != nil {
		return e
	}

	if self.Flag&FlagTC != 0 {
		self.Authority = self.Authority[0:0]
		self.Addition = self.Addition[0:0]
		return nil
	}

	if e := self.Authority.unpack(in, self.Bytes); e != nil {
		return e
	}

	if e := self.Addition.unpack(in, self.Bytes); e != nil {
		return e
	}

	return nil
}

func (self *Packet) unpackHeader(in *bytes.Reader) error {
	buf := make([]byte, 12)
	if _, e := in.Read(buf); e != nil {
		return e
	}

	self.Id = enc.Uint16(buf[0:2])
	self.Flag = enc.Uint16(buf[2:4])
	if enc.Uint16(buf[4:6]) != 1 {
		return errors.New("not one question")
	}

	self.Answer = make([]*RR, enc.Uint16(buf[6:8]))
	self.Authority = make([]*RR, enc.Uint16(buf[8:10]))
	self.Addition = make([]*RR, enc.Uint16(buf[10:12]))

	return nil
}

func (self *Packet) packHeader(out *bytes.Buffer) {
	buf := make([]byte, 12)

	enc.PutUint16(buf[0:2], self.Id)
	enc.PutUint16(buf[2:4], self.Flag)
	enc.PutUint16(buf[4:6], 1) // always have one question
	enc.PutUint16(buf[6:8], self.Answer.LenU16())
	enc.PutUint16(buf[8:10], self.Authority.LenU16())
	enc.PutUint16(buf[10:12], self.Addition.LenU16())

	out.Write(buf)
}

func (self *Packet) PackQuery() []byte {
	out := new(bytes.Buffer)

	self.packHeader(out)
	self.Question.pack(out)

	self.Bytes = out.Bytes() // swap in
	return self.Bytes
}

func Q(d *domain.Domain, t uint16) *Packet {
	return Qid(d, t, randomId())
}

func Qid(d *domain.Domain, t, id uint16) *Packet {
	m := new(Packet)

	if t == 0 {
		t = A
	}

	m.Id = id
	m.Flag = 0
	m.Question = &Question{d, t, IN}
	m.PackQuery()

	return m
}

func (self *Packet) PrintTo(p *printer.Printer) {
	p.Printf("#%d %s", self.Id, flagString(self.Flag))
	p.Printf("ques %v", self.Question)
	self.Answer.PrintNameTo(p, "answ")
	self.Authority.PrintNameTo(p, "auth")
	self.Addition.PrintNameTo(p, "addi")
}

func (self *Packet) String() string {
	return printer.String(self)
}

func (self *Packet) SelectWith(s Selector) []*RR {
	ret := make([]*RR, 0, 10)
	ret = self.Answer.selectAndAppend(s, SecAnsw, ret)
	ret = self.Authority.selectAndAppend(s, SecAuth, ret)
	ret = self.Addition.selectAndAppend(s, SecAddi, ret)

	return ret
}

func (self *Packet) SelectIPs(d *domain.Domain) []*RR {
	return self.SelectWith(&IPSelector{d})
}

func (self *Packet) SelectRedirects(z, d *domain.Domain) []*RR {
	return self.SelectWith(&RedirectSelector{z, d})
}

func (self *Packet) SelectAnswers(d *domain.Domain, t uint16) []*RR {
	return self.SelectWith(&AnswerSelector{d, t})
}

func (self *Packet) SelectRecords(d *domain.Domain, t uint16) []*RR {
	return self.SelectWith(&RecordSelector{d, t})
}
