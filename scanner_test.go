package scim_filtering

import (
	"strings"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	var tests = []struct {
		s       string
		token   Token
		literal string
	}{
		// special tokens
		{s: ``, token: EOF},
		{s: ` `, token: WS, literal: " "},
		{s: "\n", token: WS, literal: "\n"},
		{s: "\t", token: WS, literal: "\t"},

		// unknown tokens
		{s: `#`, token: UNKNOWN, literal: `#`},

		// identifiers
		{s: `identifier`, token: ID, literal: `identifier`},
		{s: `identifier.id`, token: ID, literal: `identifier.id`},

		{s: `eq`, token: EQ, literal: "eq"},
		{s: `Eq`, token: EQ, literal: "eq"},
		{s: `EQ`, token: EQ, literal: "eq"},
		{s: `eQ`, token: EQ, literal: "eq"},

		{s: `ne`, token: NE, literal: "ne"},
		{s: `co`, token: CO, literal: "co"},
		{s: `sw`, token: SW, literal: "sw"},
		{s: `ew`, token: EW, literal: "ew"},
		{s: `pr`, token: PR, literal: "pr"},
		{s: `gt`, token: GT, literal: "gt"},
		{s: `ge`, token: GE, literal: "ge"},
		{s: `lt`, token: LT, literal: "lt"},
		{s: `le`, token: LE, literal: "le"},

		{s: "and", token: AND, literal: "and"},
		{s: "or", token: OR, literal: "or"},
		{s: "not", token: NOT, literal: "not"},

		{s: "(", token: LPAR, literal: "("},
		{s: ")", token: RPAR, literal: ")"},
		{s: "[", token: LBRA, literal: "["},
		{s: "]", token: RBRA, literal: "]"},

		// values
		{s: `"john"`, token: V, literal: `john`},
	}

	for i, test := range tests {
		s := NewScanner(strings.NewReader(test.s))
		token, literal := s.Scan()
		if test.token != token {
			t.Errorf("%d. %q wrong token: exp=%q got=%q <%q>", i, test.s, test.token, token, literal)
		} else if test.literal != literal {
			t.Errorf("%d. %q wrong token: exp=%q got=%q", i, test.s, test.literal, literal)
		}
	}
}
