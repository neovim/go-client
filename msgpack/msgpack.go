package msgpack

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
	unusedCode       = 0xc1
	falseCode        = 0xc2
	trueCode         = 0xc3
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
	string8Code      = 0xd9
	string16Code     = 0xda
	string32Code     = 0xdb
	binary8Code      = 0xc4
	binary16Code     = 0xc5
	binary32Code     = 0xc6
	array16Code      = 0xdc
	array32Code      = 0xdd
	map16Code        = 0xde
	map32Code        = 0xdf
	fixext1Code      = 0xd4
	fixext2Code      = 0xd5
	fixext4Code      = 0xd6
	fixext8Code      = 0xd7
	fixext16Code     = 0xd8
	ext8Code         = 0xc7
	ext16Code        = 0xc8
	ext32Code        = 0xc9
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

type field struct {
	name      string
	omitEmpty bool
	array     bool
	index     []int
	typ       reflect.Type
	empty     reflect.Value
}

func collectFields(fields []*field, t reflect.Type, visited map[reflect.Type]bool, depth map[string]int, index []int) []*field {
	// Break recursion.
	if visited[t] {
		return fields
	}
	visited[t] = true

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if sf.PkgPath != "" && !sf.Anonymous {
			// Skip field if not exported and not anonymous.
			continue
		}

		var (
			name      string
			omitEmpty bool
			array     bool
		)
		for i, p := range strings.Split(sf.Tag.Get("msgpack"), ",") {
			if i == 0 {
				name = p
			} else if p == "omitempty" {
				omitEmpty = true
			} else if p == "array" {
				array = true
			} else {
				panic(fmt.Errorf("msgpack: unknown field tag %s for type %s", p, t.Name()))
			}
		}

		if name == "-" {
			// Skip field when field tag starts with "-".
			continue
		}

		ft := sf.Type
		if ft.Name() == "" && ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}

		if name == "" && sf.Anonymous && ft.Kind() == reflect.Struct {
			// Flatten anonymous struct field.
			fields = collectFields(fields, ft, visited, depth, append(index, i))
			continue
		}

		if name == "" {
			name = sf.Name
		}

		// Check for name collisions.
		d, found := depth[name]
		if !found {
			d = 65535
		}
		if len(index) == d {
			// There is another field with same name and same depth.
			// Remove that field and skip this field.
			j := 0
			for i := 0; i < len(fields); i++ {
				if name != fields[i].name {
					fields[j] = fields[i]
					j++
				}
			}
			fields = fields[:j]
			continue
		}
		depth[name] = len(index)

		f := &field{
			name:      name,
			omitEmpty: omitEmpty,
			array:     array,
			index:     make([]int, len(index)+1),
			typ:       sf.Type,
		}
		copy(f.index, index)
		f.index[len(index)] = i

		// Parse empty field tag.

		if e := sf.Tag.Get("empty"); e != "" {
			switch sf.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				bits := 0
				if sf.Type.Kind() != reflect.Int {
					bits = sf.Type.Bits()
				}
				v, err := strconv.ParseInt(e, 10, bits)
				if err != nil {
					panic(fmt.Errorf("msgpack: error parsing field empty field %s.%s: %v", t.Name(), sf.Name, err))
				}
				f.empty = reflect.New(sf.Type).Elem()
				f.empty.SetInt(v)
			case reflect.Bool:
				v, err := strconv.ParseBool(e)
				if err != nil {
					panic(fmt.Errorf("msgpack: error parsing field empty field %s.%s: %v", t.Name(), sf.Name, err))
				}
				f.empty = reflect.New(sf.Type).Elem()
				f.empty.SetBool(v)
			case reflect.String:
				f.empty = reflect.New(sf.Type).Elem()
				f.empty.SetString(e)
			default:
				panic(fmt.Errorf("msgpack: unsupported empty field %s.%s", t.Name(), sf.Name))
			}
		}

		fields = append(fields, f)

	}
	return fields
}

func fieldsForType(t reflect.Type) ([]*field, bool) {
	fields := collectFields(nil, t, make(map[reflect.Type]bool), make(map[string]int), nil)
	array := false
	for _, field := range fields {
		if field.array {
			array = true
			break
		}
	}
	return fields, array
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	for _, i := range index {
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return reflect.Value{}
			}
			v = v.Elem()
		}
		v = v.Field(i)
	}
	return v
}
