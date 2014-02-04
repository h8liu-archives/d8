package subcmd

type sub struct {
	Cmd         string // the sub command
	Description string // showed in help message
	Entry       func() // entry function
}
