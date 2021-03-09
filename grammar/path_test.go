package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExamplePath() {
	p, _ := ast.New([]byte("members[value eq \"2819c223-7f76-453a-919d-413861904646\"].displayName"))
	fmt.Println(Path(p))
	// Output:
	// ["Path",[["ValuePath",[["AttrPath",[["AttrName","members"]]],["AttrExp",[["AttrPath",[["AttrName","value"]]],["CompareOp","eq"],["String","\"2819c223-7f76-453a-919d-413861904646\""]]]]],["AttrName","displayName"]]] <nil>
}
