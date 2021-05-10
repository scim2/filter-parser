package filter

import (
	"encoding/json"
	"fmt"
	"github.com/di-wu/parser/ast"
	"github.com/scim2/filter-parser/v2/internal/grammar"
	"strings"
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
				t.Fatal(err)
			}
			{ // Empty config.
				i, err := config{}.parseNumber(n)
				if err != nil {
					t.Fatal(err)
				}
				if i != test.expected {
					t.Error(test.expected, i)
				}
			}
			{ // Config with useNumber = true.
				d := json.NewDecoder(strings.NewReader(test.nStr))
				d.UseNumber()
				var number json.Number
				if err := d.Decode(&number); err != nil {
					t.Error(err)
				}

				i, err := config{
					useNumber: true,
				}.parseNumber(n)
				if err != nil {
					t.Fatal(err)
				}
				if i != json.Number(test.nStr) {
					t.Error(test.nStr, i)
				}

				// Check if equal to json.Decode.
				if i != number {
					t.Error(number, i)
				}
			}
		})
	}
}
