package subcmd

import (
	"io"
)

var theSub = New()

func Add(f func(), sub, desc string) error { return theSub.Add(f, sub, desc) }
func Help(out io.Writer)                   { theSub.Help(out) }
func Main()                                { theSub.Main() }
