package filter

import (
	"fmt"
	"strings"
	"testing"
)

func ExampleParseFilter_and() {
	ast, _ := ParseFilter([]byte("title pr and userType eq \"Employee\""))

	fmt.Println(ast)
	// Output:
	// title pr and userType eq "Employee"
}

func ExampleParseFilter_attrExp() {
	fmt.Println(ParseFilter([]byte("schemas eq \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"")))
	// Output:
	// schemas eq "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User" <nil>
}

func ExampleParseFilter_caseInsensitivity() {
	fmt.Println(ParseFilter([]byte("NAME PR AND NOT (FIRST EQ \"test\") AND ANOTHER NE \"test\"")))
	// Output:
	// NAME pr and not(FIRST eq "test") and ANOTHER ne "test" <nil>
}

func ExampleParseFilter_not() {
	fmt.Println(ParseFilter([]byte("not (emails co \"example.com\" or emails.value co \"example.org\")")))
	// Output:
	// not(emails co "example.com" or emails.value co "example.org") <nil>
}

func ExampleParseFilter_or() {
	fmt.Println(ParseFilter([]byte("title pr or userType eq \"Intern\"")))
	// Output:
	// title pr or userType eq "Intern" <nil>
}

func ExampleParseFilter_parentheses() {
	fmt.Println(ParseFilter([]byte("(emails.type eq \"work\")")))
	// Output:
	// (emails.type eq "work") <nil>
}

func ExampleParseFilter_valuePath() {
	fmt.Println(ParseFilter([]byte("emails[type eq \"work\" and value co \"@example.com\"]")))
	// Output:
	// emails[type eq "work" and value co "@example.com"] <nil>
}

func Example_walk() {
	expression, _ := ParseFilter([]byte("emails[type eq \"work\" and value co \"@example.com\"] or ims[type eq \"xmpp\" and value co \"@foo.com\"]"))
	var walk func(e Expression) error
	walk = func(e Expression) error {
		switch v := e.(type) {
		case *LogicalExpression:
			_ = walk(v.Left)
			_ = walk(v.Right)
		case *ValuePath:
			_ = walk(v.ValueFilter)
		case *AttributeExpression:
			fmt.Printf("%s %s %q\n", v.AttributePath, v.Operator, v.CompareValue)
		default:
			// etc...
		}
		return nil
	}
	_ = walk(expression)
	// Output:
	// type eq "work"
	// value co "@example.com"
	// type eq "xmpp"
	// value co "@foo.com"
}

func TestParseFilter(t *testing.T) {
	for _, example := range []string{
		"userName eq \"bjensen\"",
		"userName Eq \"bjensen\"",
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

		"name pr and userName pr and title pr",
		"name pr and not (first eq \"test\") and another ne \"test\"",
		"NAME PR AND NOT (FIRST EQ \"test\") AND ANOTHER NE \"test\"",
		"name pr or userName pr or title pr",
	} {
		t.Run(example, func(t *testing.T) {

			if strings.HasPrefix(example, "userType") {
				fmt.Println("Test:" + example)
			}
			ast, err := ParseFilter([]byte(example))
			fmt.Println(ast)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
