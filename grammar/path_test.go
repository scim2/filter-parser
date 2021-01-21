package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExamplePath() {
	p, _ := ast.New([]byte("members[value eq \"2819c223-7f76-453a-919d-413861904646\"].displayName"))
	fmt.Println(Path(p))
	// Output:
	// [004] [[009] [[006] [[007] members], [005] [[006] [[007] value], [008] eq, [023] "2819c223-7f76-453a-919d-413861904646"]], [007] displayName] <nil>
}
