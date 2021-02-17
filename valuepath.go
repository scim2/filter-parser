package filter

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/scim2/filter-parser/v2/grammar"
	typ "github.com/scim2/filter-parser/v2/types"
)

// ParseValuePath parses the given raw data as an ValuePath.
func ParseValuePath(raw []byte) (ValuePath, error) {
	p, err := ast.New(raw)
	if err != nil {
		return ValuePath{}, err
	}
	node, err := grammar.ValuePath(p)
	if err != nil {
		return ValuePath{}, err
	}
	if _, err := p.Expect(parser.EOD); err != nil {
		return ValuePath{}, err
	}
	return parseValuePath(node)
}

func parseValuePath(node *ast.Node) (ValuePath, error) {
	if node.Type != typ.ValuePath {
		return ValuePath{}, invalidTypeError(typ.ValuePath, node.Type)
	}

	children := node.Children()
	if l := len(children); l != 2 {
		return ValuePath{}, invalidLengthError(typ.ValuePath, 2, l)
	}

	attrPath, err := parseAttrPath(children[0])
	if err != nil {
		return ValuePath{}, err
	}

	valueFilter, err := parseValueFilter(children[1])
	if err != nil {
		return ValuePath{}, err
	}

	return ValuePath{
		AttributePath: attrPath,
		ValueFilter:   valueFilter,
	}, nil
}

func parseValueFilter(node *ast.Node) (Expression, error) {
	switch t := node.Type; t {
	case typ.ValueLogExpOr, typ.ValueLogExpAnd:
		children := node.Children()
		if l := len(children); l != 2 {
			return nil, invalidLengthError(node.Type, 2, l)
		}

		left, err := parseAttrExp(children[0])
		if err != nil {
			return nil, err
		}
		right, err := parseAttrExp(children[1])
		if err != nil {
			return nil, err
		}

		var operator LogicalOperator
		if node.Type == typ.ValueLogExpOr {
			operator = OR
		} else {
			operator = AND
		}

		return &LogicalExpression{
			Left:     &left,
			Right:    &right,
			Operator: operator,
		}, nil
	case typ.AttrExp:
		attrExp, err := parseAttrExp(node)
		if err != nil {
			return nil, err
		}
		return &attrExp, nil
	case typ.ValueFilterNot:
		children := node.Children()
		if l := len(children); l != 1 {
			return nil, invalidLengthError(typ.ValueFilterNot, 1, l)
		}

		valueFilter, err := parseValueFilter(children[0])
		if err != nil {
			return nil, err
		}
		return &NotExpression{
			Expression: valueFilter,
		}, nil
	default:
		return nil, invalidChildTypeError(typ.ValuePath, t)
	}
}
