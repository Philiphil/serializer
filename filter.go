package serializer

import (
	"reflect"
	"strings"
)

func filterByGroups[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if !isStruct(elemType) {
		return obj
	}

	var newFields []reflect.StructField

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if isFieldExported(field) && isFieldIncluded(field, groups) {
			newFields = append(newFields, field)
		}
	}

	newStructType := reflect.StructOf(newFields)
	newValue := reflect.New(newStructType).Elem()

	for i, field := range newFields {
		fieldName := field.Name
		fieldValue := value.FieldByName(fieldName)
		newValue.Field(i).Set(fieldValue)
	}

	return newValue.Interface().(T)
}

func isFieldIncluded(field reflect.StructField, groups []string) bool {
	if len(groups) == 0 {
		return true //No filtration then
	}

	tag := field.Tag.Get("group")
	if tag == "" {
		return false
	}

	groupList := strings.Split(tag, ",")
	for _, group := range groups {
		for _, g := range groupList {
			if group == g {
				return true
			}
		}
	}

	return false
}

func isFieldExported(field reflect.StructField) bool {
	return field.PkgPath == ""
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}
