package msgpack

// Codes from the MessagePack specification.
const (
	fixIntCodeMin    = 0x00
	fixIntCodeMax    = 0x7f
	fixMapCodeMin    = 0x80
	fixMapCodeMax    = 0x8f
	fixArrayCodeMin  = 0x90
	fixArrayCodeMax  = 0x9f
	fixStringCodeMin = 0xa0
	fixStringCodeMax = 0xbf
	nilCode          = 0xc0
	unusedCode       = 0xc1 // never used
	falseCode        = 0xc2
	trueCode         = 0xc3
	binary8Code      = 0xc4
	binary16Code     = 0xc5
	binary32Code     = 0xc6
	ext8Code         = 0xc7
	ext16Code        = 0xc8
	ext32Code        = 0xc9
	float32Code      = 0xca
	float64Code      = 0xcb
	uint8Code        = 0xcc
	uint16Code       = 0xcd
	uint32Code       = 0xce
	uint64Code       = 0xcf
	int8Code         = 0xd0
	int16Code        = 0xd1
	int32Code        = 0xd2
	int64Code        = 0xd3
	fixExt1Code      = 0xd4
	fixExt2Code      = 0xd5
	fixExt4Code      = 0xd6
	fixExt8Code      = 0xd7
	fixExt16Code     = 0xd8
	string8Code      = 0xd9
	string16Code     = 0xda
	string32Code     = 0xdb
	array16Code      = 0xdc
	array32Code      = 0xdd
	map16Code        = 0xde
	map32Code        = 0xdf
	negFixIntCodeMin = 0xe0
	negFixIntCodeMax = 0xff
)

type aborted struct{ err error }

func abort(err error) { panic(aborted{err}) }

func handleAbort(err *error) {
	if r := recover(); r != nil {
		if a, ok := r.(aborted); ok {
			*err = a.err
		} else {
			panic(r)
		}
	}
}
