package serializer

import (
	"reflect"
	"strings"
)

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func filterByGroups[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()
	if !elementFilterable(elemType) {
		//panic ou log ?
		return obj
	}

	if isStruct(elemType) {
		return filterStructRecursive(obj, groups...)
	} else if elemType.Kind() == reflect.Map {
		//	return filterMap(obj, groups...)
	}
	return obj
}

func filterStructRecursive[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	var newFields []reflect.StructField

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if isFieldExported(field) && isFieldIncluded(field, groups) {
			fieldValue := value.Field(i)

			// Appliquer la fonction récursivement si le champ est de type structure
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
		newValue.Field(i).Set(fieldValue)
	}

	return newValue.Interface().(T)
}

func elementFilterable(elem reflect.Type) bool {
	if isStruct(elem) {
		return true
	} else if elem.Kind() == reflect.Map {
		//	return filterMap(obj, groups...)
	}
	return false
}

func filterStruct[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

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

func isTagIncluded(tag string, groups []string) bool {
	if tag == "" || len(groups) == 0 {
		return true
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

func isFieldIncluded(field reflect.StructField, groups []string) bool {
	tag := field.Tag.Get("group")
	return isTagIncluded(tag, groups)
}

func isKeyExported(key interface{}, groups []string) bool {
	if keyStr, ok := key.(string); ok {
		return isTagIncluded(keyStr, groups)
	}
	return true
}

func isFieldExported(field reflect.StructField) bool {
	return field.PkgPath == ""
}
