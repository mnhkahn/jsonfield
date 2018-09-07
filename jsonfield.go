package jsonfield

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	fflib "github.com/pquerna/ffjson/fflib/v1"
)

// Marshal is compatible with std json.Marshal.
// If fields is nil it acts as the std json package.
// If fields is not nil it will only return the provided fields.
func Marshal(v interface{}, fields ...string) ([]byte, error) {
	if len(fields) == 0 {
		return json.Marshal(v)
	} else {
		return marshalFields(v, fields)
	}
}

// marshalFields marshals json based on the fields parameters.
// It's based on reflect and ffjson.
func marshalFields(c interface{}, fields []string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	_, v := reflectTypeAndValue(c)

	buf.WriteString("{")
	for i, field := range fields {
		vv := v.FieldByName(field)
		if vv.Kind() != reflect.Invalid {
			buf.WriteString(`"` + field + `":`)

			var err error
			switch vv.Kind() {
			case reflect.Bool:
				zzz := strconv.AppendBool(nil, vv.Bool())
				_, err = buf.Write(zzz)
			case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				i := vv.Int()
				fflib.FormatBits2(buf, uint64(i), 10, i < 0)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
				i := vv.Uint()
				fflib.FormatBits2(buf, uint64(i), 10, i < 0)
			case reflect.Float32, reflect.Float64:
				zzz := strconv.AppendFloat(nil, vv.Float(), 'g', -1, 64)
				_, err = buf.Write(zzz)
			case reflect.String:
				_, err = buf.WriteString(vv.String())
			case reflect.Slice, reflect.Map, reflect.Array:
				if vv.IsNil() {
					buf.WriteString("null")
				} else {
					_, err = buf.WriteString(fmt.Sprintf(`"%+v"`, vv.Interface()))
				}
			case reflect.Struct:
				tt := reflect.TypeOf(vv.Interface())
				if tt.Implements(marshalerType) {
					m := vv.Interface().(json.Marshaler)
					b, err := m.MarshalJSON()
					if err != nil {
						return nil, err
					} else {
						b, err := json.Marshal(vv.Interface())
						if err != nil {
							return nil, err
						}
						_, err = buf.Write(b)
					}
					_, err = buf.Write(b)
				}
			case reflect.Ptr:
				tt := reflect.TypeOf(vv.Interface()).Elem()
				if tt.Implements(marshalerType) {
					m := vv.Interface().(json.Marshaler)
					b, err := m.MarshalJSON()
					if err != nil {
						return nil, err
					}
					_, err = buf.Write(b)
				} else {
					b, err := json.Marshal(vv.Interface())
					if err != nil {
						return nil, err
					}
					_, err = buf.Write(b)
				}
			default:
				_, err = buf.WriteString(vv.String())
			}

			if err != nil {
				return nil, err
			}
		}

		if i != len(fields)-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}

// reflectTypeAndValue ...
func reflectTypeAndValue(c interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = reflect.Indirect(v)
	}

	return t, v
}

var (
	marshalerType = reflect.TypeOf(new(json.Marshaler)).Elem()
)
