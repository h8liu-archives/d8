package wire

type Section []*RR

func (self Section) LenU16() uint16 {
	if self == nil {
		return 0
	}

	if len(self) > 0xffff {
		panic("too many rrs")
	}

	return uint16(len(self))
}
