package scim_filtering

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
func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdentifiers()
	} else if isDoubleQuote(ch) {
		s.unread()
		return s.scanValue()
	}

	switch ch {
	case eof:
		return EOF, ""
	case '(':
		return LPAR, string(ch)
	case ')':
		return RPAR, string(ch)
	case '[':
		return LBRA, string(ch)
	case ']':
		return RBRA, string(ch)
	}

	return UNKNOWN, string(ch)
}

// read reads the next character from the reader.
// Returns the rune(0) if an error occurs.
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places last rune back in the reader.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// scanWhitespace removes current rune and all the whitespace after it.
func (s *Scanner) scanWhitespace() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdentifiers removes current rune and all the identifier runes after it.
func (s *Scanner) scanIdentifiers() (Token, string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == ':' {
			// removes all runes before last colon
			buf = bytes.Buffer{}
		} else if !isLetter(ch) && !isDigit(ch) && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// buf to lower case
	lower := strings.ToLower(buf.String())

	switch lower {
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

	return ID, lower
}

// scanValue removes current rune and all the value runes after it.
func (s *Scanner) scanValue() (Token, string) {
	var buf bytes.Buffer
	_ = s.read()
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isDoubleQuote(ch) {
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	return V, buf.String()
}

// isWhitespace checks whether the given rune is a whitespace.
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\n' || ch == '\t'
}

// isLetter checks whether the given rune is a letter (a-zA-Z).
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDigit checks whether the given rune is a digit (0-9)
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// isDoubleQuote checks whether the given rune is a double quote (\").
func isDoubleQuote(ch rune) bool {
	return ch == '"'
}

var eof = rune(0)
