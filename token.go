package scim_filtering

type Token int

const (
	UNKNOWN Token = iota

	WS   // whitespace (i.e. " ")
	ID   // identifier (i.e. "userName")
	EOF  // end-of-file
	V    // value (i.e. "john")

	EQ  // equal
)
