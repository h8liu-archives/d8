package term

import (
	"printer"
)

type Task interface {
	printer.Printable
	Run(c Cursor)
}
