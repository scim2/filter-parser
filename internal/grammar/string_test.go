package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleString() {
	p, _ := ast.New([]byte("\"2819c223-7f76-453a-919d-413861904646\""))
	fmt.Println(String(p))
	// Output:
	// ["String","\"2819c223-7f76-453a-919d-413861904646\""] <nil>
}

func ExampleString_complex() {
	p, _ := ast.New([]byte("\"W/\\\"990-6468886345120203448\\\"\""))
	fmt.Println(String(p))
	// Output:
	// ["String","\"W/\\\"990-6468886345120203448\\\"\""] <nil>
}
