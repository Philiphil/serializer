package serializer

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
)

func (s *Serializer) Deserialize(data string, obj any) error {
	if !isPointer(obj) {
		return fmt.Errorf("object must be pointer")
	}
	switch s.Format {
	case JSON:
		return json.Unmarshal([]byte(data), obj)
	case XML:
		return xml.Unmarshal([]byte(data), obj)
	case CSV:
		return s.deserializeCSV(data, obj)
	default:
		return fmt.Errorf("Unsupported format: %s", s.Format)
	}
}

func (s *Serializer) MergeObjects(target interface{}, source interface{}) error {
	targetValue := reflect.ValueOf(target)
	sourceValue := reflect.ValueOf(source)

	if targetValue.Kind() != reflect.Ptr || sourceValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both target and source must be pointers")
	}

	targetValue = targetValue.Elem()
	sourceValue = sourceValue.Elem()

	for i := 0; i < targetValue.NumField(); i++ {
		targetField := targetValue.Field(i)
		sourceField := sourceValue.Field(i)

		if targetField.CanSet() && !isEmpty(sourceField) {
			targetField.Set(sourceField)
		}
	}

	return nil
}

func (s *Serializer) DeserializeAndMerge(data string, target interface{}) error {
	source := reflect.New(reflect.TypeOf(target).Elem()).Interface()

	if err := s.Deserialize(data, source); err != nil {
		return err
	}

	return s.MergeObjects(target, source)
}

func isEmpty(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

func isPointer(v interface{}) bool {
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Ptr
}

func (s *Serializer) deserializeCSV(data string, obj any) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("Invalid object type for CSV deserialization")
	}

	reader := csv.NewReader(strings.NewReader(data))
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	elemType := value.Elem().Type()
	elem := reflect.New(elemType).Elem()

	for _, row := range rows {
		for i, fieldValue := range row {
			field := elem.Field(i)
			if field.IsValid() && field.CanSet() {
				field.SetString(fieldValue)
			}
		}
	}

	value.Elem().Set(elem)
	return nil
}
