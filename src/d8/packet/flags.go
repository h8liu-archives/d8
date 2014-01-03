package packet

import (
	"fmt"
	"strings"
)

const (
	FlagResponse = 0x1 << 15

	FlagRA = 0x1 << 7
	FlagRD = 0x1 << 8
	FlagTC = 0x1 << 9
	FlagAA = 0x1 << 10

	RcodeMask = 0xf
	OpMask    = 0x3 << 11

	OpQuery  = 0 << 11
	OpIquery = 1 << 11
	OpStatus = 2 << 11

	RcodeOkay        = 0
	RcodeFormatError = iota
	RcodeServerFail
	RcodeNameError
	RcodeNotImplement
	RcodeRefused
)

type flagTags struct {
	tags []string
}

func newFlagTags() *flagTags {
	ret := new(flagTags)
	ret.tags = make([]string, 0, 10)
	return ret
}

func (self *flagTags) Tag(b bool, s string) {
	if b {
		self.tags = append(self.tags, s)
	}
}

func (self *flagTags) String() string {
	return strings.Join(self.tags, " ")
}

func rcode(flag uint16) uint16 {
	return flag & RcodeMask
}

func flagString(flag uint16) string {
	t := newFlagTags()

	t.Tag((flag&FlagResponse) == 0, "query")
	t.Tag((flag&OpMask) == OpStatus, "status")
	t.Tag((flag&OpMask) == OpIquery, "iquery")
	t.Tag((flag&FlagAA) != 0, "auth")
	t.Tag((flag&FlagTC) != 0, "trunc")
	t.Tag((flag&FlagRD) != 0, "rec-desir")
	t.Tag((flag&FlagRA) != 0, "rec-avail")

	rcode := rcode(flag)
	t.Tag(rcode == RcodeFormatError, "fmt-err")
	t.Tag(rcode == RcodeServerFail, "serv-fail")
	t.Tag(rcode == RcodeNameError, "name-err")
	t.Tag(rcode == RcodeNotImplement, "not-impl")
	t.Tag(rcode == RcodeRefused, "refused")
	t.Tag(rcode > RcodeRefused, fmt.Sprintf("rcode%d", rcode))

	return t.String()
}
