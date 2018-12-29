# Query Filter Parser for SCIM v2.0
[[scim](http://www.simplecloud.info/#Specification)] [[filtering](https://tools.ietf.org/html/rfc7644#section-3.4.2.2)]

## Implemented Operators
### Attribute Operators
- [x] eq
- [x] ne
- [x] co
- [x] sw
- [x] ew
- [x] pr
- [x] gt
- [x] ge
- [x] lt
- [x] le

### Logical Operators
- [x] and
- [x] or
- [x] not
- [ ] precedence

### Grouping Operators
- [ ] ( )
- [ ] [ ]

## Case Sensitivity
Attribute names and attribute operators used in filters are case insensitive.  
For example, the following two expressions will evaluate to the same logical value:

```json
filter=userName Eq "john"
filter=Username eq "john"
```

## Expressions Requirements
Each expression MUST contain an attribute name followed by
an attribute operator and optional value.

Multiple expressions MAY be combined using logical operators.

Expressions MAY be grouped together using round brackets "(" and ")".

Filters MUST be evaluated using the following order of operations, in
   order of precedence:

   1.  Grouping operators

   2.  Logical operators - where "not" takes precedence over "and",
       which takes precedence over "or"

   3.  Attribute operators
