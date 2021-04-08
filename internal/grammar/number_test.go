package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleNumber() {
	p, _ := ast.New([]byte("-10.0e+01"))
	fmt.Println(Number(p))
	// Output:
	// ["Number",[["Minus","-"],["Int","10"],["Frac",[["Digits","0"]]],["Exp",[["Sign","+"],["Digits","01"]]]]] <nil>
}
