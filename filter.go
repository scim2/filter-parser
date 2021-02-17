package filter

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/scim2/filter-parser/v2/grammar"
	typ "github.com/scim2/filter-parser/v2/types"
)

// ParseFilter parses the given raw data as an Expression.
func ParseFilter(raw []byte) (Expression, error) {
	p, err := ast.New(raw)
	if err != nil {
		return nil, err
	}
	node, err := grammar.Filter(p)
	if err != nil {
		return nil, err
	}
	if _, err := p.Expect(parser.EOD); err != nil {
		return nil, err
	}
	return parseFilterOr(node)
}

func parseFilterOr(node *ast.Node) (Expression, error) {
	if node.Type != typ.FilterOr {
		return nil, invalidTypeError(typ.FilterOr, node.Type)
	}

	children := node.Children()
	if len(children) == 0 {
		return nil, invalidLengthError(typ.FilterOr, 1, 0)
	}

	if len(children) == 1 {
		return parseFilterAnd(children[0])
	}

	var or LogicalExpression
	for _, node := range children {
		exp, err := parseFilterAnd(node)
		if err != nil {
			return nil, err
		}
		switch {
		case or.Left == nil:
			or.Left = exp
		case or.Right == nil:
			or.Right = exp
			or.Operator = OR
		default:
			or = LogicalExpression{
				Left:     &or,
				Right:    exp,
				Operator: OR,
			}
		}
	}
	return &or, nil
}

func parseFilterAnd(node *ast.Node) (Expression, error) {
	if node.Type != typ.FilterAnd {
		return nil, invalidTypeError(typ.FilterAnd, node.Type)
	}

	children := node.Children()
	if len(children) == 0 {
		return nil, invalidLengthError(typ.FilterAnd, 1, 0)
	}

	if len(children) == 1 {
		return parseFilterValue(children[0])
	}

	var and LogicalExpression
	for _, node := range children {
		exp, err := parseFilterValue(node)
		if err != nil {
			return nil, err
		}
		switch {
		case and.Left == nil:
			and.Left = exp
		case and.Right == nil:
			and.Right = exp
			and.Operator = AND
		default:
			and = LogicalExpression{
				Left:     &and,
				Right:    exp,
				Operator: AND,
			}
		}
	}
	return &and, nil
}

func parseFilterValue(node *ast.Node) (Expression, error) {
	switch t := node.Type; t {
	case typ.ValuePath:
		valuePath, err := parseValuePath(node)
		if err != nil {
			return nil, err
		}
		return &valuePath, nil
	case typ.AttrExp:
		attrExp, err := parseAttrExp(node)
		if err != nil {
			return nil, err
		}
		return &attrExp, nil
	case typ.FilterNot:
		children := node.Children()
		if l := len(children); l != 1 {
			return nil, invalidLengthError(typ.FilterNot, 1, l)
		}

		exp, err := parseFilterOr(children[0])
		if err != nil {
			return nil, err
		}
		return &NotExpression{
			Expression: exp,
		}, nil
	case typ.FilterOr:
		return parseFilterOr(node)
	default:
		return nil, invalidChildTypeError(typ.FilterAnd, t)
	}
}
