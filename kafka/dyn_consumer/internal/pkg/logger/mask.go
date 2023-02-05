package logger

import (
	"reflect"
	"strings"
)

// if the field has a json tag, use the name
func getName(field reflect.StructField) string {
	jsonTag, foundJsonTag := field.Tag.Lookup("json")
	if foundJsonTag {
		jsonTags := strings.SplitN(jsonTag, ",", 2)
		if len(jsonTags) > 0 {
			jsonNameTag := jsonTags[0]
			if jsonNameTag != "" && jsonNameTag != "-" {
				return jsonTags[0]
			}
		}
	}
	return field.Name
}

// s: struct to be masked
// parentHidden:
func maskInterface(s interface{}, parentHidden bool) interface{} {
	defer func() {
		// for all unhandled condition, just recover.
		// maskInterface will return nil
		recover()
	}()
	if s == nil {
		return nil
	}

	s = convert(s)

	sType := reflect.TypeOf(s)
	switch sType.Kind() {
	case reflect.Struct:
		return maskStruct(s, parentHidden)
	case reflect.Slice:
		return maskSlice(s, parentHidden)
	case reflect.Ptr:
		v := reflect.ValueOf(s).Elem().Interface()
		return maskInterface(v, parentHidden)
	default:
		if parentHidden {
			return maskedStr
		}
		return s
	}
}

func maskSlice(s interface{}, parentHidden bool) []interface{} {
	sliceVal := reflect.ValueOf(s)
	m := make([]interface{}, sliceVal.Len())

	for i := 0; i < sliceVal.Len(); i++ {
		v := sliceVal.Index(i)
		m[i] = maskInterface(v.Interface(), parentHidden)
	}
	return m

}

func maskStruct(s interface{}, parentHidden bool) map[string]interface{} {
	m := map[string]interface{}{}
	fields := structFields(s)
	for _, f := range fields {
		if !isExported(f) {
			continue
		}
		fName := getName(f)
		masked := parentHidden || f.Tag.Get("logger") == "-"
		m[fName] = maskInterface(value(s, f), masked)
	}
	return m
}

// reflect helper
func structFields(s interface{}) []reflect.StructField {
	var fields []reflect.StructField
	sType := reflect.TypeOf(s)
	for i := 0; i < sType.NumField(); i++ {
		fields = append(fields, sType.Field(i))
	}
	return fields
}

func isExported(f reflect.StructField) bool {
	// PkgPath is empty for exported fields.
	// See https://golang.org/pkg/reflect/#StructField
	return f.PkgPath == ""
}

func value(s interface{}, f reflect.StructField) interface{} {
	return reflect.ValueOf(s).FieldByIndex(f.Index).Interface()
}
