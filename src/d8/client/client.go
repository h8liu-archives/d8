package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"d8/packet"
)

type Client struct {
	conn   *net.UDPConn
	idPool *idPool

	jobs       map[uint16]*job
	newJobs    chan *job
	sendErrors chan *job
	recvs      chan *Message
	timer      <-chan time.Time
}

const (
	DnsPort    = 53
	ClientPort = 3553
)

func NewPort(port uint16) (*Client, error) {
	ret := new(Client)
	addr := &net.UDPAddr{Port: ClientPort}
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
	return NewPort(ClientPort)
}

const packetMaxSize = 512

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
			log.Print("unpack:", e)
			continue
		}

		m := &Message{
			RemoteAddr: addr,
			Packet:     p,
			Timestamp:  time.Now(),
		}
		log.Println("recv:", m)
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

var ErrTimeout = errors.New("timeout")

func (self *Client) serve() {
	for {
		select {
		case job := <-self.newJobs:
			log.Println("new:", job.id)
			id := job.id
			bugOn(self.jobs[id] != nil)
			self.jobs[id] = job
			log.Println("new done")
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
			log.Println("recved:", m)
			id := m.Packet.Id
			job := self.jobs[id]
			if job == nil {
				// this might happen with the timeout window is set too small
				log.Printf("recved zombie msg with id %d", id)
			} else {
				job.CloseRecv(m)
			}
		case now := <-self.timer:
			timeouts := make([]uint16, 0, 1024)

			for id, job := range self.jobs {
				bugOn(job.id != id)
				if job.deadline.Before(now) {
					job.CloseErr(ErrTimeout)

					// iterating the map, so delete afterwards for safty
					timeouts = append(timeouts, id)
				}
			}

			if len(timeouts) > 0 {
				log.Println("timeouts:", now, timeouts)
			}

			for _, id := range timeouts {
				self.delJob(id)
			}
		}
	}
}

const timeout = time.Second * 3

func (self *Client) Send(q *Query, c chan<- *Exchange) {
	log.Println("query:", q)

	id := self.idPool.Fetch()
	message := newMessage(q, id)
	exchange := &Exchange{
		Query: q,
		Send:  message,
	}
	job := &job{
		id:       id,
		exchange: exchange,
		deadline: time.Now().Add(timeout),
		c:        c,
	}

	self.newJobs <- job // set a place in mapping

	e := self.send(message)
	if e != nil {
		job.CloseErr(e)

		// release the spot reserved if not timed out
		self.sendErrors <- job
	}
}

func (self *Client) send(m *Message) error {
	if m.RemoteAddr.Port == 0 {
		m.RemoteAddr.Port = DnsPort
	}

	log.Println("send:", m.RemoteAddr)
	fmt.Print(hex.Dump(m.Packet.Bytes))
	_, e := self.conn.WriteToUDP(m.Packet.Bytes, m.RemoteAddr)
	return e
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
