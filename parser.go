package scim_filtering

import (
	"errors"
	"fmt"
	"io"
)

type Expression struct {
	Name     string
	Operator Token
	Value    string
}

type Statement struct {
	Operator   Token
	Expression *Expression
	Statements []*Statement
}

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

func (p *Parser) Parse() (*Statement, error) {
	token, literal := p.scanIgnoreWhitespace()
	switch token {
	case NOT:
		statement, err := p.Parse()
		if err != nil {
			return nil, err
		}
		return &Statement{
			Operator:   NOT,
			Expression: statement.Expression,
		}, nil
	case ID:
		break
	default:
		return nil, fmt.Errorf("found %q, expected identifier", literal)
	}

	expression := &Expression{Name: literal}

	token, literal = p.scanIgnoreWhitespace()
	switch token {
	case EQ, NE, CO, SW, EW, PR, GT, GE, LT, LE:
		break
	default:
		return nil, fmt.Errorf("found %q, expected operator", literal)
	}

	expression.Operator = token

	token, literal = p.scanIgnoreWhitespace()
	if token != V && literal != "" {
		return nil, fmt.Errorf("found %q, expected value", token)
	}

	expression.Value = literal

	token, literal = p.scanIgnoreWhitespace()
	if token == EOF {
		return &Statement{
			Expression: expression,
		}, nil
	}

	switch token {
	case AND, OR:
		statement, err := p.Parse()
		if err != nil {
			return nil, err
		}
		return &Statement{
			Operator: token,
			Statements: []*Statement{
				{
					Operator:   0,
					Expression: expression,
				},
				statement,
			},
		}, nil
	default:
		return nil, errors.New("not implemented")
	}
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

// scanIgnoreWhiteSpace scans the next token that is not whitespace.
func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	token, literal := p.scan()
	if token == WS {
		token, literal = p.scan()
	}
	return token, literal
}
