package filter

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePath(t *testing.T) {
	var tests = []struct {
		s    string
		expr Path
		err  string
	}{
		{
			s: `urn:ietf:params:scim:schemas:core:2.0:User:userName`,
			expr: Path{
				URIPrefix:     "urn:ietf:params:scim:schemas:core:2.0:User",
				AttributeName: "userName",
			},
		},
		{
			s:   `.familyName`,
			err: `found ".", expected identifier`,
		},
		{
			s: `members`,
			expr: Path{
				AttributeName: "members",
			},
		},
		{
			s:   `members.`,
			err: `found "members.", sub attribute can not be empty`,
		},
		{
			s: `name.familyName`,
			expr: Path{
				AttributeName: "name",
				SubAttribute:  "familyName",
			},
		},
		{
			s: `addresses[type eq "work"]`,
			expr: Path{
				AttributeName: "addresses",
				ValueExpression: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "type",
					},
					CompareOperator: EQ,
					CompareValue:    "work",
				},
			},
		},
		{
			s:   `addresses[type eq "work"]x`,
			err: `found "x", expected '.' or eof`,
		},
		{
			s: `members[value eq "id"].displayName`,
			expr: Path{
				AttributeName: "members",
				SubAttribute:  "displayName",
				ValueExpression: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "value",
					},
					CompareOperator: EQ,
					CompareValue:    "id",
				},
			},
		},
		{
			s:   `members.displayName[value eq "id"]`,
			err: `found "[", expected eof`,
		},
	}

	for i, test := range tests {
		if test.err == "" && !isPath(test.s) {
			t.Errorf("invalid path: %s", test.s)
		}
		expr, err := NewParser(strings.NewReader(test.s)).ParsePath()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParse_AttributeOperators(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		// eq operator
		{
			s: `userName Eq "john"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "userName",
				},
				CompareOperator: EQ,
				CompareValue:    "john",
			},
		},
		{
			s: `userName Eq "john"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "userName",
				},
				CompareOperator: EQ,
				CompareValue:    "john",
			},
		},
		{
			s: `name.formatted eq "john doe"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "name",
					SubAttribute:  "formatted",
				},
				CompareOperator: EQ,
				CompareValue:    "john doe",
			},
		},

		// other operators
		{
			s: `userName ne "john"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "userName",
				},
				CompareOperator: NE,
				CompareValue:    "john",
			},
		},
		{
			s: `name.familyName co "doe"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "name",
					SubAttribute:  "familyName",
				},
				CompareOperator: CO,
				CompareValue:    "doe",
			},
		},
		{
			s: `urn:ietf:params:scim:schemas:core:2.0:User:userName sw "j"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					URIPrefix:     "urn:ietf:params:scim:schemas:core:2.0:User",
					AttributeName: "userName",
				},
				CompareOperator: SW,
				CompareValue:    "j",
			},
		},
		{
			s: `username ew "n"`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "username",
				},
				CompareOperator: EW,
				CompareValue:    "n",
			},
		},
		{
			s: `title pr`,
			expr: AttributeExpression{
				AttributePath: AttributePath{
					AttributeName: "title",
				},
				CompareOperator: PR,
			},
		},

		// invalid operator
		{
			s:   `error x "value"`,
			err: `found "x", expected operator`,
		},
	}

	for i, test := range tests {
		if test.err == "" && !isFilter(test.s) {
			t.Errorf("invalid filter: %s", test.s)
		}
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParse_SubAttributes(t *testing.T) {
	var tests = []struct {
		s    string
		expr Expression
		err  string
	}{
		// invalid operator
		{
			s:   `emails.x.y pr`,
			err: `found "emails.x.y", no multiple sub attributes allowed`,
		},
	}

	for i, test := range tests {
		if test.err == "" && !isFilter(test.s) {
			t.Errorf("invalid filter: %s", test.s)
		}
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParse_LogicalOperators(t *testing.T) {
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
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
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
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
				CompareOperator: AND,
				Y: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
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
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
				CompareOperator: OR,
				Y: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
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
						AttributePath: AttributePath{
							AttributeName: "emails",
						},
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath: AttributePath{
							AttributeName: "emails",
						},
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
				CompareOperator: OR,
				Y: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.be",
				},
			},
		},
		{
			s: `emails co "example.be" or emails co "example.com" and emails co "example.org"`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.be",
				},
				CompareOperator: OR,
				Y: BinaryExpression{
					X: AttributeExpression{
						AttributePath: AttributePath{
							AttributeName: "emails",
						},
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath: AttributePath{
							AttributeName: "emails",
						},
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
			},
		},
	}

	for i, test := range tests {
		if test.err == "" && !isFilter(test.s) {
			t.Errorf("invalid filter: %s", test.s)
		}
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParse_GroupingOperators_Parenthesis(t *testing.T) {
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
						AttributePath: AttributePath{
							AttributeName: "emails",
						},
						CompareOperator: CO,
						CompareValue:    "example.be",
					},
					CompareOperator: OR,
					Y: AttributeExpression{
						AttributePath: AttributePath{
							AttributeName: "emails",
						},
						CompareOperator: CO,
						CompareValue:    "example.com",
					},
				},
				CompareOperator: AND,
				Y: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
			},
		},
		{
			s: `emails co "example.org" and not (emails co "example.be" or emails co "example.com")`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
				CompareOperator: AND,
				Y: UnaryExpression{
					X: BinaryExpression{
						X: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "emails",
							},
							CompareOperator: CO,
							CompareValue:    "example.be",
						},
						CompareOperator: OR,
						Y: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "emails",
							},
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
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.com",
				},
				CompareOperator: AND,
				Y: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
			},
		},
		{
			s: `emails co "example.org" and (emails.type eq "work" and emails.value co "example.org")`,
			expr: BinaryExpression{
				X: AttributeExpression{
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
				CompareOperator: AND,
				Y: BinaryExpression{
					X: AttributeExpression{
						AttributePath: AttributePath{
							AttributeName: "emails",
							SubAttribute:  "type",
						},
						CompareOperator: EQ,
						CompareValue:    "work",
					},
					CompareOperator: AND,
					Y: AttributeExpression{
						AttributePath: AttributePath{
							AttributeName: "emails",
							SubAttribute:  "value",
						},
						CompareOperator: CO,
						CompareValue:    "example.org",
					},
				},
			},
		},
	}

	for i, test := range tests {
		if test.err == "" && !isFilter(test.s) {
			t.Errorf("invalid filter: %s", test.s)
		}
		expr, err := NewParser(strings.NewReader(test.s)).Parse()
		if !reflect.DeepEqual(test.err, errToString(err)) {
			t.Errorf("%d. %q: wrong error:\n exp=%s\n got=%s\n\n", i, test.s, test.err, err)
		} else if test.err == "" && !reflect.DeepEqual(test.expr, expr) {
			t.Errorf("%d. %q: wrong expr:\n exp=%s\n got=%s\n\n", i, test.s, test.expr, expr)
		}
	}
}

func TestParse_GroupingOperators_Brackets(t *testing.T) {
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
					AttributePath: AttributePath{
						AttributeName: "emails",
					},
					CompareOperator: CO,
					CompareValue:    "example.org",
				},
				CompareOperator: AND,
				Y: ValuePath{
					AttributeName: "emails",
					ValueExpression: BinaryExpression{
						X: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "type",
							},
							CompareOperator: EQ,
							CompareValue:    "work",
						},
						CompareOperator: AND,
						Y: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "value",
							},
							CompareOperator: CO,
							CompareValue:    "example.org",
						},
					},
				},
			},
		},
		{
			s: `emails[type eq "work" and value co "example.org"] or emails[type eq "private" and value co "example.com"]`,
			expr: BinaryExpression{
				X: ValuePath{
					AttributeName: "emails",
					ValueExpression: BinaryExpression{
						X: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "type",
							},
							CompareOperator: EQ,
							CompareValue:    "work",
						},
						CompareOperator: AND,
						Y: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "value",
							},
							CompareOperator: CO,
							CompareValue:    "example.org",
						},
					},
				},
				CompareOperator: OR,
				Y: ValuePath{
					AttributeName: "emails",
					ValueExpression: BinaryExpression{
						X: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "type",
							},
							CompareOperator: EQ,
							CompareValue:    "private",
						},
						CompareOperator: AND,
						Y: AttributeExpression{
							AttributePath: AttributePath{
								AttributeName: "value",
							},
							CompareOperator: CO,
							CompareValue:    "example.com",
						},
					},
				},
			},
		},
	}

	for i, test := range tests {
		if test.err == "" && !isFilter(test.s) {
			t.Errorf("invalid filter: %s", test.s)
		}
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
