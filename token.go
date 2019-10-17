package filter

// Token is an int type representing tokens.
type Token int

const (
	// UNKNOWN is an unknown token
	UNKNOWN Token = iota

	// WHITESPACE is a whitespace (i.e. " ")
	WHITESPACE
	// IDENTIFIER is an identifier (i.e. "userName")
	IDENTIFIER
	// EOF represents an end-of-file
	EOF
	// VALUE is a value (i.e. "john")
	VALUE

	// EQ = equals
	EQ
	// NE = not equal
	NE
	// CO = contains
	CO
	// SW = starts with
	SW
	// EW = ends with
	EW
	// PR = present (has value)
	PR
	// GT = greater than
	GT
	// GE = greater than or equal to
	GE
	// LT = less than
	LT
	// LE = less than or equal to
	LE

	// AND = logical "and"
	AND
	// OR = logical "or"
	OR
	// NOT = "not" function
	NOT

	// LeftParenthesis = "("
	LeftParenthesis
	// RightParenthesis = ")"
	RightParenthesis
	// LeftBracket = "["
	LeftBracket
	// RightBracket = "]"
	RightBracket
)

// string representation of the tokens.
var tokens = [...]string{
	UNKNOWN: "unknown",

	WHITESPACE: " ",
	IDENTIFIER: "id",
	EOF:        "",
	VALUE:      "value",

	EQ: "equal",
	NE: "not equal",
	CO: "contains",
	SW: "starts with",
	EW: "ends with",
	PR: "present",
	GT: ">",
	GE: ">=",
	LT: "<",
	LE: "<=",

	AND: "and",
	OR:  "or",
	NOT: "not",

	LeftParenthesis:  "(",
	RightParenthesis: ")",
	LeftBracket:      "[",
	RightBracket:     "]",
}

// IsOperator returns whether the token is an operator.
func (token Token) IsOperator() bool {
	switch token {
	case EQ, NE, CO, SW, EW, PR, GT, GE, LT, LE:
		return true
	}
	return false
}

// IsAssociative return whether the token in associative.
func (token Token) IsAssociative() bool {
	switch token {
	case AND, OR:
		return true
	}
	return false
}

const (
	// highestPrecedence is the highest precedence of a token (highest integer)
	highestPrecedence = 2
	// lowestPrecedence is the lowest precedence of a token (lowest integer)
	lowestPrecedence = 0
)

// Precedence returns the precedence value of the token.
func (token Token) Precedence() int {
	switch token {
	case AND:
		return 2
	case OR:
		return 1
	}
	return 0
}

func (token Token) String() string {
	if 0 <= token && token < Token(len(tokens)) {
		return tokens[token]
	}
	return tokens[0]
}
