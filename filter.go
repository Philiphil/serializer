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
			fieldValue := value.Field(i)

			if isStruct(field.Type) {
				filteredElem := filterByGroups(fieldValue.Interface(), groups...)
				newFields = append(newFields, reflect.StructField{
					Name: field.Name,
					Type: reflect.TypeOf(filteredElem),
					Tag:  field.Tag,
				})
			} else {
				newFields = append(newFields, field)
			}
		}
	}

	newStructType := reflect.StructOf(newFields)
	newValue := reflect.New(newStructType).Elem()

	for i, field := range newFields {
		fieldName := field.Name
		fieldValue := value.FieldByName(fieldName)
		newFieldValue := newValue.Field(i)

		assignFieldValue(field, newFieldValue, fieldValue, groups...)
	}

	return newValue.Interface().(T)
}

func assignFieldValue(field reflect.StructField, destValue reflect.Value, srcValue reflect.Value, groups ...string) {
	if field.Type == srcValue.Type() {
		destValue.Set(srcValue)
	} else if field.Type.AssignableTo(srcValue.Type()) {
		destValue.Set(srcValue)
	} else if isStruct(field.Type) {
		filteredElem := filterByGroups(srcValue.Interface(), groups...)
		destValue.Set(reflect.ValueOf(filteredElem))
	}
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
