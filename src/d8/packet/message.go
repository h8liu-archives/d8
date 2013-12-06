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

type Message struct {
	Packet []byte

	Id        uint16
	Flag      uint16
	Question  *Question
	Answer    Section
	Authority Section
	Addition  Section
}

func randomId() uint16 { return uint16(rand.Uint32()) }

var enc = binary.BigEndian

func Unpack(p []byte) (*Message, error) {
	m := new(Message)
	m.Packet = p
	e := m.Unpack()

	return m, e
}

func (self *Message) Unpack() error {
	if self.Packet == nil {
		return errors.New("nil packet")
	}

	in := bytes.NewReader(self.Packet)

	if e := self.unpackHeader(in); e != nil {
		return e
	}

	if e := self.Question.unpack(in, self.Packet); e != nil {
		return e
	}

	if e := self.Answer.unpack(in, self.Packet); e != nil {
		return e
	}

	if e := self.Authority.unpack(in, self.Packet); e != nil {
		return e
	}

	if e := self.Addition.unpack(in, self.Packet); e != nil {
		return e
	}

	return nil
}

func (self *Message) unpackHeader(in *bytes.Reader) error {
	var buf [12]byte
	if _, e := in.Read(buf[:]); e != nil {
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

func (self *Message) packHeader(out *bytes.Buffer) {
	var buf [12]byte

	enc.PutUint16(buf[0:2], self.Id)
	enc.PutUint16(buf[2:4], self.Flag)
	enc.PutUint16(buf[4:6], 1) // always have one question
	enc.PutUint16(buf[6:8], self.Answer.LenU16())
	enc.PutUint16(buf[8:10], self.Authority.LenU16())
	enc.PutUint16(buf[10:12], self.Addition.LenU16())

	out.Write(buf[:])
}

func (self *Message) PackQuery() []byte {
	out := new(bytes.Buffer)

	self.packHeader(out)
	self.Question.pack(out)

	self.Packet = out.Bytes() // swap in
	return self.Packet
}

func Q(d *domain.Domain, t uint16) *Message {
	m := new(Message)

	m.Id = randomId()
	m.Flag = 0
	m.Question = &Question{d, t, IN}
	m.PackQuery()

	return m
}
