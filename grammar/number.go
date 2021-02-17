package grammar

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	"github.com/scim2/filter-parser/v2/types"
)

func Number(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type: typ.Number,
			Value: op.And{
				op.Optional(
					Minus,
				),
				Int,
				op.Optional(
					Frac,
				),
				op.Optional(
					Exp,
				),
			},
		},
	)
}

func Minus(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type:  typ.Minus,
			Value: "-",
		},
	)
}

func Exp(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type: typ.Exp,
			Value: op.And{
				op.Or{
					"e",
					"E",
				},
				op.Optional(
					Sign,
				),
				Digits,
			},
		},
	)
}

func Sign(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type: typ.Sign,
			Value: op.Or{
				"-",
				"+",
			},
		},
	)
}

func Digits(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type: typ.Digits,
			Value: op.MinOne(
				parser.CheckRuneRange('0', '9'),
			),
		},
	)
}

func Frac(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type: typ.Frac,
			Value: op.And{
				".",
				Digits,
			},
		},
	)
}

func Int(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(
		ast.Capture{
			Type: typ.Int,
			Value: op.Or{
				"0",
				op.And{
					parser.CheckRuneRange('1', '9'),
					op.MinZero(
						parser.CheckRuneRange('0', '9'),
					),
				},
			},
		},
	)
}
