package term

type Branch struct {
	Task
	Children []Node
}

var _ Node = new(Branch)

func newBranch(t Task) *Branch {
	ret := new(Branch)
	ret.Task = t
	ret.Children = make([]Node, 0, 5)

	return ret
}

func (self *Branch) add(n Node) {
	if self == nil {
		return
	}
	self.Children = append(self.Children, n)
}

func (self *Branch) IsLeaf() bool { return false }
