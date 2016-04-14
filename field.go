package reflector

import (
	"stash.abc.ee/micro/reflector/parser"
	"reflect"
	"unicode"
	"fmt"
)

type reflectionField struct {
	configField *parser.ConfigField
	hasValue    bool
	fieldIndex  int
	fieldType   reflect.Type
	isStruct    bool
	fields      []reflectionField
}

// get string information for field
func (reflectionField reflectionField) GetInfo() interface{} {
	var s string

	// detect type
	switch kind := reflectionField.fieldType.Kind(); kind {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s = "int"
	case reflect.Float32, reflect.Float64:
		s = "float"
	case reflect.String:
		s = "string"
	case reflect.Slice:
		// use slice
		value := []interface{}{}
		for _, field := range reflectionField.fields {
			value = append(value, field.GetInfo())
		}
		return value
	case reflect.Struct:
		// use map of strings
		value := map[string]interface{}{}
		for _, field := range reflectionField.fields {
			value[field.configField.Name] = field.GetInfo()
		}
		return value
	}

	// processing default value
	if reflectionField.configField.IsRequired {
		s += " required"
	}
	if v := reflectionField.configField.DefaultValue; v != nil {
		s += fmt.Sprintf(" default %v", v)
	}
	// processing is_required and depends on
	return s
}

// Processing tags
func processingTags(st reflect.Type, tagName string) []reflectionField {
	if st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	fields := []reflectionField{}
	for fieldIndex := 0; fieldIndex < st.NumField(); fieldIndex++ {
		field := st.Field(fieldIndex)
		if newField := processingField(field, tagName); newField != nil {
			for _, field := range fields {
				if newField.configField.Name == field.configField.Name {
					panic(fmt.Sprintf("there are already field with name `%s`", newField.configField.Name))
				}
			}
			newField.fieldIndex = fieldIndex
			fields = append(fields, *newField)
		}
	}
	// check depends on
	for i, f := range fields {
		if name := f.configField.DependsOn.ConfigFieldName; name != "" {
			found := false
			for ii, df := range fields {
				if ii != i {
					if df.configField.Name == name {
						found = true
						break
					}
				}
			}
			if !found {
				panic(fmt.Sprintf("field `%s` depends on `%s` which does not exists in struct", f.configField.Name, name))
			}
		}
	}
	return fields
}

// Internal processing of field
func processingField(field reflect.StructField, tagName string) *reflectionField {
	reflectionField := reflectionField{}
	if field.Type.Kind() == reflect.Ptr {
		panic("pointer is not allowed")
	}
	reflectionField.fieldType = field.Type
	name := []rune(field.Name)
	name[0] = unicode.ToLower(name[0])
	if string(name) == field.Name {
		panic(fmt.Sprintf("field %s is unexported", field.Name))
	}

	// processing tag
	p := parser.NewParser(field.Tag.Get(tagName))
	if configField, err := p.Parse(); err != nil {
		return nil // return empty field
	} else {
		reflectionField.configField = configField
		if reflectionField.configField.Name == "" {
			reflectionField.configField.Name = field.Name
		}
	}

	// Processing struct and slices
	if field.Type.Kind() == reflect.Struct {
		// processing struct
		reflectionField.isStruct = true
		reflectionField.fields = processingTags(field.Type, tagName)

	} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct {
		// processing slice of struct
		reflectionField.fields = processingTags(field.Type.Elem(), tagName)
	}

	return &reflectionField

}
