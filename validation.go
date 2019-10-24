package filter

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

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

// AttrName  = ALPHA *(nameChar)
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
	split := strings.Split(s, ":")
	if len(split) > 1 && len(split) < 9 {
		return false
	}

	s = split[len(split)-1]
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
	switch strings.ToLower(s) {
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
	case 0, 1:
		return false
	case 2:
		return attrPath(split[0]) && split[1] == "pr"
	default:
		return attrPath(split[0]) && compareOp(split[1]) && compValue(strings.Join(split[2:], " "))
	}
}

// Path = attrPath / valuePath [subAttr]
func isPath(s string) bool {
	if attrPath(s) || valuePath(s) {
		return true
	}
	idx := strings.LastIndex(s, ".")
	if idx < 1 {
		return false
	}
	first := s[:idx]
	last := s[idx:]
	return attrPath(first) && subAttr(last) || valuePath(first) && subAttr(last)
}

// isFilter = attrExp / logExp / valuePath / *1"not" "(" FILTER ")"
func isFilter(s string) bool {
	if strings.HasPrefix(s, "not") {
		return isFilter(strings.TrimPrefix(s[3:], " "))
		//return isFilter(strings.TrimLeft(s[3:], " "))
	}
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return isFilter(s[1 : len(s)-1])
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

// valFilter = attrExp / valLogExp / *1"not" "(" valFilter ")"
func valueFilter(s string) bool {
	if strings.HasPrefix(s, "not") && strings.HasSuffix(s, ")") {
		return valueFilter(s[3:])
	}
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		return valueFilter(s[1 : len(s)-1])
	}
	return attrExp(s) || valueLogExp(s) || logExp(s)
}

// valLogExp = attrExp SP ("and" / "or") SP attrExp
func valueLogExp(s string) bool {
	first, last := splitOnLog(s)
	if first == "" || last == "" {
		return false
	}
	return attrExp(first) && attrExp(last)
}

// logExp = FILTER SP ("and" / "or") SP FILTER
func logExp(s string) bool {
	first, last := splitOnLog(s)
	if first == "" || last == "" {
		return false
	}
	return isFilter(first) && isFilter(last)
}

func splitOnLog(s string) (string, string) {
	// regex for (...) and [...]
	brackets := regexp.MustCompile(`(\(.*?\))|(\[.*?\])`)
	inside := brackets.FindAllString(s, -1)

	// replace all strings inside brackets w/ "%"
	for _, v := range inside {
		s = strings.Replace(s, v, "%", 1)
	}

	// split on first and/or not inside ()/[]
	r := regexp.MustCompile(` and | or `)
	splits := r.Split(s, 2)

	if len(splits) != 2 {
		return "", ""
	}

	// repopulate string
	idx := 0
	for i, v := range splits {
		count := strings.Count(v, "%")
		for count > 0 {
			splits[i] = strings.Replace(v, "%", inside[idx], 1)
			count = count - 1
			idx = idx + 1
		}
	}

	return splits[0], splits[1]
}
