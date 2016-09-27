// Copyright 2016 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msgpack

import (
	"errors"
	"io"
	"math"
)

// Encoder writes values in MessagePack format.
type Encoder struct {
	buf         [32]byte
	w           io.Writer
	writeString func(string) (int, error)
	err         error // permanent error
}

// NewEncoder allocates and initializes a new Unpacker.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{
		w: w,
	}
	if ws, ok := w.(interface {
		WriteString(string) (int, error)
	}); ok {
		e.writeString = ws.WriteString
	} else {
		e.writeString = e.writeStringUnopt
	}
	return e
}

func (e *Encoder) writeStringUnopt(s string) (int, error) {
	if len(s) <= len(e.buf) {
		copy(e.buf[:], s)
		return e.w.Write(e.buf[:len(s)])
	}
	return e.w.Write([]byte(s))
}

type numCodes struct {
	c8, c16, c32, c64 byte
}

var (
	stringLenEncodings = &numCodes{string8Code, string16Code, string32Code, 0}
	binaryLenEncodings = &numCodes{binary8Code, binary16Code, binary32Code, 0}
	arrayLenEncodings  = &numCodes{0, array16Code, array32Code, 0}
	mapLenEncodings    = &numCodes{0, map16Code, map32Code, 0}
	extLenEncodings    = &numCodes{ext8Code, ext16Code, ext32Code, 0}
	uintEncodings      = &numCodes{uint8Code, uint16Code, uint32Code, uint64Code}
)

func (e *Encoder) encodeNum(fc *numCodes, v uint64) []byte {
	switch {
	case fc.c8 != 0 && v <= math.MaxUint8:
		e.buf[0] = fc.c8
		e.buf[1] = byte(v)
		return e.buf[:2]
	case v <= math.MaxUint16:
		e.buf[0] = fc.c16
		e.buf[1] = byte(v >> 8)
		e.buf[2] = byte(v)
		return e.buf[:3]
	case v <= math.MaxUint32:
		e.buf[0] = fc.c32
		e.buf[1] = byte(v >> 24)
		e.buf[2] = byte(v >> 16)
		e.buf[3] = byte(v >> 8)
		e.buf[4] = byte(v)
		return e.buf[:5]
	default:
		e.buf[0] = fc.c64
		e.buf[1] = byte(v >> 56)
		e.buf[2] = byte(v >> 48)
		e.buf[3] = byte(v >> 40)
		e.buf[4] = byte(v >> 32)
		e.buf[5] = byte(v >> 24)
		e.buf[6] = byte(v >> 16)
		e.buf[7] = byte(v >> 8)
		e.buf[8] = byte(v)
		return e.buf[:9]
	}
}

// PackNil writes a Nil value to the MessagePack stream.
func (e *Encoder) PackNil() error {
	e.buf[0] = nilCode
	_, err := e.w.Write(e.buf[:1])
	return err
}

// PackBool writes a Bool value to the MessagePack stream.
func (e *Encoder) PackBool(b bool) error {
	if b {
		e.buf[0] = trueCode
	} else {
		e.buf[0] = falseCode
	}
	_, err := e.w.Write(e.buf[:1])
	return err
}

// PackInt packs an Int value to the MessagePack stream.
func (e *Encoder) PackInt(v int64) error {
	var b []byte
	switch {
	case 0 <= v && v <= math.MaxInt8:
		e.buf[0] = byte(v)
		b = e.buf[:1]
	case v > 0:
		// Pack as unsigned for compatibility with other encoders.
		b = e.encodeNum(uintEncodings, uint64(v))
	case v >= -32:
		e.buf[0] = byte(v)
		b = e.buf[:1]
	case v >= math.MinInt8:
		e.buf[0] = int8Code
		e.buf[1] = byte(v)
		b = e.buf[:2]
	case v >= math.MinInt16:
		e.buf[0] = int16Code
		e.buf[1] = byte(v >> 8)
		e.buf[2] = byte(v)
		b = e.buf[:3]
	case v >= math.MinInt32:
		e.buf[0] = int32Code
		e.buf[1] = byte(v >> 24)
		e.buf[2] = byte(v >> 16)
		e.buf[3] = byte(v >> 8)
		e.buf[4] = byte(v)
		b = e.buf[:5]
	default:
		e.buf[0] = int64Code
		e.buf[1] = byte(v >> 56)
		e.buf[2] = byte(v >> 48)
		e.buf[3] = byte(v >> 40)
		e.buf[4] = byte(v >> 32)
		e.buf[5] = byte(v >> 24)
		e.buf[6] = byte(v >> 16)
		e.buf[7] = byte(v >> 8)
		e.buf[8] = byte(v)
		b = e.buf[:9]
	}
	_, err := e.w.Write(b)
	return err
}

// PackUint packs a Uint value to the message pack stream.
func (e *Encoder) PackUint(v uint64) error {
	var b []byte
	if v <= math.MaxInt8 {
		// Pack as signed for compatibility with other encoders.
		e.buf[0] = byte(v)
		b = e.buf[:1]
	} else {
		b = e.encodeNum(uintEncodings, v)
	}
	_, err := e.w.Write(b)
	return err
}

func (e *Encoder) packArrayMapLen(fixMin int, fc *numCodes, v int) error {
	if v < 0 || v > math.MaxUint32 {
		return errors.New("msgpack: illegal array or map size")
	}
	var b []byte
	if v < 16 {
		e.buf[0] = byte(fixMin + v)
		b = e.buf[:1]
	} else {
		b = e.encodeNum(fc, uint64(v))
	}
	_, err := e.w.Write(b)
	return err
}

// PackArrayLen write an Array length to the MessagePack stream. The
// application must write n objects to the stream following this call.
func (e *Encoder) PackArrayLen(n int) error {
	return e.packArrayMapLen(fixArrayCodeMin, arrayLenEncodings, n)
}

// PackMapLen write an Map length to the MessagePack stream. The application
// must write n key-value pairs to the stream following this call.
func (e *Encoder) PackMapLen(n int) error {
	return e.packArrayMapLen(fixMapCodeMin, mapLenEncodings, n)
}

// PackExtension writes an extension to the MessagePack stream.
func (e *Encoder) PackExtension(kind int, data []byte) error {
	var b []byte
	switch len(data) {
	case 1:
		e.buf[0] = fixext1Code
		e.buf[1] = byte(kind)
		b = e.buf[:2]
	case 2:
		e.buf[0] = fixext2Code
		e.buf[1] = byte(kind)
		b = e.buf[:2]
	case 4:
		e.buf[0] = fixext4Code
		e.buf[1] = byte(kind)
		b = e.buf[:2]
	case 8:
		e.buf[0] = fixext8Code
		e.buf[1] = byte(kind)
		b = e.buf[:2]
	case 16:
		e.buf[0] = fixext16Code
		e.buf[1] = byte(kind)
		b = e.buf[:2]
	default:
		b = e.encodeNum(extLenEncodings, uint64(len(data)))
		b = append(b, byte(kind))
	}
	if _, err := e.w.Write(b); err != nil {
		return err
	}
	_, err := e.w.Write(data)
	return err
}

func (e *Encoder) packStringLen(n int) error {
	var b []byte
	if n < 32 {
		e.buf[0] = byte(fixStringCodeMin + n)
		b = e.buf[:1]
	} else if n <= math.MaxUint32 {
		b = e.encodeNum(stringLenEncodings, uint64(n))
	} else {
		return errors.New("msgpack: long string or binary")
	}
	_, err := e.w.Write(b)
	return err
}

// PackString writes a String value to the MessagePack stream.
func (e *Encoder) PackString(v string) error {
	if err := e.packStringLen(len(v)); err != nil {
		return err
	}
	_, err := e.writeString(v)
	return err
}

// PackStringBytes writes a String value to the MessagePack stream.
func (e *Encoder) PackStringBytes(v []byte) error {
	if err := e.packStringLen(len(v)); err != nil {
		return err
	}
	_, err := e.w.Write(v)
	return err
}

// PackBinary writes a Binary value to the MessagePack stream.
func (e *Encoder) PackBinary(v []byte) error {
	if len(v) > math.MaxUint32 {
		return errors.New("msgpack: long string or binary")
	}
	if _, err := e.w.Write(e.encodeNum(binaryLenEncodings, uint64(len(v)))); err != nil {
		return err
	}
	_, err := e.w.Write(v)
	return err
}

// PackFloat writes a Float value to the MessagePack stream.
func (e *Encoder) PackFloat(f float64) error {
	n := math.Float64bits(f)
	e.buf[0] = float64Code
	e.buf[1] = byte(n >> 56)
	e.buf[2] = byte(n >> 48)
	e.buf[3] = byte(n >> 40)
	e.buf[4] = byte(n >> 32)
	e.buf[5] = byte(n >> 24)
	e.buf[6] = byte(n >> 16)
	e.buf[7] = byte(n >> 8)
	e.buf[8] = byte(n)
	_, err := e.w.Write(e.buf[:9])
	return err
}
