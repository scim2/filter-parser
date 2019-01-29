package scim_filter_parser

import "fmt"

type Expression interface{}

type ValueExpression struct {
	Expression
	Name     string
	Operator Token
	Value    string
}

type UnaryExpression struct {
	Expression
	Operator Token
	X        Expression
}

type BinaryExpression struct {
	Expression
	X        Expression
	Operator Token
	Y        Expression
}

func (expression ValueExpression) String() string {
	return fmt.Sprintf("'%s %s %s'", expression.Name, expression.Operator, expression.Value)
}

func (expression UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", expression.Operator, expression.X)
}

func (expression BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", expression.X, expression.Operator, expression.Y)
}
