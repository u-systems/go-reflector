package reflector

import (
	"testing"
	"reflect"
)

func TestProcessingField(t *testing.T) {
	v := struct {
		Name string `config:"name is_required has_default 'test value'"`
		Field struct{
			Label string `config:"label is_required"`
			Value int64 `config:"value has_default 10"`
		      } `config:"field"`
	}{}

	// Test fist field
	nameField := processingField(reflect.TypeOf(v).Field(0), "config")
	fieldField := processingField(reflect.TypeOf(v).Field(1), "config")

	if nameField.fieldType.Kind() != reflect.String {
		t.Error("field type is string")
	}

	// checking for field
	if fieldField.fieldType.Kind() != reflect.Struct {
		t.Error("field type is struct")
	}
}

func TestProcessingTags(t *testing.T) {

	defer func() {
		if err := recover(); err == nil {
			t.Error("there must be a panic")
		}
	}()

	type BadConfig struct {
		SslCert string `config:"ssl_cert is_required_if ssl_mode has_value true"`
	}

	typeOf :=reflect.TypeOf(BadConfig{})
	_ = processingTags(typeOf, "config")

}

func TestUsingPointer(t *testing.T) {
	type Config struct {
		Name string `config:"name is_requred"`
	}
	if fields := processingTags(reflect.TypeOf(&Config{}), "config"); len(fields) != 1 {
		t.Fatal("there must be one field")
	}
}

func TestDublicatesInTags(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("there must be a panic")
		}
	}()

	type StructWithDublicateTags struct {
		Name string `config:"name"`
		Label string `config:"name is_requred"`
	}
	processingTags(reflect.TypeOf(StructWithDublicateTags{}), "config")
}
