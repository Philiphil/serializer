package serializer

import (
	"fmt"
	"reflect"
	"strings"
)

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
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

func isList(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array || (t.Kind() == reflect.Ptr && (t.Elem().Kind() == reflect.Slice || t.Elem().Kind() == reflect.Array))
}
func isMap(t reflect.Type) bool {
	return t.Kind() == reflect.Map || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Map)
}

func assignFieldValue(field reflect.StructField, destValue reflect.Value, srcValue reflect.Value) {
	if field.Type == srcValue.Type() {
		destValue.Set(srcValue)
	} else if field.Type.AssignableTo(srcValue.Type()) {
		destValue.Set(srcValue)
	} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem() == srcValue.Type() {
		newPtr := reflect.New(field.Type.Elem())
		newPtr.Elem().Set(srcValue)
		destValue.Set(newPtr)
	} else if isStruct(srcValue.Type()) {
		destFieldType := destValue.Type()
		if destFieldType.Kind() == reflect.Ptr {
			destFieldType = destFieldType.Elem()
		}

		if destFieldType.Kind() == reflect.Struct {
			destValueConverted := reflect.New(destFieldType).Elem()

			for i := 0; i < destFieldType.NumField(); i++ {
				destField := destFieldType.Field(i)
				srcFieldValue := srcValue
				if srcFieldValue.Type().Kind() == reflect.Ptr {
					srcFieldValue = srcFieldValue.Elem()
				}

				srcFieldValue = srcFieldValue.FieldByName(destField.Name)
				destFieldValue := destValueConverted.Field(i)

				if isStruct(srcFieldValue.Type()) && destFieldValue.Kind() == reflect.Struct {
					assignFieldValue(destField, destFieldValue, srcFieldValue)
				} else {
					assignFieldValue(destField, destFieldValue, srcFieldValue)
				}
			}

			destValue.Set(destValueConverted)
		} else {
			destValue.Set(srcValue.Convert(destFieldType))
		}
	} else {
		if !srcValue.Type().ConvertibleTo(destValue.Type()) {
			fmt.Printf("Type conversion not supported from %v to %v\n", srcValue.Type(), destValue.Type())
			panic("!")
		}
		destValue.Set(srcValue.Convert(destValue.Type()))
	}
}

func filterByGroups[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if isStruct(elemType) {
		return filterByGroups_struct(obj, groups...)
	}
	if isList(elemType) {
		return filterByGroups_slice(obj, groups...)
	}
	if isMap(elemType) {
		return filterByGroups_map(obj, groups...)
	}
	return obj
}

func filterByGroups_struct[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
		value = value.Elem()
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
		assignFieldValue(field, newFieldValue, fieldValue)
	}

	return newValue.Interface().(T)
}

func filterByGroups_slice[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if value.Len() == 0 {
		return obj
	}

	firstElem := value.Index(0).Interface()

	filteredFirstElem := filterByGroups(firstElem, groups...)

	newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(filteredFirstElem)), 0, value.Len())

	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		filteredElem := filterByGroups(elem.Interface(), groups...)
		newSlice = reflect.Append(newSlice, reflect.ValueOf(filteredElem))
	}
	return newSlice.Interface().(T)
}

func filterByGroups_map[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		elemType = elemType.Elem()
	}

	mapType := reflect.MapOf(elemType.Key(), elemType.Elem())
	mapValue := reflect.MakeMap(mapType)

	iter := value.MapRange()

	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		filteredVal := filterByGroups(val.Interface(), groups...)

		filteredValValue, ok := filteredVal.(reflect.Value)
		if !ok {
			filteredValValue = reflect.ValueOf(filteredVal)
		}

		if filteredValValue.Type().AssignableTo(elemType.Elem()) {
			mapValue.SetMapIndex(key, filteredValValue)
		} else {
			destType := elemType.Elem()
			destValue := reflect.New(destType).Elem()

			assignFieldValue(destType.Field(0), destValue.Field(0), filteredValValue)

			mapValue.SetMapIndex(key, destValue)
		}
	}

	return mapValue.Interface().(T)
}
