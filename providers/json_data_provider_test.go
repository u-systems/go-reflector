package providers

import "testing"

func TestNew(t *testing.T) {
	raw := []byte(`{"field":"name"}`)
	provider := NewJsonDataProvider(raw)

	data := provider.Data()
	if v, ok := data.([]byte); !ok {
		t.Error("invalid type for data")
	} else {
		if string(v) != string(raw) {
			t.Error("invalid value")
		}
	}
}

func TestInvalidFormat(t *testing.T) {
	raw := []byte(``)
	provider := NewJsonDataProvider(raw)
	if _, err := provider.Load(); err == nil {
		t.Error("there must an error for invalid format")
	}
}

func TestLoad(t *testing.T) {
	raw := []byte(`{"host":"localhost", "port":8080}`)
	provider := NewJsonDataProvider(raw)

	if result, err := provider.Load(); err != nil {
		t.Errorf("there can not be an error: %s", err)
	} else {
		if v, found := result["host"]; !found {
			t.Error("there must be an field: host")
		} else {
			if v, ok := v.(string); !ok {
				t.Error("type for field is string")
			}  else {
				if v != "localhost" {
					t.Errorf("invalid value: %s", v)
				}
			}
		}
	}
}

func TestUnload(t *testing.T) {
	provider := NewJsonDataProvider([]byte(``))
	raw := map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	}
	if err := provider.Unload(raw); err != nil {
		t.Errorf("there can not be an error: %s", err)
	} else {
		expected := `{"host":"localhost","port":8080}`
		if string(provider.Data().([]byte)) != expected {
			t.Errorf("unexpected value: %v", provider.data)
		}
	}
	// Test unload with error?
}
