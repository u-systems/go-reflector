package providers

import (
	"encoding/json"
	"errors"
)

type JsonDataProvider struct {
	data []byte
}

func NewJsonDataProvider(data []byte) *JsonDataProvider {
	return &JsonDataProvider{
		data: data,
	}
}

func (provider *JsonDataProvider) Load() (map[string]interface{}, error) {
	var result interface{}
	if err := json.Unmarshal(provider.data, &result); err != nil {
		return nil, err
	}
	if v, ok := result.(map[string]interface{}); !ok {
		return nil, errors.New("error cast ")
	} else {
		return v, nil
	}
}

func (provider *JsonDataProvider) Unload(data map[string]interface{}) error {
	if result, err := json.Marshal(&data);  err != nil {
		return err
	} else {
		provider.data = result
	}
	return nil
}

func (provider *JsonDataProvider) Data() interface{} {
	return provider.data
}
