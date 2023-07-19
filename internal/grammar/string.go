package grammar

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"

	typ "github.com/scim2/filter-parser/v2/internal/types"
)

type quote int

const (
	singleQuote quote = iota
	doubleQuote
)

func (q quote) getValue() string {
	switch q {
	case singleQuote:
		return "'"
	case doubleQuote:
		fallthrough
	default:
		return "\""
	}
}

func Character(q quote) func(*ast.Parser) (*ast.Node, error) {
	return func(p *ast.Parser) (*ast.Node, error) {
		return p.Expect(
			op.Or{
				Unescaped(q),
				op.And{
					"\\",
					q.getValue(),
				},
				op.And{
					"\\",
					Escaped,
				},
			},
		)
	}
}

func Escaped(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		op.Or{
			"\\",
			"/",
			0x0062,
			0x0066,
			0x006E,
			0x0072,
			0x0074,
			op.And{
				"u",
				op.Repeat(4,
					op.Or{
						parser.CheckRuneRange('0', '9'),
						parser.CheckRuneRange('A', 'F'),
					},
				),
			},
		},
	)
}

func String(q quote) func(*ast.Parser) (*ast.Node, error) {
	return func(p *ast.Parser) (*ast.Node, error) {
		return p.Expect(
			ast.Capture{
				Type:        typ.String,
				TypeStrings: typ.Stringer,
				Value: op.And{
					q.getValue(),
					op.MinZero(
						Character(q),
					),
					q.getValue(),
				},
			},
		)
	}
}

func Unescaped(q quote) func(*ast.Parser) (*ast.Node, error) {
	switch q {
	case singleQuote:
		return func(p *ast.Parser) (*ast.Node, error) {
			// 0x0027 : '
			return p.Expect(
				op.Or{
					parser.CheckRuneRange(0x0020, 0x0026),
					parser.CheckRuneRange(0x0028, 0x0010FFFF),
				},
			)
		}
	case doubleQuote:
		fallthrough
	default:
		return func(p *ast.Parser) (*ast.Node, error) {
			// 0x0023 : "
			// 0x005C : \
			return p.Expect(
				op.Or{
					parser.CheckRuneRange(0x0020, 0x0021),
					parser.CheckRuneRange(0x0023, 0x005B),
					parser.CheckRuneRange(0x005D, 0x0010FFFF),
				},
			)
		}
	}
}
