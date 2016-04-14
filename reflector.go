package reflector

import (
	"reflect"
	"errors"
)

type Reflector struct  {
	source interface{}
	tagName string
	fields []reflectionField
}

// Data provider interface
type DataProvider interface {
	Load() (map[string]interface{}, error)
	Unload(map[string]interface{}) error
	Data() interface{}
}

// Create new reflector
func New(source interface{}, tagName string) (*Reflector, error) {
	// Check if source is pointer
	if reflect.TypeOf(source).Kind() != reflect.Ptr {
		return nil, errors.New("source must be a pointer")
	}
	// Check if source is
	if reflect.TypeOf(source).Elem().Kind() != reflect.Struct {
		return nil, errors.New("source must be a struct")
	}
	// check tag
	if tagName == "" {
		return nil, errors.New("tagName can not be an empty")
	}
	fields := processingTags(reflect.TypeOf(source).Elem(), tagName)
	if len(fields) == 0 {
		return nil, errors.New("source does not have configuration tags")
	}

	return &Reflector{
		source: source,
		tagName: tagName,
		fields: fields,
	}, nil
}

// Get template for reflection source
func (reflection *Reflector) Template(provider DataProvider) (interface{}, error) {

	raw := map[string]interface{}{}
	for _, field := range reflection.fields {
		raw[field.configField.Name] = field.GetInfo()
	}

	if err := provider.Unload(raw); err != nil {
		return nil, err
	}

	return provider.Data(), nil
}

// Set values and return
func (reflection *Reflector) SetValues(provider DataProvider) (interface{}, error){
	// get data from provider
	data, err := provider.Load()
	if err != nil {
		return nil, err
	}
	valueOf := reflect.ValueOf(reflection.source)
	if err := setFieldsValues(&valueOf, reflection.fields, data); err != nil {
		return nil, err
	}


	return reflection.source, nil
}



