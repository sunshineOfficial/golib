package reflect

import (
	"encoding"
	"reflect"
)

func tryUnmarshalSlice(fieldKind reflect.Kind, fieldVal reflect.Value, inputValue []string) (bool, error) {
	count := len(inputValue)
	if fieldKind != reflect.Slice || count == 0 {
		return false, nil
	}

	sliceItemType := fieldVal.Type().Elem().Kind()
	sliceValue := reflect.MakeSlice(fieldVal.Type(), count, count)

	for j := 0; j < count; j++ {
		if err := setField(sliceItemType, sliceValue.Index(j), inputValue[j]); err != nil {
			return true, err
		}
	}

	fieldVal.Set(sliceValue)
	return true, nil
}

func tryUnmarshalField(valueKind reflect.Kind, valueField reflect.Value, value string) (bool, error) {
	switch valueKind {
	case reflect.Ptr:
		return tryUnmarshalPointer(valueField, value)
	default:
		return tryUnmarshalValue(valueField, value)
	}
}

func tryUnmarshalValue(valueField reflect.Value, value string) (bool, error) {
	fieldIValue := valueField.Addr().Interface()
	if unmarshaler, ok := fieldIValue.(encoding.TextUnmarshaler); ok {
		return true, unmarshaler.UnmarshalText([]byte(value))
	}

	return false, nil
}

func tryUnmarshalPointer(fieldValue reflect.Value, value string) (bool, error) {
	if fieldValue.IsNil() {
		fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
	}

	return tryUnmarshalValue(fieldValue.Elem(), value)
}
