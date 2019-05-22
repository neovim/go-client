// +build bench

package msgpack_test

type Encoder interface {
	Encode(interface{}) error
}

type EncodeStringer interface {
	EncodeString(string) error
}

type EncodeMapper interface {
	EncodeMap(interface{}) error
}

type EncodeArrayer interface {
	EncodeArray(interface{}) error
}

type EncodeUinter interface {
	EncodeUint(uint) error
}

type EncodeUint8er interface {
	EncodeUint8(uint8) error
}

type EncodeUint16er interface {
	EncodeUint16(uint16) error
}

type EncodeUint32er interface {
	EncodeUint32(uint32) error
}

type EncodeUint64er interface {
	EncodeUint64(uint64) error
}

type EncodeFloat32er interface {
	EncodeFloat32(float32) error
}

type EncodeFloat64er interface {
	EncodeFloat64(float64) error
}

type EncodeInter interface {
	EncodeInt(int) error
}

type EncodeInt8er interface {
	EncodeInt8(int8) error
}

type EncodeInt16er interface {
	EncodeInt16(int16) error
}

type EncodeInt32er interface {
	EncodeInt32(int32) error
}

type EncodeInt64er interface {
	EncodeInt64(int64) error
}

type Decoder interface {
	Decode(interface{}) error
}

type DecodeArrayer interface {
	DecodeArray(*[]interface{}) error
}

type DecodeArrayReturner interface {
	DecodeArray() ([]interface{}, error)
}

type DecodeMapper interface {
	DecodeMap(*map[string]interface{}) error
}

type DecodeStringer interface {
	DecodeString(*string) error
}

type DecodeStringReturner interface {
	DecodeString() (string, error)
}

type DecodeUinter interface {
	DecodeUint(*uint) error
}

type DecodeUint8er interface {
	DecodeUint8(*uint8) error
}

type DecodeUint8Returner interface {
	DecodeUint8() (uint8, error)
}

type DecodeUint16er interface {
	DecodeUint16(*uint16) error
}

type DecodeUint16Returner interface {
	DecodeUint16() (uint16, error)
}

type DecodeUint32er interface {
	DecodeUint32(*uint32) error
}

type DecodeUint32Returner interface {
	DecodeUint32() (uint32, error)
}
type DecodeUint64er interface {
	DecodeUint64(*uint64) error
}

type DecodeUint64Returner interface {
	DecodeUint64() (uint64, error)
}

type DecodeInter interface {
	DecodeInt(*int) error
}

type DecodeInt8er interface {
	DecodeInt8(*int8) error
}

type DecodeInt8Returner interface {
	DecodeInt8() (int8, error)
}

type DecodeInt16er interface {
	DecodeInt16(*int16) error
}

type DecodeInt16Returner interface {
	DecodeInt16() (int16, error)
}

type DecodeInt32er interface {
	DecodeInt32(*int32) error
}

type DecodeInt32Returner interface {
	DecodeInt32() (int32, error)
}
type DecodeInt64er interface {
	DecodeInt64(*int64) error
}

type DecodeInt64Returner interface {
	DecodeInt64() (int64, error)
}

type DecodeFloat32er interface {
	DecodeFloat32(*float32) error
}

type DecodeFloat32Returner interface {
	DecodeFloat32() (float32, error)
}
type DecodeFloat64er interface {
	DecodeFloat64(*float64) error
}

type DecodeFloat64Returner interface {
	DecodeFloat64() (float64, error)
}
