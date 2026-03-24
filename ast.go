package filter

import (
	"fmt"
)

const (
	// PR is an abbreviation for 'present'.
	PR CompareOperator = "pr"
	// EQ is an abbreviation for 'equals'.
	EQ CompareOperator = "eq"
	// NE is an abbreviation for 'not equals'.
	NE CompareOperator = "ne"
	// CO is an abbreviation for 'contains'.
	CO CompareOperator = "co"
	// SW is an abbreviation for 'starts with'.
	SW CompareOperator = "sw"
	// EW an abbreviation for 'ends with'.
	EW CompareOperator = "ew"
	// GT is an abbreviation for 'greater than'.
	GT CompareOperator = "gt"
	// LT is an abbreviation for 'less than'.
	LT CompareOperator = "lt"
	// GE is an abbreviation for 'greater or equal than'.
	GE CompareOperator = "ge"
	// LE is an abbreviation for 'less or equal than'.
	LE CompareOperator = "le"

	// AND is the logical operation and (&&).
	AND LogicalOperator = "and"
	// OR is the logical operation or (||).
	OR LogicalOperator = "or"
)

// AttributeExpression represents an attribute expression/filter.
type AttributeExpression struct {
	AttributePath AttributePath
	Operator      CompareOperator
	CompareValue  any
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

func (*AttributeExpression) exprNode() {}

// AttributePath represents an attribute path with an optional URIPrefix and
// SubAttribute.
//
// Example: urn:ietf:params:scim:schemas:core:2.0:User:name.givenName
//   - URIPrefix: urn:ietf:params:scim:schemas:core:2.0:User
//   - AttributeName: name
//   - SubAttribute: givenName
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

// SubAttributeName returns the sub attribute name if present.
// Returns an empty string otherwise.
func (p *AttributePath) SubAttributeName() string {
	if p.SubAttribute != nil {
		return *p.SubAttribute
	}
	return ""
}

// URI returns the URI if present.
// Returns an empty string otherwise.
func (p *AttributePath) URI() string {
	if p.URIPrefix != nil {
		return *p.URIPrefix
	}
	return ""
}

// CompareOperator represents a compare operation.
type CompareOperator string

// Expression is a type to assign to implemented expressions.
// Valid expressions are:
//   - ValuePath
//   - AttributeExpression
//   - LogicalExpression
//   - NotExpression
type Expression interface {
	exprNode()
}

// LogicalExpression represents an 'and' / 'or' node.
type LogicalExpression struct {
	Left, Right Expression
	Operator    LogicalOperator
}

func (e LogicalExpression) String() string {
	left := fmt.Sprintf("%v", e.Left)
	if e.Operator == AND {
		if l, ok := e.Left.(*LogicalExpression); ok && l.Operator == OR {
			left = fmt.Sprintf("(%v)", e.Left)
		}
	}
	right := fmt.Sprintf("%v", e.Right)
	if e.Operator == AND {
		if r, ok := e.Right.(*LogicalExpression); ok && r.Operator == OR {
			right = fmt.Sprintf("(%v)", e.Right)
		}
	}
	return fmt.Sprintf("%s %s %s", left, e.Operator, right)
}

func (*LogicalExpression) exprNode() {}

// LogicalOperator represents a logical operation such as 'and' / 'or'.
type LogicalOperator string

// NotExpression represents an 'not' node.
type NotExpression struct {
	Expression Expression
}

func (e NotExpression) String() string {
	return fmt.Sprintf("not(%v)", e.Expression)
}

func (*NotExpression) exprNode() {}

// Path describes the target of a PATCH operation with an optional
// ValueExpression and SubAttribute.
//
// Example: members[value eq "2819c223-..."].displayName
//   - AttributePath: members
//   - ValueExpression: value eq "2819c223-..."
//   - SubAttribute: displayName
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

// SubAttributeName returns the sub attribute name if present.
// Returns an empty string otherwise.
func (p *Path) SubAttributeName() string {
	if p.SubAttribute != nil {
		return *p.SubAttribute
	}
	return ""
}

// ValuePath represents a filter on a attribute path.
type ValuePath struct {
	AttributePath AttributePath
	ValueFilter   Expression
}

func (e ValuePath) String() string {
	return fmt.Sprintf("%v[%v]", e.AttributePath, e.ValueFilter)
}

func (*ValuePath) exprNode() {}
