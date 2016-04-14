package parser

import (
	"bufio"
	"io"
	"bytes"
	"strings"
	"strconv"
	"reflect"
)

const eof = rune(0)

type Scanner struct {
	reader *bufio.Reader
}

// Create new scanner
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(reader),
	}
}

func (scanner *Scanner) Scan() (Token, TokenValue) {
	ch := scanner.read()
	if ch == eof {
		return eofToken, ""
	}

	if isWhiteSpace(ch) {
		scanner.unread()
		return scanner.scanWhiteSpace()
	} else if isApostrophe(ch) {
		scanner.unread()
		return scanner.scanString()
	} else if isSliceStart(ch) {
		scanner.unread()
		return scanner.scanSlice()
	} else if isDigit(ch) {
		scanner.unread()
		return scanner.scanNumber()
	}

	scanner.unread()
	return scanner.scanIdent()
}

func (scanner *Scanner) read() rune {
	ch, _, err := scanner.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (scanner *Scanner) unread() {
	_ = scanner.reader.UnreadRune()
}

// Scan white space
func (scanner *Scanner) scanWhiteSpace() (Token, TokenValue) {
	tokenValue := WhiteSpaceTokenValue{}
	for {
		if ch := scanner.read(); ch == eof {
			break
		} else if !isWhiteSpace(ch) {
			scanner.unread()
			break
		} else {
			tokenValue.number++
		}
	}
	return wsToken, tokenValue
}

// Scan string
func (scanner *Scanner) scanString() (Token, TokenValue) {
	var buffer bytes.Buffer
	_ = scanner.read() // skip first rune
	for {
		if ch := scanner.read(); ch == eof {
			return illegalToken, nil
		} else if isApostrophe(ch) {
			break
		} else {
			buffer.WriteRune(ch)
		}
	}

	return stringValueToken, buffer.String()
}

func (scanner *Scanner) scanIdent() (Token, TokenValue) {
	var buffer bytes.Buffer
	buffer.WriteRune(scanner.read())

	for {
		if ch := scanner.read(); ch == eof {
			break
		} else if !isLetter(ch) && ch != '_' {
			scanner.unread()
			break
		} else {
			buffer.WriteRune(ch)
		}
	}

	switch v := strings.ToLower(buffer.String()); v {
	case "is_required":
		return isRequiredToken, nil
	case "has_value":
		return hasValueToken, nil
	case "has_default":
		return hasDefaultToken, nil
	case "is_required_if":
		return isRequiredIfToken, nil
	case "true":
		return booleanValueToken, true
	case "false":
		return booleanValueToken, false
	}

	return identValueToken, buffer.String()
}

// Scan number
func (scanner *Scanner) scanNumber() (Token, TokenValue) {
	var buffer bytes.Buffer
	var isFloat bool
	buffer.WriteRune(scanner.read())

	for {
		if ch := scanner.read(); ch == eof {
			break
		} else if !isDigit(ch) && !isDot(ch) && !isLetter(ch)  {
			scanner.unread()
			break
		} else {
			buffer.WriteRune(ch)
			if isDot(ch) {
				if isFloat {
					return illegalToken, nil
				}
				isFloat = true
			}
		}
	}

	if isFloat {
		v, err := strconv.ParseFloat(buffer.String(), 64)
		if err != nil {
			return illegalToken, nil
		}
		return floatValueToken, float64(v)
	} else {
		if v, err := strconv.ParseInt(buffer.String(), 10, 64); err != nil {
			return illegalToken, nil
		} else {
			return numberValueToken, int64(v)
		}
	}

}

// Scan slice
func (scanner *Scanner) scanSlice() (Token, TokenValue) {

	_ = scanner.read() // read first brace
	slice := []interface{}{}

	for {
		ch := scanner.read()
		if ch == eof || isSliceStart(ch) {
			return illegalToken, nil

		} else if isWhiteSpace(ch) {
			scanner.unread()
			scanner.scanWhiteSpace()
		} else if isApostrophe(ch) {
			scanner.unread()
			if token, value := scanner.scanString(); token == illegalToken {
				return illegalToken, nil
			} else {
				slice = append(slice, value)
			}
		} else if ch == ',' {
			//
		} else if isDigit(ch) {
			scanner.unread()
			if token, value := scanner.scanNumber(); token == illegalToken {
				return illegalToken, nil
			} else {
				slice = append(slice, value)
			}

		} else if isSliceEnd(ch) {
			break
		}

	}

	var kind reflect.Kind = reflect.Invalid

	kindMap := map[reflect.Kind]reflect.Type{
		reflect.String: reflect.TypeOf(string("")),
		reflect.Int64: reflect.TypeOf(int64(0)),
		reflect.Float64: reflect.TypeOf(float64(0)),
	}

	for _, v := range slice {
		curKind := reflect.TypeOf(v).Kind()
		if kind == reflect.Invalid {
			kind = curKind
		} else {
			if kind != curKind {
				return sliceValueToken, slice
			}
		}
	}

	if kind == reflect.Invalid {
		return illegalToken, nil // empty slice
	}

	newSlice := reflect.MakeSlice(reflect.SliceOf(kindMap[kind]), 0, 0)
	for i := 0; i < len(slice); i++ {
		newSlice = reflect.Append(newSlice, reflect.ValueOf(slice[i]))
	}

	if kind == reflect.String {
		return sliceValueToken, newSlice.Interface().([]string)
	} else if kind == reflect.Int64 {
		return sliceValueToken, newSlice.Interface().([]int64)
	} else {
		return sliceValueToken, newSlice.Interface().([]float64)
	}

}