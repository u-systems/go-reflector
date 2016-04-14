package parser

import (
	"strings"
	"errors"
)

//
// Parser
//
type Parser struct {
	rawString string
	scanner   *Scanner
	buffer    struct {
			  token Token
			  value TokenValue
			  n     int
		  }
}

//
// Configuration field
//
type ConfigField struct {
	Name         string
	IsRequired   bool
	DefaultValue TokenValue
	DependsOn    struct {
			     ConfigFieldName string
			     Value           TokenValue
		     }
}

//
// New parser
//
func NewParser(s string) *Parser {
	return &Parser{
		rawString: s,
		scanner: NewScanner(strings.NewReader(s)),
	}
}

//
// Parse
//
func (parser *Parser) Parse() (*ConfigField, error) {
	configField := new(ConfigField)
	configField.DefaultValue = nil
	configField.DependsOn.Value = nil

	token, value := parser.scanIgnoreWhitespaces()

	if token == eofToken || token == illegalToken {
		return nil, errors.New("invalid config value")
	}

	if token == identValueToken {
		configField.Name = value.(string)
	} else {
		parser.unscan()
	}

	for {
		token, value := parser.scanIgnoreWhitespaces()
		if token == eofToken {
			break
		}

		// processing is_required
		if token == isRequiredToken {
			configField.IsRequired = true
		}

		// processing has_default
		if token == hasDefaultToken {
			// scan for value
			token, value = parser.scanIgnoreWhitespaces()
			if token != numberValueToken && token != floatValueToken &&
				token != stringValueToken && token != sliceValueToken && token != booleanValueToken {
				return nil, errors.New("has_default must has value")
			}
			configField.DefaultValue =  value
		}

		// processing if
		if token == isRequiredIfToken {
			// scan for ident
			token, value = parser.scanIgnoreWhitespaces()
			if token != identValueToken {
				return nil, errors.New("if required name of config field")
			} else {
				configField.DependsOn.ConfigFieldName = value.(string)
				token, value = parser.scanIgnoreWhitespaces()
				if token == hasValueToken {
					// scan for value
					token, value = parser.scanIgnoreWhitespaces()
					if token != stringValueToken && token != numberValueToken &&
						token != floatValueToken && token != booleanValueToken {
						return nil, errors.New("has_value needs value: string, int, float")
					} else {
						configField.DependsOn.Value = value
					}
				}
			}
		}
	}

	return configField, nil
}

// Scan ignore white spaces
func (parser *Parser) scanIgnoreWhitespaces() (token Token, value TokenValue) {
	token, value = parser.scan()
	if token == wsToken {
		token, value = parser.scan()
	}
	return
}

func (parser *Parser) scan() (token Token, value TokenValue) {
	if parser.buffer.n != 0 {
		parser.buffer.n = 0
		return parser.buffer.token, parser.buffer.value
	}
	token, value = parser.scanner.Scan()
	parser.buffer.token, parser.buffer.value = token, value
	return
}

func (parser *Parser) unscan() {
	parser.buffer.n = 1
}

