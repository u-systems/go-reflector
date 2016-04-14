package reflector

import (
	"reflect"
	"fmt"
	"errors"
)

// Set fields values
func setFieldsValues(value *reflect.Value, fields []reflectionField, data map[string]interface{}) error {

	var buffer []reflectionField

	for _, field := range fields {
		fieldValue, ok := data[field.configField.Name]
		if !ok {
			if field.configField.DefaultValue != nil {
				// if field has default value use it
				fieldValue = field.configField.DefaultValue
			} else {
				if field.configField.IsRequired {
					// check if field depends on
					if field.configField.DependsOn.ConfigFieldName != "" {
						buffer = append(buffer, field)
					} else {
						return errors.New(fmt.Sprintf("value for field `%s` is required", field.configField.Name))
					}
				}
			}
		}

		var valueField reflect.Value
		if value.Kind() == reflect.Ptr {
			valueField = value.Elem().Field(field.fieldIndex)
		} else {
			valueField = value.Field(field.fieldIndex)
		}

		if err := setFieldValue(&valueField, field, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func setFieldValue(value *reflect.Value, field reflectionField, data interface{}) error {
	switch  value.Kind() {
	case reflect.Struct:
		if data, ok := data.(map[string]interface{}); !ok {
			return errors.New(fmt.Sprintf("invalid data format for field `%s`", field.configField.Name))
		} else {
			return setFieldsValues(value, field.fields, data)
		}
	// Simple types
	// Strings
	case reflect.String:
		if data == nil {
			data = ""
		}
		if v, ok := data.(string); !ok {
			return errors.New(fmt.Sprintf("invalid type `%T` for field `%s` expected `%s`", data, field.configField.Name, value.Kind()))
		} else {
			value.SetString(v)
		}
	// Processing integers
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		if data == nil {
			data = 0
		}
		var castingError bool = true
		switch reflect.TypeOf(data).Kind() {
		case reflect.Int:
			if v, ok := data.(int); ok {
				value.SetInt(int64(v))
				castingError = false
			}
		case reflect.Int8:
			if v, ok := data.(int8); ok {
				value.SetInt(int64(v))
				castingError = false
			}
		case reflect.Int16:
			if v, ok := data.(int16); ok {
				value.SetInt(int64(v))
				castingError = false
			}
		case reflect.Int32:
			if v, ok := data.(int32); ok {
				value.SetInt(int64(v))
				castingError = false
			}
		case reflect.Int64:
			if v, ok := data.(int64); ok {
				value.SetInt(int64(v))
				castingError = false
			}
		}
		if castingError {
			return errors.New(fmt.Sprintf("invalid type `%T` for field `%s` expected %s", data,
				field.configField.Name, value.Kind()))
		}

	case reflect.Float32, reflect.Float64:
		if data == nil {
			data = 0.0
		}
		if v, ok := data.(float64); !ok {
			return errors.New(fmt.Sprintf("invalid type `%T` for field `%s`", data, field.configField.Name))
		} else {
			value.SetFloat(v)
		}
	case reflect.Bool:
		if data == nil {
			data = false
		}
		if v, ok := data.(bool); !ok {
			return errors.New(fmt.Sprintf("invalid type `%T` for field `%s`", data, field.configField.Name))
		} else {
			value.SetBool(v)
		}
	//
	// Slices
	//
	case reflect.Slice:
		fieldType := field.fieldType
		if sliceData, ok := data.([]interface{}); !ok {
			return errors.New(fmt.Sprintf("invalid type `%T` for field `%s`", data, field.configField.Name))
		} else {
			// create slice
			slice := reflect.MakeSlice(reflect.SliceOf(fieldType.Elem()), 0, 0)
			for index := 0; index < len(sliceData); index++ {
				// create new value for slice elements
				// New creates pointer to type - need to use Elem() to get value
				elemValue := reflect.New(fieldType.Elem())
				elemValue = elemValue.Elem()
				if err := setFieldValue(&elemValue, field, sliceData[index]); err != nil {
					return errors.New(fmt.Sprintf("error (%s) create slice element with index: %d",
						err, index))
				}
				slice = reflect.Append(slice, elemValue) // add element to slice
			}
			if slice.Len() == 0 && field.configField.IsRequired {
				field.hasValue = false
				//return errors.New(fmt.Sprintf("value for field `%s` is required", field.name))

			}
			value.Set(slice)
		}
	}

	return nil
}