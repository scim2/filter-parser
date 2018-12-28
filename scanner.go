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
		} else if !isLetter(ch) && ch != '.' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	switch strings.ToLower(buf.String()) {
	case "eq":
		return EQ, strings.ToLower(buf.String())
	}

	return ID, strings.ToLower(buf.String())
}

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

// isDoubleQuote checks whether the given rune is a double quote (\").
func isDoubleQuote(ch rune) bool {
	return ch == '"'
}

var eof = rune(0)
