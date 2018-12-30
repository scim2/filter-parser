package scim_filtering

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
			expr: ValueExpression{
				Name:     "username",
				Operator: EQ,
				Value:    "john",
			},
		},
		{
			s: `userName Eq "john"`,
			expr: ValueExpression{
				Name:     "username",
				Operator: EQ,
				Value:    "john",
			},
		},
		{
			s: `name.formatted eq "john doe"`,
			expr: ValueExpression{
				Name:     "name.formatted",
				Operator: EQ,
				Value:    "john doe",
			},
		},

		// other operators
		{
			s: `username ne "john"`,
			expr: ValueExpression{
				Name:     "username",
				Operator: NE,
				Value:    "john",
			},
		},
		{
			s: `name.familyName co "doe"`,
			expr: ValueExpression{
				Name:     "name.familyname",
				Operator: CO,
				Value:    "doe",
			},
		},
		{
			s: `urn:ietf:params:scim:schemas:core:2.0:User:userName sw "j"`,
			expr: ValueExpression{
				Name:     "username",
				Operator: SW,
				Value:    "j",
			},
		},
		{
			s: `username ew "n"`,
			expr: ValueExpression{
				Name:     "username",
				Operator: EW,
				Value:    "n",
			},
		},
		{
			s: `title pr`,
			expr: ValueExpression{
				Name:     `title`,
				Operator: PR,
			},
		},


		// empty value
		{
			s: `Username eq`,
			expr: ValueExpression{
				Name:     "username",
				Operator: EQ,
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
				Operator: NOT,
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.com",
				},
			},
		},

		// and operator
		{
			s: `emails co "example.com" and emails co "example.org"`,
			expr: BinaryExpression{
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.com",
				},
				Operator: AND,
				Y: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.org",
				},
			},
		},

		// or operator
		{
			s: `emails co "example.com" or emails co "example.org"`,
			expr: BinaryExpression{
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.com",
				},
				Operator: OR,
				Y: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.org",
				},
			},
		},

		// precedence
		{
			s: `emails co "example.com" and emails co "example.org" or emails co "example.be"`,
			expr: BinaryExpression{
				X: BinaryExpression{
					X: ValueExpression{
						Name:     "emails",
						Operator: CO,
						Value:    "example.com",
					},
					Operator: AND,
					Y: ValueExpression{
						Name:     "emails",
						Operator: CO,
						Value:    "example.org",
					},
				},
				Operator: OR,
				Y: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.be",
				},
			},
		},
		{
			s: `emails co "example.be" or emails co "example.com" and emails co "example.org"`,
			expr: BinaryExpression{
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.be",
				},
				Operator: OR,
				Y: BinaryExpression{
					X: ValueExpression{
						Name:     "emails",
						Operator: CO,
						Value:    "example.com",
					},
					Operator: AND,
					Y: ValueExpression{
						Name:     "emails",
						Operator: CO,
						Value:    "example.org",
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

func TestParser_GroupingOperators(t *testing.T) {
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
					X: ValueExpression{
						Name:     "emails",
						Operator: CO,
						Value:    "example.be",
					},
					Operator: OR,
					Y: ValueExpression{
						Name:     "emails",
						Operator: CO,
						Value:    "example.com",
					},
				},
				Operator: AND,
				Y: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.org",
				},
			},
		},
		{
			s: `emails co "example.org" and not (emails co "example.be" or emails co "example.com")`,
			expr: BinaryExpression{
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.org",
				},
				Operator: AND,
				Y: UnaryExpression{
					X: BinaryExpression{
						X: ValueExpression{
							Name:     "emails",
							Operator: CO,
							Value:    "example.be",
						},
						Operator: OR,
						Y: ValueExpression{
							Name:     "emails",
							Operator: CO,
							Value:    "example.com",
						},
					},
					Operator: NOT,
				},
			},
		},
		{
			s: `emails co "example.com" and (emails co "example.org")`,
			expr: BinaryExpression{
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.com",
				},
				Operator: AND,
				Y: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    "example.org",
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
