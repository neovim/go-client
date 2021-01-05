package msgpack

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Unmarshaler is the interface implemented by objects that can decode
// themselves from a MessagePack stream.
type Unmarshaler interface {
	UnmarshalMsgPack(d *Decoder) error
}

// ErrInvalidDecodeArg is the invalid argument error.
var ErrInvalidDecodeArg = errors.New("msgpack: argument to Decode must be non-nil pointer, slice or map")

// DecodeConvertError describes a MessagePack value that was not appropriate
// for a value of a specific Go type.
type DecodeConvertError struct {
	// The MessagePack type of the value.
	SrcType Type
	// Option value.
	SrcValue interface{}
	// Type of the Go value that could not be assigned to.
	DestType reflect.Type
}

// Error implements the error interface.
func (e *DecodeConvertError) Error() string {
	if e.SrcValue == nil {
		return fmt.Sprintf("msgpack: cannot convert %s to %s", e.SrcType, e.DestType)
	}
	return fmt.Sprintf("msgpack: cannot convert %s(%v) to %s", e.SrcType, e.SrcValue, e.DestType)
}

func decodeUnsupportedType(ds *decodeState, v reflect.Value) {
	ds.saveErrorAndSkip(v, nil)
}

// decodeState represents the state while decoding value.
type decodeState struct {
	*Decoder
	errSaved error
}

func (ds *decodeState) unpack() {
	if err := ds.Decoder.Unpack(); err != nil {
		abort(err)
	}
}

func (ds *decodeState) skip() {
	if err := ds.Decoder.Skip(); err != nil {
		abort(err)
	}
}

func (ds *decodeState) saveErrorAndSkip(destValue reflect.Value, srcValue interface{}) {
	if ds.errSaved == nil {
		ds.errSaved = &DecodeConvertError{
			SrcType:  ds.Type(),
			SrcValue: srcValue,
			DestType: destValue.Type(),
		}
	}
	ds.skip()
}

// Decode decodes the next value in the stream to v.
//
// Decode uses the inverse of the encodings that Encoder.Encode uses,
// allocating maps, slices, and pointers as necessary, with the following
// additional rules:
//
// To decode into a pointer, Decode first handles the case of a MessagePack
// nil. In that case, Decode sets the pointer to nil. Otherwise, Decode decodes
// the stream into the value pointed at by the pointer. If the pointer is nil,
// Decode allocates a new value for it to point to.
//
// To decode a MessagePack array into a slice, Decode sets the slice length to
// the length of the MessagePack array or reallocates the slice if there is
// insufficient capaicity. Slice elments are not cleared before decoding the
// element.
//
// To decode a MessagePack array into a Go array, Decode decodes the
// MessagePack array elements into corresponding Go array elements.  If the Go
// array is smaller than the MessagePack array, the additional MessagePack
// array elements are discarded. If the MessagePack array is smaller than the
// Go array, the additional Go array elements are set to zero values.
//
// If a MessagePack value is not appropriate for a given target type, or if a
// MessagePack number overflows the target type, Decode skips that field and
// completes the decoding as best it can.  If no more serious errors are
// encountered, Decode returns an DecodeConvertError describing the earliest
// such error.
func (d *Decoder) Decode(v interface{}) (err error) {
	defer handleAbort(&err)
	ds := &decodeState{
		Decoder: d,
	}
	ds.unpack()

	rv := reflect.ValueOf(v)
	if (rv.Kind() != reflect.Ptr && rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map) || rv.IsNil() {
		ds.skip()
		return ErrInvalidDecodeArg
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	decoderForType(rv.Type(), nil)(ds, rv)

	return ds.errSaved
}

var decodeFuncCache struct {
	sync.RWMutex
	m map[reflect.Type]decodeFunc
}

type decodeFunc func(*decodeState, reflect.Value)

type decodeBuilder struct {
	m map[reflect.Type]decodeFunc
}

func decoderForType(t reflect.Type, b *decodeBuilder) decodeFunc {
	decodeFuncCache.RLock()
	f, ok := decodeFuncCache.m[t]
	decodeFuncCache.RUnlock()
	if ok {
		return f
	}

	save := false
	if b == nil {
		b = &decodeBuilder{m: make(map[reflect.Type]decodeFunc)}
		save = true
	} else if f, ok := b.m[t]; ok {
		return f
	}

	// Add temporary entry to break recursion
	b.m[t] = func(ds *decodeState, v reflect.Value) {
		f(ds, v)
	}
	f = b.decoder(t)
	b.m[t] = f

	if save {
		decodeFuncCache.Lock()

		if decodeFuncCache.m == nil {
			decodeFuncCache.m = make(map[reflect.Type]decodeFunc)
		}
		for t, f := range b.m {
			decodeFuncCache.m[t] = f
		}

		decodeFuncCache.Unlock()
	}
	return f
}

func (b *decodeBuilder) decoder(t reflect.Type) decodeFunc {
	if t.Kind() == reflect.Ptr && t.Implements(unmarshalerType) {
		return unmarshalDecoder
	}

	var f decodeFunc
	switch t.Kind() {
	case reflect.Bool:
		f = boolDecoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f = intDecoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		f = uintDecoder
	case reflect.Float32, reflect.Float64:
		f = floatDecoder
	case reflect.String:
		f = stringDecoder
	case reflect.Array:
		f = b.arrayDecoder(t)
	case reflect.Slice:
		f = b.sliceDecoder(t)
	case reflect.Map:
		f = b.mapDecoder(t)
	case reflect.Interface:
		f = interfaceDecoder
	case reflect.Struct:
		f = b.structDecoder(t)
	case reflect.Ptr:
		f = b.ptrDecoder(t)
	default:
		f = decodeUnsupportedType
	}

	if t.Kind() != reflect.Ptr && reflect.PtrTo(t).Implements(unmarshalerType) {
		f = unmarshalAddrDecoder{f}.decode
	}

	return f
}

func boolDecoder(ds *decodeState, v reflect.Value) {
	var x bool

	switch ds.Type() {
	case Bool:
		x = ds.Bool()
	case Int:
		x = ds.Int() != 0
	case Uint:
		x = ds.Uint() != 0
	default:
		ds.saveErrorAndSkip(v, nil)
		return
	}

	v.SetBool(x)
}

func intDecoder(ds *decodeState, v reflect.Value) {
	var x int64

	switch ds.Type() {
	case Int:
		x = ds.Int()
	case Uint:
		n := ds.Uint()
		x = int64(n)
		if x < 0 {
			ds.saveErrorAndSkip(v, n)
			return
		}
	case Float:
		f := ds.Float()
		x = int64(f)
		if float64(x) != f {
			ds.saveErrorAndSkip(v, f)
			return
		}
	default:
		ds.saveErrorAndSkip(v, nil)
		return
	}

	if v.OverflowInt(x) {
		ds.saveErrorAndSkip(v, x)
		return
	}

	v.SetInt(x)
}

func uintDecoder(ds *decodeState, v reflect.Value) {
	var x uint64

	switch ds.Type() {
	case Uint:
		x = ds.Uint()
	case Int:
		i := ds.Int()
		if i < 0 {
			ds.saveErrorAndSkip(v, i)
			return
		}
		x = uint64(i)
	case Float:
		f := ds.Float()
		x = uint64(f)
		if float64(x) != f {
			ds.saveErrorAndSkip(v, f)
			return
		}
	default:
		ds.saveErrorAndSkip(v, nil)
		return
	}

	if v.OverflowUint(x) {
		ds.saveErrorAndSkip(v, x)
		return
	}

	v.SetUint(x)
}

func floatDecoder(ds *decodeState, v reflect.Value) {
	var x float64

	switch ds.Type() {
	case Int:
		i := ds.Int()
		x = float64(i)
		if int64(x) != i {
			ds.saveErrorAndSkip(v, i)
			return
		}
	case Uint:
		n := ds.Uint()
		x = float64(n)
		if uint64(x) != n {
			ds.saveErrorAndSkip(v, n)
			return
		}
	case Float:
		x = ds.Float()
	default:
		ds.saveErrorAndSkip(v, nil)
		return
	}

	v.SetFloat(x)
}

func stringDecoder(ds *decodeState, v reflect.Value) {
	var x string

	switch ds.Type() {
	case Binary, String:
		x = ds.String()
	default:
		ds.saveErrorAndSkip(v, nil)
		return
	}

	v.SetString(x)
}

func byteSliceDecoder(ds *decodeState, v reflect.Value) {
	var x []byte

	switch ds.Type() {
	case Nil:
		// Nothing to do
	case Binary, String:
		// TODO: check if OK to set?
		x = ds.Bytes()
	default:
		ds.saveErrorAndSkip(v, nil)
		return
	}

	v.SetBytes(x)
}

func interfaceDecoder(ds *decodeState, v reflect.Value) {
	if ds.Type() == Nil {
		v.Set(reflect.Zero(v.Type()))
		return
	}

	if v.IsNil() {
		if v.NumMethod() > 0 {
			// We don't know how to make an object of this interface type.
			ds.saveErrorAndSkip(v, nil)
			return
		}
		v.Set(reflect.Value(reflect.ValueOf(decodeNoReflect(ds))))
		return
	}

	v = v.Elem()
	if (v.Kind() == reflect.Ptr ||
		v.Kind() == reflect.Map ||
		v.Kind() == reflect.Slice) && !v.IsNil() {
		decoderForType(v.Type(), nil)(ds, v)
		return
	}

	ds.saveErrorAndSkip(v, nil)
}

type sliceArrayDecoder struct {
	elem decodeFunc
}

func (dec sliceArrayDecoder) decodeArray(ds *decodeState, v reflect.Value) {
	n := ds.Len()
	for i := 0; i < n; i++ {
		ds.unpack()
		if i < v.Len() {
			dec.elem(ds, v.Index(i))
		} else {
			ds.skip()
		}
	}

	if n < v.Len() {
		z := reflect.Zero(v.Type().Elem())
		for i := n; i < v.Len(); i++ {
			v.Index(i).Set(z)
		}
	}
}

func (b *decodeBuilder) arrayDecoder(t reflect.Type) decodeFunc {
	return sliceArrayDecoder{elem: decoderForType(t.Elem(), b)}.decodeArray
}

func (dec sliceArrayDecoder) decodeSlice(ds *decodeState, v reflect.Value) {
	if !v.CanAddr() {
		dec.decodeArray(ds, v)
		return
	}

	n := ds.Len()
	if n > v.Cap() {
		newv := reflect.MakeSlice(v.Type(), n, n)
		reflect.Copy(newv, v)
		v.Set(newv)
	} else {
		v.SetLen(n)
	}

	for i := 0; i < n; i++ {
		ds.unpack()
		dec.elem(ds, v.Index(i))
	}
}

func (b *decodeBuilder) sliceDecoder(t reflect.Type) decodeFunc {
	if t.Elem().Kind() == reflect.Uint8 {
		return byteSliceDecoder
	}

	return sliceArrayDecoder{elem: decoderForType(t.Elem(), b)}.decodeSlice
}

type mapDecoder struct {
	key  decodeFunc
	elem decodeFunc
}

func (dec *mapDecoder) decode(ds *decodeState, v reflect.Value) {
	if ds.Type() != MapLen {
		ds.saveErrorAndSkip(v, nil)
		return
	}

	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	n := ds.Len()
	for i := 0; i < n; i++ {
		ds.unpack()
		key := reflect.New(v.Type().Key()).Elem()
		dec.key(ds, key)

		ds.unpack()
		elem := reflect.New(v.Type().Elem()).Elem()
		dec.elem(ds, elem)

		v.SetMapIndex(key, elem)
	}
}

func (b *decodeBuilder) mapDecoder(t reflect.Type) decodeFunc {
	dec := &mapDecoder{
		key:  decoderForType(t.Key(), b),
		elem: decoderForType(t.Elem(), b),
	}
	return dec.decode
}

type fieldDec struct {
	index []int
	f     decodeFunc
	empty reflect.Value
}

func (fd *fieldDec) setEmpty(v reflect.Value) {
	if !fd.empty.IsValid() {
		return
	}

	fv := fieldByIndex(v, fd.index)
	fv.Set(fd.empty)
}

type structArrayDecoder []*fieldDec

func (dec structArrayDecoder) decode(ds *decodeState, v reflect.Value) {
	for _, fd := range dec {
		fd.setEmpty(v)
	}

	if ds.Type() != ArrayLen {
		ds.saveErrorAndSkip(v, nil)
		return
	}

	n := ds.Len()
	for i := 0; i < n; i++ {
		ds.unpack()
		if i < len(dec) {
			fd := dec[i]
			fv := fieldByIndex(v, fd.index)
			fd.f(ds, fv)
		} else {
			ds.skip()
		}
	}
}

type structDecoder map[string]*fieldDec

func (dec structDecoder) decode(ds *decodeState, v reflect.Value) {
	for _, fd := range dec {
		fd.setEmpty(v)
	}

	if ds.Type() != MapLen {
		ds.saveErrorAndSkip(v, nil)
		return
	}

	n := ds.Len()
	for i := 0; i < n; i++ {
		// Key
		ds.unpack()

		var fd *fieldDec
		if ds.Type() == String || ds.Type() == Binary {
			fd = dec[string(ds.BytesNoCopy())]
		} else {
			ds.saveErrorAndSkip(reflect.ValueOf(""), nil)
		}

		// Value
		ds.unpack()

		if fd != nil {
			fv := fieldByIndex(v, fd.index)
			fd.f(ds, fv)
		} else {
			ds.skip()
		}
	}
}

func (b *decodeBuilder) structDecoder(t reflect.Type) decodeFunc {
	fields, array := fieldsForType(t)

	if array {
		var dec structArrayDecoder
		for _, field := range fields {
			dec = append(dec, &fieldDec{
				index: field.index,
				f:     decoderForType(field.typ, b),
			})
		}
		return dec.decode
	}

	dec := make(structDecoder)
	for _, field := range fields {
		dec[field.name] = &fieldDec{
			index: field.index,
			f:     decoderForType(field.typ, b),
			empty: field.empty,
		}
	}

	return dec.decode
}

type ptrDecoder struct {
	elem decodeFunc
}

func (dec ptrDecoder) decode(ds *decodeState, v reflect.Value) {
	if ds.Type() == Nil {
		v.Set(reflect.Zero(v.Type()))
		return
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	dec.elem(ds, v.Elem())
}

func (b *decodeBuilder) ptrDecoder(t reflect.Type) decodeFunc {
	return ptrDecoder{elem: decoderForType(t.Elem(), b)}.decode
}

var unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()

func unmarshalDecoder(ds *decodeState, v reflect.Value) {
	if ds.Type() == Nil {
		v.Set(reflect.Zero(v.Type()))
		return
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	m := v.Interface().(Unmarshaler)
	err := m.UnmarshalMsgPack(ds.Decoder)
	if e, ok := err.(*DecodeConvertError); ok {
		if ds.errSaved != nil {
			ds.errSaved = e
		}
	} else if err != nil {
		abort(err)
	}
}

type unmarshalAddrDecoder struct{ f decodeFunc }

func (dec unmarshalAddrDecoder) decode(ds *decodeState, v reflect.Value) {
	if !v.CanAddr() {
		dec.f(ds, v)
		return
	}

	unmarshalDecoder(ds, v.Addr())
}

type extensionValue struct {
	kind int
	data []byte
}

func (ev extensionValue) MarshalMsgPack(e *Encoder) error {
	return e.PackExtension(ev.kind, ev.data)
}

func decodeNoReflect(ds *decodeState) (x interface{}) {
	switch ds.Type() {
	case Int:
		return ds.Int()
	case Uint:
		return ds.Uint()
	case Float:
		return ds.Float()
	case Bool:
		return ds.Bool()
	case Nil:
		return nil
	case String:
		return ds.String()
	case Binary:
		return ds.Bytes()
	case ArrayLen:
		n := ds.Len()
		a := make([]interface{}, n)
		for i := 0; i < n; i++ {
			ds.unpack()
			a[i] = decodeNoReflect(ds)
		}
		return a

	case MapLen:
		n := ds.Len()
		m := make(map[string]interface{})
		for i := 0; i < n; i++ {
			ds.unpack()

			if ds.Type() != String && ds.Type() != Binary {
				ds.saveErrorAndSkip(reflect.ValueOf(""), nil)
				ds.unpack()
				ds.skip()
				continue
			}

			key := ds.String()
			ds.unpack()
			m[key] = decodeNoReflect(ds)
		}
		return m

	case Extension:
		if f := ds.extensions[ds.Extension()]; f != nil {
			v, err := f(ds.Bytes())
			if e, ok := err.(*DecodeConvertError); ok {
				if ds.errSaved != nil {
					ds.errSaved = e
				}
			} else if err != nil {
				abort(err)
			}
			return v
		}
		return extensionValue{ds.Extension(), ds.Bytes()}

	default:
		return nil
	}
}
