package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleURI() {
	p, _ := ast.New([]byte("urn:ietf:params:scim:schemas:core:2.0:User:userName"))
	fmt.Println(URI(p))
	// Output:
	// [024] urn:ietf:params:scim:schemas:core:2.0:User: <nil>
}
