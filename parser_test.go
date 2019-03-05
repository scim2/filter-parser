package scim

import (
	"reflect"
	"strings"
	"testing"
)

func TestParser_AttributeOperators(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		// eq operator
		{
			s: `userName Eq "john"`,
			expr: AttributeExpression{
				AttributePath:   "username",
				CompareOperator: EQ,
				CompareValue:    "john",
			},
		},
		{
			s: `userName Eq "john"`,
			expr: AttributeExpression{
				AttributePath:   "username",
				CompareOperator: EQ,
				CompareValue:    "john",
			},
		},
		{
			s: `name.formatted eq "john doe"`,
			expr: AttributeExpression{
				AttributePath:   "name.formatted",
				CompareOperator: EQ,
				CompareValue:    "john doe",
			},
		},

		// other operators
		{
			s: `username ne "john"`,
			expr: AttributeExpression{
				AttributePath:   "username",
				CompareOperator: NE,
				CompareValue:    "john",
			},
		},
		{
			s: `name.familyName co "doe"`,
			expr: AttributeExpression{
				AttributePath:   "name.familyname",
				CompareOperator: CO,
				CompareValue:    "doe",
			},
		},
		{
			s: `urn:ietf:params:scim:schemas:core:2.0:User:userName sw "j"`,
			expr: AttributeExpression{
				AttributePath:   "username",
				CompareOperator: SW,
				CompareValue:    "j",
			},
		},
		{
			s: `username ew "n"`,
			expr: AttributeExpression{
				AttributePath:   "username",
				CompareOperator: EW,
				CompareValue:    "n",
			},
		},
		{
			s: `title pr`,
			expr: AttributeExpression{
				AttributePath:   `title`,
				CompareOperator: PR,
			},
		},

		// empty value
		{
			s: `Username eq`,
			expr: AttributeExpression{
				AttributePath:   "username",
				CompareOperator: EQ,
			},
		},

		// invalid operator
		{
			s:   `error x "value"`,
			err: `found "x", expected operator`,
		},
	}

	for i, test := range tests {
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParser_LogicalOperators(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		// not operator
		{
			s: `not emails co "example.com"`,
			expr: UnaryExpression{
				CompareOperator: NOT,
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
			},
		},

		// and operator
		{
			s: `emails co "example.com" and emails co "example.org"`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
				CompareOperator: AND,
				Y: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
			},
		},

		// or operator
		{
			s: `emails co "example.com" or emails co "example.org"`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
				CompareOperator: OR,
				Y: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
			},
		},

		// precedence
		{
			s: `emails co "example.com" and emails co "example.org" or emails co "example.be"`,
			expr: BinaryExpression{
				X: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails",
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath:   "emails",
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
				CompareOperator: OR,
				Y: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.be",
				},
			},
		},
		{
			s: `emails co "example.be" or emails co "example.com" and emails co "example.org"`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.be",
				},
				CompareOperator: OR,
				Y: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails",
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath:   "emails",
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
			},
		},
	}

	for i, test := range tests {
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParser_GroupingOperators_Parenthesis(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		// parenthesis
		{
			s: `(emails co "example.be" or emails co "example.com") and emails co "example.org"`,
			expr: BinaryExpression{
				X: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails",
						CompareOperator: CO,
						CompareValue:    "example.be",
					},
					CompareOperator: OR,
					Y: AttributeExpression{
						AttributePath:   "emails",
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
				},
				CompareOperator: AND,
				Y: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
			},
		},
		{
			s: `emails co "example.org" and not (emails co "example.be" or emails co "example.com")`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
				CompareOperator: AND,
				Y: UnaryExpression{
					X: BinaryExpression{
						X: AttributeExpression{
							AttributePath:   "emails",
							CompareOperator: CO,
							CompareValue:    "example.be",
						},
						CompareOperator: OR,
						Y: AttributeExpression{
							AttributePath:   "emails",
							CompareOperator: CO,
							CompareValue:    "example.com",
						},
					},
					CompareOperator: NOT,
				},
			},
		},
		{
			s: `emails co "example.com" and (emails co "example.org")`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
				CompareOperator: AND,
				Y: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
			},
		},
		{
			s: `emails co "example.org" and (emails.type eq "work" and emails.value co "example.org")`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
				CompareOperator: AND,
				Y: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails.type",
						CompareOperator: EQ,
						CompareValue:    "work",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath:   "emails.value",
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
			},
		},
	}

	for i, test := range tests {
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParser_GroupingOperators_Brackets(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		// brackets
		{
			s: `emails co "example.org" and emails[type eq "work" and value co "example.org"]`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
				CompareOperator: AND,
				Y: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails.type",
						CompareOperator: EQ,
						CompareValue:    "work",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath:   "emails.value",
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
			},
		},
		{
			s: `emails[type eq "work" and value co "example.org"] or emails[type eq "private" and value co "example.com"]`,
			expr: BinaryExpression{
				X: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails.type",
						CompareOperator: EQ,
						CompareValue:    "work",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath:   "emails.value",
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
				CompareOperator: OR,
				Y: BinaryExpression{
					X: AttributeExpression{
						AttributePath:   "emails.type",
						CompareOperator: EQ,
						CompareValue:    "private",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath:   "emails.value",
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
				},
			},
		},
	}

	for i, test := range tests {
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func errToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
