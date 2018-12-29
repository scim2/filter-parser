package scim_filtering

import (
	"fmt"
	"io"
)

// Parser is a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		token   Token  // last read token
		literal string // last read literal
		n       int    // buffer size (max = 1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (Expression, error) {
	return p.expression(LowestPrecedence)
}

func (p *Parser) expression(precedence int) (Expression, error) {
	var left interface{}
	token, literal := p.scanIgnoreWhitespace()
	switch token {
	case UNKNOWN:
		return nil, fmt.Errorf("unknown token: %q", literal)
	case ID:
		operator, operatorLiteral := p.scanIgnoreWhitespace()
		if !operator.IsOperator() {
			return nil, fmt.Errorf("found %q, expected operator", operatorLiteral)
		}

		value, valueLiteral := p.scanIgnoreWhitespace()
		if value != V && valueLiteral != "" {
			return nil, fmt.Errorf("found %q, expected value", token)
		}

		left = ValueExpression{
			Name:     literal,
			Operator: operator,
			Value:    valueLiteral,
		}
	case NOT:
		expression, err := p.expression(HighestPrecedence)
		if err != nil {
			return nil, err
		}
		left = UnaryExpression{
			X:        expression,
			Operator: NOT,
		}
	}

	for precedence < p.peek().Precedence() {
		token, _ := p.scanIgnoreWhitespace()
		if token.IsAssociative() {
			expression, err := p.expression(token.Precedence())
			if err != nil {
				return nil, err
			}
			left = BinaryExpression{
				X:        left,
				Operator: token,
				Y:        expression,
			}
		}
	}

	return left, nil
}

// scan returns the next token in the scanner.
func (p *Parser) scan() (Token, string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.token, p.buf.literal
	}

	token, literal := p.s.Scan()
	p.buf.token, p.buf.literal = token, literal

	return token, literal
}

// unscan places the last read token back in the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) peek() Token {
	token, _ := p.scan()
	if token == WS {
		token, _ = p.scan()
		p.unscan()
	}
	p.unscan()
	return token
}

// scanIgnoreWhiteSpace scans the next token that is not whitespace.
func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	token, literal := p.scan()
	if token == WS {
		token, literal = p.scan()
	}
	return token, literal
}
