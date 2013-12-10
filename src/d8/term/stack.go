package term

type stack struct {
	s []*Branch
}

func newStack() *stack {
	ret := new(stack)
	ret.s = make([]*Branch, 0, 20)
	return ret
}

func (self *stack) Push(b *Branch) {
	self.s = append(self.s, b)
}

func (self *stack) Pop() *Branch {
	n := len(self.s)
	if n == 0 {
		return nil
	}
	ret := self.s[n-1]
	self.s = self.s[:n-1]
	return ret
}

func (self *stack) Len() int {
	return len(self.s)
}

func (self *stack) Top() *Branch {
	n := len(self.s)
	if n == 0 {
		return nil
	}
	return self.s[n-1]
}

func (self *stack) TopAdd(n Node) {
	t := self.Top()
	if t != nil {
		t.Add(n)
	}
}
