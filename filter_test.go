package filter

import (
	"fmt"
	"testing"
)

func ExampleParseFilter_attrExp() {
	fmt.Println(ParseFilter([]byte("schemas eq \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"")))
	// Output:
	// schemas eq "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User" <nil>
}

func ExampleParseFilter_valuePath() {
	fmt.Println(ParseFilter([]byte("emails[type eq \"work\" and value co \"@example.com\"]")))
	// Output:
	// emails[type eq "work" and value co "@example.com"] <nil>
}

func ExampleParseFilter_not() {
	fmt.Println(ParseFilter([]byte("not (emails co \"example.com\" or emails.value co \"example.org\")")))
	// Output:
	// not(emails co "example.com" or emails.value co "example.org") <nil>
}

func ExampleParseFilter_and() {
	fmt.Println(ParseFilter([]byte("title pr and userType eq \"Employee\"")))
	// Output:
	// title pr and userType eq "Employee" <nil>
}

func ExampleParseFilter_or() {
	fmt.Println(ParseFilter([]byte("title pr or userType eq \"Intern\"")))
	// Output:
	// title pr or userType eq "Intern" <nil>
}

func ExampleParseFilter_parentheses() {
	fmt.Println(ParseFilter([]byte("(emails.type eq \"work\")")))
	// Output:
	// emails.type eq "work" <nil>
}

func TestParseFilter(t *testing.T) {
	for _, example := range []string{
		"userName eq \"bjensen\"",
		"name.familyName co \"O'Malley\"",
		"userName sw \"J\"",
		"urn:ietf:params:scim:schemas:core:2.0:User:userName sw \"J\"",
		"title pr",
		"meta.lastModified gt \"2011-05-13T04:42:34Z\"",
		"meta.lastModified ge \"2011-05-13T04:42:34Z\"",
		"meta.lastModified lt \"2011-05-13T04:42:34Z\"",
		"meta.lastModified le \"2011-05-13T04:42:34Z\"",
		"title pr and userType eq \"Employee\"",
		"title pr or userType eq \"Intern\"",
		"schemas eq \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"",
		"userType eq \"Employee\" and (emails co \"example.com\" or emails.value co \"example.org\")",
		"userType ne \"Employee\" and not (emails co \"example.com\" or emails.value co \"example.org\")",
		"userType eq \"Employee\" and (emails.type eq \"work\")",
		"userType eq \"Employee\" and emails[type eq \"work\" and value co \"@example.com\"]",
		"emails[type eq \"work\" and value co \"@example.com\"] or ims[type eq \"xmpp\" and value co \"@foo.com\"]",
	} {
		t.Run(example, func(t *testing.T) {
			if _, err := ParseFilter([]byte(example)); err != nil {
				t.Error(err)
			}
		})
	}
}
