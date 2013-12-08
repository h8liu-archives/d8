package client

import (
	"printer"
	"time"
)

type job struct {
	id       uint16
	exchange *Exchange
	deadline time.Time
	printer  *printer.Printer
	c        chan<- *Exchange
}

func (self *job) Close() {
	if self.printer != nil {
		self.exchange.printRecv(self.printer)
	}
	self.c <- self.exchange
}

func (self *job) CloseErr(e error) {
	self.exchange.Error = e
	self.Close()
}

func (self *job) CloseRecv(m *Message) {
	self.exchange.Recv = m
	bugOn(self.exchange.Error != nil)
	self.Close()
}
