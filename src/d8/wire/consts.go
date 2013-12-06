package wire

// rdata type
const (
	A     = 1
	NS    = 2
	MD    = 3
	MF    = 4
	CNAME = 5
	SOA   = 6
	MB    = 7
	MG    = 8
	MR    = 9
	NULL  = 10
	WKS   = 11
	PTR   = 12
	HINFO = 13
	MINFO = 14
	MX    = 15
	TXT   = 16
	AAAA  = 28
)

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

// class code
const (
	IN = 1
	CS = 2
	CH = 3
	HS = 4
)

const DnsPort = 53
