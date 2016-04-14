package parser

type Token int

const (
	illegalToken Token = iota
	wsToken // whitespace
	eofToken //

	identValueToken // ident value
	stringValueToken // string value
	numberValueToken // number value
	floatValueToken // float value
	sliceValueToken // slice
	booleanValueToken // boolean value true|false

	isRequiredToken // is_required
	isRequiredIfToken // is_required
	hasDefaultToken // has_default
	hasValueToken // has_value

)

var names = map[Token]string{
	illegalToken: "illegal",
	wsToken: "white space",
	eofToken: "eof",
	identValueToken: "ident",
	stringValueToken: "string",
	numberValueToken: "number",
	floatValueToken: "float",
	sliceValueToken: "slice",
	isRequiredToken: "is_required",
	isRequiredIfToken: "is_required_if",
	hasDefaultToken: "has_default ...",
	hasValueToken: "has_value ...",
}

func (token Token) String() string {
	if name, ok := names[token]; !ok {
		return "<?>"
	} else {
		return "<" + name + ">"
	}
}

func isWhiteSpace(ch rune) bool {
	return ch == ' '
}

func isDot(ch rune) bool {
	return ch == '.'
}

func isApostrophe(ch rune) bool {
	return ch == '\''
}

func isSliceStart(ch rune) bool {
	return ch == '['
}

func isSliceEnd(ch rune) bool {
	return ch == ']'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}