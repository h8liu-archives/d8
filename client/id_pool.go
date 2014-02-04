package client

import (
	"math/rand"
	"time"
)

const idCount = 65536
const nprepare = idCount / 4

type idPool struct {
	using    []bool
	nusing   int
	returns  chan uint16
	prepared chan uint16
	rand     *rand.Rand
}

func newIdPool() *idPool {
	ret := new(idPool)
	ret.using = make([]bool, idCount)
	ret.returns = make(chan uint16, 10)
	ret.prepared = make(chan uint16, nprepare)
	ret.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	go ret.serve()

	return ret
}

func (self *idPool) pick() uint16 {
	for {
		ret := uint16(self.rand.Uint32())
		if !self.using[ret] {
			return ret
		}
	}
}

func (self *idPool) prepare() {
	id := self.pick()
	self.using[id] = true
	self.nusing++
	self.prepared <- id
}

func (self *idPool) return_(id uint16) bool {
	if !self.using[id] {
		return false
	}

	self.nusing--
	self.using[id] = false
	return true
}

func (self *idPool) serve() {
	for i := 0; i < nprepare; i++ {
		self.prepare()
	}

	for r := range self.returns {
		if self.return_(r) {
			self.prepare()
		}
		bugOn(self.nusing != nprepare)
	}
}

func (self *idPool) Fetch() uint16 {
	return <-self.prepared
}

func (self *idPool) Return(id uint16) {
	self.returns <- id
}
