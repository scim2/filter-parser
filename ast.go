package filter

import "fmt"

// Expression is a type to assign to implemented expressions.
type Expression interface{}

// AttributeExpression is an Expression with a name, operator and value.
type AttributeExpression struct {
	Expression
	AttributePath   AttributePath
	CompareOperator Token
	CompareValue    string
}

type AttributePath struct {
	AttributeName string
	SubAttribute  string
}

type ValuePath struct {
	Expression
	AttributeName   string
	ValueExpression Expression
}

// UnaryExpression is an expression with a token bound to a (child) expression X.
type UnaryExpression struct {
	Expression
	CompareOperator Token
	X               Expression
}

// BinaryExpression is an expression with a token bound to two (child) expressions X and Y.
type BinaryExpression struct {
	Expression
	X               Expression
	CompareOperator Token
	Y               Expression
}

func (expression AttributeExpression) String() string {
	return fmt.Sprintf("'%s %s %s'", expression.AttributePath, expression.CompareOperator, expression.CompareValue)
}

func (attributePath AttributePath) String() string {
	if attributePath.SubAttribute != "" {
		return fmt.Sprintf("%s.%s", attributePath.AttributeName, attributePath.SubAttribute)
	}
	return attributePath.AttributeName
}

func (valuePath ValuePath) String() string {
	return fmt.Sprintf("%s[%s]", valuePath.AttributeName, valuePath.ValueExpression)
}

func (expression UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", expression.CompareOperator, expression.X)
}

func (expression BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", expression.X, expression.CompareOperator, expression.Y)
}
