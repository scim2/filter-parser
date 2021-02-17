package grammar

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	typ "github.com/scim2/filter-parser/v2/types"
)

func ValuePath(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.ValuePath,
		Value: op.And{
			AttrPath,
			op.MinZero(SP),
			'[',
			op.MinZero(SP),
			ValueFilterAll,
			op.MinZero(SP),
			']',
		},
	})
}

func ValueFilterAll(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		ValueFilter,
		ValueFilterNot,
	})
}

func ValueFilter(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		ValueLogExpOr,
		ValueLogExpAnd,
		AttrExp,
	})
}

func ValueLogExpOr(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.ValueLogExpOr,
		Value: op.And{
			AttrExp,
			op.MinZero(SP),
			"or",
			op.MinZero(SP),
			AttrExp,
		},
	})
}

func ValueLogExpAnd(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.ValueLogExpAnd,
		Value: op.And{
			AttrExp,
			op.MinZero(SP),
			"and",
			op.MinZero(SP),
			AttrExp,
		},
	})
}

func ValueFilterNot(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.ValueFilterNot,
		Value: op.And{
			"not",
			op.MinZero(SP),
			'(',
			op.MinZero(SP),
			ValueFilter,
			op.MinZero(SP),
			')',
		},
	})
}
