package parser

import (
	"strings"
	"testing"
)

func initScanner(s string) *Scanner {
	return NewScanner(strings.NewReader(s))
}

func TestNewScanner(t *testing.T) {
	_ = NewScanner(strings.NewReader("raw_string"))
}

func TestScanMethod(t *testing.T) {
	scanner := initScanner("")
	if token, _ := scanner.Scan(); token != eofToken {
		t.Error("there must be eof token")
	}

	scanner = initScanner(" ")
	if token, _ := scanner.Scan(); token != wsToken {
		t.Error("there must be white space")
	}

	scanner = initScanner("'string'")
	if token, _ := scanner.Scan(); token != stringValueToken {
		t.Error("this is string value")
	}

	// test slice
	scanner = initScanner("[1,2,3]")
	if token, _ := scanner.Scan(); token != sliceValueToken {
		t.Error("this is slice")
	}

	// test digits
	scanner = initScanner("1")
	if token, _ := scanner.Scan(); token != numberValueToken {
		t.Error("this is number")
	}

	// scan ident
	scanner = initScanner("value")
	if token, value := scanner.Scan(); token != identValueToken {
		t.Error("this is ident")
	} else {
		if value != "value"{
			t.Errorf("invalid value `%s`", value)
		}
	}
}

func TestWhiteSpace(t *testing.T) {
	scanner := initScanner("    ")
	if token, value := scanner.scanWhiteSpace(); token != wsToken {
		t.Error("this is whitespace")
	} else {

		if v, ok := value.(WhiteSpaceTokenValue); !ok {
			t.Error("token value must be WhiteSpaceTokenValue")
		} else {
			if v.number != 4 {
				t.Error("there must be 4 spaces")
			}
		}
	}
}

func TestString(t *testing.T) {
	scanner := initScanner("'string'")
	if token, value := scanner.scanString(); token != stringValueToken {
		t.Error("invalid token")
	} else {
		if v, ok := value.(string); !ok {
			t.Error("invalid token value type")
		} else if v != "string" {
			t.Error("invalid token value")
		}
	}

	// Test invalid
	scanner = initScanner("'string")
	if token, _ := scanner.scanString(); token != illegalToken {
		t.Error("there must be an illegal token")
	}

}

func TestNumbers(t *testing.T) {
	scanner := initScanner("10")
	if token, value := scanner.scanNumber(); token != numberValueToken {
		t.Error("this is number")
	} else {
		if value != int64(10) {
			t.Errorf("invalid value %T", value)
		}
	}
	// test invalid
	scanner = initScanner("10s")
	if token, value := scanner.scanNumber(); token != illegalToken {
		t.Errorf("token is illegal - %v", value)
	}
	scanner = initScanner("10.10.10")
	if token, _ :=scanner.scanNumber(); token != illegalToken {
		t.Error("value is invalid")
	}
	scanner = initScanner("10.s")
	if token, _ := scanner.scanNumber(); token != illegalToken {
		t.Error("value is invalid")
	}
	// test float
	scanner = initScanner("10.10")
	if token, value := scanner.scanNumber(); token != floatValueToken {
		t.Error("invalid token")
	} else {
		if value != float64(10.10) {
			t.Error("invalid value")
		}
	}

}

func TestIdents(t *testing.T) {
	tests := []struct {
		input string
		token Token
	}{
		{
			input: "is_required",
			token: isRequiredToken,
		},
		{
			input: "has_default",
			token: hasDefaultToken,
		},
		{
			input: "has_value",
			token: hasValueToken,
		},
		{
			input: "is_required_if",
			token: isRequiredIfToken,
		},
		{
			input: "some_value",
			token: identValueToken,
		},
		{
			input: "some_value is_required",
			token: identValueToken,
		},
		{
			input: "false",
			token: booleanValueToken,
		},
	}
	for _, test := range tests {
		scanner := initScanner(test.input)
		if token, _ := scanner.scanIdent(); token != test.token {
			t.Errorf("invalid token for input: %s", test.input)
		}
	}
}

func TestNumberSlice(t *testing.T) {
	scanner := initScanner("[1,2,3]")
	if token, value := scanner.scanSlice(); token != sliceValueToken {
		t.Errorf("illegal token %q", token)
	} else {
		if v, ok := value.([]int64); !ok {
			t.Errorf("invalid slice type %T", value)
		} else {
			expected := []int64{int64(1), int64(2), int64(3)}
			for i := 0; i < len(v); i++ {
				if v[i] != expected[i] {
					t.Error("invalid value")
				}
			}
		}
	}
}

func TestFloatSlice(t *testing.T) {
	scanner := initScanner("[1.0,2.0,3.1]")
	if token, value := scanner.scanSlice(); token != sliceValueToken {
		t.Errorf("invalid token: %q", token)
	} else {
		if v, ok := value.([]float64); !ok {
			t.Errorf("invalid slice type %T", value)
		} else {
			expected := []float64{1.0, 2.0, 3.1}
			for i := 0; i < len(expected); i++ {
				if expected[i] != v[i] {
					t.Errorf("expected: %f, actual: %f", expected[i], v[i])
				}
			}
		}
	}
}

func TestStringSlice(t *testing.T) {
	scanner := initScanner("['1', '2', '3']")
	if token, value := scanner.scanSlice(); token != sliceValueToken {
		t.Errorf("invalid token: %s", token)
	} else {
		if v, ok := value.([]string); !ok {
			t.Errorf("invalid slice type: %T", value)
		} else {
			ex := []string{"1", "2", "3"}
			for i := 0; i < len(ex); i++ {
				if ex[i] != v[i] {
					t.Errorf("expected: %s, real: %s", ex[i], v[i])
				}
			}
		}
	}

	scanner = initScanner("['some value, 'another value']")
	if token, _ := scanner.scanSlice(); token != illegalToken {
		t.Error("slice values are illegal")
	}

}

func TestInterfaceSlice(t *testing.T) {
	scanner := initScanner("['1', 2, 3.0]")
	if token, value := scanner.scanSlice(); token != sliceValueToken {
		t.Errorf("invalid token: %s", token)
	} else {
		if v, ok := value.([]interface{}); !ok {
			t.Error("invalid type for slice")
		} else {

			if (v[0]).(string) != "1" {
				t.Error("invalid value")
			}
		}

	}
}

func TestInvalidSlices(t *testing.T) {

	tests := []string{
		"[1,2,3",
		"['as, 'sdsd']",
		"[10,10s], ",
		"[[1,2], 'trest']",
		"[]",
	}
	for _, test :=range tests {
		scanner := initScanner(test)
		if token, _ := scanner.scanSlice(); token != illegalToken {
			t.Errorf("%s - slice is illegal", test)
		}
	}

}

func TestBooleanValue(t *testing.T) {
	scanner := initScanner("true")
	if token, value := scanner.scanIdent(); token != booleanValueToken {
		t.Errorf("it is boolean value")
	} else {
		if !value.(bool) {
			t.Error("value is true")
		}
	}
}