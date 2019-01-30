package scim

import "fmt"

// Expression is a type to assign to implemented expressions.
type Expression interface{}

// ValueExpression is an Expression with a name, operator and value.
type ValueExpression struct {
	Expression
	Name     string
	Operator Token
	Value    string
}

// UnaryExpression is an expression with a token bound to a (child) expression X.
type UnaryExpression struct {
	Expression
	Operator Token
	X        Expression
}

// BinaryExpression is an expression with a token bound to two (child) expressions X and Y.
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
