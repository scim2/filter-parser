package filter

import (
	"fmt"
	"io"
	"strings"
)

// Parser is a parser.
type Parser struct {
	s      *Scanner
	prefix string
	buf    struct {
		token   Token  // last read token
		literal string // last read literal
		n       int    // buffer size (max = 1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// ParsePath parses the "path" attribute value in the scanner.
// RFC: https://tools.ietf.org/html/rfc7644#section-3.5.2
func (parser *Parser) ParsePath() (Path, error) {
	return parser.parsePath()
}

func (parser *Parser) parsePath() (Path, error) {
	token, literal := parser.scanIgnoreWhitespace()

	if token != IDENTIFIER {
		return Path{}, fmt.Errorf("found %q, expected identifier", literal)
	}

	var uriPrefix string
	uriPrefix, literal = splitURIPrefix(literal)

	if sub := strings.Split(literal, "."); len(sub) > 1 {
		if len(sub) > 2 {
			return Path{}, fmt.Errorf("found %q, no multiple sub attributes allowed", literal)
		}
		if sub[1] == "" {
			return Path{}, fmt.Errorf("found %q, sub attribute can not be empty", literal)
		}
		if parser.peek() != EOF {
			_, literal = parser.scanIgnoreWhitespace()
			return Path{}, fmt.Errorf("found %q, expected eof", literal)
		}
		return Path{
			URIPrefix:     uriPrefix,
			AttributeName: sub[0],
			SubAttribute:  sub[1],
		}, nil
	}

	name := literal

	token, literal = parser.scanIgnoreWhitespace()
	if token != LeftBracket {
		if token != EOF {
			return Path{}, fmt.Errorf("found %q, expected eof", literal)
		}
		return Path{
			URIPrefix:     uriPrefix,
			AttributeName: name,
		}, nil
	}

	expression, err := parser.parse(lowestPrecedence)
	if err != nil {
		return Path{}, err
	}

	token, literal = parser.scanIgnoreWhitespace()
	if token != RightBracket {
		return Path{}, fmt.Errorf("found %q, expected '['", literal)
	}

	token, literal = parser.scanIgnoreWhitespace()
	if token == EOF {
		return Path{
			URIPrefix:       uriPrefix,
			AttributeName:   name,
			ValueExpression: expression,
		}, nil
	}
	if token != Dot {
		return Path{}, fmt.Errorf("found %q, expected '.' or eof", literal)
	}

	token, literal = parser.scanIgnoreWhitespace()
	if token != IDENTIFIER {
		return Path{}, fmt.Errorf("found %q, expected identifier", literal)
	}
	sub := literal

	token, literal = parser.scanIgnoreWhitespace()
	if token != EOF {
		return Path{}, fmt.Errorf("found %q, expected eof", literal)
	}

	return Path{
		URIPrefix:       uriPrefix,
		AttributeName:   name,
		ValueExpression: expression,
		SubAttribute:    sub,
	}, nil
}

// Parse returns an abstract syntax tree of the string in the scanner.
// https://tools.ietf.org/html/rfc7644#section-3.4.2.2
func (parser *Parser) Parse() (Expression, error) {
	return parser.parse(lowestPrecedence)
}

func (parser *Parser) parse(precedence int) (Expression, error) {
	var filter interface{}
	token, literal := parser.scanIgnoreWhitespace()

	if parser.peek() == LeftBracket {
		parser.prefix = literal
		token, literal = parser.scanIgnoreWhitespace()
	}

	switch token {
	case UNKNOWN:
		return nil, fmt.Errorf("unknown token: %q", literal)
	case LeftParenthesis, LeftBracket:
		expression, err := parser.parse(lowestPrecedence)
		if err != nil {
			return nil, err
		}
		parenthesis, parenthesisLiteral := parser.scanIgnoreWhitespace()
		if parenthesis != RightParenthesis && parenthesis != RightBracket {
			return nil, fmt.Errorf("found %q, expected '[' or '('", parenthesisLiteral)
		}

		if parenthesis == RightBracket {
			var uriPrefix string
			uriPrefix, parser.prefix = splitURIPrefix(parser.prefix)

			filter = ValuePath{
				URIPrefix:       uriPrefix,
				AttributeName:   parser.prefix,
				ValueExpression: expression,
			}
		} else {
			filter = expression
		}
	case IDENTIFIER:
		var err error
		filter, err = parser.parseAttributeExpression(token, literal)
		if err != nil {
			return nil, err
		}
	case NOT:
		expression, err := parser.parse(highestPrecedence)
		if err != nil {
			return nil, err
		}
		filter = UnaryExpression{
			X:               expression,
			CompareOperator: NOT,
		}
	}

	for precedence < parser.peek().Precedence() {
		token, _ := parser.scanIgnoreWhitespace()
		if token.IsAssociative() {
			expression, err := parser.parse(token.Precedence())
			if err != nil {
				return nil, err
			}
			filter = BinaryExpression{
				X:               filter,
				CompareOperator: token,
				Y:               expression,
			}
		}
	}

	return filter, nil
}

// parseAttributeExpression returns a value parse with the remaining operator and value of preceding identifier.
func (parser *Parser) parseAttributeExpression(token Token, literal string) (AttributeExpression, error) {
	operator, operatorLiteral := parser.scanIgnoreWhitespace()
	if !operator.IsOperator() {
		return AttributeExpression{}, fmt.Errorf("found %q, expected operator", operatorLiteral)
	}

	var valueLiteral string
	if operator != PR {
		var value Token
		value, valueLiteral = parser.scanIgnoreWhitespace()
		if value != VALUE && valueLiteral != "" {
			return AttributeExpression{}, fmt.Errorf("found %q, expected value", token)
		}
	}

	var uriPrefix string
	uriPrefix, literal = splitURIPrefix(literal)

	if sub := strings.Split(literal, "."); len(sub) > 1 {
		if len(sub) > 2 {
			return AttributeExpression{}, fmt.Errorf("found %q, no multiple sub attributes allowed", literal)
		}
		return AttributeExpression{
			AttributePath: AttributePath{
				URIPrefix:     uriPrefix,
				AttributeName: sub[0],
				SubAttribute:  sub[1],
			},
			CompareOperator: operator,
			CompareValue:    valueLiteral,
		}, nil
	}

	return AttributeExpression{
		AttributePath: AttributePath{
			URIPrefix:     uriPrefix,
			AttributeName: literal,
		},
		CompareOperator: operator,
		CompareValue:    valueLiteral,
	}, nil
}

// scan returns the next token in the scanner.
func (parser *Parser) scan() (Token, string) {
	if parser.buf.n != 0 {
		parser.buf.n = 0
		return parser.buf.token, parser.buf.literal
	}

	token, literal := parser.s.Scan()
	parser.buf.token, parser.buf.literal = token, literal

	return token, literal
}

// unscan places the last read token back in the buffer.
func (parser *Parser) unscan() {
	parser.buf.n = 1
}

// peek returns the next token in the scanner that is not whitespace.
func (parser *Parser) peek() Token {
	token, _ := parser.scan()
	if token == WHITESPACE {
		token, _ = parser.scan()
		parser.unscan()
	}
	parser.unscan()
	return token
}

// scanIgnoreWhiteSpace scans the next token that is not whitespace.
func (parser *Parser) scanIgnoreWhitespace() (Token, string) {
	token, literal := parser.scan()
	if token == WHITESPACE {
		token, literal = parser.scan()
	}
	return token, literal
}
