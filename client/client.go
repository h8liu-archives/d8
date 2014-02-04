package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/h8liu/d8/packet"
)

type Client struct {
	conn	*net.UDPConn
	idPool	*idPool

	jobs		map[uint16]*job
	newJobs		chan *job
	sendErrors	chan *job
	recvs		chan *Message
	timer		<-chan time.Time
}

const (
	DNSPort = 53
)

func NewPort(port uint16) (*Client, error) {
	ret := new(Client)

	addr := &net.UDPAddr{Port: int(port)}
	if port == 0 {
		addr = nil
	}

	var e error
	ret.conn, e = net.ListenUDP("udp4", addr)
	if e != nil {
		return nil, e
	}

	ret.newJobs = make(chan *job, 0)
	ret.sendErrors = make(chan *job, 10)
	ret.recvs = make(chan *Message, 10)
	ret.timer = time.Tick(time.Millisecond * 100)
	ret.idPool = newIdPool()
	ret.jobs = make(map[uint16]*job)

	go ret.recv()
	go ret.serve()

	return ret, nil
}

func New() (*Client, error) {
	return NewPort(0)
}

const packetMaxSize = 1600

func newRecvBuf() []byte {
	return make([]byte, packetMaxSize)
}

func (self *Client) recv() {
	buf := newRecvBuf()

	for {
		n, addr, e := self.conn.ReadFromUDP(buf)
		if e != nil {
			log.Print("recv:", e)
			continue
		}

		p, e := packet.Unpack(buf[:n])
		if e != nil {
			log.Print("unpack: ", e)
			fmt.Println(hex.Dump(buf[:n]))

			continue
		}

		m := &Message{
			RemoteAddr:	addr,
			Packet:		p,
			Timestamp:	time.Now(),
		}
		self.recvs <- m

		buf = newRecvBuf()
	}
}

func (self *Client) delJob(id uint16) {
	if self.jobs[id] != nil {
		delete(self.jobs, id)
		self.idPool.Return(id)
	}
}

var errTimeout = errors.New("timeout")

func (self *Client) serve() {
	for {
		select {
		case job := <-self.newJobs:
			id := job.id
			bugOn(self.jobs[id] != nil)
			self.jobs[id] = job
		case job := <-self.sendErrors:
			/*
				Need to check if it is still the same job. In some rare racing
				cases, sendErrors will be delayed (like by a send that takes
				too long), and timeout might trigger first, hence reallocate
				the id to another job.
			*/
			if self.jobs[job.id] == job {
				// still the same job
				self.delJob(job.id)
			}
		case m := <-self.recvs:
			id := m.Packet.Id
			job := self.jobs[id]
			if job == nil {
				// this might happen with the timeout window is set too small
				log.Printf("recved zombie msg with id %d", id)
			} else {
				bugOn(job.id != id)
				job.CloseRecv(m)
				self.delJob(id)
			}
		case now := <-self.timer:
			timeouts := make([]uint16, 0, 1024)

			for id, job := range self.jobs {
				bugOn(job.id != id)
				if job.deadline.Before(now) {
					job.CloseErr(errTimeout)

					// iterating the map, so delete afterwards for safty
					timeouts = append(timeouts, id)
				}
			}

			for _, id := range timeouts {
				self.delJob(id)
			}
		}
	}
}

const timeout = time.Second * 3

func (self *Client) Send(q *QueryPrinter, c chan<- *Exchange) {
	id := self.idPool.Fetch()
	message := newMessage(q.Query, id)
	if message.RemoteAddr.Port == 0 {
		message.RemoteAddr.Port = DNSPort
	}

	exchange := &Exchange{
		Query:		q.Query,
		Send:		message,
		PrintFlag:	q.PrintFlag,
	}
	job := &job{
		id:		id,
		exchange:	exchange,
		deadline:	time.Now().Add(timeout),
		printer:	q.Printer,
		c:		c,
	}

	self.newJobs <- job	// set a place in mapping

	if q.Printer != nil {
		exchange.printSend(q.Printer)
	}

	e := self.send(message)
	if e != nil {
		job.CloseErr(e)

		// release the spot reserved if not timed out
		self.sendErrors <- job
	}
}

func (self *Client) send(m *Message) error {
	_, e := self.conn.WriteToUDP(m.Packet.Bytes, m.RemoteAddr)
	return e
}

func (self *Client) AsyncQuery(q *QueryPrinter, f func(*Exchange)) {
	c := make(chan *Exchange)
	go func() {
		f(<-c)
	}()

	self.Send(q, c)
}

func (self *Client) Query(q *QueryPrinter) *Exchange {
	c := make(chan *Exchange, 1)	// we need a slot in case of send error
	self.Send(q, c)

	return <-c
}
