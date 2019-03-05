package scim

import (
	"encoding/json"
	"strconv"
	"strings"
)

var URI = "urn:ietf:params:scim:schemas:core:2.0"

// ALPHA = A–Z, a–z
func alpha(r rune) bool {
	i := int32(r)
	return (65 <= i && i <= 90) || (97 <= i && i <= 122)
}

// DIGIT = 0-9
func digit(r rune) bool {
	i := int32(r)
	return 48 <= i && i <= 57
}

// nameChar = "-" / "_" / DIGIT / ALPHA
func nameChar(r rune) bool {
	return r == '-' || r == '_' || digit(r) || alpha(r)
}

// ATTRNAME  = ALPHA *(nameChar)
func AttrName(s string) bool {
	for idx, r := range s {
		if idx == 0 && !alpha(r) {
			return false
		}

		if !nameChar(r) {
			return false
		}
	}
	return true
}

// subAttr   = "." ATTRNAME
func subAttr(s string) bool {
	if s[0] != '.' {
		return false
	}
	return AttrName(s[1:])
}

// attrPath  = [URI ":"] ATTRNAME *1subAttr
func attrPath(s string) bool {
	s = strings.TrimPrefix(s, URI)
	s = strings.TrimPrefix(s, ":")
	if len(s) < 1 {
		return false
	}

	for idx, attributeName := range strings.Split(s, ".") {
		if !AttrName(attributeName) {
			return false
		}

		// one sub attribute
		if idx > 1 {
			return false
		}
	}
	return true
}

// compValue = false / null / true / number / string
func compareOp(s string) bool {
	switch s {
	case "eq", "ne", "co", "sw", "ew", "gt", "lt", "ge", "le":
		return true
	default:
		return false
	}
}

func boolean(s string) bool {
	if _, err := strconv.ParseBool(s); err != nil {
		return false
	}
	return true
}

func null(s string) bool {
	if s != "null" {
		return false
	}
	return true
}

func number(s string) bool {
	var number json.Number
	if err := json.Unmarshal([]byte(s), &number); err != nil {
		return false
	}
	return true
}

// word = string
func word(s string) bool {
	if s[0] != '"' && s[len(s)-1] != '"' {
		return false
	}
	return json.Unmarshal([]byte(s), &s) == nil
}

// compValue = false / null / true / number / string
func compValue(s string) bool {
	return boolean(s) || null(s) || number(s) || word(s)
}

// space = SP
func space(s string) bool {
	return s == " "
}

// attrExp = (attrPath SP "pr") / (attrPath SP compareOp SP compValue)
func attrExp(s string) bool {
	split := strings.Split(s, " ")
	switch len(split) {
	case 2:
		return attrPath(split[0]) && split[1] == "pr"
	case 3:
		return attrPath(split[0]) && compareOp(split[1]) && compValue(split[2])
	default:
		return false
	}
}

// FILTER = attrExp / logExp / valuePath / *1"not" "(" FILTER ")"
func Filter(s string) bool {
	if strings.HasPrefix(s, "not") {
		return Filter(s[3:])
	}
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return Filter(s[1 : len(s)-1])
	}
	return attrExp(s) || logExp(s) || valuePath(s)
}

// valuePath = attrPath "[" valFilter "]"
func valuePath(s string) bool {
	if !strings.HasSuffix(s, "]") {
		return false
	}
	s = s[:len(s)-1]
	split := strings.Split(s, "[")
	switch len(split) {
	case 2:
		return attrPath(split[0]) && valueFilter(split[1])
	default:
		return false
	}
}

// valFilter = attrExp / logExp / *1"not" "(" valFilter ")"
func valueFilter(s string) bool {
	if strings.HasPrefix(s, "not") && strings.HasSuffix(s, ")") {
		return valueFilter(s[3:])
	}
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return valueFilter(s[1 : len(s)-1])
	}
	return attrExp(s) || logExp(s)
}

// logExp = FILTER SP ("and" / "or") SP FILTER
func logExp(s string) bool {
	split := strings.Split(s, " and ")
	var double []string
	for _, s := range split {
		for _, elem := range strings.Split(s, " or ") {
			double = append(double, elem)
		}
	}
	if len(double) <= 1 {
		return false
	}

	for _, s := range double {
		if !Filter(s) {
			return false
		}
	}
	return true
}
