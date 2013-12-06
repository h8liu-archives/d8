// Package wire parses DNS messages into records and pack DNS queries
package packet

import (
	"bytes"
	"encoding/binary"
	"math/rand"

	"d8/domain"
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

func u16(buf []byte) uint16       { return enc.Uint16(buf) }
func u32(buf []byte) uint32       { return enc.Uint32(buf) }
func putU16(buf []byte, i uint16) { enc.PutUint16(buf, i) }
func putU32(buf []byte, i uint32) { enc.PutUint32(buf, i) }

func Unpack(p []byte) *Message {
	m := new(Message)
	m.Packet = p
	m.Unpack()

	return m
}

func (self *Message) Unpack() {
	// TODO
}

func (self *Message) packHeader(out *bytes.Buffer) {
	var buf [12]byte

	putU16(buf[0:2], self.Id)
	putU16(buf[2:4], self.Flag)
	putU16(buf[4:6], 1) // always have one question
	putU16(buf[6:8], self.Answer.LenU16())
	putU16(buf[8:10], self.Authority.LenU16())
	putU16(buf[10:12], self.Addition.LenU16())

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
