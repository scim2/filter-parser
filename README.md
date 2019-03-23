[![Go Report Card](https://goreportcard.com/badge/github.com/di-wu/scim-filter-parser)](https://goreportcard.com/report/github.com/di-wu/scim-filter-parser)
[![GoDoc](https://godoc.org/github.com/di-wu/scim-filter-parser?status.svg)](https://godoc.org/github.com/di-wu/scim-filter-parser)
# Query Filter Parser for SCIM v2.0
links: [[scim](http://www.simplecloud.info/#Specification)] [[filtering](https://tools.ietf.org/html/rfc7644#section-3.4.2.2)] \
other: [[pratt parser](https://en.wikipedia.org/wiki/Pratt_parser)] [[json](http://json.org/)]
## Implemented Operators
### Attribute Operators
- [x] eq, ne, co, sw, ew, gt, ge, lt, le
- [x] pr

### Logical Operators
- [x] and, or
- [x] not
- [x] precedence

### Grouping Operators
- [x] ( )
- [x] [ ]

## Case Sensitivity
Attribute names and attribute operators used in filters are case insensitive.  
For example, the following two expressions will evaluate to the same logical value:

```
filter=userName Eq "john"
filter=Username eq "john"
```

## Expressions Requirements
Each expression MUST contain an attribute name followed by an attribute operator and optional value.

Multiple expressions MAY be combined using logical operators.

Expressions MAY be grouped together using round brackets "(" and ")".

Filters MUST be evaluated using the following order of operations, in order of precedence:

   1.  Grouping operators
   2.  Logical operators - where "not" takes precedence over "and",
       which takes precedence over "or"
   3.  Attribute operators
