package parser_test

import (
	"testing"
	"stash.abc.ee/micro/reflector/parser"
)

func TestNew(t *testing.T) {
	_ = parser.NewParser("")
}

func TestEmptyConfig(t *testing.T) {
	parser := parser.NewParser("")
	if _, err := parser.Parse(); err == nil {
		t.Error("there must be an error for empty value")
	}
}

func TestConfigWithoutName(t *testing.T) {
	parser := parser.NewParser("is_required")
	if _, err := parser.Parse(); err != nil {
		t.Errorf("there no error %s", err)
	}
}

func TestHasDefault(t *testing.T) {
	parser := parser.NewParser("name has_default 'default value'")
	if configField, err := parser.Parse(); err != nil {
		t.Errorf("there can not be an error: %s", err)
	} else {
		if configField.Name != "name" {
			t.Errorf("invalid value for name: %s", configField.Name)
		}
		if v, ok := configField.DefaultValue.(string); !ok {
			t.Error("default value is string")
		} else {
			if v != "default value" {
				t.Errorf("invalid default value: `%s`", v)
			}
		}
	}
}

func TestHasDefaultError(t *testing.T) {
	parser := parser.NewParser("name has_default")
	if _, err := parser.Parse(); err == nil {
		t.Errorf("there must be an error")
	}
}

func TestIfStatement(t *testing.T) {
	parser := parser.NewParser("name is_required_if field has_value true")
	if configField, err := parser.Parse(); err != nil {
		t.Errorf("there can not be an error: %s", err)
	} else {
		if configField.DependsOn.ConfigFieldName != "field" {
			t.Error("field depens on `field`")
		}
		if v, ok := configField.DependsOn.Value.(bool); !ok {
			t.Error("value type - bool")
		} else {
			if !v {
				t.Errorf("value is true")
			}
		}
	}
}

func TestIfWithoutValue(t *testing.T) {
	parser := parser.NewParser("is_required_if field has_value")
	if _, err := parser.Parse(); err == nil {
		t.Error("there must an error")
	}
}

func TestIfStatementError(t *testing.T) {
	parser := parser.NewParser("name is_required_if")
	if _, err := parser.Parse(); err == nil {
		t.Error("there must be ab error")
	}
}


