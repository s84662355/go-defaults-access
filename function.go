package access

import (
	"errors"
	_ "fmt"
	"reflect"
	"strconv"
)

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:

		if valueKind == reflect.Slice {
			vals := strings.Split(val, ",")
			switch structField.Type() {
			case reflect.TypeOf([]string{}):
				structField.Set(reflect.ValueOf(vals))
			case reflect.TypeOf([]int{}):
				t := make([]int, len(vals))
				for i, v := range vals {
					val, err := strconv.ParseInt(v, 10, 32)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = int(val)
				}
				structField.Set(reflect.ValueOf(t))
			case reflect.TypeOf([]int64{}):
				t := make([]int64, len(vals))
				for i, v := range vals {
					val, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = val
				}
				structField.Set(reflect.ValueOf(t))
			case reflect.TypeOf([]uint{}):
				t := make([]uint, len(vals))
				for i, v := range vals {
					val, err := strconv.ParseUint(v, 10, 32)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = uint(val)
				}
				structField.Set(reflect.ValueOf(t))
			case reflect.TypeOf([]uint64{}):
				t := make([]uint64, len(vals))
				for i, v := range vals {
					val, err := strconv.ParseUint(v, 10, 64)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = val
				}
				structField.Set(reflect.ValueOf(t))
			case reflect.TypeOf([]float32{}):
				t := make([]float32, len(vals))
				for i, v := range vals {
					val, err := strconv.ParseFloat(v, 32)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = float32(val)
				}
				structField.Set(reflect.ValueOf(t))
			case reflect.TypeOf([]float64{}):
				t := make([]float64, len(vals))
				for i, v := range vals {
					val, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = val
				}
				structField.Set(reflect.ValueOf(t))
			case reflect.TypeOf([]bool{}):
				t := make([]bool, len(vals))
				for i, v := range vals {
					val, err := parseBool(v)
					if err != nil {
						return fmt.Errorf("%s:%s", prefix, err)
					}
					t[i] = val
				}
				structField.Set(reflect.ValueOf(t))
			}
		}

		return errors.New("unknown type")
	}
	return nil
}

func set(iface interface{}) error {
	typ := reflect.TypeOf(iface).Elem()  //reflect.Type
	val := reflect.ValueOf(iface).Elem() //reflect.Value
	if typ.Kind() != reflect.Struct {
		return errors.New("binding element must be a struct")
	}
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}

		if !structField.IsNil() {
			continue
		}

		defaultValue, ok := typeField.Tag.Lookup("default")
		if !ok {
			continue
		}
		kind := structField.Kind()
		var err error
		switch kind {
		case reflect.Ptr:
			err = setPointer(structField, defaultValue)
			if err == nil {
				continue
			}
		default:
			err = setWithProperType(kind, defaultValue, structField)
			if err == nil {
				continue
			}
		}

		if structField.Type().NumMethod() > 1 && structField.CanInterface() {
			if u, ok := structField.Interface().(Unmarshaler); ok {
				if !u.IsNil() {
					continue
				}
				err = u.Default(defaultValue)
				if err == nil {
					continue
				}
			}
		}

		return err
	}
	return nil

}

func setPointer(v reflect.Value, defaultValue string) error {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
		return setPointer(v, defaultValue)
	} else {
		elem := v.Elem()
		return setWithProperType(elem.Kind(), defaultValue, elem)
	}
}

func parseBool(v string) (bool, error) {
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}
