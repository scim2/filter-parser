package grammar

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	"github.com/scim2/filter-parser/v2/internal/types"
)

// URI parses a URN as defined in RFC 8141 Section 2.
// https://datatracker.ietf.org/doc/html/rfc8141#section-2
//
// Each segment between colons may contain alphanumeric characters,
// hyphens, and periods:
//   - NID allows alphanum and "-" (ldh production).
//   - NSS allows pchar (RFC 3986), which includes unreserved chars
//     (ALPHA / DIGIT / "-" / "." / "_" / "~"), but in practice SCIM
//     schema URNs only use alphanum, "-", and ".".
func URI(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ.URI,
		TypeStrings: typ.Stringer,
		Value: op.MinOne(op.And{
			op.MinOne(op.Or{
				parser.CheckRuneRange('a', 'z'),
				parser.CheckRuneRange('A', 'Z'),
				parser.CheckRuneRange('0', '9'),
				'-',
				'.',
			}),
			":",
		}),
	})
}
