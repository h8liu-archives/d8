package packet

// flags
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

	RcodeOkay = iota
	RcodeFormatError
	RcodeNameError
	RcodeNotImplement
	RcodeRefused
)
