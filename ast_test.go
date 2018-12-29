package scim_filtering

import (
	"fmt"
	"testing"
)

func TestAST(t *testing.T) {
	expression := BinaryExpression{
		X: ValueExpression{
			Name:     "emails",
			Operator: CO,
			Value:    ".com",
		},
		Operator: OR,
		Y: BinaryExpression{
			X: ValueExpression{
				Name:     "emails",
				Operator: CO,
				Value:    ".org",
			},
			Operator: AND,
			Y: UnaryExpression{
				X: ValueExpression{
					Name:     "emails",
					Operator: CO,
					Value:    ".be",
				},
				Operator: NOT,
			},
		},
	}

	expected := "('emails contains .com' or ('emails contains .org' and not 'emails contains .be'))"
	actual := fmt.Sprintf("%s", expression)

	if actual != expected {
		t.Errorf("expected %s got %s", expected, actual)
	}
}
