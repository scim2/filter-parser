package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleFilterOr() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("title pr or userType eq \"Intern\"")
	p("emails[type eq \"work\" and value co \"@example.com\"] or ims[type eq \"xmpp\" and value co \"@foo.com\"]")
	// Output:
	// [001] [[002] [[005] [[006] [[007] title]]], [002] [[005] [[006] [[007] userType], [008] eq, [023] "Intern"]]] <nil>
	// [001] [[002] [[009] [[006] [[007] emails], [011] [[005] [[006] [[007] type], [008] eq, [023] "work"], [005] [[006] [[007] value], [008] co, [023] "@example.com"]]]], [002] [[009] [[006] [[007] ims], [011] [[005] [[006] [[007] type], [008] eq, [023] "xmpp"], [005] [[006] [[007] value], [008] co, [023] "@foo.com"]]]]] <nil>
}

func ExampleFilterAnd() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("title pr and userType eq \"Employee\"")
	p("userType eq \"Employee\" and emails[type eq \"work\" and value co \"@example.com\"]")
	// Output:
	// [001] [[002] [[005] [[006] [[007] title]], [005] [[006] [[007] userType], [008] eq, [023] "Employee"]]] <nil>
	// [001] [[002] [[005] [[006] [[007] userType], [008] eq, [023] "Employee"], [009] [[006] [[007] emails], [011] [[005] [[006] [[007] type], [008] eq, [023] "work"], [005] [[006] [[007] value], [008] co, [023] "@example.com"]]]]] <nil>
}

func ExampleFilterNot() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("userType ne \"Employee\" and not (emails co \"example.com\" or emails.value co \"example.org\")")
	// Output:
	// [001] [[002] [[005] [[006] [[007] userType], [008] ne, [023] "Employee"], [003] [[001] [[002] [[005] [[006] [[007] emails], [008] co, [023] "example.com"]], [002] [[005] [[006] [[007] emails, [007] value], [008] co, [023] "example.org"]]]]]] <nil>
}

func ExampleFilterParentheses() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("userType eq \"Employee\" and (emails.type eq \"work\")")
	p("userType eq \"Employee\" and (emails co \"example.com\" or emails.value co \"example.org\")")
	// Output:
	// [001] [[002] [[005] [[006] [[007] userType], [008] eq, [023] "Employee"], [001] [[002] [[005] [[006] [[007] emails, [007] type], [008] eq, [023] "work"]]]]] <nil>
	// [001] [[002] [[005] [[006] [[007] userType], [008] eq, [023] "Employee"], [001] [[002] [[005] [[006] [[007] emails], [008] co, [023] "example.com"]], [002] [[005] [[006] [[007] emails, [007] value], [008] co, [023] "example.org"]]]]] <nil>
}
