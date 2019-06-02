package scim

import (
	"testing"
)

func TestAlpha(t *testing.T) {
	for _, r := range "abcdefghijklmnopqrstuvwxyz" {
		if !alpha(r) {
			t.Errorf("rune should be an alpha value: %s, %d", string(r), r)
		}
	}

	for _, r := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if !alpha(r) {
			t.Errorf("rune should be an alpha value: %s, %d", string(r), r)
		}
	}

	for _, r := range "@[`{" {
		if alpha(r) {
			ErrorRune(t, false, "alpha", r)
		}
	}
}

func TestDigit(t *testing.T) {
	for _, r := range "0123456789" {
		if !digit(r) {
			t.Errorf("rune should be a digit value: %s, %d", string(r), r)
		}
	}

	for _, r := range "/:" {
		if digit(r) {
			ErrorRune(t, false, "digit", r)
		}
	}
}

func TestNameChar(t *testing.T) {
	for _, r := range "azAZ09-_" {
		if !nameChar(r) {
			ErrorRune(t, true, "name char", r)
		}
	}

	for _, r := range " ÿ\n" {
		if nameChar(r) {
			ErrorRune(t, false, "name char", r)
		}
	}
}

func ErrorRune(t *testing.T, equal bool, kind string, r rune) {
	if equal {
		t.Errorf("rune should be a(n) %s value: %s, %d", kind, string(r), r)
	} else {
		t.Errorf("rune should NOT be a(n) %s value: %s, %d", kind, string(r), r)
	}
}

func TestAttrName(t *testing.T) {
	for _, s := range []string{"attrName"} {
		if !AttrName(s) {
			ErrorString(t, true, "attribute name", s)
		}
	}

	for _, s := range []string{"attribute.Name", "0attrName"} {
		if AttrName(s) {
			ErrorString(t, false, "attribute name", s)
		}
	}
}

func TestSubAttr(t *testing.T) {
	for _, s := range []string{".subAttr"} {
		if !subAttr(s) {
			ErrorString(t, true, "sub attribute", s)
		}
	}

	for _, s := range []string{"subAttr", "..subAttr", ".0subAttr"} {
		if subAttr(s) {
			ErrorString(t, false, "sub attribute", s)
		}
	}
}

func TestAttrPath(t *testing.T) {
	for _, s := range []string{
		"urn:ietf:params:scim:schemas:core:2.0:User:name",
		"urn:ietf:params:scim:schemas:core:2.0:User:name.givenName",
		"id",
		"name.givenName",
	} {
		if !attrPath(s) {
			ErrorString(t, true, "attribute path", s)
		}
	}

	for _, s := range []string{
		"urn:ietf:params:scim:schemas:core:2.0:User",
		"name.given.name",
	} {
		if attrPath(s) {
			ErrorString(t, false, "attribute path", s)
		}
	}
}

func TestCompareOp(t *testing.T) {
	for _, s := range []string{"eq", "ne", "co", "sw", "ew", "gt", "lt", "ge", "le"} {
		if !compareOp(s) {
			t.Errorf("string should be a compare operator: %s", s)
		}
	}

	if compareOp("pr") {
		t.Errorf("string should NOT be a compare operator: pr")
	}
}

func TestBoolean(t *testing.T) {
	for _, s := range []string{"true", "false"} {
		if !boolean(s) {
			ErrorString(t, true, "boolean", s)
		}
	}

	for _, s := range []string{"maybe"} {
		if boolean(s) {
			ErrorString(t, false, "boolean", s)
		}
	}
}

func TestNull(t *testing.T) {
	for _, s := range []string{"null"} {
		if !null(s) {
			ErrorString(t, true, "null", s)
		}
	}

	for _, s := range []string{"nil", "Null"} {
		if null(s) {
			ErrorString(t, false, "null", s)
		}
	}
}

func TestNumber(t *testing.T) {
	for _, s := range []string{"1", "-1", "1e7"} {
		if !number(s) {
			ErrorString(t, true, "number", s)
		}
	}

	for _, s := range []string{"inf", "--5"} {
		if number(s) {
			ErrorString(t, false, "number", s)
		}
	}
}

func TestString(t *testing.T) {
	for _, s := range []string{"\"\"", "\"string\""} {
		if !word(s) {
			ErrorString(t, true, "number", s)
		}
	}

	for _, s := range []string{"string"} {
		if word(s) {
			ErrorString(t, false, "number", s)
		}
	}
}

func TestCompValue(t *testing.T) {
	for _, s := range []string{"true", "null", "\"\\\"\"", "2e7", "-5"} {
		if !compValue(s) {
			ErrorString(t, true, "comparison", s)
		}
	}

	for _, s := range []string{"\""} {
		if compValue(s) {
			ErrorString(t, false, "comparison", s)
		}
	}
}

func TestSP(t *testing.T) {
	if !space(" ") {
		ErrorString(t, true, "space", "\" \"")
	}

	if space("\t") {
		ErrorString(t, false, "space", "\\t")
	}
}

func TestAttrExp(t *testing.T) {
	for _, s := range []string{"name.givenName pr", "name.givenName eq \"givenName\""} {
		if !attrExp(s) {
			ErrorString(t, true, "attribute expression", s)
		}
	}

	for _, s := range []string{"name.givenName pr ", "name.givenName eq givenName"} {
		if attrExp(s) {
			ErrorString(t, false, "attribute expression", s)
		}
	}
}

func TestFilter(t *testing.T) {
	for _, s := range []string{
		"name.givenName pr", "name.givenName eq \"givenName\"",
		"not(name.givenName pr) and name.givenName eq \"givenName\"",
		"emails[type eq true]",
		"not(name.givenName pr)",
	} {
		if !Filter(s) {
			ErrorString(t, true, "filter", s)
		}
	}

	for _, s := range []string{
		"name.givenName pr ", "name.givenName eq givenName",
		"not  name.givenName pr and name.givenName eq \"givenName\"",
		"emails[type eq maybe]",
		"not  name.givenName pr",
	} {
		if Filter(s) {
			ErrorString(t, false, "filter", s)
		}
	}
}

func TestValuePath(t *testing.T) {
	for _, s := range []string{
		"emails[type eq \"work\"]",
		"emails[type eq true]",
	} {
		if !valuePath(s) {
			ErrorString(t, true, "value path", s)
		}
	}

	for _, s := range []string{
		"emails[type pr \"work\"]",
		"emails[type eq work]",
	} {
		if valuePath(s) {
			ErrorString(t, false, "value path", s)
		}
	}
}

func TestValueFilter(t *testing.T) {
	for _, s := range []string{
		"name.givenName pr", "name.givenName eq \"givenName\"",
		"not(name.givenName pr) and name.givenName eq \"givenName\"",
		"not(not(name.givenName pr) and name.givenName eq \"givenName\")",
	} {
		if !valueFilter(s) {
			ErrorString(t, true, "value filter", s)
		}
	}

	for _, s := range []string{
		"name.givenName pr ", "name.givenName eq givenName",
		"not  name.givenName pr and name.givenName eq \"givenName\"",
		"not not(name.givenName pr) and name.givenName eq \"givenName\")",
	} {
		if valueFilter(s) {
			ErrorString(t, false, "value filter", s)
		}
	}
}

func TestValueLogExp(t *testing.T) {
	for _, s := range []string{
		"name.givenName eq \"givenName\" and name.givenName eq \"givenName\"",
	} {
		if !valueLogExp(s) {
			ErrorString(t, true, "value filter", s)
		}
	}
}

func TestLogExp(t *testing.T) {
	for _, s := range []string{
		"name.givenName pr and name.givenName eq \"givenName\"",
	} {
		if !logExp(s) {
			ErrorString(t, true, "logical expression", s)
		}
	}

	for _, s := range []string{
		"name.givenName pr andor name.givenName eq givenName",
		"name[givenName eq \"givenName\" and name[familyName eq \"familyName\"]]",
	} {
		if logExp(s) {
			ErrorString(t, false, "logical expression", s)
		}
	}
}

func ErrorString(t *testing.T, equal bool, kind string, s string) {
	if equal {
		t.Errorf("string should be a(n) %s value: %s", kind, s)
	} else {
		t.Errorf("string should NOT be a(n) %s value: %s", kind, s)
	}
}
