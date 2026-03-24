package filter

import "fmt"

func ExampleParseAttrPath() {
	fmt.Println(ParseAttrPath([]byte("urn:ietf:params:scim:schemas:core:2.0:User:name.familyName")))
	// Output:
	// urn:ietf:params:scim:schemas:core:2.0:User:name.familyName <nil>
}

func ExampleParseAttrPath_dash() {
	fmt.Println(ParseAttrPath([]byte("urn:example:scim:schemas:extension:my-custom-ext:1.0:User:name.familyName")))
	// Output:
	// urn:example:scim:schemas:extension:my-custom-ext:1.0:User:name.familyName <nil>
}
