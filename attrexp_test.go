package filter

import (
	"fmt"
	"github.com/di-wu/parser/ast"
	"github.com/scim2/filter-parser/v2/grammar"
	"testing"
)

func ExampleParseAttrExp_pr() {
	fmt.Println(ParseAttrExp([]byte("userName pr")))
	// Output:
	// userName pr <nil>
}

func ExampleParseAttrExp_sw() {
	fmt.Println(ParseAttrExp([]byte("userName sw \"J\"")))
	// Output:
	// userName sw "J" <nil>
}

func TestParseNumber(t *testing.T) {
	for _, test := range []struct {
		nStr     string
		expected interface{}
	}{
		{
			nStr:     "-5.1e-2",
			expected: -0.051,
		},
		{
			nStr:     "-5.1e2",
			expected: float64(-510),
		},
		{
			nStr:     "-510",
			expected: -510,
		},
	} {
		t.Run(test.nStr, func(t *testing.T) {
			p, _ := ast.New([]byte(test.nStr))
			n, err := grammar.Number(p)
			if err != nil {
				t.Error(err)
				return
			}
			i, err := parseNumber(n)
			if err != nil {
				t.Error(err)
				return
			}
			if i != test.expected {
				t.Error(test.expected, i)
			}
		})
	}
}
