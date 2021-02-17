package grammar

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	"github.com/scim2/filter-parser/v2/types"
)

func Filter(p *ast.Parser) (*ast.Node, error) {
	return FilterOr(p)
}

func FilterOr(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.FilterOr,
		Value: op.And{
			FilterAnd,
			op.MinZero(op.And{
				op.MinOne(SP),
				"or",
				op.MinOne(SP),
				FilterAnd,
			}),
		},
	})
}

func FilterAnd(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.FilterAnd,
		Value: op.And{
			FilterValue,
			op.MinZero(op.And{
				op.MinOne(SP),
				"and",
				op.MinOne(SP),
				FilterValue,
			}),
		},
	})
}

func FilterNot(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.FilterNot,
		Value: op.And{
			"not",
			op.MinZero(SP),
			FilterParentheses,
		},
	})
}

func FilterValue(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		ValuePath,
		AttrExp,
		FilterNot,
		FilterParentheses,
	})
}

func FilterParentheses(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{
		'(',
		op.MinZero(SP),
		FilterOr,
		op.MinZero(SP),
		')',
	})
}
