package scim_filtering

type Token int

const (
	UNKNOWN Token = iota

	WS   // whitespace (i.e. " ")
	ID   // identifier (i.e. "userName")
	EOF  // end-of-file
	V    // value (i.e. "john")

	// Attribute Operators
	EQ  // equal
	NE  // not equal
	CO  // contains
	SW  // starts with
	EW  // ends with
	PR  // present (has value)
	GT  // greater than
	GE  // greater than or equal to
	LT  // less than
	LE  // less than or equal to

	// Logical Operators
	AND  // logical "and
	OR   // logical "or"
	NOT  // "not" function
)
