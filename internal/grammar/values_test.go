package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleFalse() {
	p, _ := ast.New([]byte("FaLSe"))
	fmt.Println(False(p))
	// Output:
	// ["False","FaLSe"] <nil>
}

func ExampleTrue() {
	p, _ := ast.New([]byte("TRue"))
	fmt.Println(True(p))
	// Output:
	// ["True","TRue"] <nil>
}
