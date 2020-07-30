package filter

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// Scanner is a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (scanner *Scanner) Scan() (Token, string) {
	ch := scanner.read()

	if isWhitespace(ch) {
		scanner.unread()
		return scanner.scanWhitespace()
	} else if isLetter(ch) {
		scanner.unread()
		return scanner.scanIdentifiers()
	} else if isDoubleQuote(ch) {
		scanner.unread()
		return scanner.scanValue()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '(':
		return LeftParenthesis, string(ch)
	case ')':
		return RightParenthesis, string(ch)
	case '[':
		return LeftBracket, string(ch)
	case ']':
		return RightBracket, string(ch)
	case '.':
		return Dot, string(ch)

	}

	return UNKNOWN, string(ch)
}

// read reads the next character from the reader.
// Returns the rune(0) if an error occurs.
func (scanner *Scanner) read() rune {
	ch, _, err := scanner.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places last rune back in the reader.
func (scanner *Scanner) unread() {
	_ = scanner.r.UnreadRune()
}

// scanWhitespace removes current rune and all the whitespace after it.
func (scanner *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(scanner.read())

	for {
		if ch := scanner.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			scanner.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WHITESPACE, buf.String()
}

// scanIdentifiers removes current rune and all the identifier runes after it.
func (scanner *Scanner) scanIdentifiers() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(scanner.read())

	for {
		if ch := scanner.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '.' && ch != ':' {
			scanner.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return checkIdentifier(buf.String())
}

// checkIdentifier checks whether given literal is a token or identifier.
func checkIdentifier(literal string) (Token, string) {
	switch lower := strings.ToLower(literal); lower {
	case "eq":
		return EQ, lower
	case "ne":
		return NE, lower
	case "co":
		return CO, lower
	case "sw":
		return SW, lower
	case "ew":
		return EW, lower
	case "pr":
		return PR, lower
	case "gt":
		return GT, lower
	case "ge":
		return GE, lower
	case "lt":
		return LT, lower
	case "le":
		return LE, lower

	case "and":
		return AND, lower
	case "or":
		return OR, lower
	case "not":
		return NOT, lower
	}

	return IDENTIFIER, literal
}

// scanValue removes current rune and all the value runes after it.
func (scanner *Scanner) scanValue() (Token, string) {
	var buf bytes.Buffer
	_ = scanner.read()
	buf.WriteRune(scanner.read())

	for {
		if ch := scanner.read(); ch == eof {
			break
		} else if isDoubleQuote(ch) {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return VALUE, buf.String()
}

// isWhitespace checks whether the given rune is a whitespace.
func isWhitespace(char rune) bool {
	return char == ' ' || char == '\n' || char == '\t'
}

// isLetter checks whether the given rune is a letter (a-zA-Z).
func isLetter(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

// isDigit checks whether the given rune is a digit (0-9)
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

// isDoubleQuote checks whether the given rune is a double quote (\").
func isDoubleQuote(char rune) bool {
	return char == '"'
}

var eof = rune(0)
