package client

import (
	"net"
	"time"

	"d8/domain"
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

type Query struct {
	Domain *domain.Domain
	Type   uint16
}

type Message struct {
	RemoteAddr *net.UDPAddr
	Message    *packet.Message
	Timestamp  time.Time
}

type Exchange struct {
	Query *Query
	Send  *Message
	Recv  *Message
	Error error
}

func (self *Client) Send(q *Query, c <-chan *Exchange) {
	// will send the message out via network
	// if any send error, then put an error exchange into c

	panic("todo")
}

func (self *Client) AsyncQuery(q *Query, f func(*Exchange)) {
	c := make(chan *Exchange)
	go func() {
		f(<-c)
	}()

	self.Send(q, c)
}

func (self *Client) Query(q *Query) *Exchange {
	c := make(chan *Exchange, 1) // we need a slot in case of send error
	self.Send(q, c)

	return <-c
}
