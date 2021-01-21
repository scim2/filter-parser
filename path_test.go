package filter

import "fmt"

func ExampleParsePath_attrPath() {
	fmt.Println(ParsePath([]byte("members")))
	fmt.Println(ParsePath([]byte("name.familyName")))
	// Output:
	// members <nil>
	// name.familyName <nil>
}

func ExampleParsePath_valuePath() {
	fmt.Println(ParsePath([]byte("members[value eq \"2819c223-7f76-453a-919d-413861904646\"]")))
	fmt.Println(ParsePath([]byte("members[value eq \"2819c223-7f76-453a-919d-413861904646\"].displayName")))
	// Output:
	// members[value eq "2819c223-7f76-453a-919d-413861904646"] <nil>
	// members[value eq "2819c223-7f76-453a-919d-413861904646"].displayName <nil>
}
