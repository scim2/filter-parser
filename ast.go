package scim_filtering

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

func (e ValueExpression) String() string {
	return fmt.Sprintf("'%s %s %s'", e.Name, e.Operator, e.Value)
}

func (e UnaryExpression) String() string {
	return fmt.Sprintf("%s %s", e.Operator, e.X)
}

func (e BinaryExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.X, e.Operator, e.Y)
}
