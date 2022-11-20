package grammar

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleFilterAnd() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("title pr and userType eq \"Employee\"")
	p("userType eq \"Employee\" and emails[type eq \"work\" and value co \"@example.com\"]")
	// Output:
	// ["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","title"]]]]],["AttrExp",[["AttrPath",[["AttrName","userType"]]],["CompareOp","eq"],["String","\"Employee\""]]]]]]] <nil>
	// ["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","userType"]]],["CompareOp","eq"],["String","\"Employee\""]]],["ValuePath",[["AttrPath",[["AttrName","emails"]]],["ValueLogExpAnd",[["AttrExp",[["AttrPath",[["AttrName","type"]]],["CompareOp","eq"],["String","\"work\""]]],["AttrExp",[["AttrPath",[["AttrName","value"]]],["CompareOp","co"],["String","\"@example.com\""]]]]]]]]]]] <nil>
}

func ExampleFilterNot() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("userType ne \"Employee\" and not (emails co \"example.com\" or emails.value co \"example.org\")")
	// Output:
	// ["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","userType"]]],["CompareOp","ne"],["String","\"Employee\""]]],["FilterNot",[["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","emails"]]],["CompareOp","co"],["String","\"example.com\""]]]]],["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","emails"],["AttrName","value"]]],["CompareOp","co"],["String","\"example.org\""]]]]]]]]]]]]] <nil>
}

func ExampleFilterOr() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("title pr or userType eq \"Intern\"")
	p("emails[type eq \"work\" and value co \"@example.com\"] or ims[type eq \"xmpp\" and value co \"@foo.com\"]")
	// Output:
	// ["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","title"]]]]]]],["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","userType"]]],["CompareOp","eq"],["String","\"Intern\""]]]]]]] <nil>
	// ["FilterOr",[["FilterAnd",[["ValuePath",[["AttrPath",[["AttrName","emails"]]],["ValueLogExpAnd",[["AttrExp",[["AttrPath",[["AttrName","type"]]],["CompareOp","eq"],["String","\"work\""]]],["AttrExp",[["AttrPath",[["AttrName","value"]]],["CompareOp","co"],["String","\"@example.com\""]]]]]]]]],["FilterAnd",[["ValuePath",[["AttrPath",[["AttrName","ims"]]],["ValueLogExpAnd",[["AttrExp",[["AttrPath",[["AttrName","type"]]],["CompareOp","eq"],["String","\"xmpp\""]]],["AttrExp",[["AttrPath",[["AttrName","value"]]],["CompareOp","co"],["String","\"@foo.com\""]]]]]]]]]]] <nil>
}

func ExampleFilterParentheses() {
	p := func(s string) {
		p, _ := ast.New([]byte(s))
		fmt.Println(Filter(p))
	}
	p("userType eq \"Employee\" and (emails.type eq \"work\")")
	p("userType eq \"Employee\" and (emails co \"example.com\" or emails.value co \"example.org\")")
	// Output:
	// ["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","userType"]]],["CompareOp","eq"],["String","\"Employee\""]]],["FilterPrecedence",[["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","emails"],["AttrName","type"]]],["CompareOp","eq"],["String","\"work\""]]]]]]]]]]]]] <nil>
	// ["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","userType"]]],["CompareOp","eq"],["String","\"Employee\""]]],["FilterPrecedence",[["FilterOr",[["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","emails"]]],["CompareOp","co"],["String","\"example.com\""]]]]],["FilterAnd",[["AttrExp",[["AttrPath",[["AttrName","emails"],["AttrName","value"]]],["CompareOp","co"],["String","\"example.org\""]]]]]]]]]]]]] <nil>
}
