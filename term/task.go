package term

import (
	"github.com/h8liu/d8/printer"
)

type Task interface {
	printer.Printable
	Run(c Cursor)
}
