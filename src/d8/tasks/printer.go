package tasks

import (
	"d8/term"
)

func ShiftOutWith(c term.Cursor, s string) {
	c.ShiftOut()
	c.Print(s)
}
