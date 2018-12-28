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
