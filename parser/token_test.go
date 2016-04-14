package parser

import "testing"

func TestTokenString(t *testing.T) {
	var token = illegalToken
	if token.String() != "<illegal>" {
		t.Error("invalid string value for token")
	}
	// test unknown value

	token = 100
	if token.String() != "<?>" {
		t.Error("invalid value for unknown token")
	}
}

func TestWhitespace(t *testing.T) {
	if !isWhiteSpace(' ') {
		t.Error("is whitespace rune")
	}
}

func TestDot(t *testing.T) {
	if !isDot('.') {
		t.Error("is dot char rune")
	}
}

func TestDigits(t *testing.T) {
	var runes []rune
	for i := '0'; i <= '9'; i++ {
		runes = append(runes, i)
	}
	for _, ch := range runes {
		if !isDigit(ch) {
			t.Errorf("%c is digit", ch)
		}
	}
}

func TestLetters(t *testing.T) {
	var runes []rune
	for ch := 'a'; ch <= 'z'; ch++ {
		runes = append(runes, ch)
	}
	for ch := 'Z'; ch <= 'Z'; ch++ {
		runes = append(runes, ch)
	}
	for _, ch := range runes {
		if !isLetter(ch) {
			t.Errorf("%c is letter")
		}
	}
}

func TestApostrophe(t *testing.T) {
	if !isApostrophe('\'') {
		t.Error("this is apostrophe")
	}
}

func TestSliceStartEnd(t *testing.T) {
	if !isSliceStart('[') {
		t.Error("slice starts")
	}
	if !isSliceEnd(']') {
		t.Error("slice ends")
	}
}