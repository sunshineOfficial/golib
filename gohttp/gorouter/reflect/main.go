package reflect

import (
	"reflect"
	"strings"
)

func SetValuesToItem(values map[string][]string, tagType TagType, item any) error {
	if item == nil {
		return nil
	}

	itemType := reflect.TypeOf(item).Elem()
	itemValue := reflect.ValueOf(item).Elem()

	if itemType.Kind() == reflect.Map {
		for k, v := range values {
			itemValue.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v[0]))
		}

		return nil
	}

	if itemType.Kind() != reflect.Struct {
		return ErrStructOrMapRequired
	}

	for i := 0; i < itemType.NumField(); i++ {
		fieldTyp := itemType.Field(i)
		fieldVal := itemValue.Field(i)

		if fieldTyp.Anonymous && fieldVal.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}
		if !fieldVal.CanSet() {
			continue
		}

		if err := unmarshalStructField(values, tagType, fieldVal, fieldTyp); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalStructField(values map[string][]string, tagType TagType, fieldVal reflect.Value, fieldTyp reflect.StructField) error {
	fieldKind := fieldVal.Kind()
	fieldName := fieldTyp.Tag.Get(string(tagType))
	if fieldTyp.Anonymous && fieldVal.Kind() == reflect.Struct && fieldName != "" {
		return ErrNoAnonymous
	}

	if len(fieldName) == 0 {
		return nil
	}

	var required bool
	fieldName, required = parseFieldName(tagType, fieldName)

	inputValue, exists := getValue(values, fieldName)
	switch {
	case !exists && required:
		return NewRequiredFieldError(tagType, fieldName)

	case !exists:
		return nil
	}

	if ok, err := tryUnmarshalField(fieldTyp.Type.Kind(), fieldVal, inputValue[0]); ok {
		if err != nil {
			return err
		}

		return nil
	}

	if ok, err := tryUnmarshalSlice(fieldKind, fieldVal, inputValue); ok {
		if err != nil {
			return err
		}

		return nil
	}

	if err := setField(fieldTyp.Type.Kind(), fieldVal, inputValue[0]); err != nil {
		return err
	}

	return nil
}

func parseFieldName(tagType TagType, fieldName string) (name string, required bool) {
	parts := strings.Split(fieldName, ",")
	name = parts[0]

	if tagType == TagTypePath {
		required = true
	}

	for i := 1; i < len(parts); i++ {
		if strings.EqualFold(parts[i], "required") {
			required = true
		}
	}

	return name, required
}

func getValue(values map[string][]string, fieldName string) ([]string, bool) {
	inputValue, exists := values[fieldName]
	if !exists {
		// чтение через values[fieldName] - наиболее быстрый вариант и он должен покрывать большинство кейсов
		// однако для поддержки нечувствительности к регистру (как у encoding/json) нужно пройтись по values вручную

		for k, v := range values {
			if strings.EqualFold(k, fieldName) {
				inputValue = v
				exists = true
				break
			}
		}
	}

	return inputValue, exists
}
