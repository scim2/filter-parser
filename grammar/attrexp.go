package grammar

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	typ "github.com/scim2/filter-parser/v2/types"
)

func AttrExp(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.AttrExp,
		Value: op.And{
			AttrPath,
			op.MinOne(SP),
			op.Or{
				"pr",
				op.And{
					CompareOp,
					op.MinOne(SP),
					CompareValue,
				},
			},
		},
	})
}

func AttrPath(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.AttrPath,
		Value: op.And{
			op.Optional(URI),
			AttrName,
			op.Optional(SubAttr),
		},
	})
}

func AttrName(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: typ.AttrName,
		Value: op.And{
			op.Optional('$'),
			Alpha,
			op.MinZero(NameChar),
		},
	})
}

func NameChar(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{'-', '_', Digit, Alpha})
}

func SubAttr(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{'.', AttrName})
}

func CompareOp(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:  typ.CompareOp,
		Value: op.Or{"eq", "ne", "co", "sw", "ew", "gt", "lt", "ge", "le"},
	})
}

func CompareValue(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{False, Null, True, Number, String})
}
