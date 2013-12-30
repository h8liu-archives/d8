package subcmd

import (
	"fmt"
	"io"
	"os"
	"sort"
)

type SubCmd struct {
	subs map[string]*sub
	def  func()
}

func New() *SubCmd {
	return &SubCmd{make(map[string]*sub), nil}
}

func (self *SubCmd) add(sub *sub) error {
	name := sub.Cmd
	if self.subs[name] != nil {
		return fmt.Errorf("subcmd %s already exists", name)
	}

	self.subs[name] = sub
	return nil
}

func (self *SubCmd) Add(f func(), c, desc string) error {
	return self.add(&sub{c, desc, f})
}

func (self *SubCmd) Help(out io.Writer) {
	cmds := make([]string, 0, len(self.subs))
	for c, _ := range self.subs {
		cmds = append(cmds, c)
	}

	sort.Strings(cmds)
	fmt.Fprintf(out, "usage: %s <sub command>\n", os.Args[0])
	for _, c := range cmds {
		v := self.subs[c]
		cmd := v.Cmd
		if cmd == "" {
			cmd = "<nothing>"
		}

		fmt.Fprintf(out, "    %-11s %s\n", cmd, v.Description)
	}
}

func (self *SubCmd) Default(f func()) {
	self.def = f
}

func (self *SubCmd) Main() {
	args := os.Args

	subcmd := ""
	if len(args) > 1 {
		subcmd = args[1]
	}
	sub := self.subs[subcmd]
	if sub == nil {
		if subcmd == "help" {
			self.Help(os.Stdout)
			os.Exit(0)
			return
		}

		if self.def != nil {
			self.def()
			return
		}

		fmt.Fprintf(os.Stderr, "error: unknown subcmd '%s'\n", subcmd)
		os.Exit(1)
		return
	}

	os.Args = os.Args[1:]
	sub.Entry()
}
