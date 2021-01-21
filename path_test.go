package filter

import (
	"fmt"
	"testing"
)

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

func TestParsePath(t *testing.T) {
	for _, example := range []string{
		"members",
		"name.familyName",
		"addresses[type eq \"work\"]",
		"members[value eq \"2819c223-7f76-453a-919d-413861904646\"]",
		"members[value eq \"2819c223-7f76-453a-919d-413861904646\"].displayName",
	} {
		t.Run(example, func(t *testing.T) {
			if _, err := ParsePath([]byte(example)); err != nil {
				t.Error(err)
			}
		})
	}
}
