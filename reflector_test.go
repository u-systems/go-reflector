package reflector

import (
	"testing"
	"stash.abc.ee/micro/reflector/providers"
)

type TestDataProvider struct {
	data []byte
}

func TestNew(t *testing.T) {

	testStruct := &struct {
		Name string `config:"name is_required"`
	}{}

	if reflector, err := New(testStruct, "config"); err != nil {
		t.Errorf("there can not be an error", err)
	} else {
		if reflector.tagName != "config" {
			t.Error("reflector tag name must be `config`")
		}
	}

	// Test empty tag
	if _, err := New(testStruct, ""); err == nil {
		t.Error("There must be an error for empty tag")
	}

}

func TestProcessingStructWithoutTags(t *testing.T) {
	testStruct := &struct {
		Name string
	}{}

	if _, err := New(testStruct, "config"); err == nil {
		t.Error("there must an error for struct without tag")
	}

}

func TestPointerAndStructError(t *testing.T) {
	notAStruct := "error_value"
	if _, err := New(notAStruct, "config"); err == nil {
		t.Error("there must be an error for not pointer value")
	}
	if _, err := New(&notAStruct, "config"); err == nil {
		t.Error("there must be an error beacause of non struct use")
	}
}

func TestTemplate(t *testing.T) {
	config := &struct {
		Host string `config:"host is_required"`
		Port int64 `config:"port has_default 8080"`
		Percent float64 `config:"percent"`
		Server struct {
			Name string `config:"name is_required"`
			Params []struct{
				Label string `config:"label"`
			} `config:"params"`
		       } `config:"server"`
		Parameters []string `config:"parameters"`
	}{}
	if r, err := New(config, "config"); err != nil {
		t.Errorf("there can not be an error: %s", err)
	} else {
		provider := providers.NewJsonDataProvider(nil)
		if template, err := r.Template(provider); err != nil {
			t.Errorf("there can not be error: %s", err)
		} else {
			t.Log(string(template.([]byte)))
		}
	}
}

func TestSetValues(t *testing.T) {

	type Config struct {
		Host string `config:"host is_required"`
		Port int `config:"port has_default 8080"`
	}

	r, err := New(&Config{}, "config")
	if err != nil {
		t.Fatalf("there can not be an error: %s", err)
	}

	provider := providers.NewJsonDataProvider([]byte(`{"host":"127.0.0.1"}`))
	if filledConfig, err := r.SetValues(provider); err != nil {
		t.Errorf("there can not be an error - %s", err)
	} else {
		if config, ok := filledConfig.(*Config); !ok {
			t.Logf("%T\n", filledConfig)
			t.Error("failed to cast config to proper type")
		} else {
			t.Logf("%+v\n", config)
			if config.Host != "127.0.0.1" {
				t.Error("invalid value for Host")
			}
			if config.Port != 8080 {
				t.Error("invalid value for port")
			}
		}
	}
}

func TestSetFloatValue(t *testing.T) {
	type Config struct {
		Percent float64 `config:"percent"`
	}
	r, err := New(&Config{}, "config")
	if err != nil {
		t.Fatalf("fatal error - %s", err)
	}

	provider := providers.NewJsonDataProvider([]byte(`{"percent":1.2}`))
	if config, err := r.SetValues(provider); err != nil {
		t.Errorf("error - %s", err)
	} else {
		if config, ok := config.(*Config); !ok {
			t.Errorf("failed to convert to proper type")
		} else {
			if config.Percent != 1.2 {
				t.Errorf("invalid value: %v", config.Percent)
			}
		}
	}
}

func TestSetStructValue(t *testing.T) {
	type Config struct {
		Server struct{
			Name string `config:"name"`
		       } `config:"server"`
	}

	r, err := New(&Config{}, "config")
	if err != nil {
		t.Fatalf("fatal error - %s", err)
	}

	provider := providers.NewJsonDataProvider([]byte(`{"server":{"name":"test_server"}}`))
	if config, err := r.SetValues(provider); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v\n", config)
	}
}

func TestBoolValue(t *testing.T) {
	type Config struct {
		Active bool `config:"is_active"`
	}
	r, err := New(&Config{}, "config")
	if err != nil {
		t.Fatal(err)
	}
	provider := providers.NewJsonDataProvider([]byte(`{"is_active":false}`))
	if config, err := r.SetValues(provider); err != nil {
		t.Error(err)
	} else {
		if config.(*Config).Active {
			t.Error("invalid value for active")
		}
	}
}

func TestSliceValue(t *testing.T) {
	type Config struct {
		Params []string `config:"params"`
	}
	r, err := New(&Config{}, "config")
	if err != nil {
		t.Fatal(err)
	}
	provider := providers.NewJsonDataProvider([]byte(`{"params":["one","two","three"]}`))
	if config, err := r.SetValues(provider); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v\n", config)
	}
}



