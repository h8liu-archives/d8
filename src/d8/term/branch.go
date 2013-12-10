package term

type Branch struct {
	Task
	Children []Node
}

func newBranch(t Task) *Branch {
	ret := new(Branch)
	ret.Task = t
	ret.Children = make([]Node, 0, 5)

	return ret
}

func (self *Branch) Add(n Node) {
	if self == nil {
		return
	}
	self.Children = append(self.Children, n)
}
