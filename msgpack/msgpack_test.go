package msgpack

import (
	"reflect"
)

type (
	arrayLen uint32
)

type (
	mapLen uint32
)

type extension struct {
	k int
	d string
}

type testExtension1 struct {
	data []byte
}

func (x *testExtension1) UnmarshalMsgPack(dec *Decoder) error {
	if dec.Type() != Extension || dec.Extension() != 1 {
		err := &DecodeConvertError{
			SrcType:  dec.Type(),
			DestType: reflect.TypeOf(x),
		}
		dec.Skip()
		return err
	}
	x.data = dec.Bytes()
	return nil
}

func (x testExtension1) MarshalMsgPack(enc *Encoder) error {
	return enc.PackExtension(1, x.data)
}

var testExtensionMap = ExtensionMap{
	1: func(data []byte) (interface{}, error) { return testExtension1{data}, nil },
}
