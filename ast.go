package filter

import (
	"fmt"
)

type CompareOperator string
type LogicalOperator string

const (
	PR CompareOperator = "pr"
	EQ CompareOperator = "eq"
	NE CompareOperator = "ne"
	CO CompareOperator = "co"
	SW CompareOperator = "sw"
	EW CompareOperator = "ew"
	GT CompareOperator = "gt"
	LT CompareOperator = "lt"
	GE CompareOperator = "ge"
	LE CompareOperator = "le"

	AND LogicalOperator = "and"
	OR  LogicalOperator = "or"
)

// Expression is a type to assign to implemented expressions.
type Expression interface {
	exprNode()
}

func (*ValuePath) exprNode()           {}
func (*AttributeExpression) exprNode() {}
func (*LogicalExpression) exprNode()   {}
func (*NotExpression) exprNode()       {}

// AttributeExpression is an Expression with a name, operator and value.
type AttributeExpression struct {
	AttributePath AttributePath
	Operator      CompareOperator
	CompareValue  interface{}
}

func (e AttributeExpression) String() string {
	s := fmt.Sprintf("%v %s", e.AttributePath, e.Operator)
	if e.CompareValue != nil {
		switch e.CompareValue.(type) {
		case string:
			s += fmt.Sprintf(" %q", e.CompareValue)
		default:
			s += fmt.Sprintf(" %v", e.CompareValue)
		}
	}
	return s
}

type LogicalExpression struct {
	Left, Right Expression
	Operator    LogicalOperator
}

func (e LogicalExpression) String() string {
	return fmt.Sprintf("%v %s %v", e.Left, e.Operator, e.Right)
}

type ValuePath struct {
	AttributePath AttributePath
	ValueFilter   Expression
}

func (e ValuePath) String() string {
	return fmt.Sprintf("%v[%v]", e.AttributePath, e.ValueFilter)
}

type NotExpression struct {
	Expression Expression
}

func (e NotExpression) String() string {
	return fmt.Sprintf("not(%v)", e.Expression)
}

// AttributePath represents an attribute path. Both URIPrefix and SubAttr are
// optional values and can be nil.
// e.g. urn:ietf:params:scim:schemas:core:2.0:User:name.givenName
//      ^                                          ^    ^
//      URIPrefix                                  |    SubAttribute
//                                                 AttributeName
type AttributePath struct {
	URIPrefix     *string
	AttributeName string
	SubAttribute  *string
}

func (p AttributePath) String() string {
	s := p.AttributeName
	if p.URIPrefix != nil {
		s = fmt.Sprintf("%s:%s", p.URI(), s)
	}
	if p.SubAttribute != nil {
		s = fmt.Sprintf("%s.%s", s, p.SubAttributeName())
	}
	return s
}

// URI returns the URI is present. Also removes the trailing ':'.
// Returns an empty string otherwise.
func (p *AttributePath) URI() string {
	if p.URIPrefix != nil {
		return *p.URIPrefix
	}
	return ""
}

// SubAttributeName returns the sub attribute name is present.
// Returns an empty string otherwise.
func (p *AttributePath) SubAttributeName() string {
	if p.SubAttribute != nil {
		return *p.SubAttribute
	}
	return ""
}

// Path describes the target of a PATCH operation. Path can have an optional
// ValueExpression and SubAttribute.
// e.g. members[value eq "2819c223-7f76-453a-919d-413861904646"].displayName
//      ^       ^                                                ^
//      |       ValueExpression                                  SubAttribute
//      AttributePath
type Path struct {
	AttributePath   AttributePath
	ValueExpression Expression
	SubAttribute    *string
}

func (p Path) String() string {
	s := p.AttributePath.String()
	if p.ValueExpression != nil {
		s += fmt.Sprintf("[%s]", p.ValueExpression)
	}
	if p.SubAttribute != nil {
		s += fmt.Sprintf(".%s", *p.SubAttribute)
	}
	return s
}
