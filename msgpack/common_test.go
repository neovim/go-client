// Copyright 2015 Gary Burd. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package msgpack

import (
	"bytes"
	"fmt"
	"io"
)

type arrayLen uint32
type mapLen uint32
type extension struct {
	k int
	d string
}

// unpack unpacks a byte slice to the following types.
//
//   Type      Go
//   Nil       nil
//   Bool      bool
//   Int       int
//   Uint      int
//   Float     float64
//   ArrayLen  arrayLen
//   MapLen    mapLen
//   String    string
//   Binary    []byte
//   Extension extension
//
// This function is not suitable for unpack tests because the integer and float
// types are mapped to int and float64 respectively.
func unpack(p []byte) ([]interface{}, error) {
	var data []interface{}
	u := NewDecoder(bytes.NewReader(p))
	for {
		err := u.Unpack()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		var v interface{}
		switch u.Type() {
		case Nil:
			v = nil
		case Bool:
			v = u.Bool()
		case Int:
			v = int(u.Int())
		case Uint:
			v = int(u.Uint())
		case Float:
			v = u.Float()
		case Binary:
			v = u.Bytes()
		case String:
			v = u.String()
		case ArrayLen:
			v = arrayLen(u.Int())
		case MapLen:
			v = mapLen(u.Int())
		case Extension:
			v = extension{u.Extension(), u.String()}
		default:
			return nil, fmt.Errorf("unpack %d not handled", u.Type())
		}
		data = append(data, v)
	}
	return data, nil
}

// pack packs the values vs and returns the result.
//
//  Go Type     Encoder method
//  nil         PackNil
//  bool        PackBool
//  int64       PackInt
//  uint64      PackUint
//  float64     PackFloat
//  arrayLen    PackArrayLen
//  mapLen      PackMapLen
//  string      PackString(s, false)
//  []byte      PackBytes(s, true)
//  extension   PackExtension(k, d)
func pack(vs ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for _, v := range vs {
		var err error
		switch v := v.(type) {
		case int64:
			err = enc.PackInt(v)
		case uint64:
			err = enc.PackUint(v)
		case bool:
			err = enc.PackBool(v)
		case float64:
			err = enc.PackFloat(v)
		case arrayLen:
			err = enc.PackArrayLen(int(v))
		case mapLen:
			err = enc.PackMapLen(int(v))
		case string:
			err = enc.PackString(v)
		case []byte:
			err = enc.PackBinary(v)
		case extension:
			err = enc.PackExtension(v.k, []byte(v.d))
		case nil:
			err = enc.PackNil()
		default:
			err = fmt.Errorf("no pack for type %T", v)
		}
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
