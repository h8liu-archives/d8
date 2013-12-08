package client

import (
	"net"
	"time"

	"d8/packet"
)

type Message struct {
	RemoteAddr *net.UDPAddr
	Packet     *packet.Packet
	Timestamp  time.Time
}

func newMessage(q *Query, id uint16) *Message {
	return &Message{
		RemoteAddr: q.Server,
		Packet:     packet.Qid(q.Domain, q.Type, id),
		Timestamp:  time.Now(),
	}
}
