// Package wire parses DNS messages into records and pack DNS queries
package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand"

	"d8/domain"
	. "d8/packet/consts"
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

func randomId() uint16 { return uint16(rand.Uint32()) }

var enc = binary.BigEndian

func Unpack(p []byte) (*Packet, error) {
	m := new(Packet)
	m.Bytes = p
	e := m.Unpack()

	return m, e
}

func (self *Packet) Unpack() error {
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

	m.Id = id
	m.Flag = 0
	m.Question = &Question{d, t, IN}
	m.PackQuery()

	return m
}
