package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleURI() {
	p, _ := ast.New([]byte("urn:ietf:params:scim:schemas:core:2.0:User:userName"))
	fmt.Println(URI(p))
	// Output:
	// ["URI","urn:ietf:params:scim:schemas:core:2.0:User:"] <nil>
}

func ExampleURI_dash() {
	p, _ := ast.New([]byte("urn:example:scim:schemas:extension:my-custom-ext:1.0:User:userName"))
	fmt.Println(URI(p))
	// Output:
	// ["URI","urn:example:scim:schemas:extension:my-custom-ext:1.0:User:"] <nil>
}
