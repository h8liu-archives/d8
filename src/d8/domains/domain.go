package domain

type Domain struct {
	name string
	labels []string
}

func (self *Domain) Equals(other *Domain) bool {
	return self.name == other.name
}

func (self *Domain) String() string {
	return self.name
}

func New(s string) *Domain {

	return nil
}
