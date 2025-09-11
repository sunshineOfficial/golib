package reflect

import (
	"reflect"
	"strconv"
)

func setField(valueKind reflect.Kind, valueField reflect.Value, value string) error {
	if ok, err := tryUnmarshalField(valueKind, valueField, value); ok {
		return err
	}

	switch valueKind {
	case reflect.Ptr:
		return setField(valueField.Elem().Kind(), valueField.Elem(), value)
	case reflect.Int:
		return setIntField(valueField, value, 0)
	case reflect.Int8:
		return setIntField(valueField, value, 8)
	case reflect.Int16:
		return setIntField(valueField, value, 16)
	case reflect.Int32:
		return setIntField(valueField, value, 32)
	case reflect.Int64:
		return setIntField(valueField, value, 64)
	case reflect.Uint:
		return setUIntField(valueField, value, 0)
	case reflect.Uint8:
		return setUIntField(valueField, value, 8)
	case reflect.Uint16:
		return setUIntField(valueField, value, 16)
	case reflect.Uint32:
		return setUIntField(valueField, value, 32)
	case reflect.Uint64:
		return setUIntField(valueField, value, 64)
	case reflect.Bool:
		return setBoolField(valueField, value)
	case reflect.Float32:
		return setFloatField(valueField, value, 32)
	case reflect.Float64:
		return setFloatField(valueField, value, 64)
	case reflect.String:
		valueField.SetString(value)
	default:
		return NewUnsupportedTypeError(valueKind)
	}

	return nil
}

func setIntField(fieldValue reflect.Value, value string, bitSize int) error {
	if value == "" {
		fieldValue.SetInt(0)
		return nil
	}

	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		fieldValue.SetInt(intVal)
	}

	return err
}

func setUIntField(fieldValue reflect.Value, value string, bitSize int) error {
	if value == "" {
		fieldValue.SetUint(0)
		return nil
	}

	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		fieldValue.SetUint(uintVal)
	}

	return err
}

func setBoolField(fieldValue reflect.Value, value string) error {
	if value == "" {
		fieldValue.SetBool(false)
		return nil
	}

	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		fieldValue.SetBool(boolVal)
	}

	return err
}

func setFloatField(fieldValue reflect.Value, value string, bitSize int) error {
	if value == "" {
		fieldValue.SetFloat(0)
		return nil
	}

	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		fieldValue.SetFloat(floatVal)
	}

	return err
}
