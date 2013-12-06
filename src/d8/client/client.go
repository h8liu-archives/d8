package client

import (
	"net"
	"time"

	"d8/packet"
)

type Client struct {
	conn *net.UDPConn
}

const DnsPort = 53

func New() *Client {
	return NewPort(DnsPort)
}

func NewPort(port uint16) *Client {
	panic("todo")
}

type Message struct {
	RemoteAddr *net.UDPAddr
	Message    *packet.Message
	Timestamp  time.Time
}

type Exchange struct {
	Question *Message
	Reply    *Message
	Error    error
}

func (self *Client) Send(m *Message, c <-chan *Exchange) {
	// will send the message out via network
	// if any send error, then put an error exchange into c

	panic("todo")
}

func (self *Client) AsyncQuery(m *Message, f func(*Exchange)) {
	c := make(chan *Exchange)
	go func() {
		f(<-c)
	}()

	self.Send(m, c)
}

func (self *Client) Query(m *Message) (*Message, error) {
	c := make(chan *Exchange, 1) // we need a slot in case of send error
	self.Send(m, c)

	exchange := <-c
	return exchange.Reply, exchange.Error
}
