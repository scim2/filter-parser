package filter

import (
	"fmt"
	"testing"
)

func TestAST(t *testing.T) {
	expression := BinaryExpression{
		X: AttributeExpression{
			AttributePath:   "emails",
			CompareOperator: CO,
			CompareValue:    ".com",
		},
		CompareOperator: OR,
		Y: BinaryExpression{
			X: AttributeExpression{
				AttributePath:   "emails",
				CompareOperator: CO,
				CompareValue:    ".org",
			},
			CompareOperator: AND,
			Y: UnaryExpression{
				X: AttributeExpression{
					AttributePath:   "emails",
					CompareOperator: CO,
					CompareValue:    ".be",
				},
				CompareOperator: NOT,
			},
		},
	}

	expected := "('emails contains .com' or ('emails contains .org' and not 'emails contains .be'))"
	actual := fmt.Sprintf("%s", expression)

	if actual != expected {
		t.Errorf("expected %s got %s", expected, actual)
	}
}
