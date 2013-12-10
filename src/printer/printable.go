package printer

type Printable interface {
	PrintTo(p *Printer)
}
