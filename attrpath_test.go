package filter

import "fmt"

func ExampleParseAttrPath() {
	fmt.Println(ParseAttrPath([]byte("urn:ietf:params:scim:schemas:core:2.0:User:name.familyName")))
	// Output:
	// urn:ietf:params:scim:schemas:core:2.0:User:name.familyName <nil>
}
