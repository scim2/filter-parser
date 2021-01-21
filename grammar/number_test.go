package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleNumber() {
	p, _ := ast.New([]byte("-10.0e+01"))
	fmt.Println(Number(p))
	// Output:
	// [016] [[017] -, [022] 10, [021] [[020] 0], [018] [[019] +, [020] 01]] <nil>
}
